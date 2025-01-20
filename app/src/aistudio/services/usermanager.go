package services

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"

	dmstores "github.com/vapusdata-oss/aistudio/aistudio/datastoreops"
	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	utils "github.com/vapusdata-oss/aistudio/aistudio/utils"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	globals "github.com/vapusdata-oss/aistudio/core/globals"
	models "github.com/vapusdata-oss/aistudio/core/models"
	coreutils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

type UserManagerAgent struct {
	managerRequest *pb.UserManagerRequest
	getterRequest  *pb.UserGetterRequest
	dbStore        *dmstores.DMStore
	*StudioServices
	user         *models.Users
	result       *pb.UserResponse
	organization *models.Organization
	*models.VapusInterfaceAgentBase
}

func (x *UserManagerAgent) GetResult() *pb.UserResponse {
	x.FinishAt = coreutils.GetEpochTime()
	return x.result
}

func (x *UserManagerAgent) LogAgent() {
	x.Logger.Info().Msgf("UserManagerAgent - %v action started at %v and finished at %v with status %v", x.AgentId, x.InitAt, x.FinishAt, x.Status)
}

func (x *StudioServices) NewUserManagerAgent(ctx context.Context, managerRequest *pb.UserManagerRequest, getterRequest *pb.UserGetterRequest) (*UserManagerAgent, error) {
	vapusStudioClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		x.logger.Error().Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}

	organization, err := x.DMStore.GetOrganization(ctx, vapusStudioClaim[encryption.ClaimOrganizationKey], vapusStudioClaim)
	if err != nil {
		return nil, dmerrors.DMError(utils.ErrOrganization404, err)
	}
	agent := &UserManagerAgent{
		result:         &pb.UserResponse{Output: &pb.UserResponse_VapusUser{}},
		managerRequest: managerRequest,
		getterRequest:  getterRequest,
		dbStore:        x.DMStore,
		StudioServices: x,
		organization:   organization,
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
		if getterRequest.GetUserId() == "" {
			agent.Action = getterRequest.GetAction().String()
		} else {
			agent.Action = pb.UserGetterActions_GET_USER.String()
		}
	} else {
		agent.Action = ""
	}
	agent.Logger = pkgs.GetSubDMLogger(globals.DATAPRODUCTAGENT.String(), agent.AgentId)
	return agent, nil
}

func (x *UserManagerAgent) Act(action string) error {
	if action != "" {
		x.Action = action
	}
	log.Println("Action", x.Action)
	switch x.Action {
	case pb.UserAgentActions_INVITE_USERS.String():
		return x.InviteUser()
	case pb.UserAgentActions_PATCH_USER.String():
		return x.PatchUser()
	case pb.UserGetterActions_GET_USER.String():
		return x.GetUser(x.getterRequest.GetUserId())
	case pb.UserGetterActions_LIST_USERS.String():
		return x.ListUsers()
	case pb.UserGetterActions_LIST_PLATFORM_USERS.String():
		return x.ListStudioUsers()
	default:
		return dmerrors.DMError(utils.ErrInvalidUserManagerAction, nil)
	}
}

