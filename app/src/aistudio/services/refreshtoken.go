package services

import (
	"context"
	"log"
	"time"

	utils "github.com/vapusdata-oss/aistudio/aistudio/utils"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	"github.com/vapusdata-oss/aistudio/core/globals"
	"github.com/vapusdata-oss/aistudio/core/serviceops/grpcops"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
	grpccodes "google.golang.org/grpc/codes"
)

func (dms *StudioServices) RefreshTokenAgentHandler(ctx context.Context, managerRequest *pb.RefreshTokenManagerRequest, getterRequest *pb.RefreshTokenGetterRequest) (*pb.RefreshTokenResponse, error) {
	validTill := time.Now().Add(globals.DEFAULT_AT_VALIDITY)
	log.Println("validTill: ", validTill)
	vapusStudioClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		dms.logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, encryption.ErrInvalidJWTClaims
	}
	if managerRequest != nil {
		switch managerRequest.GetUtility() {
		case pb.RefreshTokenAgentUtility_GENERATE_REFRESH_TOKEN:
			_, _, err := dms.GenerateStudioAccessToken(ctx, "", "", "request.GetOrganization()", vapusStudioClaim)
			if err != nil {
				return nil, grpcops.HandleGrpcError(dmerrors.DMError(utils.Err401, err), grpccodes.Unauthenticated)
			}
			return nil, nil
		case pb.RefreshTokenAgentUtility_REVOKE_REFRESH_TOKEN:
			_, _, err := dms.GenerateStudioAccessToken(ctx, "", "", "request.GetDomai n()", vapusStudioClaim)
			if err != nil {
				return nil, grpcops.HandleGrpcError(dmerrors.DMError(utils.Err401, err), grpccodes.Unauthenticated)
			}
			return nil, nil
		default:
			return nil, dmerrors.DMError(utils.ErrInvalidAccessTokenAgentUtility, nil)
		}
	}
	return nil, nil
}
