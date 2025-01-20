package guardrails

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	"github.com/vapusdata-oss/aistudio/core/aistudio/prompts"
	aimodels "github.com/vapusdata-oss/aistudio/core/aistudio/providers"
	"github.com/vapusdata-oss/aistudio/core/models"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	aipb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

type GuardRailClient struct {
	Guardrail *models.AIGuardrails
	modelPool []*GuardModelNodePool
}

type GuardModelNodePool struct {
	Connection aimodels.AIModelNodeInterface
	IsAccount  bool
	Model      string
}

type GuardrailScanner struct {
	ContentGuard    []string
	TopicGuard      []string
	WordGuard       []string
	SenstivityGuard []string
	client          *GuardRailClient
	ID              string
}

type GuardrailsFunc func(*GuardRailClient)

func WithSpec(spec *models.AIGuardrails) GuardrailsFunc {
	return func(g *GuardRailClient) {
		g.Guardrail = spec
	}
}

func WithModelPool(pool []*GuardModelNodePool) GuardrailsFunc {
	return func(g *GuardRailClient) {
		g.modelPool = pool
	}
}

func New(opts ...GuardrailsFunc) *GuardRailClient {
	cl := &GuardRailClient{}
	for _, opt := range opts {
		opt(cl)
	}
	return cl
}

func (g *GuardRailClient) Scan(ctx context.Context, message string, logger zerolog.Logger) *GuardrailScanner {
	scanner := &GuardrailScanner{
		client: g,
	}
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		scanner.ScanContent(ctx, message, logger)
	}()
	go func() {
		defer wg.Done()
		scanner.ScanTopic(ctx, message, logger)
	}()
	go func() {
		defer wg.Done()
		scanner.ScanWords(ctx, message, logger)
	}()
	go func() {
		defer wg.Done()
		scanner.SensitivityDataAction(ctx, message, logger)
	}()
	wg.Wait()
	return scanner
}

func (g *GuardrailScanner) ScanContent(ctx context.Context, message string, logger zerolog.Logger) {}
func (g *GuardrailScanner) ScanTopic(ctx context.Context, message string, logger zerolog.Logger) {
	log.Println("Scanning topic-----------", message, "+++++++++++++++++++++++++", g.client.Guardrail.Schema)
	var contentHeatMap = map[string]string{
		"sexual":     g.client.Guardrail.Contents.Sexual,
		"hateSpeech": g.client.Guardrail.Contents.HateSpeech,
		"threats":    g.client.Guardrail.Contents.Threats,
		"insults":    g.client.Guardrail.Contents.Insults,
		"misconduct": g.client.Guardrail.Contents.Misconduct,
	}

	payload := prompts.NewPrompter(&aipb.ChatRequest{
		Model: g.client.modelPool[0].Model,
		Tools: func() []*mpb.ToolCall {
			var toolCalls []*mpb.ToolCall
			f := models.GetFunctionCallFromString(g.client.Guardrail.Schema)
			toolCalls = append(toolCalls, &mpb.ToolCall{
				Type: strings.ToLower(mpb.AIToolCallType_FUNCTION.String()),
				FunctionSchema: &mpb.FunctionCall{
					Name:           f.Name,
					Arguments:      f.GetStringParamSchema(),
					Description:    f.Description,
					RequiredFields: f.RequiredFields,
				},
			})
			return toolCalls
		}(),
		Messages: []*aipb.ChatMessageObject{
			{
				Role:    aipb.AIMessageRoles_SYSTEM,
				Content: "You are an AI guardrail inspector, please scan the user input based and provide the tool call response. Strictly, do not generate any other dataset.",
			},
			{
				Role:    aipb.AIMessageRoles_USER,
				Content: message,
			},
		},
	}, nil, nil, logger)
	payload.RenderPrompt()
	err := g.client.modelPool[0].Connection.GenerateContent(ctx, payload)
	if err != nil {
		logger.Error().Err(err).Msg("error while generating content")
		return
	}
	for _, obj := range payload.ToolCallResponse {
		argsresult := obj.FunctionSchema.GetArguments()
		result := map[string]interface{}{}
		log.Println("Scanning resultddd-----------", argsresult)
		log.Println("map vals-----------", contentHeatMap)
		err = json.Unmarshal([]byte(argsresult), &result)
		if err != nil {
			logger.Error().Err(err).Msg("error while unmarshalling function call arguments for topic")
			continue
		}
		ct, ok := result["content_guardrails"]
		if ok {
			cg, ok := ct.(map[string]interface{})
			if ok {
				for k, v := range cg {
					val := v.(string)
					acceptedVal, ok := mpb.GuardRailLevels_value[contentHeatMap[k]]
					if !ok {
						logger.Error().Msg("error while getting guardrail level")
						continue
					}
					foundLevel, ok := mpb.GuardRailLevels_value[strings.ToUpper(val)]
					if !ok {
						logger.Error().Msg("error while getting guardrail level")
						continue
					}
					log.Println("Scanning topic-----------", k, "acceptedVal", acceptedVal, "foundLevel", foundLevel)
					if acceptedVal != 0 {
						if foundLevel >= acceptedVal {
							g.ContentGuard = append(g.ContentGuard, fmt.Sprintf("%s is %v for input provided by the user, guardrail check failed", k, v))
						}
					}
				}
			}
		}
		tg, ok := result["topic_guardrails"]
		log.Println("Scanning tg-----------", tg)
		if ok {
			topic, ok := tg.(map[string]interface{})
			log.Println("Scanning tg topic-----------", topic)
			if ok {
				log.Println("Scanning topic-----------", topic)
				for k, v := range topic {
					if v.(bool) == true {
						log.Println("Scanning topic-----------", k)
						g.TopicGuard = append(g.TopicGuard, "user message contains topic "+k)
					}
				}
			}
		}
	}
}
func (g *GuardrailScanner) ScanWords(ctx context.Context, message string, logger zerolog.Logger) {
	for _, rule := range g.client.Guardrail.Words {
		for _, word := range rule.Words {
			log.Println("Scanning word-----------", word, " in message ", message)
			if strings.Contains(message, word) {
				g.WordGuard = append(g.WordGuard, word)
			}
		}
	}
}
func (g *GuardrailScanner) SensitivityDataAction(ctx context.Context, message string, logger zerolog.Logger) {
}
