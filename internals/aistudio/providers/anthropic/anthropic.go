package anthropic

import (
	"bufio"
	"context"
	"encoding/json"
	"html"
	"log"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
	aicore "github.com/vapusdata-oss/aistudio/core/aistudio/core"
	"github.com/vapusdata-oss/aistudio/core/aistudio/prompts"
	httpCls "github.com/vapusdata-oss/aistudio/core/http"
	dmlogger "github.com/vapusdata-oss/aistudio/core/logger"
	"github.com/vapusdata-oss/aistudio/core/models"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

type AnthropicInterface interface {
	GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload, model string) error
	GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	CrawlModels(ctx context.Context) ([]*models.AIModelBase, error)
}

type Anthropic struct {
	client     *httpCls.RestHttp
	log        zerolog.Logger
	modelNode  *models.AIModelNode
	maxRetries int
	params     map[string]interface{}
}

func New(node *models.AIModelNode, retries int, logger zerolog.Logger) (AnthropicInterface, error) {
	token := ""
	if node.GetCredentials("default") != nil {
		token = node.GetCredentials("default").ApiToken
	}
	httpCl, err := httpCls.New(logger,
		httpCls.WithAddress(defaultEndpoint),
		httpCls.WithBasePath(baseAPIPath),
		httpCls.WithCustomHeaders(map[string]string{
			versionHeaderKey: defaultAPIVersion,
			apiKeyHeader:     token,
		}),
	)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating http client for anthropic")
		return nil, err
	}
	return &Anthropic{
		client:    httpCl,
		log:       dmlogger.GetSubDMLogger(logger, "ailogger", "Anthropic"),
		modelNode: node,
	}, nil
}

func (x *Anthropic) buildRequest(model string, payload *prompts.GenerativePrompterPayload, stream bool) []byte {
	messages := make([]message, 0)
	sysMessage := ""
	if payload.Params.Mode == pb.AIInterfaceMode_CHAT_MODE {
		for _, msg := range payload.SessionContext {
			if msg.Role == prompts.USER {
				messages = append(messages, message{
					Role: aicore.USER,
					Content: []content{{
						Type: messageTypeText,
						Text: &msg.Message,
					},
					},
				})
			} else {
				messages = append(messages, message{
					Role: aicore.ASSISTANT,
					Content: []content{{
						Type: messageTypeText,
						Text: &msg.Message,
					},
					},
				})
			}
		}
	}
	for _, c := range payload.Params.Messages {
		if c.Role == pb.AIMessageRoles_USER {
			messages = append(messages, message{
				Role: aicore.USER,
				Content: []content{{
					Type: messageTypeText,
					Text: &c.Content,
				},
				},
			})
		} else if c.Role == pb.AIMessageRoles_SYSTEM {
			sysMessage = sysMessage + "\n" + c.Content
		}

	}

	tools := make([]Tool, 0)
	if len(payload.ToolCalls) > 0 {
		for _, tool := range payload.ToolCalls {
			if tool != nil && tool.FunctionSchema != nil {
				tools = append(tools, Tool{
					Name:        tool.FunctionSchema.Name,
					Description: tool.FunctionSchema.Description,
					InputSchema: tool.FunctionSchema.Arguments,
				})
			}
		}
		payload.ParseToolCallResponse()
	}
	reqObj := &GenerativeRequest{
		System:      sysMessage,
		Model:       model,
		Messages:    messages,
		Stream:      stream,
		Temperature: float64(payload.Params.Temperature),
		MaxTokens:   int(payload.Params.MaxOutputTokens),
	}
	reqBytes, err := json.Marshal(reqObj)
	if err != nil {
		x.log.Error().Err(err).Msg("Error marshalling request object")
		return nil
	}
	return reqBytes
}

func (x *Anthropic) buildResponse(resp *GenerativeResponse, payload *prompts.GenerativePrompterPayload) {
	dt := ""
	for _, content := range resp.Content {
		if content.Type == "tool_use" {
			payload.ToolCallResponse = append(payload.ToolCallResponse, &mpb.ToolCall{
				Type: mpb.AIToolCallType_FUNCTION.String(),
				FunctionSchema: &mpb.FunctionCall{
					Name:      *content.Name,
					Arguments: *content.Input,
				},
			})
		} else {
			dt = dt + *content.Text
		}
	}
	payload.ParseOutput(&prompts.PayloadgenericResponse{
		FinishReason: string(*resp.StopReason),
		Data:         dt,
		Role:         resp.Role,
	})
	return
}

