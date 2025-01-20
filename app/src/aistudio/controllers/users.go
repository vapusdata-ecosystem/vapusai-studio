package dmcontrollers

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
	grpccodes "google.golang.org/grpc/codes"

	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	dmsvc "github.com/vapusdata-oss/aistudio/aistudio/services"
	utils "github.com/vapusdata-oss/aistudio/aistudio/utils"
	dmauthn "github.com/vapusdata-oss/aistudio/core/authn"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	"github.com/vapusdata-oss/aistudio/core/globals"
	grpcops "github.com/vapusdata-oss/aistudio/core/serviceops/grpcops"
	dmutils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

type VapusDataUsers struct {
	pb.UnimplementedUserManagementServiceServer
	validator      *dmutils.DMValidator
	StudioServices *dmsvc.StudioServices
	logger         zerolog.Logger
}

var VapusDataUsersManager *VapusDataUsers

func NewVapusDataUsers() *VapusDataUsers {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "VapusDataUsers")
	validator, err := dmutils.NewDMValidator()
	if err != nil {
		l.Panic().Err(err).Msg("Error while loading validator")
	}

	l.Info().Msg("VapusDataUsers Controller initialized")
	return &VapusDataUsers{
		validator:      validator,
		logger:         l,
		StudioServices: dmsvc.StudioServicesManager,
	}
}

func InitVapusDataUsers() {
	if VapusDataUsersManager == nil {
		VapusDataUsersManager = NewVapusDataUsers()
	}
}

func (dmc *VapusDataUsers) AccessTokenInterface(ctx context.Context, request *pb.AccessTokenInterfaceRequest) (*pb.AccessTokenResponse, error) {
	dmc.logger.Info().Msg("Generating platform access token.......")
	return dmc.StudioServices.AccessTokenAgentHandler(ctx, request)
}

func (dmc *VapusDataUsers) RegisterUser(ctx context.Context, request *pb.RegisterUserRequest) (*pb.AccessTokenResponse, error) {
	dmc.logger.Info().Msg("Registering user.......")
	return dmc.StudioServices.SignupHandler(ctx, request)
}

