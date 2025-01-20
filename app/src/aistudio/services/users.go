package services

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	dmstores "github.com/vapusdata-oss/aistudio/aistudio/datastoreops"
	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	utils "github.com/vapusdata-oss/aistudio/aistudio/utils"
	authn "github.com/vapusdata-oss/aistudio/core/authn"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	globals "github.com/vapusdata-oss/aistudio/core/globals"
	"github.com/vapusdata-oss/aistudio/core/models"
	"github.com/vapusdata-oss/aistudio/core/serviceops/grpcops"
	coreutils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
	grpccodes "google.golang.org/grpc/codes"
)

func (dms *StudioServices) AccessTokenAgentHandler(ctx context.Context, request *pb.AccessTokenInterfaceRequest) (*pb.AccessTokenResponse, error) {
	validTill := time.Now().Add(globals.DEFAULT_AT_VALIDITY)
	vapusStudioClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		dms.logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, encryption.ErrInvalidJWTClaims
	}
	switch request.GetUtility() {
	case pb.AccessTokenAgentUtility_ORG_LOGIN:
		accessToken, scope, err := dms.GenerateStudioAccessToken(ctx, request.GetIdToken(), "", request.GetOrganization(), vapusStudioClaim)
		if err != nil {
			log.Println("error while generating token", err, err == utils.ErrUser404)
			if errors.Is(err, utils.ErrUser404) {
				log.Println("user not found,++++++++++++++")
				return nil, grpcops.HandleGrpcError(dmerrors.DMError(utils.Err401, err), grpccodes.NotFound)
			}
			dms.logger.Err(err).Msg("error while generating token")
			return nil, grpcops.HandleGrpcError(dmerrors.DMError(utils.Err401, err), grpccodes.Unauthenticated)
		}
		return &pb.AccessTokenResponse{
			Token: &pb.AccessToken{
				AccessToken: accessToken,
				ValidTill:   validTill.Unix(),
			},
			TokenScope: scope,
		}, nil
	case pb.AccessTokenAgentUtility_REFRESH_TOKEN_LOGIN:
		accessToken, scope, err := dms.GenerateStudioAccessToken(ctx, "", request.GetRefreshToken(), request.GetOrganization(), vapusStudioClaim)
		if err != nil {
			dms.logger.Err(err).Msg("error while generating token")
			return nil, grpcops.HandleGrpcError(dmerrors.DMError(utils.Err401, err), grpccodes.Unauthenticated)
		}
		return &pb.AccessTokenResponse{
			Token: &pb.AccessToken{
				AccessToken: accessToken,
				ValidTill:   validTill.Unix(),
			},
			TokenScope: scope,
		}, nil
	default:
		return nil, dmerrors.DMError(utils.ErrInvalidAccessTokenAgentUtility, nil)
	}
}

func (dms *StudioServices) SignupHandler(ctx context.Context, request *pb.RegisterUserRequest) (*pb.AccessTokenResponse, error) {
	validTill := time.Now().Add(globals.DEFAULT_AT_VALIDITY)
	vapusStudioClaim := map[string]string{
		encryption.ClaimAccountKey: dmstores.AccountPool.VapusID,
	}
	claims, err := authn.ValidateOIDCAuth(request.GetIdToken(), dms.logger)
	if err != nil {
		dms.logger.Err(err).Msg("invalid token, validation failed")
		return nil, dmerrors.DMError(err, nil)
	}
	lu, err := claimToLocal(claims)
	if err != nil {
		dms.logger.Err(err).Msg("invalid token, relevant claims are missing")
		return nil, dmerrors.DMError(err, nil)
	}
	userObj, err := dms.DMStore.CreateUser(ctx, lu, []string{
		mpb.StudioUserRoles_PLATFORM_USERS.String(),
	}, nil, vapusStudioClaim)
	if err != nil {
		dms.logger.Err(err).Msg("error while creating user")
		return nil, dmerrors.DMError(err, nil)
	}
	if request.GetOrganization() == "" {
		dm := strings.Split(userObj.Email, "@")
		request.Organization = strings.ToTitle(dm[0]) + "." + strings.ToTitle(dm[1]) + " organization"
	}
	vapusStudioClaim[encryption.ClaimUserIdKey] = userObj.UserId
	organization, err := organizationConfigureTool(ctx, &models.Organization{
		Name: request.GetOrganization(),
	}, dms.DMStore, dms.logger, vapusStudioClaim)
	if err != nil {
		dms.logger.Err(err).Msg("error while configuring organization")
		return nil, dmerrors.DMError(err, nil)
	}
	userObj.OrganizationRoles = []*models.UserOrganizationRole{
		{
			OrganizationId: organization.VapusID,
			RoleArns:       []string{mpb.StudioUserRoles_PLATFORM_USERS.String()},
			InvitedOn:      coreutils.GetEpochTime(),
		},
	}
	dms.DMStore.PutUser(ctx, userObj, vapusStudioClaim)
	accessToken, scope, err := dms.GenerateStudioAccessToken(ctx, request.GetIdToken(), "", organization.VapusID, vapusStudioClaim)
	if err != nil {
		return nil, grpcops.HandleGrpcError(dmerrors.DMError(utils.Err401, err), grpccodes.Unauthenticated)
	}
	return &pb.AccessTokenResponse{
		Token: &pb.AccessToken{
			AccessToken: accessToken,
			ValidTill:   validTill.Unix(),
		},
		TokenScope: scope,
	}, nil
}

