package services

import (
	"context"

	dmstores "github.com/vapusdata-oss/aistudio/aistudio/datastoreops"
	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	utils "github.com/vapusdata-oss/aistudio/aistudio/utils"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	globals "github.com/vapusdata-oss/aistudio/core/globals"
	models "github.com/vapusdata-oss/aistudio/core/models"
	coreutils "github.com/vapusdata-oss/aistudio/core/utils"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

type AuthzManagerAgent struct {
	managerRequest *pb.AuthzManagerRequest
	getterRequest  *pb.AuthzGetterRequest
	dbStore        *dmstores.DMStore
	*StudioServices
	user   *models.Users
	result *pb.AuthzResponse
	org    *models.Organization
	*models.VapusInterfaceAgentBase
}

func (x *AuthzManagerAgent) GetResult() *pb.AuthzResponse {
	x.FinishAt = coreutils.GetEpochTime()
	return x.result
}

func (x *AuthzManagerAgent) LogAgent() {
	x.Logger.Info().Msgf("AuthzManagerAgent - %v action started at %v and finished at %v with status %v", x.AgentId, x.InitAt, x.FinishAt, x.Status)
}

func (x *StudioServices) NewAuthzManagerAgent(ctx context.Context, managerRequest *pb.AuthzManagerRequest, getterRequest *pb.AuthzGetterRequest) (*AuthzManagerAgent, error) {
	vapusStudioClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		x.logger.Error().Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}

	organization, err := x.DMStore.GetOrganization(ctx, vapusStudioClaim[encryption.ClaimOrganizationKey], vapusStudioClaim)
	if err != nil {
		return nil, dmerrors.DMError(utils.ErrOrganization404, err)
	}
	agent := &AuthzManagerAgent{
		result:         &pb.AuthzResponse{Output: &pb.AuthzResponse_AuthzRoles{}},
		managerRequest: managerRequest,
		getterRequest:  getterRequest,
		dbStore:        x.DMStore,
		StudioServices: x,
		org:            organization,
		VapusInterfaceAgentBase: &models.VapusInterfaceAgentBase{
			InitAt:    coreutils.GetEpochTime(),
			Ctx:       ctx,
			CtxClaim:  vapusStudioClaim,
			AgentType: globals.USERMANAGERAGENT.String(),
		},
	}
	agent.SetAgentId()
	if managerRequest != nil {
		agent.Action = managerRequest.GetAction().String()
	} else if getterRequest != nil {
		if getterRequest.GetRoleArn() == "" {
			agent.Action = "pb.AuthzManagerRequest_GET_AUTHZ_ROLE.String()"
		} else {
			agent.Action = "pb.AuthzManagerRequest_LIST_AUTHZ_ROLES.String()"
		}
	} else {
		agent.Action = ""
	}
	agent.Logger = pkgs.GetSubDMLogger(globals.DATAPRODUCTAGENT.String(), agent.AgentId)
	return agent, nil
}

func (x *AuthzManagerAgent) Act(action string) error {
	// switch x.Action {
	// case pb.AuthzManagerRequest_CREATE_AUTHZ_ROLE.String():
	// 	return x.CreateAuthzRole()
	// case pb.AuthzManagerRequest_GET_AUTHZ_ROLE.String():
	// 	return x.GetAuthzRole()
	// case pb.AuthzManagerRequest_UPDATE_AUTHZ_ROLE.String():
	// 	return x.UpdateAuthzRole()
	// case pb.AuthzManagerRequest_DELETE_AUTHZ_ROLE.String():
	// 	return x.DeleteAuthzRole()
	// case pb.AuthzManagerRequest_LIST_AUTHZ_ROLES.String():
	// 	return x.ListAuthzRoles()
	// default:
	// 	return dmerrors.DMError(utils.ErrInvalidAction, nil)
	// }
	return nil
}
