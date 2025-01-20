package clients

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	serviceops "github.com/vapusdata-oss/aistudio/core/serviceops"
	"github.com/vapusdata-oss/aistudio/webapp/pkgs"
	"github.com/vapusdata-oss/aistudio/webapp/utils"
)

type GrpcClient struct {
	*serviceops.VapusSvcInternalClients
	logger zerolog.Logger
}

var GrpcClientManager *GrpcClient

func NewGrpcClient() *GrpcClient {
	logger := pkgs.GetSubDMLogger("webapp", "grpcClients")
	cl, err := serviceops.SetupVapusSvcInternalClients(context.Background(), pkgs.NetworkConfigManager, "", logger)
	if err != nil {
		logger.Err(err).Msg("error while initializing vapus svc internal clients.")
	}
	return &GrpcClient{
		VapusSvcInternalClients: cl,
		logger:                  logger,
	}
}

func InitGrpcClient() {
	if GrpcClientManager == nil {
		GrpcClientManager = NewGrpcClient()
	}
}

func (s *GrpcClient) Close() {
	s.VapusSvcInternalClients.Close()
}

func (s *GrpcClient) SetAuthCtx(eCtx echo.Context) context.Context {
	token, err := utils.GetCookie(eCtx, pkgs.ACCESS_TOKEN)
	if err != nil || token == "" {
		return eCtx.Request().Context()
	}
	return utils.GetBearerCtx(context.Background(), token)
}
