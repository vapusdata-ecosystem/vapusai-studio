package services

import (
	"context"
	"fmt"
	"slices"

	dmstores "github.com/vapusdata-oss/aistudio/aistudio/datastoreops"
	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	utils "github.com/vapusdata-oss/aistudio/aistudio/utils"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	globals "github.com/vapusdata-oss/aistudio/core/globals"
	"github.com/vapusdata-oss/aistudio/core/models"
	coreutils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

type OrganizationAgent struct {
	result         *pb.OrganizationResponse
	dbStore        *dmstores.DMStore
	StudioServices *StudioServices
	managerRequest *pb.OrganizationManagerRequest
	getterRequest  *pb.OrganizationGetterRequest
	organization   *models.Organization
	errors         []error
	*models.VapusInterfaceAgentBase
}

func (x *OrganizationAgent) GetResult() *pb.OrganizationResponse {
	x.FinishAt = coreutils.GetEpochTime()
	return x.result
}

func (x *OrganizationAgent) LogAgent() {
	x.Logger.Info().Msgf("OrganizationAgent - %v action started at %v and finished at %v with status %v", x.AgentId, x.InitAt, x.FinishAt, x.Status)
}

func (x *StudioServices) NewOrganizationAgent(ctx context.Context, managerRequest *pb.OrganizationManagerRequest, getterRequest *pb.OrganizationGetterRequest) (*OrganizationAgent, error) {
	vapusStudioClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		x.logger.Error().Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}

	organization, err := x.DMStore.GetOrganization(ctx, vapusStudioClaim[encryption.ClaimOrganizationKey], vapusStudioClaim)
	if err != nil {
		return nil, dmerrors.DMError(utils.ErrOrganization404, err)
	}

	agent := &OrganizationAgent{
		managerRequest: managerRequest,
		getterRequest:  getterRequest,
		dbStore:        x.DMStore,
		StudioServices: x,
		organization:   organization,
		result:         &pb.OrganizationResponse{Output: &pb.OrganizationResponse_OrganizationOutput{Users: make([]*pb.OrganizationResponse_OrganizationOutput_OrganizationUsers, 0)}},
		VapusInterfaceAgentBase: &models.VapusInterfaceAgentBase{
			InitAt:    coreutils.GetEpochTime(),
			Ctx:       ctx,
			CtxClaim:  vapusStudioClaim,
			Action:    managerRequest.GetActions().String(),
			AgentType: globals.DOMAINAGENT.String(),
		},
	}
	if managerRequest != nil {
		agent.Action = managerRequest.GetActions().String()
	} else if getterRequest != nil {
		if getterRequest.GetOrgId() == "" {
			agent.Action = pb.OrganizationAgentActions_LIST_ORGS.String()
		} else {
			agent.Action = pb.OrganizationAgentActions_DESCRIBE_ORG.String()
		}
	} else {
		agent.Action = ""
	}
	agent.SetAgentId()
	agent.Logger = pkgs.GetSubDMLogger(globals.DOMAINAGENT.String(), agent.AgentId)
	return agent, nil
}

func (x *OrganizationAgent) Act(action string) error {

	if action != "" {
		x.Action = action
	}
	switch x.Action {
	case pb.OrganizationAgentActions_CONFIGURE_ORG.String():
		x.organization = utils.DmNodeToObj(x.managerRequest)
		return x.configureOrganization()
	case pb.OrganizationAgentActions_LIST_ORGS.String():
		return x.listOrganizations()
	case pb.OrganizationAgentActions_PATCH_ORG.String():
		x.organization = utils.DmNodeToObj(x.managerRequest)
		return x.patchOrganization()
	case pb.OrganizationAgentActions_ADD_ORG_USER.String():
		x.organization = utils.DmNodeToObj(x.managerRequest)
		return x.addOrganizationUsers()
	case pb.OrganizationAgentActions_DESCRIBE_ORG.String():
		return x.describeOrganization()
	default:
		return dmerrors.DMError(utils.ErrInvalidOrganizationAction, nil) //nolint:wrapcheck
	}
}

