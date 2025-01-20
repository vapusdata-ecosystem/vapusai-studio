package dmcontrollers

import (
	"context"

	"github.com/rs/zerolog"

	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	dmsvc "github.com/vapusdata-oss/aistudio/aistudio/services"
	dmutils "github.com/vapusdata-oss/aistudio/core/utils"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

type AIModelStudio struct {
	pb.UnimplementedAIModelStudioServer
	validator      *dmutils.DMValidator
	StudioServices *dmsvc.StudioServices
	Logger         zerolog.Logger
}

var AIModelInterfaceManager *AIModelStudio

func NewAIModelStudio() *AIModelStudio {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "AIModelInterface")
	validator, err := dmutils.NewDMValidator()
	if err != nil {
		l.Panic().Err(err).Msg("Error while loading validator")
	}

	l.Info().Msg("AIModelInterface Controller initialized")
	return &AIModelStudio{
		validator:      validator,
		Logger:         l,
		StudioServices: dmsvc.StudioServicesManager,
	}
}

func InitNewAIModelStudio() {
	if AIModelInterfaceManager == nil {
		AIModelInterfaceManager = NewAIModelStudio()
	}
}

func (v *AIModelStudio) GenerateEmbeddings(ctx context.Context, req *pb.EmbeddingsInterface) (*pb.EmbeddingsResponse, error) {
	agent, err := v.StudioServices.NewEmbeddingAgent(ctx, req)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while creating new NewEmbeddingAgent Agent request")
		return nil, err
	}
	err = agent.Act()
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while processing NewEmbeddingAgent Agent request")
		return nil, err
	}
	return &pb.EmbeddingsResponse{
		Output: agent.GetEmbeddings(),
	}, nil
}

func (v *AIModelStudio) Chat(ctx context.Context, req *pb.ChatRequest) (*pb.ChatResponse, error) {
	agent, err := v.StudioServices.NewAIModelAgent(ctx, req, nil)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while creating new AIModelNodeInterface Agent request")
		return nil, err
	}
	err = agent.Act()
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while processing AIModelNodeInterface Agent request")
		return nil, err
	}
	return agent.GetResult(), nil
}

func (v *AIModelStudio) ChatStream(req *pb.ChatRequest, stream pb.AIModelStudio_ChatStreamServer) error {
	ctx := stream.Context()
	agent, err := v.StudioServices.NewAIModelAgent(ctx, req, stream)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while creating new AIModelNodeInterface Agent request")
		return err
	}
	err = agent.Act()
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while processing AIModelNodeInterface Agent request")
		return err
	}
	return nil
}
