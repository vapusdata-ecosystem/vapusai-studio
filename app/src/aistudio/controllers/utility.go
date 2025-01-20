package dmcontrollers

import (
	"context"

	"github.com/rs/zerolog"
	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	services "github.com/vapusdata-oss/aistudio/aistudio/services"
	utils "github.com/vapusdata-oss/aistudio/aistudio/utils"
	grpcops "github.com/vapusdata-oss/aistudio/core/serviceops/grpcops"
	dpb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

type UtilityController struct {
	dpb.UnimplementedUtilityServiceServer
	StudioServices *services.StudioServices
	logger         zerolog.Logger
}

var UtilityControllerManager *UtilityController

func NewUtilityController() *UtilityController {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "UtilityController")
	l.Info().Msg("UtilityController initialized")
	return &UtilityController{
		StudioServices: services.StudioServicesManager,
		logger:         l,
	}
}

func InitUtilityController() {
	if UtilityControllerManager == nil {
		UtilityControllerManager = NewUtilityController()
	}
}

func (x *UtilityController) Upload(ctx context.Context, request *dpb.UploadRequest) (*dpb.UploadResponse, error) {
	utAgent, err := x.StudioServices.NewUtilityAgent(ctx, request, nil)
	if err != nil {
		return nil, grpcops.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	err = utAgent.Act()
	if err != nil {
		return nil, grpcops.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	resp := utAgent.GetUploadedResult()
	return resp, nil
}

func (x *UtilityController) UploadStream(stream dpb.UtilityService_UploadStreamServer) error {
	utAgent, err := x.StudioServices.NewUtilityAgent(stream.Context(), nil, stream)
	if err != nil {
		return grpcops.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	err = utAgent.Act()
	if err != nil {
		return grpcops.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	return nil
}

func (dmc *UtilityController) StoreDMSecrets(ctx context.Context, request *pb.StoreDMSecretsRequest) (*pb.StoreDMSecretsResponse, error) {
	err := dmc.StudioServices.StoreCredential(ctx, request)
	if err != nil {
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	response := &pb.StoreDMSecretsResponse{
		Name:   request.GetName(),
		VPath:  request.GetVPath(),
		DmResp: grpcops.HandleDMResponse(ctx, utils.CREDDENTIALS_STORED, "201"),
	}
	return response, nil
}
