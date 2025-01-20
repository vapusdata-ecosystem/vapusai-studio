package dmcontrollers

// Package controllers provides the implementation of domain controllers.
// These controllers handle the business logic for the domain package.

import (
	"context"

	"github.com/rs/zerolog"
	dmstores "github.com/vapusdata-oss/aistudio/aistudio/datastoreops"
	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	services "github.com/vapusdata-oss/aistudio/aistudio/services"
	utils "github.com/vapusdata-oss/aistudio/aistudio/utils"
	grpcops "github.com/vapusdata-oss/aistudio/core/serviceops/grpcops"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	dpb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

type OrganizationController struct {
	dpb.UnimplementedStudioServiceServer
	StudioServices *services.StudioServices
	Logger         zerolog.Logger
}

var OrganizationControllerManager *OrganizationController

func NewOrganizationController() *OrganizationController {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "OrganizationController")
	l.Info().Msg("Organization Controller initialized")
	return &OrganizationController{
		StudioServices: services.StudioServicesManager,
		Logger:         l,
	}
}

func InitOrganizationController() {
	if OrganizationControllerManager == nil {
		OrganizationControllerManager = NewOrganizationController()
	}
}

func (nc *OrganizationController) OrganizationManager(ctx context.Context, request *dpb.OrganizationManagerRequest) (*dpb.OrganizationResponse, error) {
	agent, err := nc.StudioServices.NewOrganizationAgent(ctx, request, nil)
	if err != nil {
		return nil, grpcops.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	err = agent.Act("")
	if err != nil {
		return nil, grpcops.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response := agent.GetResult()
	agent.LogAgent()
	response.DmResp = grpcops.HandleDMResponse(ctx, "VdcDeployment action executed successfully", "200")
	return response, nil
}

func (nc *OrganizationController) OrganizationGetter(ctx context.Context, request *dpb.OrganizationGetterRequest) (*dpb.OrganizationResponse, error) {
	agent, err := nc.StudioServices.NewOrganizationAgent(ctx, nil, request)
	if err != nil {
		return nil, grpcops.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	err = agent.Act("")
	if err != nil {
		return nil, grpcops.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response := agent.GetResult()
	agent.LogAgent()
	response.DmResp = grpcops.HandleDMResponse(ctx, "VdcDeployment action executed successfully", "200")
	return response, nil
}

func (dmc *OrganizationController) StudioPublicInfo(ctx context.Context, request *mpb.EmptyRequest) (*pb.StudioPublicInfoResponse, error) {
	accountInfo := dmstores.AccountPool
	if accountInfo == nil {
		return &pb.StudioPublicInfoResponse{}, grpcstatus.Error(grpccodes.NotFound, "Account details not found")
	}
	return &pb.StudioPublicInfoResponse{
		Logo:        accountInfo.Profile.Logo,
		AccountName: accountInfo.Name,
		Favicon:     accountInfo.Profile.Favicon,
	}, nil
}

func (dmc *OrganizationController) AccountManager(ctx context.Context, request *pb.AccountManagerRequest) (*pb.AccountResponse, error) {
	agent, err := dmc.StudioServices.NewAccountAgent(ctx, request)
	if err != nil {
		dmc.Logger.Err(err).Ctx(ctx).Msg("Error while initializing AccountManager")
		return nil, grpcops.HandleGrpcError(err, grpccodes.Internal)
	}
	err = agent.Act("")
	if err != nil {
		dmc.Logger.Err(err).Ctx(ctx).Msg("Error while performing the action on AccountManager")
		return nil, grpcops.HandleGrpcError(err, grpccodes.Internal)
	}
	agent.LogAgent()
	return &pb.AccountResponse{
		Output: agent.GetResponse().ConvertToPb(),
		DmResp: grpcops.HandleDMResponse(ctx, utils.ACCOUNT_CREATED, "200"),
	}, nil
}

func (dmc *OrganizationController) AccountGetter(ctx context.Context, request *mpb.EmptyRequest) (*pb.AccountResponse, error) {
	accountInfo, err := dmc.StudioServices.GetAccount(ctx)
	if err != nil {
		dmc.Logger.Err(err).Ctx(ctx).Msg("Error while getting account info")
		return nil, grpcops.HandleGrpcError(err, grpccodes.Internal)
	}
	return &pb.AccountResponse{
		Output: accountInfo.ConvertToPb(),
		DmResp: grpcops.HandleDMResponse(ctx, utils.ACCOUNT_CREATED, "200"),
	}, nil
}