func (x *UserManagerAgent) InviteUser() error {
	var user *models.Users
	var err error
	if x.managerRequest == nil {
		x.Logger.Error().Msg("error while getting user invite request")
		return dmerrors.DMError(utils.ErrUserInviteCreateFailed, nil)
	}
	userObj := utils.DmUPbToObj(x.managerRequest.GetSpec())
	if userObj == nil {
		x.Logger.Error().Msg("error while converting user object from pb")
		return dmerrors.DMError(utils.ErrUserInviteCreateFailed, nil)
	}
	for _, organizationRole := range userObj.OrganizationRoles {
		organizationRole.InvitedOn = coreutils.GetEpochTime()
	}
	exists := x.DMStore.UserInviteExists(x.Ctx, userObj.Email, x.CtxClaim)
	if exists {
		user, err := x.DMStore.GetUser(x.Ctx, userObj.Email, x.CtxClaim)
		if err != nil {
			x.Logger.Error().Msgf("invite already exists for email - %v, but not registered user", userObj.Email)
			return dmerrors.DMError(utils.ErrUserInviteExists, nil)
		}
		for _, organizationRole := range userObj.OrganizationRoles {
			for _, existOrganizationRole := range user.OrganizationRoles {
				if organizationRole.OrganizationId == existOrganizationRole.OrganizationId {
					existOrganizationRole.RoleArns = append(existOrganizationRole.RoleArns, organizationRole.RoleArns...)
				} else {

					user.OrganizationRoles = append(user.OrganizationRoles, organizationRole)
				}
			}
		}
		user.PreSaveUpdate(x.CtxClaim[encryption.ClaimUserIdKey])
		err = x.DMStore.PutUser(x.Ctx, user, x.CtxClaim)
		if err != nil {
			x.Logger.Error().Msgf("error while updating user - %v", err)
			return dmerrors.DMError(utils.ErrUserInviteCreateFailed, err)
		}
	} else {
		userObj.InvitedType = mpb.UserInviteType_INVITE_ACCESS.String()
		userObj.SetUserId()
		userObj.StudioRoles = []string{mpb.StudioUserRoles_PLATFORM_USERS.String()}
		userObj.Status = mpb.CommonStatus_INVITED.String()
		userObj.PreSaveInvite(x.CtxClaim, globals.DEFAULT_AT_VALIDITY) // fetch user rom MD
		user, err = x.DMStore.CreateUser(x.Ctx, nil, userObj.StudioRoles, userObj, x.CtxClaim)
		if err != nil {
			x.Logger.Error().Msgf("error while creating user while inviting - %v", err)
			return dmerrors.DMError(utils.ErrUserInviteCreateFailed, err)
		}
		emailer := dmstores.PluginPool.StudioPlugins.Emailer
		body := globals.UserInviteEmailTemplate
		body = strings.ReplaceAll(body, "{Account}", dmstores.AccountPool.Name)
		body = strings.ReplaceAll(body, "{Name}", userObj.FirstName+" "+userObj.LastName)
		body = strings.ReplaceAll(body, "{Link}", pkgs.NetworkConfigManager.ExternalURL+"/login")
		err = emailer.SendRawEmail(x.Ctx, &models.EmailOpts{
			To:               []string{userObj.Email},
			Subject:          "Welcome to VapusData Studio",
			HtmlTemplateBody: body,
		}, x.AgentId)
		if err != nil {
			x.Logger.Error().Msgf("error while sending email on user invite - %v", err)
			return dmerrors.DMError(utils.ErrUserInviteCreateFailed, err)
		}
	}
	x.result.Output.Users = utils.DmUToPb([]*models.Users{user}, x.CtxClaim[encryption.ClaimOrganizationKey])
	return nil
}

func (x *UserManagerAgent) GetUser(userId string) error {
	if userId == "" {
		userId = x.CtxClaim[encryption.ClaimUserIdKey]
	}
	user, err := x.DMStore.GetUser(x.Ctx, userId, x.CtxClaim)
	if err != nil {
		x.Logger.Error().Ctx(x.Ctx).Msgf("error while getting user - %v", err)
		return dmerrors.DMError(utils.ErrUser404, err)
	}
	organizationMap := make(map[string]string)
	organizationIds := ""
	for _, organizationRole := range user.OrganizationRoles {
		organizationIds = fmt.Sprintf("%s'%s',", organizationIds, organizationRole.OrganizationId)
	}
	organizationIds = strings.TrimRight(organizationIds, ",")
	organizations, err := x.DMStore.ListOrganizations(x.Ctx,
		"vapus_id in ("+organizationIds+")", x.CtxClaim)
	if err != nil {
		x.Logger.Error().Ctx(x.Ctx).Msgf("error while getting organization - %v", err)
		return dmerrors.DMError(utils.ErrUser404, err)
	}
	for _, organization := range organizations {
		organizationMap[organization.VapusID] = organization.Name
	}
	x.result.OrganizationMap = organizationMap
	if user.IsValidUserByOrganization(x.CtxClaim[encryption.ClaimOrganizationKey]) {
		x.result.Output.Users = utils.DmUToPb([]*models.Users{user}, x.CtxClaim[encryption.ClaimOrganizationKey])
		return nil
	} else if user.Organization == x.CtxClaim[encryption.ClaimOrganizationKey] {
		dmRole := user.GetOrganizationRole(x.CtxClaim[encryption.ClaimOrganizationKey])
		if len(dmRole) > 0 && slices.Contains(dmRole[0].RoleArns, mpb.StudioUserRoles_DOMAIN_OWNERS.String()) {
			x.result.Output.Users = utils.DmUToPb([]*models.Users{user}, x.CtxClaim[encryption.ClaimOrganizationKey])
			return nil
		} else {
			x.Logger.Error().Msg("error while getting organization role from user object")
			return dmerrors.DMError(utils.ErrUser404, nil)
		}
	} else {
		x.Logger.Error().Ctx(x.Ctx).Msgf("error while getting user - %v", err)
		return dmerrors.DMError(utils.ErrUser404, err)
	}
}

