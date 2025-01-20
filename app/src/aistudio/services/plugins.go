package services

import (
	"context"
	"fmt"
	"slices"
	"strings"

	guuid "github.com/google/uuid"
	dmstores "github.com/vapusdata-oss/aistudio/aistudio/datastoreops"
	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	"github.com/vapusdata-oss/aistudio/aistudio/utils"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	"github.com/vapusdata-oss/aistudio/core/globals"
	models "github.com/vapusdata-oss/aistudio/core/models"
	"github.com/vapusdata-oss/aistudio/core/plugins"
	coreutils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

type PluginManagerAgent struct {
	*models.VapusInterfaceAgentBase
	managerRequest *pb.PluginManagerRequest
	getterRequest  *pb.PluginGetterRequest
	result         []*mpb.Plugin
	organization   *models.Organization
	*dmstores.DMStore
}

func (s *StudioServices) NewPluginManagerAgent(ctx context.Context, managerRequest *pb.PluginManagerRequest, getterRequest *pb.PluginGetterRequest) (*PluginManagerAgent, error) {
	vapusStudioClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		s.logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, encryption.ErrInvalidJWTClaims
	}
	organization, err := s.DMStore.GetOrganization(ctx, vapusStudioClaim[encryption.ClaimOrganizationKey], vapusStudioClaim)
	if err != nil {
		s.logger.Error().Err(err).Ctx(ctx).Msg("error while getting organization from datastore")
		return nil, dmerrors.DMError(utils.ErrOrganization404, err)
	}
	agent := &PluginManagerAgent{
		managerRequest: managerRequest,
		getterRequest:  getterRequest,
		result:         make([]*mpb.Plugin, 0),
		DMStore:        s.DMStore,
		organization:   organization,
		VapusInterfaceAgentBase: &models.VapusInterfaceAgentBase{
			CtxClaim: vapusStudioClaim,
			Ctx:      ctx,
			Action:   managerRequest.GetAction().String(),
			InitAt:   coreutils.GetEpochTime(),
		},
	}
	agent.SetAgentId()
	if managerRequest != nil {
		agent.Action = managerRequest.GetAction().String()
	} else {
		agent.Action = ""
	}
	agent.Logger = pkgs.GetSubDMLogger(globals.AIPROMPTAGENT.String(), agent.AgentId)
	return agent, nil
}

func (v *PluginManagerAgent) GetAgentId() string {
	return v.AgentId
}

func (v *PluginManagerAgent) GetResult() []*mpb.Plugin {
	v.FinishAt = coreutils.GetEpochTime()
	v.FinalLog()
	return v.result
}

func (v *PluginManagerAgent) Act() error {
	switch v.GetAction() {
	case pb.PluginAgentAction_CONFIGURE_PLUGIN.String():
		return v.configurePlugin()
	case pb.PluginAgentAction_PATCH_PLUGIN.String():
		return v.patchPlugin()
	default:
		if v.getterRequest != nil {
			if v.getterRequest.GetPluginId() != "" {
				return v.describePlugin()
			} else {
				return v.listPlugins()
			}
		}
		v.Logger.Error().Msg("invalid action")
		return utils.ErrInvalidAction
	}
}

func (v *PluginManagerAgent) configurePlugin() error {
	cQ := fmt.Sprintf("status = 'ACTIVE' AND deleted_at IS NULL AND plugin_type = '%s'", v.managerRequest.GetSpec().GetPluginType())
	switch v.managerRequest.GetSpec().GetScope() {
	case mpb.ResourceScope_ORG_SCOPE.String():
		cQ += fmt.Sprintf(" AND scope = '%s' AND organization = '%s'", mpb.ResourceScope_USER_SCOPE.String(), v.CtxClaim[encryption.ClaimOrganizationKey])
	case mpb.ResourceScope_USER_SCOPE.String():
		cQ += fmt.Sprintf(" AND scope = '%s' AND created_by = '%s' AND organization = '%s'", mpb.ResourceScope_USER_SCOPE.String(), v.CtxClaim[encryption.ClaimUserIdKey], v.CtxClaim[encryption.ClaimOrganizationKey])
	case mpb.ResourceScope_ACCOUNT_SCOPE.String():
		cQ += fmt.Sprintf(" AND scope = '%s'", mpb.ResourceScope_ACCOUNT_SCOPE.String())
	default:
		return dmerrors.DMError(utils.ErrInvalidResourceScope, nil)
	}
	count, err := v.CountPlugins(v.Ctx,
		cQ,
		v.CtxClaim)
	if count > 0 {
		v.Logger.Error().Msg("plugin already exists")
		return dmerrors.DMError(utils.ErrPluginServiceScopeExists, nil)
	}
	validScope, ok := plugins.PluginTypeScopeMap[v.managerRequest.GetSpec().GetPluginType().String()]
	if !ok {
		v.Logger.Error().Msg("invalid plugin type")
		return dmerrors.DMError(utils.ErrUnSupportedPluginType, nil)
	}
	if !slices.Contains(validScope, v.managerRequest.GetSpec().GetScope()) {
		v.Logger.Error().Msg("invalid plugin scope")
		return dmerrors.DMError(utils.ErrPluginScope403, nil)
	}
	plugin := (&models.Plugin{}).ConvertFromPb(v.managerRequest.GetSpec())
	plugin.PreSaveCreate(v.CtxClaim)
	if plugin.NetworkParams.SecretName == coreutils.EMPTYSTR {
		plugin.NetworkParams.SecretName = coreutils.GetSecretName("plugin", "", "creds-"+guuid.New().String())
	}
	err = v.StoreCredsInStore(v.Ctx, plugin.NetworkParams.SecretName, plugin.NetworkParams.Credentials, v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while saving plugin creds")
		return err
	}

	if plugin.Scope == mpb.ResourceScope_ORG_SCOPE.String() {
		if !strings.Contains(v.CtxClaim[encryption.ClaimOrganizationRolesKey], mpb.StudioUserRoles_DOMAIN_OWNERS.String()) {
			return dmerrors.DMError(utils.ErrPluginOrganizationScope403, nil)
		}
	}
	plugin.Organization = v.CtxClaim[encryption.ClaimOrganizationKey]
	plugin.NetworkParams.Credentials = nil
	plugin.NetworkParams.IsAlreadyInSecretBS = true
	err = v.ConfigurePlugin(v.Ctx, plugin, v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while configuring plugin")
		return dmerrors.DMError(utils.ErrPluginCreate400, err)
	}
	v.result = []*mpb.Plugin{plugin.ConvertToPb()}
	return nil
}

