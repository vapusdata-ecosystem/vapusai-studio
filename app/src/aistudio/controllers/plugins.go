package dmcontrollers

import (
	"context"

	grpccodes "google.golang.org/grpc/codes"

	"github.com/rs/zerolog"
	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	services "github.com/vapusdata-oss/aistudio/aistudio/services"
	grpcops "github.com/vapusdata-oss/aistudio/core/serviceops/grpcops"
	dpb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

type PluginController struct {
	dpb.UnimplementedPluginServiceServer
	StudioServices *services.StudioServices
	logger         zerolog.Logger
}

var PluginControllerManager *PluginController

func NewPluginController() *PluginController {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "PluginController")
	l.Info().Msg("Organization Controller initialized")
	return &PluginController{
		StudioServices: services.StudioServicesManager,
		logger:         l,
	}
}

func InitPluginController() {
	if PluginControllerManager == nil {
		PluginControllerManager = NewPluginController()
	}
}

func (x *PluginController) PluginManager(ctx context.Context, request *dpb.PluginManagerRequest) (*dpb.PluginResponse, error) {
	agent, err := x.StudioServices.NewPluginManagerAgent(ctx, request, nil)
	if err != nil {
		return nil, grpcops.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	err = agent.Act()
	if err != nil {
		return nil, grpcops.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response := &dpb.PluginResponse{
		Output: agent.GetResult(),
	}
	response.DmResp = grpcops.HandleDMResponse(ctx, "PluginManager action executed successfully", "201")
	return response, nil
}

func (x *PluginController) PluginGetter(ctx context.Context, request *dpb.PluginGetterRequest) (*dpb.PluginResponse, error) {
	agent, err := x.StudioServices.NewPluginManagerAgent(ctx, nil, request)
	if err != nil {
		return nil, grpcops.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	err = agent.Act()
	if err != nil {
		return nil, grpcops.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response := &dpb.PluginResponse{
		Output: agent.GetResult(),
	}
	response.DmResp = grpcops.HandleDMResponse(ctx, "PluginGetter action executed successfully", "201")
	return response, nil
}

func (x *PluginController) PluginAction(ctx context.Context, request *dpb.PluginActionRequest) (*dpb.PluginActionResponse, error) {
	agent, err := x.StudioServices.NewPluginActionsAgent(ctx, request)
	if err != nil {
		return nil, grpcops.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	err = agent.Act()
	if err != nil {
		return nil, grpcops.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	return &dpb.PluginActionResponse{
		DmResp: grpcops.HandleDMResponse(ctx, "Email sent successfully", "201"),
	}, nil
}
