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

type VapusAIGuardrails struct {
	pb.UnimplementedAIGuardrailsServer
	validator      *dmutils.DMValidator
	StudioServices *dmsvc.StudioServices
	Logger         zerolog.Logger
}

var VapusAIGuardrailsManager *VapusAIGuardrails

func NewVapusAIGuardrails() *VapusAIGuardrails {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "VapusAIGuardrails")
	validator, err := dmutils.NewDMValidator()
	if err != nil {
		l.Panic().Err(err).Msg("Error while loading validator")
	}

	l.Info().Msg("VapusAIGuardrails Controller initialized")
	return &VapusAIGuardrails{
		validator:      validator,
		Logger:         l,
		StudioServices: dmsvc.StudioServicesManager,
	}
}

func InitVapusAIGuardrails() {
	if VapusAIGuardrailsManager == nil {
		VapusAIGuardrailsManager = NewVapusAIGuardrails()
	}
}

func (v *VapusAIGuardrails) Manager(ctx context.Context, req *pb.GuardrailsManagerRequest) (*pb.GuardrailsResponse, error) {
	agent, err := v.StudioServices.NewVapusAIGuardrailManagerAgent(ctx, req, nil)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while creating new AI guardrails manager request")
		return nil, err
	}
	err = agent.Act()
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while processing AI guardrails manager request")
		return nil, err
	}
	response := &pb.GuardrailsResponse{
		Output: agent.GetResult(),
	}
	response.DmResp = grpcops.HandleDMResponse(ctx, "AIModelPromptConfigAgent action executed successfully", "200")
	return response, nil
}

func (v *VapusAIGuardrails) Getter(ctx context.Context, req *pb.GuardrailsGetterRequest) (*pb.GuardrailsResponse, error) {
	agent, err := v.StudioServices.NewVapusAIGuardrailManagerAgent(ctx, nil, req)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while creating new AI guardrails agent request")
		return nil, err
	}
	err = agent.Act()
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while processing AI guardrails agent request")
		return nil, err
	}
	response := &pb.GuardrailsResponse{
		Output: agent.GetResult(),
	}
	response.DmResp = grpcops.HandleDMResponse(ctx, "AIModelPromptConfigAgent action executed successfully", "200")
	return response, nil
}