func (v *PluginManagerAgent) patchPlugin() error {
	existingObj, err := v.GetPlugin(v.Ctx, v.managerRequest.GetSpec().GetPluginId(), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting plugin from datastore")
		return dmerrors.DMError(utils.ErrPlugin404, err)
	}
	if !existingObj.Editable {
		v.Logger.Error().Msg("plugin not editable")
		return dmerrors.DMError(utils.ErrPluginNotEditable, nil)
	} else if existingObj.Organization != v.CtxClaim[encryption.ClaimOrganizationKey] && existingObj.CreatedBy != v.CtxClaim[encryption.ClaimUserIdKey] {
		v.Logger.Error().Msg("plugin not in organization")
		return dmerrors.DMError(utils.ErrPlugin403, nil)
	}
	plugin := (&models.Plugin{}).ConvertFromPb(v.managerRequest.GetSpec())
	existingObj.Name = plugin.Name
	existingObj.PreSaveUpdate(v.CtxClaim[encryption.ClaimUserIdKey])
	if plugin.NetworkParams.SecretName == coreutils.EMPTYSTR {
		existingObj.NetworkParams = plugin.NetworkParams
		existingObj.NetworkParams.SecretName = coreutils.GetSecretName("plugin", "", "creds-"+guuid.New().String())
		err = v.StoreCredsInStore(v.Ctx, existingObj.NetworkParams.SecretName, existingObj.NetworkParams.Credentials, v.CtxClaim)
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while saving plugin creds")
			return err
		}
		existingObj.NetworkParams.Credentials = nil
		existingObj.NetworkParams.IsAlreadyInSecretBS = true
	}
	existingObj.Status = mpb.CommonStatus_ACTIVE.String()
	err = v.PutPlugin(v.Ctx, existingObj, v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while configuring plugin")
		return dmerrors.DMError(utils.ErrPluginPatch400, err)
	}
	v.result = []*mpb.Plugin{existingObj.ConvertToPb()}
	return nil
}

func (v *PluginManagerAgent) describePlugin() error {
	result, err := v.GetPlugin(v.Ctx, v.getterRequest.GetPluginId(), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting plugin from datastore")
		return dmerrors.DMError(utils.ErrPlugin404, err)
	}
	if !result.Editable {
		v.Logger.Error().Msg("plugin not editable")
		return dmerrors.DMError(utils.ErrPluginNotEditable, nil)
	} else if result.Organization != v.CtxClaim[encryption.ClaimOrganizationKey] && result.CreatedBy != v.CtxClaim[encryption.ClaimUserIdKey] {
		v.Logger.Error().Msg("plugin not in organization")
		return dmerrors.DMError(utils.ErrPlugin403, nil)
	}
	v.result = []*mpb.Plugin{result.ConvertToPb()}
	return nil
}

func (v *PluginManagerAgent) listPlugins() error {
	v.Logger.Info().Msg("getting ai agent list from datastore")
	result, err := v.ListPlugins(v.Ctx,
		fmt.Sprintf(`status = 'ACTIVE' AND deleted_at IS NULL AND 
		(
	(scope='PLATFORM_SCOPE') OR (scope='DOMAIN_SCOPE' AND organization='%s') OR (scope='USER_SCOPE' AND organization='%s' AND created_by='%s')
		)
		 ORDER BY created_at DESC`, v.CtxClaim[encryption.ClaimOrganizationKey], v.CtxClaim[encryption.ClaimOrganizationKey], v.CtxClaim[encryption.ClaimUserIdKey]),
		v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting plugin from datastore")
		return dmerrors.DMError(utils.ErrPlugin404, err)
	}
	v.result = utils.DPPLObj2Pb(result)
	return nil
}

func (v *PluginManagerAgent) archivePlugins() error {
	v.Logger.Info().Msg("getting ai agent list from datastore")

	result, err := v.GetPlugin(v.Ctx, v.getterRequest.GetPluginId(), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting plugin from datastore")
		return dmerrors.DMError(utils.ErrPlugin404, err)
	}
	if result.Organization != v.CtxClaim[encryption.ClaimOrganizationKey] && result.CreatedBy != v.CtxClaim[encryption.ClaimUserIdKey] {
		v.Logger.Error().Msg("plugin not onwed by current user")
		return dmerrors.DMError(utils.ErrPlugin403, nil)
	}
	result.Status = mpb.CommonStatus_DELETED.String()
	result.DeletedAt = coreutils.GetEpochTime()
	result.DeletedBy = v.CtxClaim[encryption.ClaimUserIdKey]
	err = v.PutPlugin(v.Ctx, result, v.CtxClaim)
	v.result = nil
	return nil
}