func (dms *StudioServices) GenerateStudioAccessToken(ctx context.Context, idToken, refreshToken, organization string, ctxClaim map[string]string) (string, mpb.AccessTokenScope, error) {
	var userObj *models.Users
	var organizationObj *models.Organization
	var err error
	var lu *pkgs.LocalUserM
	if idToken != "" {
		claims, err := authn.ValidateOIDCAuth(idToken, dms.logger)
		if err != nil {
			dms.logger.Err(err).Msg("invalid token, validation failed")
			return "", 0, dmerrors.DMError(err, nil)
		}
		lu, err = claimToLocal(claims)
		if err != nil {
			dms.logger.Err(err).Msg("invalid token, relevant claims are missing")
			return "", 0, dmerrors.DMError(err, nil)
		}
	} else if refreshToken != "" {
		refreshToken = encryption.GenerateSHA256(refreshToken, "")
		rtObj, err := dms.DMStore.GetStudioRTinfo(ctx, refreshToken, ctxClaim)
		if err != nil {
			dms.logger.Err(err).Msg("invalid token, refresh token doesn't exists")
			return "", 0, dmerrors.DMError(utils.ErrRefreshToken404, err)
		}
		if rtObj.ValidTill < coreutils.GetEpochTime() {
			dms.logger.Err(err).Msg("invalid token, refresh token expired")
			return "", 0, dmerrors.DMError(utils.ErrRefreshTokenExpired, nil)
		}
		if rtObj.Status != mpb.CommonStatus_ACTIVE.String() {
			dms.logger.Err(err).Msg("invalid token, refresh token is not active")
			return "", 0, dmerrors.DMError(utils.ErrRefreshTokenInactive, nil)
		}
		lu = &pkgs.LocalUserM{
			Email: rtObj.UserId,
		}
		organization = rtObj.Organization
	} else {
		lu = &pkgs.LocalUserM{
			Email: ctxClaim[encryption.ClaimUserIdKey],
		}
	}
	validTill := time.Now().Add(globals.DEFAULT_AT_VALIDITY)
	if organization == "" {
		userObj, err = dms.DMStore.GetOrUpdateUser(ctx, lu, true, true, ctxClaim)
		if err != nil {
			dms.logger.Err(err).Msg("invalid token, user doesn't exists")
			return idToken, 0, utils.ErrUser404
		}
		organizationObj, err = dms.DMStore.GetOrganization(ctx, userObj.OrganizationRoles[0].OrganizationId, nil)
		if err != nil {
			dms.logger.Err(err).Msg("invalid token, organization doesn't exists")
			return "", 0, dmerrors.DMError(utils.ErrOrganization404, err)
		}
		dms.logger.Info().Msgf("User obtained for generating token with default organization- %v ", userObj)
	} else {
		lu.Organization = organization
		organizationObj, err = dms.DMStore.GetOrganization(ctx, lu.Organization, nil)
		if err != nil {
			dms.logger.Err(err).Msg("invalid token, organization doesn't exists")
			return "", 0, dmerrors.DMError(utils.ErrOrganization404, err)
		}
		userObj, err = dms.DMStore.GetOrUpdateUser(ctx, lu, false, false, ctxClaim)
		if err != nil {
			dms.logger.Err(err).Msg("invalid token, user doesn't exists")
			return "", 0, dmerrors.DMError(err, nil)
		}

		if !userObj.IsValidUserByOrganization(organizationObj.VapusID) {
			dms.logger.Err(err).Msg("invalid request, user is not present in this organization")
			return "", 0, dmerrors.DMError(err, nil)
		}
		dms.logger.Info().Msgf("User obtained - %v ", userObj)
	}
	dms.logger.Info().Msgf("Organization obtained - %v ", organizationObj)
	dms.logger.Info().Msgf("User obtained here - %v ", userObj)
	tokenId, token, err := generateOrganizationAccessToken(userObj, organizationObj.VapusID, validTill)
	if err != nil {
		return "", 0, dmerrors.DMError(err, nil)
	}

	newCtx := context.TODO()
	go dms.DMStore.LogStudioJwtinfo(newCtx, &models.JwtLog{
		JwtId:        tokenId,
		UserId:       userObj.UserId,
		Organization: organizationObj.VapusID,
		Scope:        mpb.AccessTokenScope_DOMAIN_TOKEN.String(),
	}, ctxClaim)
	return token, mpb.AccessTokenScope_DOMAIN_TOKEN, nil
}
