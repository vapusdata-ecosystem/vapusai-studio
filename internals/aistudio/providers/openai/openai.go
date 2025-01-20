package openaistd

import (
	"context"
	"encoding/json"
	"errors"
	"html"
	"io"
	"log"

	"github.com/rs/zerolog"
	openai "github.com/sashabaranov/go-openai"
	aicore "github.com/vapusdata-oss/aistudio/core/aistudio/core"
	"github.com/vapusdata-oss/aistudio/core/aistudio/prompts"
	dmlogger "github.com/vapusdata-oss/aistudio/core/logger"
	"github.com/vapusdata-oss/aistudio/core/models"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

var (
	OpenAILLMMap = map[string]string{
		"gpt-4o":      openai.GPT4o,
		"gpt-4":       openai.GPT4,
		"gpt-4o-mini": openai.GPT4oMini,
	}
	OpenAIEmbeddingMap = map[string]openai.EmbeddingModel{
		"text-embedding-3-large":  openai.LargeEmbedding3,
		"text-embedding-3-small":  openai.SmallEmbedding3,
		"ttext-embedding-3-small": openai.AdaEmbeddingV2,
	}
)

var ResponseFormatMap = map[string]openai.ChatCompletionResponseFormatType{
	mpb.AIResponseFormat_TEXT.String():        openai.ChatCompletionResponseFormatTypeText,
	mpb.AIResponseFormat_JSON_SCHEMA.String(): openai.ChatCompletionResponseFormatTypeJSONObject,
}

type OpenAIInterface interface {
	GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload, model string) error
	GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	CrawlModels(ctx context.Context) ([]*models.AIModelBase, error)
}

type OpenAI struct {
	client    *openai.Client
	log       zerolog.Logger
	modelNode *models.AIModelNode
	params    map[string]interface{}
}

func New(node *models.AIModelNode, logger zerolog.Logger) OpenAIInterface {
	token := ""
	if node.GetCredentials("default") != nil {
		token = node.GetCredentials("default").ApiToken
	}
	return &OpenAI{
		client:    openai.NewClient(token),
		log:       dmlogger.GetSubDMLogger(logger, "ailogger", "openai"),
		modelNode: node,
	}
}

func (o *OpenAI) GenerateEmbeddings(ctx context.Context, payload *prompts.AIEmbeddingPayload, model string) error {
	o.log.Info().Msgf("Generating embeddings from openai with model %s", model)
	var rModel openai.EmbeddingModel
	rModel, ok := OpenAIEmbeddingMap[model]
	if !ok {
		o.log.Err(nil).Msg("invalid model name in request")
		return aicore.ErrInvalidAIModel
	}
	o.log.Info().Msgf("Generating embeddings from openai with model %s", model)
	resp, err := o.client.CreateEmbeddings(
		ctx,
		openai.EmbeddingRequest{
			Model:      rModel,
			Input:      payload.Input,
			User:       openai.ChatMessageRoleUser,
			Dimensions: payload.Dimensions,
		},
	)
	if err != nil {
		o.log.Err(err).Msg("error while generating embeddings from openai")
		return err
	}
	payload.Embeddings = &models.VectorEmbeddings{
		Vectors32: resp.Data[0].Embedding,
	}
	return nil
}

func (o *OpenAI) buildRequest(payload *prompts.GenerativePrompterPayload) ([]openai.ChatCompletionMessage, []openai.Tool) {
	request := []openai.ChatCompletionMessage{}
	tools := []openai.Tool{}
	if payload.Params.Mode == pb.AIInterfaceMode_CHAT_MODE {
		for _, c := range payload.SessionContext {
			if c.Role == prompts.USER {
				request = append(request, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleUser,
					Content: c.Message,
				})
			} else {
				request = append(request, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleAssistant,
					Content: c.Message,
				})
			}
		}
	}
	for _, c := range payload.Params.Messages {
		mess := openai.ChatCompletionMessage{
			Role:      openai.ChatMessageRoleUser,
			Content:   c.Content,
			ToolCalls: []openai.ToolCall{},
		}
		if c.Role == pb.AIMessageRoles_SYSTEM {
			mess.Role = openai.ChatMessageRoleSystem
		}
		request = append(request, mess)
	}
	if len(payload.ToolCalls) > 0 {
		for _, tool := range payload.ToolCalls {
			log.Println("Tool call++++++++++++++++++++++++++++++++++++++++++++++++++++", tool.FunctionSchema.Name)
			if tool != nil && tool.FunctionSchema != nil {
				paramObj := map[string]interface{}{}
				err := json.Unmarshal([]byte(tool.FunctionSchema.Arguments), &paramObj)
				if err != nil {
					o.log.Err(err).Msg("error while unmarshalling tool call arguments")
					continue
				}
				tools = append(tools, openai.Tool{
					Type: openai.ToolType(tool.Type),
					Function: &openai.FunctionDefinition{
						Name:        tool.FunctionSchema.Name,
						Parameters:  paramObj,
						Description: tool.FunctionSchema.Description,
					},
				})
			}
		}
		// request = append(request, tc)
	}
	return request, tools
}

