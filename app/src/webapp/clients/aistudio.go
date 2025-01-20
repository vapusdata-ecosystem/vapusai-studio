package clients

import (
	"log"

	"github.com/labstack/echo/v4"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

func (x *GrpcClient) AIModelNodes(eCtx echo.Context) []*mpb.AIModelNode {
	result, err := x.AIModelClient.Getter(x.SetAuthCtx(eCtx), &pb.AIModelNodeGetterRequest{})
	if err != nil {
		x.logger.Err(err).Msg("error while getting AI Model Nodes")
		return []*mpb.AIModelNode{}
	}
	return result.Output.AiModelNodes
}

func (x *GrpcClient) AIModelNodesDetails(eCtx echo.Context, id string) *mpb.AIModelNode {
	log.Println("id", id)
	result, err := x.AIModelClient.Getter(x.SetAuthCtx(eCtx), &pb.AIModelNodeGetterRequest{
		AiModelNodeId: id,
	})
	if err != nil || result == nil || result.Output == nil || len(result.Output.AiModelNodes) == 0 {
		x.logger.Err(err).Msg("error while getting AI Model Nodes")
		return nil
	}
	return result.Output.AiModelNodes[0]
}

func (x *GrpcClient) AIModelPrompts(eCtx echo.Context) []*mpb.AIModelPrompt {
	result, err := x.AIPromptClient.Getter(x.SetAuthCtx(eCtx), &pb.PromptGetterRequest{})
	if err != nil {
		x.logger.Err(err).Msg("error while getting AI Prompt List")
		return []*mpb.AIModelPrompt{}
	}
	return result.Output
}

func (x *GrpcClient) AIModelPromptDetails(eCtx echo.Context, id string) *mpb.AIModelPrompt {
	result, err := x.AIPromptClient.Getter(x.SetAuthCtx(eCtx), &pb.PromptGetterRequest{
		PromptId: id,
	})
	if err != nil || len(result.Output) == 0 {
		x.logger.Err(err).Msg("error while getting AI Prompt details")
		return nil
	}
	return result.Output[0]
}

func (x *GrpcClient) ListAIAgents(eCtx echo.Context) []*mpb.VapusAIAgent {
	result, err := x.AIAgentClient.Getter(x.SetAuthCtx(eCtx), &pb.AgentGetterRequest{})
	if err != nil {
		x.logger.Err(err).Msg("error while getting AI Prompt List")
		return []*mpb.VapusAIAgent{}
	}
	return result.Output
}

func (x *GrpcClient) AIAgentDetails(eCtx echo.Context, id string) *mpb.VapusAIAgent {
	result, err := x.AIAgentClient.Getter(x.SetAuthCtx(eCtx), &pb.AgentGetterRequest{
		AgentId: id,
	})
	if err != nil || len(result.Output) == 0 {
		x.logger.Err(err).Msg("error while getting AI Prompt details")
		return nil
	}
	return result.Output[0]
}

func (x *GrpcClient) ListAIGuardrails(eCtx echo.Context) []*mpb.AIGuardrails {
	result, err := x.AIGurdrailsClient.Getter(x.SetAuthCtx(eCtx), &pb.GuardrailsGetterRequest{})
	if err != nil {
		x.logger.Err(err).Msg("error while getting AI Prompt List")
		return []*mpb.AIGuardrails{}
	}
	return result.Output
}

func (x *GrpcClient) DescribeAIGuardrail(eCtx echo.Context, id string) *mpb.AIGuardrails {
	result, err := x.AIGurdrailsClient.Getter(x.SetAuthCtx(eCtx), &pb.GuardrailsGetterRequest{
		GuardrailId: id,
	})
	if err != nil || len(result.Output) == 0 {
		x.logger.Err(err).Msg("error while getting AI Prompt details")
		return nil
	}
	return result.Output[0]
}
