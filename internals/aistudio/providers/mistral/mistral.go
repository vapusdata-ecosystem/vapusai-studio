package mistral

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
	aicore "github.com/vapusdata-oss/aistudio/core/aistudio/core"
	"github.com/vapusdata-oss/aistudio/core/aistudio/prompts"
	httpCls "github.com/vapusdata-oss/aistudio/core/http"
	dmlogger "github.com/vapusdata-oss/aistudio/core/logger"
	"github.com/vapusdata-oss/aistudio/core/models"
	utils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

type MistralInterface interface {
	GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload, model string) error
	GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	CrawlModels(ctx context.Context) ([]*models.AIModelBase, error)
	FIM(ctx context.Context, payload *prompts.GenerativePrompterPayload, model string) error
}

type Mistral struct {
	client     *httpCls.RestHttp
	log        zerolog.Logger
	modelNode  *models.AIModelNode
	maxRetries int
	params     map[string]interface{}
}

func New(node *models.AIModelNode, retries int, logger zerolog.Logger) (MistralInterface, error) {
	token := ""
	if node.GetCredentials("default") != nil {
		token = node.GetCredentials("default").ApiToken
	}
	httpCl, err := httpCls.New(logger,
		httpCls.WithAddress(defaultEndpoint),
		httpCls.WithBasePath(baseAPIPath),
		httpCls.WithBearerAuth(token),
	)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating http client for Mistral")
		return nil, err
	}
	return &Mistral{
		client:    httpCl,
		log:       dmlogger.GetSubDMLogger(logger, "ailogger", "Mistral"),
		modelNode: node,
	}, nil
}

func (x *Mistral) buildRequest(model string, payload *prompts.GenerativePrompterPayload, stream bool) []byte {
	messages := make([]*Message, 0)

	// messages = append(messages, &Message{
	// 	Role:    aicore.ASSISTANT,
	// 	Content: &payload.Prompter.Prompt.As.AssistantMessage,
	// })
	if payload.Params.Mode == pb.AIInterfaceMode_CHAT_MODE {
		for _, msg := range payload.SessionContext {
			if msg.Role == prompts.USER {
				messages = append(messages, &Message{
					Role:    aicore.USER,
					Content: &msg.Message,
				})
			} else {
				messages = append(messages, &Message{
					Role:    aicore.ASSISTANT,
					Content: &msg.Message,
				})
			}
		}
	}
	for _, c := range payload.Params.Messages {
		mess := &Message{
			Role:    aicore.USER,
			Content: &c.Content,
		}
		if c.Role == pb.AIMessageRoles_SYSTEM {
			mess.Role = aicore.SYSTEM
		}
		messages = append(messages, mess)
	}

	tools := make([]Tool, 0)
	if len(payload.ToolCalls) > 0 {
		for _, tool := range payload.ToolCalls {
			if tool != nil && tool.FunctionSchema != nil {
				tools = append(tools, Tool{
					Type: ToolType(getToolType(tool.Type)),
					Function: Function{
						Name:        tool.FunctionSchema.Name,
						Description: tool.FunctionSchema.Description,
						Parameters:  tool.FunctionSchema.Arguments,
					},
				})
			}
		}
	}
	reqObj := &GenerativeRequest{
		Model:       model,
		Messages:    messages,
		Stream:      stream,
		Temperature: float64(payload.Params.Temperature),
		ToolChoice:  "auto",
		MaxTokens:   int(payload.Params.MaxOutputTokens),
		TopP: func() float64 {
			if payload.Params.TopP == 0.0 {
				return defaultTopP
			}
			return payload.Params.TopP
		}(),
	}
	if len(tools) > 0 {
		reqObj.Tools = tools
		reqObj.ToolChoice = "any"
	}
	reqBytes, err := json.Marshal(reqObj)
	if err != nil {
		x.log.Error().Err(err).Msg("Error marshalling request object")
		return nil
	}
	return reqBytes
}

func (x *Mistral) buildFimRequest(model string, payload *prompts.GenerativePrompterPayload) []byte {
	ip := ""
	for _, c := range payload.Params.Messages {
		if c.Role == pb.AIMessageRoles_USER {
			ip = ip + "\n" + c.Content
		}
	}
	req := &FIMRequests{
		Model:       model,
		Prompt:      ip,
		Suffix:      payload.Suffix,
		Temperature: float64(payload.Params.Temperature),
		MaxTokens:   int(payload.Params.MaxOutputTokens),
	}
	bbytes, err := json.Marshal(req)
	if err != nil {
		x.log.Error().Err(err).Msg("Error marshalling request object for Mistral")
		return nil
	}
	return bbytes
}

func (x *Mistral) buildResponse(resp *GenerativeResponse, payload *prompts.GenerativePrompterPayload) error {
	if len(resp.Choices) == 0 {
		x.log.Warn().Msg("No choices found in response")
		return aicore.ErrNoResponseFromAIModel
	}
	for _, choice := range resp.Choices {
		if choice.FinishReason == "tool_calls" {
			for _, tool := range choice.Message.ToolCalls {
				payload.ToolCallResponse = append(payload.ToolCallResponse, &mpb.ToolCall{
					Type: tool.Type.String(),
					FunctionSchema: &mpb.FunctionCall{
						Name:      tool.Function.Name,
						Arguments: tool.Function.Arguments,
					},
				})
			}
			payload.ParseToolCallResponse()
			continue
		} else {
			payload.ParseOutput(&prompts.PayloadgenericResponse{
				FinishReason: string(resp.Choices[0].FinishReason),
				Data:         resp.Object,
				Role:         resp.Choices[0].Message.Role,
			})
		}
	}
	return nil
}