func (dmc *VapusDataUsers) UserManager(ctx context.Context, request *pb.UserManagerRequest) (*pb.UserResponse, error) {
	agent, err := dmc.StudioServices.NewUserManagerAgent(ctx, request, nil)
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
func (dmc *VapusDataUsers) UserGetter(ctx context.Context, request *pb.UserGetterRequest) (*pb.UserResponse, error) {
	agent, err := dmc.StudioServices.NewUserManagerAgent(ctx, nil, request)
	if err != nil {
		return nil, grpcops.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	err = agent.Act("")
	if err != nil {
		return nil, grpcops.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response := agent.GetResult()
	agent.LogAgent()
	log.Println("response -------------------------||||||||||||||||||||", response.OrganizationMap)
	response.DmResp = grpcops.HandleDMResponse(ctx, "VdcDeployment action executed successfully", "200")
	return response, nil
}

func (dmc *VapusDataUsers) RefreshTokenManager(ctx context.Context, request *pb.RefreshTokenManagerRequest) (*pb.RefreshTokenResponse, error) {
	response, err := dmc.StudioServices.RefreshTokenAgentHandler(ctx, request, nil)
	if err != nil {
		return nil, grpcops.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response.DmResp = grpcops.HandleDMResponse(ctx, "VdcDeployment action executed successfully", "200")
	return response, nil
}
func (dmc *VapusDataUsers) RefreshTokenGetter(ctx context.Context, request *pb.RefreshTokenGetterRequest) (*pb.RefreshTokenResponse, error) {
	dmc.logger.Info().Msg("Generating platform access token.......")
	response, err := dmc.StudioServices.RefreshTokenAgentHandler(ctx, nil, request)
	if err != nil {
		return nil, grpcops.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response.DmResp = grpcops.HandleDMResponse(ctx, "VdcDeployment action executed successfully", "200")
	return response, nil
}

func (dmc *VapusDataUsers) AuthzManager(ctx context.Context, request *pb.AuthzManagerRequest) (*pb.AuthzResponse, error) {
	dmc.logger.Info().Msg("Generating platform access token.......")
	agent, err := dmc.StudioServices.NewAuthzManagerAgent(ctx, request, nil)
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

func (dmc *VapusDataUsers) AuthzGetter(ctx context.Context, request *pb.AuthzGetterRequest) (*pb.AuthzResponse, error) {
	dmc.logger.Info().Msg("Generating platform access token.......")
	agent, err := dmc.StudioServices.NewAuthzManagerAgent(ctx, nil, request)
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

func (dmc *VapusDataUsers) LoginHandler(ctx context.Context, request *mpb.EmptyRequest) (*pb.LoginHandlerResponse, error) {
	state, err := dmutils.GenerateRandomState()
	if err != nil {
		dmc.logger.Err(err).Msg(utils.ErrAuthenticatorInitFailed.Error())
		return nil, grpcops.HandleGrpcError(dmerrors.DMError(utils.ErrLoginFailed, nil), grpccodes.Internal)
	}

	return &pb.LoginHandlerResponse{
		LoginUrl:    pkgs.AuthnManager.Authenticator.AuthCodeURL(state, oauth2.SetAuthURLParam("prompt", "login")),
		CallbackUrl: pkgs.AuthnManager.RedirectURL,
		RedirectUri: pkgs.AuthnManager.Authenticator.Organization,
	}, nil

}

func (dmc *VapusDataUsers) LoginCallback(ctx context.Context, request *pb.LoginCallBackRequest) (*pb.AccessTokenResponse, error) {
	validTill := time.Now().Add(globals.DEFAULT_AT_VALIDITY)
	dmc.logger.Info().Msg("Exchanging code for token.......")
	dmc.logger.Info().Msgf("Callback URL: %v", request.GetHost())
	token, err := pkgs.AuthnManager.Authenticator.Exchange(ctx, request.GetCode(), oauth2.SetAuthURLParam("redirect_uri", request.GetHost()))
	if err != nil {
		dmc.logger.Err(err).Msg(dmauthn.ErrTokenExchangeFailed.Error())
		return nil, grpcops.HandleGrpcError(dmerrors.DMError(dmauthn.ErrTokenExchangeFailed, nil), grpccodes.Unauthenticated)

	}
	_, err = pkgs.AuthnManager.Authenticator.VerifyIDToken(ctx, token)

	if err != nil {
		dmc.logger.Err(err).Msg(dmauthn.ErrIDTokenVerificationFailed.Error())
		return nil, grpcops.HandleGrpcError(dmerrors.DMError(dmauthn.ErrIDTokenVerificationFailed, nil), grpccodes.Unauthenticated)

	}
	accessToken, scope, err := dmc.StudioServices.GenerateStudioAccessToken(ctx, token.Extra("id_token").(string), "", "", nil)
	if err != nil {
		log.Println("error while generating token", err, err == utils.ErrUser404)
		if errors.Is(err, utils.ErrUser404) {
			if pkgs.ServiceConfigManager.SelfSignup {
				log.Println("user not found,++++++++++++++>>>>>>>>>>>>>", token.Extra("id_token").(string))
				return dmc.StudioServices.SignupHandler(ctx, &pb.RegisterUserRequest{
					IdToken: token.Extra("id_token").(string),
				})
			}
		}
		dmc.logger.Err(err).Msg("error while generating token")
		return nil, grpcops.HandleGrpcError(dmerrors.DMError(utils.Err401, err), grpccodes.Unauthenticated)
	}
	return &pb.AccessTokenResponse{
		Token: &pb.AccessToken{
			AccessToken: accessToken,
			ValidTill:   validTill.Unix(),
			ValidFrom:   dmutils.GetEpochTime(),
			IdToken:     token.Extra("id_token").(string),
		},
		TokenScope: scope,
		DmResp:     grpcops.HandleDMResponse(ctx, utils.ACCESS_TOKEN_CREATED, "201"),
	}, nil
}
