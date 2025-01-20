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

type AIModels struct {
	pb.UnimplementedAIModelsServer
	validator      *dmutils.DMValidator
	StudioServices *dmsvc.StudioServices
	Logger         zerolog.Logger
}

var AIModelsManager *AIModels

func NewAIModels() *AIModels {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "AIModels")
	validator, err := dmutils.NewDMValidator()
	if err != nil {
		l.Panic().Err(err).Msg("Error while loading validator")
	}

	l.Info().Msg("AIModels Controller initialized")
	return &AIModels{
		validator:      validator,
		Logger:         l,
		StudioServices: dmsvc.StudioServicesManager,
	}
}

func InitAIModels() {
	if AIModelsManager == nil {
		AIModelsManager = NewAIModels()
	}
}

func (v *AIModels) Manager(ctx context.Context, req *pb.AIModelNodeConfiguratorRequest) (*pb.AIModelNodeResponse, error) {
	agent, err := v.StudioServices.NewVapusAINodeManagerAgent(ctx, req, nil)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while creating new AIStudio Node Agent request")
		return nil, err
	}
	err = agent.Act()
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while processing AIStudio Node Agent request")
		return nil, err
	}
	response := agent.GetResult()
	response.DmResp = grpcops.HandleDMResponse(ctx, "AIModelNodeConfigAgent action executed successfully", "200")
	return response, nil
}

func (v *AIModels) Getter(ctx context.Context, req *pb.AIModelNodeGetterRequest) (*pb.AIModelNodeResponse, error) {
	agent, err := v.StudioServices.NewVapusAINodeManagerAgent(ctx, nil, req)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while creating new AIStudio Node Agent request")
		return nil, err
	}
	err = agent.Act()
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while processing AIStudio Node Agent request")
		return nil, err
	}
	response := agent.GetResult()
	response.DmResp = grpcops.HandleDMResponse(ctx, "AIModelNodeConfigAgent action executed successfully", "200")
	return response, nil
}
