package clients

import (
	"github.com/labstack/echo/v4"
	"github.com/vapusdata-oss/aistudio/webapp/pkgs"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

func (s *GrpcClient) GetAccountInfo(eCtx echo.Context) (*pb.AccountResponse, error) {
	return pkgs.VapusSvcInternalClientManager.SvcConn.AccountGetter(s.SetAuthCtx(eCtx), &mpb.EmptyRequest{})
}

func (x *GrpcClient) GetDashboard(eCtx echo.Context) *mpb.EmptyRequest {
	return &mpb.EmptyRequest{}
}

func (x *GrpcClient) GetUserInfo(eCtx echo.Context, userId string) (*mpb.User, map[string]string, error) {
	req := &pb.UserGetterRequest{
		Action: pb.UserGetterActions_GET_USER,
	}
	if userId != "" {
		req.UserId = userId
	}
	result, err := x.UserConn.UserGetter(x.SetAuthCtx(eCtx), req)

	if err != nil || result.Output == nil || len(result.Output.Users) == 0 {
		x.logger.Err(err).Msg("error while getting user info")
		return nil, nil, err
	}
	return result.Output.Users[0], result.OrganizationMap, nil
}

func (x *GrpcClient) GetMyOrganizationUsers(eCtx echo.Context) []*mpb.User {
	result, err := x.UserConn.UserGetter(x.SetAuthCtx(eCtx), &pb.UserGetterRequest{
		Action: pb.UserGetterActions_LIST_USERS,
	})

	if err != nil || result.Output == nil || len(result.Output.Users) == 0 {
		x.logger.Err(err).Msg("error while getting user info")
		return []*mpb.User{}
	}
	return result.Output.Users
}

func (x *GrpcClient) GetStudioUsers(eCtx echo.Context) []*mpb.User {
	result, err := x.UserConn.UserGetter(x.SetAuthCtx(eCtx), &pb.UserGetterRequest{
		Action: pb.UserGetterActions_LIST_STUDIO_USERS,
	})

	if err != nil || result.Output == nil || len(result.Output.Users) == 0 {
		x.logger.Err(err).Msg("error while getting user info")
		return []*mpb.User{}
	}
	return result.Output.Users
}

func (x *GrpcClient) ListPlugins(eCtx echo.Context) []*mpb.Plugin {
	result, err := x.PluginServiceClient.PluginGetter(x.SetAuthCtx(eCtx), &pb.PluginGetterRequest{})
	if err != nil {
		x.logger.Err(err).Msg("error while getting plugin list")
		return []*mpb.Plugin{}
	}
	return result.Output
}

func (x *GrpcClient) GetPlugin(eCtx echo.Context, id string) *mpb.Plugin {
	result, err := x.PluginServiceClient.PluginGetter(x.SetAuthCtx(eCtx), &pb.PluginGetterRequest{
		PluginId: id,
	})
	if err != nil || len(result.Output) == 0 {
		x.logger.Err(err).Msg("error while getting plugin info")
		return nil
	}
	return result.Output[0]
}

func (x *GrpcClient) StudioPublicInfo(eCtx echo.Context, id string) *pb.StudioPublicInfoResponse {
	result, err := x.SvcConn.StudioPublicInfo(x.SetAuthCtx(eCtx), &mpb.EmptyRequest{})
	if err != nil {
		x.logger.Err(err).Msg("error while getting Studio public info")
		return result
	}
	return result
}
