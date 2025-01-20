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

type VapusAIAgents struct {
	pb.UnimplementedAIAgentsServer
	validator      *dmutils.DMValidator
	StudioServices *dmsvc.StudioServices
	Logger         zerolog.Logger
}

var VapusAIAgentsManager *VapusAIAgents

func NewVapusAIAgents() *VapusAIAgents {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "VapusAIAgents")
	validator, err := dmutils.NewDMValidator()
	if err != nil {
		l.Panic().Err(err).Msg("Error while loading validator")
	}

	l.Info().Msg("VapusAIAgents Controller initialized")
	return &VapusAIAgents{
		validator:      validator,
		Logger:         l,
		StudioServices: dmsvc.StudioServicesManager,
	}
}

func InitVapusAIAgents() {
	if VapusAIAgentsManager == nil {
		VapusAIAgentsManager = NewVapusAIAgents()
	}
}

func (v *VapusAIAgents) Manager(ctx context.Context, req *pb.AgentManagerRequest) (*pb.AgentResponse, error) {
	agent, err := v.StudioServices.NewVapusAIAgentManagerAgent(ctx, req, nil)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while creating new AI Agent request")
		return nil, err
	}
	err = agent.Act()
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while processing AI Agent request")
		return nil, err
	}
	response := &pb.AgentResponse{
		Output: agent.GetResult(),
	}
	response.DmResp = grpcops.HandleDMResponse(ctx, "AIModelPromptConfigAgent action executed successfully", "200")
	return response, nil
}

func (v *VapusAIAgents) Getter(ctx context.Context, req *pb.AgentGetterRequest) (*pb.AgentResponse, error) {
	agent, err := v.StudioServices.NewVapusAIAgentManagerAgent(ctx, nil, req)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while creating new AI Agent request")
		return nil, err
	}
	err = agent.Act()
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while processing AI Agent request")
		return nil, err
	}
	response := &pb.AgentResponse{
		Output: agent.GetResult(),
	}
	response.DmResp = grpcops.HandleDMResponse(ctx, "AIModelPromptConfigAgent action executed successfully", "200")
	return response, nil
}
