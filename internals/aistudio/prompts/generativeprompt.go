package prompts

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/rs/zerolog"
	aicore "github.com/vapusdata-oss/aistudio/core/aistudio/core"
	"github.com/vapusdata-oss/aistudio/core/models"
	utils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

// var Baseregex = `{%s}\s*(.*?)\s*{/%s}`
var Baseregex = `(?s)\{%s\}(.*?)\{/%s\}`

type AIEmbeddingPayload struct {
	Input          string
	Dimensions     int
	EmbeddingModel string
	Embeddings     *models.VectorEmbeddings
	SystemMessage  string
	UserMessage    string
}

type GenerativePrompterPayload struct {
	Params           *pb.ChatRequest
	Prompt           *models.AIModelPrompt
	Response         *pb.ChatResponse
	Error            error
	ParsedOutput     string
	Context          []*mpb.Mapper
	ResultMetadata   map[string]any
	Suffix           string
	isRendered       bool
	SessionID        string
	SessionContext   []*SessionMessage
	SummaryOutput    string
	ToolCallResponse []*mpb.ToolCall
	opStart          string
	opEnd            string
	Stream           pb.AIModelStudio_ChatStreamServer
	streamLog        string
	opStartReacted   bool
	opEndReacted     bool
	GuardrailsFailed bool
	GuardrailsResult map[string][]string
	ToolCalls        []*mpb.ToolCall
	logger           zerolog.Logger
}

type PayloadgenericResponse struct {
	Data         string
	Role         string
	FinishReason string
	Usage        any
}

type SessionMessage struct {
	Message string
	Role    string
}

func NewPrompter(params *pb.ChatRequest, prompt *models.AIModelPrompt, stream pb.AIModelStudio_ChatStreamServer, logger zerolog.Logger) *GenerativePrompterPayload {
	return &GenerativePrompterPayload{
		Params:           params,
		Prompt:           prompt,
		ResultMetadata:   map[string]any{},
		Stream:           stream,
		GuardrailsFailed: false,
		GuardrailsResult: map[string][]string{},
		ToolCalls:        []*mpb.ToolCall{},
	}
}

func (p *GenerativePrompterPayload) ParseOutput(opts *PayloadgenericResponse) {
	outPut := opts.Data
	if p.Prompt != nil {
		// outPut = strings.Replace(outPut, "\n", "", -1)
		log.Println("Output from ai model before parse", outPut)
		if p.Prompt.Prompt.OutputTag == "" {
			p.ParsedOutput = outPut
		} else {
			outPut, _ = utils.ExtractBetweenDelimiters(outPut, p.opStart, p.opEnd)
			if outPut == "" {
				outPut = opts.Data
			}
			p.ParsedOutput = outPut
		}
	} else {
		p.ParsedOutput = outPut
	}
	opts.Data = p.ParsedOutput
	p.BuildResponseOP(aicore.StreamEventData.String(), opts, false)
	return
}

func (p *GenerativePrompterPayload) ParseToolCallResponse() {
	if len(p.ToolCallResponse) > 0 {
		p.Response.Choices = append(p.Response.Choices, &pb.ChatResponseChoice{
			Messages: &pb.ChatMessageObject{
				Role:      pb.AIMessageRoles_ASSISTANT,
				ToolCalls: p.ToolCallResponse,
			},
		})
	}
}

func (p *GenerativePrompterPayload) RenderPrompt() {
	if p.isRendered {
		return
	}
	contextsMap := make(map[string]interface{})
	contexts := ""
	if len(p.Context) > 0 {
		for _, context := range p.Context {
			contextsMap[context.Key] = context.Value
		}
		contextBytes, err := json.Marshal(contextsMap)
		if err != nil {
			p.logger.Err(err).Msg("error while marshalling context")
			contexts = ""
		} else {
			contexts = string(contextBytes)
		}
	}
	if p.Prompt != nil {
		if p.Params.MaxOutputTokens == 0 {
			p.Params.MaxOutputTokens = DefaultMaxOPTokenLength
		}
		p.isRendered = true
		p.opStart = fmt.Sprintf("{%s}", p.Prompt.Prompt.OutputTag)
		p.opEnd = fmt.Sprintf("{/%s}", p.Prompt.Prompt.OutputTag)
		p.opEndReacted = false
		p.opStartReacted = false
		systemMessRendered := false
		for _, mess := range p.Params.Messages {
			if mess.Role == pb.AIMessageRoles_SYSTEM {
				mess.Content = p.Prompt.Prompt.SystemMessage + "\n" + mess.Content
				systemMessRendered = true
			}
			if mess.Role == pb.AIMessageRoles_USER {
				mes := mess.Content
				mess.Content = strings.Replace(p.Prompt.UserTemplate, "["+p.Prompt.Prompt.InputTag+"]", mes, -1)
				if len(contexts) > 0 && contexts != "{}" && contexts != "" {
					mess.Content = strings.Replace(mess.Content, "["+p.Prompt.Prompt.ContextTag+"]", contexts, -1)
				}
			}
		}
		if len(p.Prompt.Prompt.Tools) > 0 {
			for _, tool := range p.Prompt.Prompt.Tools {
				mf := models.GetFunctionCallFromString(tool.ToolSchema)
				argBytes, err := json.MarshalIndent(mf.Parameters, "", "  ")
				if err != nil {
					p.logger.Err(err).Msg("error while marshalling tool arguments")
					continue
				}
				p.ToolCalls = append(p.ToolCalls, &mpb.ToolCall{
					Type: tool.Type,
					FunctionSchema: &mpb.FunctionCall{
						Name:           mf.Name,
						RequiredFields: mf.RequiredFields,
						Description:    mf.Description,
						Arguments:      string(argBytes),
					},
				})
			}
		}
		if !systemMessRendered {
			p.Params.Messages = append(p.Params.Messages, &pb.ChatMessageObject{
				Role:    pb.AIMessageRoles_SYSTEM,
				Content: p.Prompt.Prompt.SystemMessage,
			})
		}
	} else {
		if len(contexts) > 0 && contexts != "{}" && contexts != "" {
			p.Params.Messages = append(p.Params.Messages, &pb.ChatMessageObject{
				Role:    pb.AIMessageRoles_USER,
				Content: "Context: " + contexts,
			})
		}
	}
	if p.Params.MaxOutputTokens == 0 {
		p.Params.MaxOutputTokens = DefaultMaxOPTokenLength
	}
	if len(p.Params.Tools) > 0 {
		p.ToolCalls = append(p.ToolCalls, p.Params.Tools...)
	}
	p.Response = &pb.ChatResponse{}
	log.Println("------------------------------------", p.ToolCalls)
	log.Println("Rendered MaxOutputTokens", p.Params.MaxOutputTokens)
	return

}

