package dmstores

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	utils "github.com/vapusdata-oss/aistudio/aistudio/utils"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	"github.com/vapusdata-oss/aistudio/core/globals"
	"github.com/vapusdata-oss/aistudio/core/models"
	coreutils "github.com/vapusdata-oss/aistudio/core/utils"
)

func (ds *DMStore) CacheFilter(ctx context.Context, action, key string, value ...string) (interface{}, error) {
	switch action {
	case globals.LIST:
		return ds.Cacher.RedisClient.Client.Do(ctx, "CF.EXISTS", key, value[0]).Result()
	case globals.ADD:
		return ds.Cacher.RedisClient.Client.Do(ctx, "CF.ADD", key, value[0]).Result()
	case globals.EXISTS:
		return ds.Cacher.RedisClient.Client.Do(ctx, "CF.EXISTS", key, value[0]).Result()
	case globals.COUNT:
		return ds.Cacher.RedisClient.Client.Do(ctx, "CF.CARD", key).Result()
	case globals.MADD:
		return ds.Cacher.RedisClient.Client.Do(ctx, "CF.MADD", key, value).Result()
	case globals.DEL:
		return ds.Cacher.RedisClient.Client.Do(ctx, "CF.DEL", key, value[0]).Result()
	default:
		return ds.Cacher.RedisClient.Client.Do(ctx, "CF.EXISTS", key, value[0]).Result()
	}
}

func (n *DMStore) GetAIModelNodeNetworkParams(ctx context.Context, aiModelNode *models.AIModelNode) (*models.AIModelNodeNetworkParams, error) {
	secrets := &models.GenericCredentialModel{}

	secretStr, err := n.SecretStore.ReadSecret(ctx, aiModelNode.NetworkParams.SecretName)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("Error while reading the secret from the vault")
		return nil, dmerrors.DMError(utils.ErrDataSourceCredsSecretGet, err)
	}
	err = json.Unmarshal([]byte(secretStr), secrets)
	log.Println("secrets -------->", string(secretStr))
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while unmarshelling creds from secret store.")
		return nil, dmerrors.DMError(utils.ErrDataSourceCredsSecretGet, err)
	}

	return &models.AIModelNodeNetworkParams{
		Url:                 aiModelNode.NetworkParams.GetUrl(),
		Credentials:         secrets,
		ApiVersion:          aiModelNode.NetworkParams.GetApiVersion(),
		LocalPath:           aiModelNode.NetworkParams.GetLocalPath(),
		SecretName:          aiModelNode.NetworkParams.SecretName,
		IsAlreadyInSecretBs: aiModelNode.NetworkParams.IsAlreadyInSecretBs,
	}, nil
}

func (n *DMStore) ReadCredentialFromStore(ctx context.Context, secretName string) (*models.GenericCredentialModel, error) {
	secrets := &models.GenericCredentialModel{}
	secretStr, err := n.SecretStore.ReadSecret(ctx, secretName)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("Error while reading the secret from the vault")
		return nil, dmerrors.DMError(utils.ErrDataSourceCredsSecretGet, err)
	}
	err = json.Unmarshal([]byte(secretStr), secrets)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while unmarshelling creds from secret store.")
		return nil, dmerrors.DMError(utils.ErrDataSourceCredsSecretGet, err)
	}
	return secrets, nil
}

func (ds *DMStore) StoreCredsInStore(ctx context.Context, secretName string, creds *models.GenericCredentialModel, ctxClaim map[string]string) error {
	result, err := coreutils.StructToMap(creds)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msgf("error while converting struct to map")
		return err
	}

	err = ds.SecretStore.WriteSecret(ctx, result, secretName)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msgf("error while writing secret %v", secretName)
		return err
	}
	return nil
}

func GetAccountFilter(ctxClaim map[string]string, condition string) string {
	if ctxClaim != nil && ctxClaim[encryption.ClaimAccountKey] != "" {
		if condition == "" {
			return " owner_account = '" + ctxClaim[encryption.ClaimAccountKey] + "' "
		} else {
			return " owner_account = '" + ctxClaim[encryption.ClaimAccountKey] + "'" + " AND " + condition
		}
	}
	return condition
}

func GetByIdFilter(fieldId, val string, ctxClaim map[string]string) string {
	if fieldId == "" && val == "" {
		return GetAccountFilter(ctxClaim, "")
	}
	if fieldId == "" {
		fieldId = "vapus_id"
	}
	if ctxClaim != nil && ctxClaim[encryption.ClaimAccountKey] != "" {
		return fieldId + " = '" + val + "'" + " AND " + GetAccountFilter(ctxClaim, "")
	}
	return fieldId + " = '" + val + "'"
}

func VapusIdFilter() string {
	return "vapus_id = ?"
}

func GetOrderByQuery(field, order string) string {
	return fmt.Sprintf(" ORDER BY %s %s", field, order)
}
