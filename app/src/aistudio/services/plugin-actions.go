package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	dmstores "github.com/vapusdata-oss/aistudio/aistudio/datastoreops"
	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	"github.com/vapusdata-oss/aistudio/aistudio/utils"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	filemanager "github.com/vapusdata-oss/aistudio/core/filemanager"
	models "github.com/vapusdata-oss/aistudio/core/models"
	coreutils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

type PluginActionsAgent struct {
	*models.VapusInterfaceAgentBase
	request *pb.PluginActionRequest
	result  []*mpb.Plugin
	*StudioServices
}

func (s *StudioServices) NewPluginActionsAgent(ctx context.Context, request *pb.PluginActionRequest) (*PluginActionsAgent, error) {
	vapusStudioClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		s.logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, encryption.ErrInvalidJWTClaims
	}
	agent := &PluginActionsAgent{
		request:        request,
		result:         make([]*mpb.Plugin, 0),
		StudioServices: s,
		VapusInterfaceAgentBase: &models.VapusInterfaceAgentBase{
			CtxClaim: vapusStudioClaim,
			Ctx:      ctx,
			InitAt:   coreutils.GetEpochTime(),
		},
	}
	agent.SetAgentId()

	agent.Logger = pkgs.GetSubDMLogger(request.PluginType, agent.AgentId)
	return agent, nil
}

func (v *PluginActionsAgent) GetAgentId() string {
	return v.AgentId
}

func (v *PluginActionsAgent) GetResult() []*mpb.Plugin {
	v.FinishAt = coreutils.GetEpochTime()
	v.FinalLog()
	return v.result
}

func (v *PluginActionsAgent) Act() error {
	switch v.request.PluginType {
	case mpb.IntegrationPluginTypes_EMAIL.String():
		return v.sendEmail()
	case mpb.IntegrationPluginTypes_FILE_STORE.String():
		return v.fileStoreAction()
	default:
		v.logger.Error().Ctx(v.Ctx).Msg("invalid plugin type for action")
		return utils.ErrInvalidPluginTypeForAction
	}
}

func (v *PluginActionsAgent) sendEmail() error {
	if dmstores.PluginPool == nil {
		dmstores.NewPluginPool(v.Ctx, v.DMStore)
	}
	emailer := dmstores.PluginPool.StudioPlugins.Emailer
	reqObj := &models.EmailOpts{}
	if err := json.Unmarshal(v.request.GetSpec(), reqObj); err != nil {
		v.Logger.Error().Err(err).Msg("error while unmarshalling request")
		return err
	}
	return emailer.SendRawEmail(v.Ctx, reqObj, v.AgentId)
}

func (v *PluginActionsAgent) fileStoreAction() error {
	plQ := fmt.Sprintf("deleted_at is null and status = 'ACTIVE' AND plugin_type='%s' AND scope='PLATFORM_SCOPE'",
		mpb.IntegrationPluginTypes_FILE_STORE.String())
	q := fmt.Sprintf("organization='%s' AND created_by='%s' AND plugin_type='%s'",
		v.CtxClaim[encryption.ClaimOrganizationKey],
		v.CtxClaim[encryption.ClaimUserIdKey],
		mpb.IntegrationPluginTypes_FILE_STORE.String())
	var uPl *models.Plugin
	var plPL *models.Plugin
	var sourceCreds *models.GenericCredentialModel
	var err error

	var wg sync.WaitGroup
	wg.Add(2)
	var errChan = make(chan error, 2)
	go func() {
		defer wg.Done()
		res, err := v.DMStore.ListPlugins(v.Ctx, q, v.CtxClaim)
		if err != nil || len(res) == 0 {
			v.Logger.Error().Err(err).Msg("error while listing file store plguin for user")
			errChan <- dmerrors.DMError(utils.ErrFileStorePlugin404, err)
			return
		}
		uPl = res[0]
	}()
	go func() {
		defer wg.Done()
		res, err := v.DMStore.ListPlugins(v.Ctx, plQ, v.CtxClaim)
		if err != nil || len(res) == 0 {
			v.Logger.Error().Err(err).Msg("error while listing file store plguin for platform")
			errChan <- dmerrors.DMError(utils.ErrFileStorePlugin404, err)
			return
		}
		plPL = res[0]
		sourceCreds, err = v.DMStore.ReadCredentialFromStore(v.Ctx, plPL.NetworkParams.SecretName)
		if err != nil || sourceCreds == nil {
			v.Logger.Error().Err(err).Msg("error while getting filestore creds from secret store")
			errChan <- dmerrors.DMError(utils.ErrPlugin404, err)
		}
	}()
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			v.logger.Error().Err(err).Msg("error while getting file store plugin")
			return dmerrors.DMError(utils.ErrFileStorePlugin404, err)
		}
	}
	if uPl == nil {
		return dmerrors.DMError(utils.ErrFileStorePlugin404, err)
	}
	sourceCreds.Username = v.CtxClaim[encryption.ClaimUserIdKey]
	fileStoreClient, err := filemanager.NewFileManagerClient(
		v.Ctx, v.logger,
		filemanager.WithService(plPL.PluginService),
		filemanager.WithCredentials(sourceCreds),
	)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while creating file store client")
		return dmerrors.DMError(utils.ErrFileStorePlugin400, err)
	}
	reqObj := &models.FileManageOpts{}
	if err := json.Unmarshal(v.request.GetSpec(), reqObj); err != nil {
		v.Logger.Error().Err(err).Msg("error while unmarshalling request")
		return err
	}
	return fileStoreClient.Upload(v.Ctx,
		reqObj)
}