func (x *Anthropic) buildStreamResponse(resp *http.Response, payload *prompts.GenerativePrompterPayload) {
	func() {
		reader := bufio.NewReader(resp.Body)
		var currentEvent streamResponse
		var errThreshold int = 0
		defer resp.Body.Close()
		for {
			ed := &streamEventData{}
			line, err := reader.ReadString('\n')
			if err != nil {
				x.log.Error().Err(err).Msg("Error reading response from stream for anthropic completion")
				if errThreshold > 5 {
					break
				} else {
					errThreshold++
					continue
				}
			}

			line = strings.TrimSpace(line)
			// Check if this is a new event
			if strings.HasPrefix(line, "event:") {
				currentEvent.Event = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			} else if strings.HasPrefix(line, "data:") {
				data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
				if currentEvent.Event == "content_block_delta" {
					bbytes := []byte(data)
					err = json.Unmarshal(bbytes, ed)
					if err != nil {
						x.log.Error().Err(err).Msg("Error parsing JSON data from stream response for anthropic completion stream")
						continue
					}
					err = payload.SendStreamData(html.UnescapeString(ed.Delta.Text), false, false)
					if err != nil {
						x.log.Err(err).Msg("error while sending stream response")
						_ = payload.SendStreamData("", false, true)
						continue
					}
				} else if currentEvent.Event == "message_stop" {
					_ = payload.SendStreamData("", true, false)
					break
				}
			} else if line == "" {
				continue
			}
		}
	}()
	return
}

func (x *Anthropic) GenerateEmbeddings(ctx context.Context, payload *prompts.AIEmbeddingPayload, model string) error {
	return aicore.ErrEmbeddingNotSupported
}

func (x *Anthropic) GenerateContent(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
	if payload.Params.Model == "" {
		x.log.Warn().Msg("Model name is empty, using default model")
		payload.Params.Model = defaultModel
	}
	req := x.buildRequest(payload.Params.Model, payload, false)
	if req == nil {
		x.log.Error().Msg("Error building request object for anthropic")
		return aicore.ErrInvalidAIModelRequest
	}
	resp := &GenerativeResponse{}
	err := x.client.Post(ctx, messagePath, req, resp, jsonContentType)
	if err != nil {
		x.log.Error().Err(err).Msg("Error generating content from anthropic completion")
		return err
	}
	x.buildResponse(resp, payload)
	return nil
}

func (x *Anthropic) GenerateContentStream(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
	if payload.Params.Model == "" {
		x.log.Warn().Msg("Model name is empty, using default model")
		payload.Params.Model = defaultModel
	}
	req := x.buildRequest(payload.Params.Model, payload, true)
	if req == nil {
		x.log.Error().Msg("Error building request object for anthropic")
		return aicore.ErrInvalidAIModelRequest
	}
	log.Println(string(req))
	response, err := x.client.StreamPost(ctx, messagePath, req, jsonContentType)
	if err != nil {
		x.log.Error().Err(err).Msg("Error generating content from anthropic completion")
		return err
	}
	if response == nil || response.Body == nil || response.StatusCode != http.StatusOK {
		x.log.Error().Msg("Error getting response from anthropic")
		return aicore.ErrInvalidAIModelRequest
	}
	x.buildStreamResponse(response, payload)
	return nil
}

func (x *Anthropic) CrawlModels(ctx context.Context) ([]*models.AIModelBase, error) {
	return []*models.AIModelBase{
		{
			ModelName: defaultModel,
			ModelId:   defaultModel,
			ModelType: mpb.AIModelType_LLM.String(),
		},
		{
			ModelName: "claude-3-5-haiku-20241022",
			ModelId:   "claude-3-5-haiku-20241022",
			ModelType: mpb.AIModelType_LLM.String(),
		},
		{
			ModelName: "claude-3-opus-20240229",
			ModelId:   "claude-3-opus-20240229",
			ModelType: mpb.AIModelType_LLM.String(),
		},
	}, nil
}
