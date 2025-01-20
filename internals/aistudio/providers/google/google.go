package googlegenai

import (
	"context"
	"fmt"
	"html"
	"log"

	"github.com/google/generative-ai-go/genai"
	"github.com/invopop/jsonschema"
	"github.com/rs/zerolog"
	aicore "github.com/vapusdata-oss/aistudio/core/aistudio/core"
	"github.com/vapusdata-oss/aistudio/core/aistudio/prompts"
	dmlogger "github.com/vapusdata-oss/aistudio/core/logger"
	"github.com/vapusdata-oss/aistudio/core/models"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var defaultModel = "gemini-2.0-flash"
var defaultEmbeddingModel = "gpt-3.5-turbo"

type GoogleGenAIInterface interface {
	GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload, model string) error
	GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	CrawlModels(ctx context.Context) ([]*models.AIModelBase, error)
}

type GoogleGenAI struct {
	client     *genai.Client
	log        zerolog.Logger
	modelNode  *models.AIModelNode
	maxRetries int
	params     map[string]interface{}
}

func New(ctx context.Context, node *models.AIModelNode, retries int, logger zerolog.Logger) (GoogleGenAIInterface, error) {
	token := ""
	if node.GetCredentials("default") != nil {
		token = node.GetCredentials("default").ApiToken
	}
	genAiCl, err := genai.NewClient(ctx, option.WithAPIKey(token))
	if err != nil {
		logger.Error().Err(err).Msg("Error creating google gen ai client")
		return nil, err
	}
	return &GoogleGenAI{
		client:    genAiCl,
		log:       dmlogger.GetSubDMLogger(logger, "ailogger", "Google Gen AI"),
		modelNode: node,
	}, nil
}

func (x *GoogleGenAI) buildTextRequest(payload *prompts.GenerativePrompterPayload, stream bool) []genai.Part {
	req := []genai.Part{}
	for _, msg := range payload.Params.Messages {
		if msg.Role == pb.AIMessageRoles_USER {
			req = append(req, genai.Text(msg.Content))
		}
	}
	return req
}

func (x *GoogleGenAI) buildToolRequest(payload *prompts.GenerativePrompterPayload, stream bool) []*genai.Tool {
	result := []*genai.Tool{}
	reflector := jsonschema.Reflector{}
	if len(payload.ToolCalls) > 0 {
		for _, toolCall := range payload.ToolCalls {
			tool := &genai.Tool{
				FunctionDeclarations: []*genai.FunctionDeclaration{},
			}
			hSchema := reflector.Reflect(toolCall.FunctionSchema.Arguments)
			ggSchema := &genai.Schema{}
			if hSchema.Type == "object" {
				ggSchema.Type = genai.TypeObject
			} else if hSchema.Type == "array" {
				ggSchema.Type = genai.TypeArray
			} else if hSchema.Type == "string" {
				ggSchema.Type = genai.TypeString
			} else if hSchema.Type == "number" {
				ggSchema.Type = genai.TypeNumber
			} else if hSchema.Type == "boolean" {
				ggSchema.Type = genai.TypeBoolean
			}
			ggSchema.Description = toolCall.FunctionSchema.Description
			ggSchema.Items = &genai.Schema{}
			// ggSchema.Enum = hSchema.Enum
			ggSchema.Required = hSchema.Required
			// ggSchema.Properties = hSchema.Properties
			// TODO: Add support for nested objects in response also
			fc := &genai.FunctionDeclaration{
				Name:        toolCall.FunctionSchema.Name,
				Description: toolCall.FunctionSchema.Description,
				Parameters:  ggSchema,
			}
			tool.FunctionDeclarations = append(tool.FunctionDeclarations, fc)
			result = append(result, tool)
		}
	}

	return result
}

func (x *GoogleGenAI) buildResponse(resp *genai.GenerateContentResponse, payload *prompts.GenerativePrompterPayload, parseOP bool) string {
	if resp == nil {
		return ""
	}
	var result string = ""
	for _, cand := range resp.Candidates {
		lc := ""
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				funcCall, ok := part.(genai.FunctionCall)
				if !ok && funcCall.Name == "" {
					continue
				}
				log.Print("===================================stream response==============================", part)
				result = result + " " + fmt.Sprintf("%v", part)
				lc = lc + " " + fmt.Sprintf("%v", part)

			}
			payload.ParseToolCallResponse()
		}
		if parseOP {
			payload.ParseOutput(&prompts.PayloadgenericResponse{
				FinishReason: cand.FinishReason.String(),
				Data:         result,
				Role:         pb.AIMessageRoles_ASSISTANT.String(),
			})
		}
	}
	return result
}

