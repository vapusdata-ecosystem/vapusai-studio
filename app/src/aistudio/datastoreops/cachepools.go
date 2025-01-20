package dmstores

import (
	"context"

	"github.com/vapusdata-oss/aistudio/core/models"
	pluginsstore "github.com/vapusdata-oss/aistudio/core/plugins"
	serviceops "github.com/vapusdata-oss/aistudio/core/serviceops"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

var AccountPool *models.Account

var PluginPool *pluginsstore.VapusPlugins

func InitPluginPool(ctx context.Context, ds *DMStore) {
	plugins, err := ds.ListPlugins(ctx, "deleted_at IS NULL AND scope = 'PLATFORM_SCOPE'", map[string]string{})
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while getting the list of plugins")
		return
	}
	for _, plugin := range plugins {
		if plugin.Status == mpb.CommonStatus_ACTIVE.String() {
			creds, err := ds.ReadCredentialFromStore(ctx, plugin.NetworkParams.SecretName)
			if err != nil {
				logger.Err(err).Ctx(ctx).Msg("error while reading the secret from the vault")
				continue
			}
			plugin.NetworkParams.Credentials = creds
		}
	}
	PluginPool = pluginsstore.NewVapusPlugins(ctx, plugins, []string{}, logger)
	logger.Info().Ctx(ctx).Msg("PluginPool initialized")
	return
}

func NewPluginPool(ctx context.Context, ds *DMStore) {
	if PluginPool == nil {
		InitPluginPool(ctx, ds)
	}
}

func InitAccountPool(ctx context.Context, ds *DMStore) {
	result := make([]*models.Account, 0)
	query := "select * from " + serviceops.AccountsTable
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		logger.Err(err).Ctx(ctx).Msg("error while getting the list of accounts")
	}
	AccountPool = result[0]
	needUpdate := false
	logger.Info().Ctx(ctx).Msg("AccountPool initializeing...checking for base os artifacts")

	if needUpdate {
		err = ds.PutAccount(ctx, AccountPool, map[string]string{})
		if err != nil {
			logger.Fatal().Err(err).Ctx(ctx).Msg("error while updating the account with new artifacts for base os")
		} else {
			logger.Info().Ctx(ctx).Msg("AccountPool updated with new artifacts for base os")
		}
	}
	logger.Info().Ctx(ctx).Msg("AccountPool initialized")
	return
}

func GetAccountFromPool(ctx context.Context, ds *DMStore, ctxClaim map[string]string) *models.Account {

	if AccountPool != nil {
		return AccountPool
	} else {
		InitAccountPool(ctx, ds)
		return AccountPool
	}
}