func (x *OrganizationAgent) addOrganizationUsers() error {
	if x.organization.VapusID != x.CtxClaim[encryption.ClaimOrganizationKey] {
		return dmerrors.DMError(utils.ErrInvalidOrganizationRequested, nil)
	}
	for _, user := range x.managerRequest.GetUsers() {
		if !slices.Contains(x.organization.Users, user.GetUserId()) {
			err := x.attachOrganization2User(user.GetUserId(), user.InviteIfNotFound, &models.UserOrganizationRole{
				OrganizationId: x.organization.VapusID,
				RoleArns:       user.GetRole(),
			})
			if err != nil {
				x.Logger.Err(err).Ctx(x.Ctx).Msgf("error while mapping user %s to this organization %v", user, x.organization)
				return dmerrors.DMError(utils.ErrUserOrganizationMapping, err) //nolint:wrapcheck
			}
		}
	}
	return nil
}

func (x *OrganizationAgent) listOrganizationUsers() error {
	var err error
	organizationObj, err := x.dbStore.GetOrganization(x.Ctx, x.organization.VapusID, x.CtxClaim)
	if err != nil {
		return dmerrors.DMError(utils.ErrInvalidOrganizationRequested, err) //nolint:wrapcheck
	}
	filter := fmt.Sprintf(`organization_roles @> '[{"organizationId": "%s"}]'`, organizationObj.VapusID)
	users, err := x.dbStore.ListUsers(x.Ctx, filter, x.CtxClaim)
	if err != nil {
		return dmerrors.DMError(utils.ErrInvalidOrganizationRequested, err) //nolint:wrapcheck
	}
	x.result.Output.Users = []*pb.OrganizationResponse_OrganizationOutput_OrganizationUsers{{
		Users: utils.DmUToPb(users, x.organization.VapusID),
		Org:   organizationObj.VapusID,
	}}
	x.result.Output.Orgs = utils.DmNToPb([]*models.Organization{organizationObj})
	return nil
}

func (x *OrganizationAgent) configureOrganization() error {
	var err error
	if x.organization == nil {
		return dmerrors.DMError(utils.ErrInvalidAddOrganizationRequest, nil)
	}
	if x.organization.OrganizationType == mpb.OrganizationTypes_SERVICE_ORG.String() {
		return dmerrors.DMError(utils.ErrCannotCreateServiceOrg, nil)
	}
	x.organization, err = organizationConfigureTool(x.Ctx, x.organization, x.dbStore, x.Logger, x.CtxClaim)
	if err != nil {
		return dmerrors.DMError(err, nil)
	}
	err = x.attachOrganization2User(x.CtxClaim[encryption.ClaimUserIdKey], true, &models.UserOrganizationRole{
		OrganizationId: x.organization.VapusID,
		RoleArns:       []string{mpb.StudioUserRoles_DOMAIN_OWNERS.String()},
	})
	if err != nil {
		x.Logger.Err(err).Ctx(x.Ctx).Msgf("error while mapping user as organization owner to this organization %v", x.organization)
		return dmerrors.DMError(utils.ErrUserOrganizationMapping, err) //nolint:wrapcheck
	}
	x.result.Output.Orgs = utils.DmNToPb([]*models.Organization{x.organization})
	return nil
}