func (x *GoogleGenAI) buildStreamResponse(resp *genai.GenerateContentResponseIterator, payload *prompts.GenerativePrompterPayload) {
	func() {
		errCounter := 0
		for {
			resp, err := resp.Next()
			if err == iterator.Done {
				x.log.Info().Msg("Stream response done")
				_ = payload.SendStreamData("", true, false)
				break
			}
			if err != nil {
				x.log.Error().Err(err).Msg("Error reading response from stream for google gen AI completion")
				_ = payload.SendStreamData("", false, true)
				if errCounter > 3 {
					x.log.Err(err).Msg("error while reading stream response")
					_ = payload.SendStreamData("", true, false)
					break
				} else {
					errCounter++
					continue
				}
			}
			err = payload.SendStreamData(html.UnescapeString(x.buildResponse(resp, payload, false)), false, false)
			if err != nil {
				x.log.Err(err).Msg("error while sending stream response")
				_ = payload.SendStreamData("", false, true)
				continue
			}
		}
	}()
	return
}

func (x *GoogleGenAI) GenerateEmbeddings(ctx context.Context, payload *prompts.AIEmbeddingPayload, model string) error {
	return aicore.ErrEmbeddingNotSupported
}

func (x *GoogleGenAI) GenerateContent(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
	if payload.Params.Model == "" {
		x.log.Warn().Msg("Model name is empty, using default model")
		payload.Params.Model = defaultModel
	}
	if payload.Params.MaxOutputTokens == 0 {
		payload.Params.MaxOutputTokens = prompts.DefaultMaxOPTokenLength
	}
	modelCl := x.client.GenerativeModel(payload.Params.Model)
	modelCl.SetTemperature(payload.Params.Temperature)
	modelCl.SetTopP(float32(payload.Params.TopP))
	modelCl.SetMaxOutputTokens(payload.Params.MaxOutputTokens)
	tools := x.buildToolRequest(payload, false)
	if len(tools) > 0 {
		modelCl.Tools = tools
	}
	input := x.buildTextRequest(payload, false)
	log.Println("Input to Google Gen AI", input, "mode", payload.Params.Mode, "---------------||||||||||||||||||")
	switch payload.Params.Mode {
	case pb.AIInterfaceMode_CHAT_MODE:
		return x.Chat(ctx, modelCl, payload)
	case pb.AIInterfaceMode_P2P:
		response, err := modelCl.GenerateContent(ctx, input...)
		if err != nil {
			x.log.Error().Err(err).Msg("Error generating content from google gen ai")
			return err
		}
		x.buildResponse(response, payload, true)
	}
	return nil
}

func (x *GoogleGenAI) GenerateContentStream(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
	if payload.Params.Model == "" {
		x.log.Warn().Msg("Model name is empty, using default model")
		payload.Params.Model = defaultModel
	}
	log.Println("payload.Params.MaxOutputTokens google gemini", payload.Params.MaxOutputTokens)
	modelCl := x.client.GenerativeModel(payload.Params.Model)
	modelCl.SetTemperature(payload.Params.Temperature)
	modelCl.SetTopP(float32(payload.Params.TopP))
	modelCl.SetMaxOutputTokens(payload.Params.MaxOutputTokens)
	input := x.buildTextRequest(payload, false)
	tools := x.buildToolRequest(payload, false)
	if len(tools) > 0 {
		modelCl.Tools = tools
	}
	log.Println("Input to Google Gen AI stream", input, "mode", payload.Params.Mode, "---------------||||||||||||||||||")
	switch payload.Params.Mode {
	case pb.AIInterfaceMode_CHAT_MODE:
		return x.ChatStream(ctx, modelCl, payload)
	case pb.AIInterfaceMode_P2P:
		response := modelCl.GenerateContentStream(ctx, input...)
		x.buildStreamResponse(response, payload)
	}
	return nil
}

func (x *GoogleGenAI) Chat(ctx context.Context, modelCl *genai.GenerativeModel, payload *prompts.GenerativePrompterPayload) error {
	modelCl.ResponseMIMEType = "text/plain"
	session := modelCl.StartChat()
	for _, msg := range payload.SessionContext {
		session.History = append(session.History, &genai.Content{
			Parts: []genai.Part{
				genai.Text(msg.Message),
			},
			Role: msg.Role,
		})
	}
	response, err := session.SendMessage(ctx, x.buildTextRequest(payload, false)...)
	if err != nil {
		x.log.Error().Err(err).Msg("Error generating content from google gen ai")
		return err
	}
	log.Println("Output from Google Gen AI", response.PromptFeedback.BlockReason)
	x.buildResponse(response, payload, true)
	return nil
}

func (x *GoogleGenAI) ChatStream(ctx context.Context, modelCl *genai.GenerativeModel, payload *prompts.GenerativePrompterPayload) error {
	modelCl.ResponseMIMEType = "text/plain"
	session := modelCl.StartChat()
	for _, msg := range payload.SessionContext {
		session.History = append(session.History, &genai.Content{
			Parts: []genai.Part{
				genai.Text(msg.Message),
			},
			Role: msg.Role,
		})
	}
	response := session.SendMessageStream(ctx, x.buildTextRequest(payload, false)...)
	x.buildStreamResponse(response, payload)
	return nil
}
