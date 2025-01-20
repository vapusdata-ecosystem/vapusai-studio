package services

import (
	"context"
	"log"
	"slices"
	"strings"

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

type AccountAgent struct {
	request      *pb.AccountManagerRequest
	response     *models.Account
	organization *models.Organization
	*StudioServices
	*models.VapusInterfaceAgentBase
}

func (dms *StudioServices) GetAccount(ctx context.Context) (*models.Account, error) {
	vapusStudioClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		dms.logger.Error().Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}
	_, err := dms.DMStore.GetOrganization(ctx, vapusStudioClaim[encryption.ClaimOrganizationKey], vapusStudioClaim)
	if err != nil {
		return nil, dmerrors.DMError(utils.ErrOrganization404, err)
	}
	resp, err := dms.DMStore.GetAccount(ctx, vapusStudioClaim)
	if err != nil {
		return nil, dmerrors.DMError(utils.ErrGetAccount, err)
	}
	return resp, nil
}

func (dms *StudioServices) NewAccountAgent(ctx context.Context, request *pb.AccountManagerRequest) (*AccountAgent, error) {
	vapusStudioClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		dms.logger.Error().Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}
	organization, err := dms.DMStore.GetOrganization(ctx, vapusStudioClaim[encryption.ClaimOrganizationKey], vapusStudioClaim)
	if err != nil {
		return nil, dmerrors.DMError(utils.ErrOrganization404, err)
	}
	log.Println("NewAccountAgent", request.GetActions().String())
	if organization.OrganizationType != mpb.OrganizationTypes_SERVICE_ORG.String() {
		return nil, dmerrors.DMError(utils.ErrNotServiceOrganization, nil)
	}
	agent := &AccountAgent{
		request:        request,
		organization:   organization,
		StudioServices: dms,
		VapusInterfaceAgentBase: &models.VapusInterfaceAgentBase{
			InitAt:    coreutils.GetEpochTime(),
			Ctx:       ctx,
			CtxClaim:  vapusStudioClaim,
			Action:    request.GetActions().String(),
			AgentType: globals.ACCOUNTAGENT.String(),
		},
	}
	agent.SetAgentId()
	agent.Logger = pkgs.GetSubDMLogger(globals.ACCOUNTAGENT.String(), agent.AgentId)
	return agent, nil
}

func (x *AccountAgent) Act(action string) error {
	if action != "" {
		x.Action = action
	}

	switch x.Action {
	case pb.AccountAgentActions_CONFIGURE_AISTUDIO_MODEL.String():
		response, err := x.configureAIAttributes()
		if err != nil {
			return err
		}
		x.response = response
		return nil
	case pb.AccountAgentActions_UPDATE_PROFILE.String():
		response, err := x.updateAccount()
		if err != nil {
			return err
		}
		x.response = response
		return nil
	default:
		return dmerrors.DMError(utils.ErrInvalidAccountAction, nil)
	}
}

func (x *AccountAgent) GetResponse() *models.Account {
	return x.response
}

func (x *AccountAgent) configureAIAttributes() (*models.Account, error) {
	userRoles := strings.Split(x.CtxClaim[encryption.ClaimRoleKey], ",")
	if !slices.Contains(userRoles, mpb.StudioUserRoles_PLATFORM_OWNERS.String()) || !slices.Contains(userRoles, mpb.StudioUserRoles_PLATFORM_OPERATORS.String()) {
		return nil, dmerrors.DMError(utils.ErrAccountOps403, nil)

	}
	var err error
	reqObj := (&models.Account{}).ConvertFromPb(x.request.GetSpec())
	account, err := x.DMStore.GetAccount(x.Ctx, x.CtxClaim)
	if err != nil {
		return nil, dmerrors.DMError(utils.ErrGetAccount, err)
	}
	account.AIAttributes = reqObj.GetAiAttributes()
	account.PreSaveUpdate(x.CtxClaim[encryption.ClaimUserIdKey])
	err = x.DMStore.PutAccount(x.Ctx, account, x.CtxClaim)
	if err != nil {
		return nil, dmerrors.DMError(utils.ErrConfigureAIStudioModel, err)
	}
	dmstores.InitAccountPool(x.Ctx, x.DMStore)
	return account, nil
}

func (x *AccountAgent) updateAccount() (*models.Account, error) {
	userRoles := strings.Split(x.CtxClaim[encryption.ClaimRoleKey], ",")
	if !slices.Contains(userRoles, mpb.StudioUserRoles_PLATFORM_OWNERS.String()) {
		return nil, dmerrors.DMError(utils.ErrAccountOps403, nil)
	}
	var err error
	reqObj := (&models.Account{}).ConvertFromPb(x.request.GetSpec())
	account, err := x.DMStore.GetAccount(x.Ctx, x.CtxClaim)
	if err != nil {
		return nil, dmerrors.DMError(utils.ErrGetAccount, err)
	}
	account.Profile = reqObj.Profile
	account.AIAttributes = reqObj.GetAiAttributes()
	account.PreSaveUpdate(x.CtxClaim[encryption.ClaimUserIdKey])
	err = x.DMStore.PutAccount(x.Ctx, account, x.CtxClaim)
	if err != nil {
		return nil, dmerrors.DMError(utils.ErrConfigureAIStudioModel, err)
	}
	dmstores.InitAccountPool(x.Ctx, x.DMStore)
	return account, nil
}

func (x *AccountAgent) LogAgent() {
	x.Logger.Info().Msgf("AccountAgent - %v action started at %v and finished at %v with status %v", x.AgentId, x.InitAt, x.FinishAt, x.Status)
}