func (p *GenerativePrompterPayload) GetUserMessage() string {
	return p.Prompt.Prompt.UserMessage
}
func (p *GenerativePrompterPayload) FilterTag(text string) string {
	if !p.opStartReacted {
		if strings.Contains(text, p.opStart) {
			p.opStartReacted = true
			return strings.ReplaceAll(text, p.opStart, "")
		} else {
			return text
		}
	} else if !p.opEndReacted {
		if strings.Contains(text, p.opEnd) {
			p.opEndReacted = true
			return strings.ReplaceAll(text, p.opEnd, "")
		} else {
			return text
		}
	} else {
		return text
	}
}

func (p *GenerativePrompterPayload) SendStreamData(content string, isEnd bool, eventErr bool) error {
	var err error
	p.streamLog = p.FilterTag(p.streamLog)
	if len(p.streamLog) < 25 {
		if isEnd {
			p.EndStream(aicore.StreamEventEnd.String())
		} else if eventErr {
			p.EndStream(aicore.StreamEventError.String())
		} else {
			p.streamLog = p.streamLog + content
			return nil
		}
	} else {
		if isEnd {
			p.EndStream(aicore.StreamEventEnd.String())
		} else if eventErr {
			p.EndStream(aicore.StreamEventError.String())
		} else {
			err = p.Stream.Send(p.BuildResponseOP(aicore.StreamEventData.String(), &PayloadgenericResponse{
				Data: p.streamLog,
				Role: pb.AIMessageRoles_ASSISTANT.String(),
			}, true),
			)
			p.streamLog = ""
			return err
		}
		return err
	}
	return err
}

func (p *GenerativePrompterPayload) EndStream(event string) error {
	if len(p.streamLog) > 0 {
		p.Stream.Send(p.BuildResponseOP(aicore.StreamEventData.String(), &PayloadgenericResponse{
			Data: p.streamLog,
			Role: pb.AIMessageRoles_ASSISTANT.String(),
		}, true),
		)
		p.streamLog = ""
	}
	return p.Stream.SendMsg(&pb.StreamChatResponse{
		Output: &pb.ChatResponse{
			Event: event,
			Model: p.Params.Model,
		},
	})
}

func (p *GenerativePrompterPayload) BuildResponseOP(event string, opts *PayloadgenericResponse, isStream bool) *pb.ChatResponse {
	var role pb.AIMessageRoles
	for k := range pb.AIMessageRoles_value {
		if strings.ToLower(k) == strings.ToLower(opts.Role) {
			role = pb.AIMessageRoles(pb.AIMessageRoles_value[k])
			break
		} else {
			role = pb.AIMessageRoles_ASSISTANT
		}
	}
	if isStream {
		return &pb.ChatResponse{
			Choices: []*pb.ChatResponseChoice{
				{
					Messages: &pb.ChatMessageObject{
						Role:    role,
						Content: opts.Data,
					},
					FinishReason: opts.FinishReason,
				},
			},
			Model:   p.Params.Model,
			Created: utils.GetEpochTime(),
			Event:   event,
		}
	}
	if p.Response == nil {
		p.Response = &pb.ChatResponse{
			Choices: make([]*pb.ChatResponseChoice, 0),
			Model:   p.Params.Model,
			Created: utils.GetEpochTime(),
		}
	}
	ch := &pb.ChatResponseChoice{
		Messages: &pb.ChatMessageObject{
			Role:    role,
			Content: opts.Data,
		},
		FinishReason: opts.FinishReason,
	}
	p.Response.Choices = append(p.Response.Choices, ch)
	return p.Response
}