func (x *UserManagerAgent) ListUsers() error {
	users, err := x.DMStore.GetOrganizationUsers(x.Ctx, x.CtxClaim[encryption.ClaimOrganizationKey], x.CtxClaim)
	if err != nil {
		x.Logger.Error().Ctx(x.Ctx).Msgf("error while getting user - %v", err)
		return dmerrors.DMError(utils.ErrUser404, err)
	}
	log.Println("users-----------------", users)
	x.result.Output.Users = utils.DmUToPb(users, x.CtxClaim[encryption.ClaimOrganizationKey])
	return nil
}

func (x *UserManagerAgent) ListStudioUsers() error {
	users, err := x.DMStore.ListUsers(x.Ctx, "", x.CtxClaim)
	if err != nil {
		x.Logger.Error().Ctx(x.Ctx).Msgf("error while getting user - %v", err)
		return dmerrors.DMError(utils.ErrUser404, err)
	}
	log.Println("users-----------------", users)
	x.result.Output.Users = utils.DmUToPb(users, x.CtxClaim[encryption.ClaimOrganizationKey])
	return nil
}

func (x *UserManagerAgent) PatchUser() error {
	userObj := utils.DmUPbToObj(x.managerRequest.GetSpec())
	if userObj == nil || userObj.UserId == "" {
		x.Logger.Error().Msg("error while converting user object from pb")
		return dmerrors.DMError(utils.ErrUser404, nil)
	}
	exUser, err := x.DMStore.GetUser(x.Ctx, userObj.UserId, x.CtxClaim)
	if err != nil {
		x.Logger.Error().Ctx(x.Ctx).Msgf("error while getting user - %v", err)
		return dmerrors.DMError(utils.ErrUser404, err)
	}
	if x.CtxClaim[encryption.ClaimUserIdKey] != exUser.UserId {
		log.Println("x.CtxClaim[encryption.ClaimUserIdKey]-----------------", x.CtxClaim[encryption.ClaimUserIdKey])
		if exUser.GetOrganizationRole(x.CtxClaim[encryption.ClaimOrganizationKey]) == nil {
			x.Logger.Error().Msg("error while getting organization role from user object")
			return dmerrors.DMError(utils.ErrUser404, nil)
		}
		if !strings.Contains(x.CtxClaim[encryption.ClaimOrganizationRolesKey], mpb.StudioUserRoles_DOMAIN_OWNERS.String()) {
			x.Logger.Error().Msg("error while getting organization role from user object")
			return dmerrors.DMError(utils.ErrUser404, nil)
		}
		log.Println("Passed -------------------")
	}
	log.Println("userObj.OrganizationRoles -------------------", userObj.OrganizationRoles)
	updatedOrganizationRoles := []string{}
	for _, organizationRole := range userObj.OrganizationRoles {
		if organizationRole.OrganizationId == x.CtxClaim[encryption.ClaimOrganizationKey] {
			updatedOrganizationRoles = append(updatedOrganizationRoles, organizationRole.RoleArns...)
		}
	}
	log.Println("updatedOrganizationRoles-------------------", updatedOrganizationRoles)
	for _, organizationRole := range exUser.OrganizationRoles {
		if organizationRole.OrganizationId == x.CtxClaim[encryption.ClaimOrganizationKey] {
			organizationRole.RoleArns = updatedOrganizationRoles
		}
	}
	exUser.PreSaveUpdate(x.CtxClaim[encryption.ClaimUserIdKey])
	exUser.FirstName = userObj.FirstName
	exUser.LastName = userObj.LastName
	exUser.DisplayName = userObj.DisplayName
	log.Println("userObj.Profile-------------------", userObj.Profile)
	exUser.Profile = userObj.Profile
	err = x.DMStore.PutUser(x.Ctx, exUser, x.CtxClaim)
	log.Println("user updated----------------->>>>>>>>>", exUser)
	x.result.Output.Users = utils.DmUToPb([]*models.Users{exUser}, x.CtxClaim[encryption.ClaimOrganizationKey])
	return nil
}