func (o *OpenAI) GenerateContent(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
	if payload.Params.Model == "" {
		o.log.Err(nil).Msg("invalid model name in request")
		payload.Params.Model = openai.GPT4o
	}
	log.Println("Max tokens", payload.Params.MaxOutputTokens)
	o.log.Info().Msgf("Generating content from openai with model %s", payload.Params.Model)
	messages, tool := o.buildRequest(payload)
	request := openai.ChatCompletionRequest{
		Model:       payload.Params.Model,
		Temperature: payload.Params.Temperature,
		MaxTokens:   int(payload.Params.MaxOutputTokens),
	}
	if len(tool) > 0 {
		request.Tools = tool
	}
	request.Messages = messages
	vbytes, _ := json.MarshalIndent(request, "", "  ")
	log.Println("request------------------->>>>>>>>>>>>>>>>>>>>>|||||||||||||||||||||||||", string(vbytes))
	resp, err := o.client.CreateChatCompletion(
		ctx,
		request,
	)

	if err != nil {
		o.log.Err(err).Msg("error while generating content from openai")
		return err
	}
	for _, c := range resp.Choices {
		if len(c.Message.ToolCalls) > 0 {
			for _, t := range c.Message.ToolCalls {
				payload.ToolCallResponse = append(payload.ToolCallResponse, &mpb.ToolCall{
					Type: string(t.Type),
					FunctionSchema: &mpb.FunctionCall{
						Name:      t.Function.Name,
						Arguments: t.Function.Arguments,
					},
				})
			}
			payload.ParseToolCallResponse()
		}
		if len(c.Message.Content) > 0 {
			payload.ParseOutput(&prompts.PayloadgenericResponse{
				FinishReason: string(resp.Choices[0].FinishReason),
				Data:         c.Message.Content,
				Role:         c.Message.Role,
			})
		}
	}
	return nil
}

func (o *OpenAI) GenerateContentStream(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
	if payload.Params.Model == "" {
		o.log.Err(nil).Msg("invalid model name in request")
		payload.Params.Model = openai.GPT4o
	}
	o.log.Info().Msgf("Generating content from openai stream with model %s", payload.Params.Model)
	messages, tool := o.buildRequest(payload)
	request := openai.ChatCompletionRequest{
		Model:       payload.Params.Model,
		Temperature: payload.Params.Temperature,
		MaxTokens:   int(payload.Params.MaxOutputTokens),
	}
	if len(tool) > 0 {
		request.Tools = tool
	}
	request.Messages = messages
	stream, err := o.client.CreateChatCompletionStream(
		ctx,
		request,
	)

	if err != nil {
		o.log.Err(err).Msg("error while generating content from openai")
		return err
	}
	func() {
		defer stream.Close()
		o.log.Info().Msg("Stream created successfully, reading response")

		for {
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				o.log.Info().Msg("EOF reached, breaking stream")
				_ = payload.SendStreamData("", true, false)
				break
			}
			if err != nil {
				o.log.Err(err).Msg("error while reading stream")
				_ = payload.SendStreamData("", false, true)
				break
			}
			if len(resp.Choices) == 0 {
				continue
			}
			err = payload.SendStreamData(html.UnescapeString(resp.Choices[0].Delta.Content), false, false)
			if err != nil {
				_ = payload.SendStreamData("", false, true)
				continue
			}
		}
	}()
	return nil
}
