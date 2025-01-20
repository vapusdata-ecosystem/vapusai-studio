package ai

import (
	"context"

	"github.com/rs/zerolog"
	aicore "github.com/vapusdata-oss/aistudio/core/aistudio/core"
	"github.com/vapusdata-oss/aistudio/core/aistudio/prompts"
	anthropic "github.com/vapusdata-oss/aistudio/core/aistudio/providers/anthropic"
	google "github.com/vapusdata-oss/aistudio/core/aistudio/providers/google"
	mistral "github.com/vapusdata-oss/aistudio/core/aistudio/providers/mistral"
	openai "github.com/vapusdata-oss/aistudio/core/aistudio/providers/openai"
	dmlogger "github.com/vapusdata-oss/aistudio/core/logger"
	"github.com/vapusdata-oss/aistudio/core/models"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

const (
	defaultRetries = 3
)

var AILogger zerolog.Logger

type AIModelNodeInterface interface {
	GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload, model string) error
	GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	CrawlModels(ctx context.Context) ([]*models.AIModelBase, error)
}

type AIModelNodeClient struct {
	node   *models.AIModelNode
	logger zerolog.Logger
}

type AiModelOpts func(*AIModelNodeClient)

func WithAIModelNode(c *models.AIModelNode) AiModelOpts {
	return func(opts *AIModelNodeClient) {
		opts.node = c
	}
}

func WithLogger(logger zerolog.Logger) AiModelOpts {
	return func(opts *AIModelNodeClient) {
		opts.logger = logger
	}
}

func NewAIModelNode(opts ...AiModelOpts) (AIModelNodeInterface, error) {
	configurator := &AIModelNodeClient{}
	for _, opt := range opts {
		opt(configurator)
	}
	AILogger = dmlogger.GetSubDMLogger(configurator.logger, "ailogger", "base")
	switch configurator.node.ServiceProvider {
	case mpb.LLMServiceProvider_OPENAI.String():
		return openai.New(configurator.node, AILogger), nil
	case mpb.LLMServiceProvider_ANTHROPIC.String():
		return anthropic.New(configurator.node, defaultRetries, AILogger)
	case mpb.LLMServiceProvider_MISTRAL.String():
		return mistral.New(configurator.node, defaultRetries, AILogger)
	case mpb.LLMServiceProvider_GEMINI.String():
		return google.New(context.TODO(), configurator.node, defaultRetries, AILogger)
	default:
		configurator.logger.Error().Msgf("Unknown service provider: %s", configurator.node.ServiceProvider)
		return nil, aicore.ErrUnknownServiceProvider
	}
}