func (x *Mistral) buildStreamResponse(resp *http.Response, payload *prompts.GenerativePrompterPayload) {

	// go func() {
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				_ = payload.SendStreamData("", true, false)
				break
			} else {
				_ = payload.SendStreamData("", false, true)
				break
			}
		}

		if bytes.Equal(line, []byte("\n")) {
			continue
		}

		// Check if the line starts with "data: ".
		if bytes.HasPrefix(line, []byte("data: ")) {
			// Trim the prefix and any leading or trailing whitespace.
			jsonLine := bytes.TrimSpace(bytes.TrimPrefix(line, []byte("data: ")))
			// Check for the special "[DONE]" message.
			if bytes.Equal(jsonLine, []byte("[DONE]")) {
				_ = payload.SendStreamData("", true, false)
				return
			}

			streamResponse := &GenerativeResponseStream{}
			if err := json.Unmarshal(jsonLine, streamResponse); err != nil {
				continue
			}
			if len(streamResponse.Choices) == 0 {
				continue
			}
			err = payload.SendStreamData(streamResponse.Choices[0].Delta.Content, false, false)
			if err != nil {
				x.log.Err(err).Msg("error while sending stream response")
				_ = payload.SendStreamData("", false, true)
				continue
			}
		}
	}
	// return
	// }()
	return
}

func (x *Mistral) GenerateEmbeddings(ctx context.Context, payload *prompts.AIEmbeddingPayload, model string) error {
	if model == "" {
		x.log.Warn().Msg("Model name is empty, using default model")
		model = defaultEmbedModel
	}
	req := &EmbeddingRequest{
		Model: defaultEmbedModel,
		Input: payload.Input,
	}
	bbytes, err := json.Marshal(req)
	if err != nil {
		x.log.Error().Err(err).Msg("Error marshalling request object for Mistral")
	}
	resp := &EmbeddingResponse{}
	err = x.client.Post(ctx, embeddingsPath, bbytes, resp, jsonContentType)
	if err != nil {
		x.log.Error().Err(err).Msg("Error generating content from Mistral completion")
		return err
	}
	vectors := []float64{}
	for _, data := range resp.Data {
		vectors = append(vectors, data.Embedding...)
	}
	payload.Embeddings = &models.VectorEmbeddings{
		Vectors64: vectors,
	}
	return nil
}

func (x *Mistral) GenerateContent(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
	if payload.Params.Model == "" {
		x.log.Warn().Msg("Model name is empty, using default model")
		payload.Params.Model = defaultModel
	}
	req := x.buildRequest(payload.Params.Model, payload, false)
	if req == nil {
		x.log.Error().Msg("Error building request object for Mistral")
		return aicore.ErrInvalidAIModelRequest
	}
	resp := &GenerativeResponse{}
	err := x.client.Post(ctx, generatePath, req, resp, jsonContentType)
	if err != nil {
		x.log.Error().Err(err).Msg("Error generating content from Mistral completion")
		return err
	}
	x.buildResponse(resp, payload)
	return nil
}

func (x *Mistral) GenerateContentStream(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
	if payload.Params.Model == "" {
		x.log.Warn().Msg("Model name is empty, using default model")
		payload.Params.Model = defaultModel
	}
	req := x.buildRequest(payload.Params.Model, payload, true)
	if req == nil {
		x.log.Error().Msg("Error building request object for mistral")
		return aicore.ErrInvalidAIModelRequest
	}
	response, err := x.client.StreamPost(ctx, generatePath, req, jsonContentType)
	if err != nil {
		x.log.Error().Err(err).Msg("Error generating content from mistral completion")
		return err
	}
	if response == nil || response.Body == nil || response.StatusCode != http.StatusOK {
		x.log.Error().Msg("Error getting response from mistral")
		return aicore.ErrInvalidAIModelRequest
	}
	log.Println(utils.GetEpochTime(), "====================1")
	log.Println("Response from mistral stream", response)
	x.buildStreamResponse(response, payload)
	log.Println(utils.GetEpochTime(), "====================2")
	return nil
}

func (x *Mistral) CrawlModels(ctx context.Context) ([]*models.AIModelBase, error) {
	resp := &ModelList{}
	result := []*models.AIModelBase{}
	err := x.client.Get(ctx, modelsPath, nil, resp, jsonContentType)
	if err != nil {
		x.log.Error().Err(err).Msg("Error crawling models from Mistral")
		return nil, err
	}
	log.Println("Response from Mistral", resp)
	for _, model := range resp.Data {
		result = append(result, &models.AIModelBase{
			ModelId:   model.ID,
			OwnedBy:   model.OwnedBy,
			ModelName: model.Name,
			ModelType: func() string {
				if strings.Contains(model.Name, "embed") {
					return mpb.AIModelType_EMBEDDING.String()
				}
				return mpb.AIModelType_LLM.String()

			}(),
		})
	}
	return result, nil
}

func (x *Mistral) FIM(ctx context.Context, payload *prompts.GenerativePrompterPayload, model string) error {
	var err error
	if model == "" {
		x.log.Warn().Msg("Model name is empty, using default model")
		model = defaultModel
	}
	req := x.buildFimRequest(model, payload)
	if req == nil {
		x.log.Error().Msg("Error marshalling request object for Mistral")
	}
	resp := &FIMResponse{}
	err = x.client.Post(ctx, fimPath, req, resp, jsonContentType)
	if err != nil {
		x.log.Error().Err(err).Msg("Error generating FIM content from Mistral completion")
		return err
	}
	payload.ParseOutput(&prompts.PayloadgenericResponse{
		FinishReason: string(resp.Choices[0].FinishReason),
		Data:         resp.Object,
		Role:         resp.Choices[0].Message.Role,
	})
	return nil
}
