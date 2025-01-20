package svcops

import (
	"context"
	"log"

	"github.com/rs/zerolog"
	"github.com/vapusdata-oss/aistudio/core/models"
	grpcops "github.com/vapusdata-oss/aistudio/core/serviceops/grpcops"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (x *VapusSvcInternalClients) GetUser(ctx context.Context, userId string, logger zerolog.Logger, retryCount int) (*models.Users, error) {
	if x == nil {
		return nil, ErrUserConnNotInitialized
	}
	if x.UserConn == nil {
		err := x.PingTestAndReconnect(ctx, x.PlDns, logger)
		if err != nil {
			return nil, ErrUserConnNotInitialized
		}
	}
	resp, err := x.UserConn.UserGetter(grpcops.SwapNewContextWithAuthToken(ctx), &pb.UserGetterRequest{
		UserId: userId,
	})
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while connecting to Studio Service, retrying with new connection.")
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if status.Code(err).String() == codes.Unavailable.String() {
				log.Println("Retry count", retryCount)
				if retryCount > 3 {
					return nil, Erruser404
				}
				retryCount++
				logger.Err(err).Ctx(ctx).Msgf("error while calling platform server toget users, retrying with new connection. Count = %v", retryCount)
				return x.GetUser(ctx, userId, logger, retryCount)
			}
		}
	}
	if err != nil {
		return nil, Erruser404
	}
	if resp.Output == nil || len(resp.Output.Users) == 0 {
		return nil, status.Error(codes.NotFound, "User not found")
	}
	return (&models.Users{}).ConvertFromPb(resp.Output.Users[0]), nil

}
