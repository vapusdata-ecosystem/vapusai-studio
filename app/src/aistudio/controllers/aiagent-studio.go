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

type VapusAIAgentStudio struct {
	pb.UnimplementedAIAgentStudioServer
	validator      *dmutils.DMValidator
	StudioServices *dmsvc.StudioServices
	Logger         zerolog.Logger
}

var VapusAIAgentStudioManager *VapusAIAgentStudio

func NewVapusAIAgentStudio() *VapusAIAgentStudio {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "VapusAIAgentStudio")
	validator, err := dmutils.NewDMValidator()
	if err != nil {
		l.Panic().Err(err).Msg("Error while loading validator")
	}

	l.Info().Msg("VapusAIAgentStudio Controller initialized")
	return &VapusAIAgentStudio{
		validator:      validator,
		Logger:         l,
		StudioServices: dmsvc.StudioServicesManager,
	}
}

func InitVapusAIAgentStudio() {
	if VapusAIAgentStudioManager == nil {
		VapusAIAgentStudioManager = NewVapusAIAgentStudio()
	}
}

func (v *VapusAIAgentStudio) ChatStream(req *pb.AgentInvokeRequest, stream pb.AIAgentStudio_ChatStreamServer) error {
	ctx := stream.Context()
	agent, err := v.StudioServices.NewVapusAIAgentThread(ctx, req, stream)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while creating new AI Agent studio thread")
		return err
	}
	err = agent.Act()
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while processing AI Agent studio request")
		return err
	}
	return nil
}

func (v *VapusAIAgentStudio) Chat(ctx context.Context, req *pb.AgentInvokeRequest) (*pb.AgentInvokeResponse, error) {
	agent, err := v.StudioServices.NewVapusAIAgentThread(ctx, req, nil)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while creating new AI Agent studio thread")
		return nil, err
	}
	err = agent.Act()
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while processing AI Agent studio request")
		return nil, err
	}
	response := &pb.AgentInvokeResponse{
		Output: agent.GetResult(),
		DmResp: grpcops.HandleDMResponse(ctx, "Agent non stream chat action executed successfully", "200"),
	}
	return response, nil
}