func (x *OrganizationAgent) attachOrganization2User(userId string, invite bool, obj *models.UserOrganizationRole) error {
	var user *models.Users
	var err error
	user, err = x.dbStore.GetUser(x.Ctx, userId, x.CtxClaim)
	if err != nil {
		if !invite {
			return dmerrors.DMError(utils.ErrInvalidUserRequested, err) //nolint:wrapcheck
		}
		x.Logger.Err(err).Ctx(x.Ctx).Msgf("User with id %s not found, inviting.", userId)
		userAgent, err := x.StudioServices.NewUserManagerAgent(x.Ctx, &pb.UserManagerRequest{
			Action: pb.UserAgentActions_INVITE_USERS,
			Spec: &mpb.User{
				Email: userId,
			},
			Organization: obj.OrganizationId,
			RoleArn:      obj.RoleArns,
		}, nil)
		if err != nil {
			return dmerrors.DMError(utils.ErrInvalidUserRequested, err) //nolint:wrapcheck
		}
		err = userAgent.Act(pb.UserAgentActions_INVITE_USERS.String())
		if err != nil {
			return dmerrors.DMError(utils.ErrInvalidUserRequested, err) //nolint:wrapcheck
		}
		userAgent.LogAgent()
		result := userAgent.GetResult()
		if len(result.GetOutput().Users) > 0 {
			user, err = x.dbStore.GetUser(x.Ctx, result.GetOutput().Users[0].UserId, x.CtxClaim)
			if err != nil {
				return dmerrors.DMError(utils.ErrInvalidUserRequested, err) //nolint:wrapcheck
			}
		}
	}
	if len(user.OrganizationRoles) == 0 {
		user.OrganizationRoles = make([]*models.UserOrganizationRole, 0)
	}
	user.OrganizationRoles = append(user.OrganizationRoles, obj)
	return x.dbStore.PutUser(x.Ctx, user, x.CtxClaim)
}

func (x *OrganizationAgent) patchOrganization() error {
	var err error
	newObj := utils.DmNodeToObj(x.managerRequest)
	if x.organization.OrganizationType == mpb.OrganizationTypes_SERVICE_ORG.String() {
		return dmerrors.DMError(utils.ErrCannotCreateServiceOrg, nil)
	}
	x.organization.PreSaveUpdate(x.CtxClaim[encryption.ClaimUserIdKey])
	for _, nwUser := range x.managerRequest.GetUsers() {
		if !slices.Contains(x.organization.Users, nwUser.GetUserId()) {
			err := x.attachOrganization2User(nwUser.GetUserId(), nwUser.InviteIfNotFound, &models.UserOrganizationRole{
				OrganizationId: x.organization.VapusID,
				RoleArns:       nwUser.GetRole(),
			})
			if err != nil {
				x.Logger.Err(err).Ctx(x.Ctx).Msgf("error while mapping user %s to this organization %v", nwUser, x.organization)
			}
		}
	}

	x.organization.DisplayName = newObj.DisplayName
	err = x.dbStore.PutOrganization(x.Ctx, x.organization, x.CtxClaim)
	if err != nil {
		x.Logger.Err(err).Ctx(x.Ctx).Msgf("error while configuring organization %v", x.organization)
		return dmerrors.DMError(utils.ErrCreateOrganization, err) //nolint:wrapcheck
	}

	x.result.Output.Orgs = utils.DmNToPb([]*models.Organization{x.organization})
	return nil
}

func (x *OrganizationAgent) listOrganizations() error {
	var filter string = ""
	dmIds := utils.GetFilterParams(x.getterRequest.GetSearchParam(), globals.OrganizationSK.String())
	if len(dmIds) > 0 {
		filter = fmt.Sprintf("vapus_id IN (%s)", dmIds)
	} else {
		filter = ""
	}
	organizations, err := x.dbStore.ListOrganizations(x.Ctx, filter, x.CtxClaim)
	if err != nil {
		return dmerrors.DMError(utils.ErrInvalidOrganizationRequested, err) //nolint:wrapcheck
	}
	x.result.Output.Orgs = utils.DmNToPb(organizations)
	return nil
}

func (x *OrganizationAgent) describeOrganization() error {
	organization, err := x.dbStore.GetOrganization(x.Ctx, x.CtxClaim[encryption.ClaimOrganizationKey], x.CtxClaim)
	if err != nil {
		return dmerrors.DMError(utils.ErrInvalidOrganizationRequested, err) //nolint:wrapcheck
	}
	x.result.Output.Orgs = utils.DmNToPb([]*models.Organization{organization})
	return nil
}
