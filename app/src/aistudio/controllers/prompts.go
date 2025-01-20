package dmcontrollers

import (
	"context"

	"github.com/rs/zerolog"
	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	dmsvc "github.com/vapusdata-oss/aistudio/aistudio/services"
	grpcops "github.com/vapusdata-oss/aistudio/core/serviceops/grpcops"
	dmutils "github.com/vapusdata-oss/aistudio/core/utils"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

type AIPrompts struct {
	pb.UnimplementedAIPromptsServer
	validator      *dmutils.DMValidator
	StudioServices *dmsvc.StudioServices
	Logger         zerolog.Logger
}

var AIPromptsManager *AIPrompts

func NewAIPrompts() *AIPrompts {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "AIPrompts")
	validator, err := dmutils.NewDMValidator()
	if err != nil {
		l.Panic().Err(err).Msg("Error while loading validator")
	}

	l.Info().Msg("AIPrompts Controller initialized")
	return &AIPrompts{
		validator:      validator,
		Logger:         l,
		StudioServices: dmsvc.StudioServicesManager,
	}
}

func InitAIPrompts() {
	if AIPromptsManager == nil {
		AIPromptsManager = NewAIPrompts()
	}
}

func (v *AIPrompts) Manager(ctx context.Context, req *pb.PromptManagerRequest) (*pb.PromptResponse, error) {
	agent, err := v.StudioServices.NewVapusAIPromptManagerAgent(ctx, req, nil)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while creating new AIStudio Prompt Agent request")
		return nil, err
	}
	err = agent.Act()
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while processing AIStudio Prompt Agent request")
		return nil, err
	}
	response := &pb.PromptResponse{
		Output: agent.GetResult(),
	}
	response.DmResp = grpcops.HandleDMResponse(ctx, "AIModelPromptConfigAgent action executed successfully", "200")
	return response, nil
}

func (v *AIPrompts) Getter(ctx context.Context, req *pb.PromptGetterRequest) (*pb.PromptResponse, error) {
	agent, err := v.StudioServices.NewVapusAIPromptManagerAgent(ctx, nil, req)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while creating new AIStudio Prompt Agent request")
		return nil, err
	}
	err = agent.Act()
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while processing AIStudio Prompt Agent request")
		return nil, err
	}
	response := &pb.PromptResponse{
		Output: agent.GetResult(),
	}
	response.DmResp = grpcops.HandleDMResponse(ctx, "AIModelPromptConfigAgent action executed successfully", "200")
	return response, nil
}
