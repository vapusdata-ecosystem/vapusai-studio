package datasvc

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/rs/zerolog"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	"github.com/vapusdata-oss/aistudio/core/globals"
	models "github.com/vapusdata-oss/aistudio/core/models"
	tpaws "github.com/vapusdata-oss/aistudio/core/thirdparty/aws"
	tpazure "github.com/vapusdata-oss/aistudio/core/thirdparty/azure"
	tpgcp "github.com/vapusdata-oss/aistudio/core/thirdparty/gcp"
	tphcvault "github.com/vapusdata-oss/aistudio/core/thirdparty/hcvault"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

type SecretStoreOps interface {
	WriteSecret(ctx context.Context, data any, name string) error
	ReadSecret(ctx context.Context, secretId string) (any, error)
	DeleteSecret(ctx context.Context, secretId string) error
}

type SecretStoreClient struct {
	SecretStoreOps
	logger zerolog.Logger
	creds  *models.StoreParams
}

func (d *SecretStoreClient) Close() {
}

var SecretParentKey = "secret_value"

func NewSecretStoreClient(ctx context.Context, opts *models.StoreParams, logger zerolog.Logger) (*SecretStoreClient, error) {
	log.Debug().Msg("Creating new secret store client")
	resultCl := &SecretStoreClient{
		logger: logger,
	}
	if opts == nil || opts.DataSourceEngine == "" || opts.Creds.DsCreds == nil {
		return nil, dmerrors.DMError(ErrInvalidDataStoreEngine, ErrDataStoreConn)
	}
	switch opts.DataSourceEngine {
	case mpb.StorageEngine_HASHICORPVAULT.String():
		var se string
		if val, ok := opts.Params[globals.SECRETENGINE]; ok {
			se = val.(string)
		}
		client, err := tphcvault.NewHcVaultManager(ctx, &tphcvault.Vault{
			URL: opts.Creds.Address,
			// AuthAppRole:     conf.HashicorpVault.AppRoleAuthnEnabled,
			// ApproleRoleID:   conf.HashicorpVault.AppRoleID,
			// ApproleSecretID: conf.HashicorpVault.AppRoleSecret,
			Token:        opts.Creds.DsCreds.ApiToken,
			SecretEngine: se,
		})
		if err != nil {
			log.Err(err).Msg("Error creating vault client")
			return nil, err
		}
		resultCl.SecretStoreOps = client
		return resultCl, nil
	case mpb.StorageEngine_AWS_VAULT.String():
		client, err := tpaws.NewAwsSmClient(ctx, &tpaws.AWSConfig{
			Region:          opts.Creds.DsCreds.AwsCreds.Region,
			AccessKeyId:     opts.Creds.DsCreds.AwsCreds.AccessKeyId,
			SecretAccessKey: opts.Creds.DsCreds.AwsCreds.SecretAccessKey,
		})
		if err != nil {
			log.Err(err).Msg("Error creating aws secret manager client")
			return nil, err
		}
		resultCl.SecretStoreOps = client
		return resultCl, nil
	case mpb.StorageEngine_GCP_VAULT.String():
		decodeData, err := base64.StdEncoding.DecodeString(opts.Creds.DsCreds.GcpCreds.ServiceAccountKey)
		if err != nil {
			log.Err(err).Msg("Error decoding gcp service account key")
			return nil, err
		}
		client, err := tpgcp.NewGcpSMStore(ctx, &tpgcp.GcpConfig{
			ServiceAccountKey: []byte(decodeData),
			ProjectID:         opts.Creds.DsCreds.GcpCreds.ProjectId,
			Region:            opts.Creds.DsCreds.GcpCreds.Region,
		})
		if err != nil {
			log.Err(err).Msg("Error creating gcp secret manager client")
			return nil, err
		}
		resultCl.SecretStoreOps = client
		return resultCl, nil
	case mpb.StorageEngine_AZURE_VAULT.String():
		client, err := tpazure.NewAzureKeyVault(ctx, &tpazure.AzureConfig{
			TenantID:     opts.Creds.DsCreds.AzureCreds.TenantId,
			ClientID:     opts.Creds.DsCreds.AzureCreds.ClientId,
			ClientSecret: opts.Creds.DsCreds.AzureCreds.ClientSecret,
		}, opts.Creds.Address)
		if err != nil {
			log.Err(err).Msg("Error creating azure key vault client")
			return nil, err
		}
		resultCl.SecretStoreOps = client
		return resultCl, nil
	default:
		return nil, dmerrors.DMError(ErrInvalidDataStoreEngine, ErrDataStoreConn)
	}
}

// WriteSecret writes the secret to the secret store
func (be *SecretStoreClient) WriteSecret(ctx context.Context, secrdData any, name string) error {
	bytes, err := json.Marshal(secrdData)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrJsonMarshel, err)
	}
	return be.SecretStoreOps.WriteSecret(ctx, map[string]interface{}{SecretParentKey: string(bytes)}, name)
}

// ReadSecret reads the secret from the secret store
func (be *SecretStoreClient) ReadSecret(ctx context.Context, secretId string) (string, error) {
	origVal, err := be.SecretStoreOps.ReadSecret(ctx, secretId)
	if err != nil {
		return "", err
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(origVal.([]byte), &result)
	if err != nil {
		be.logger.Err(err).Msg("Error unmarshalling secret data")
		return "", dmerrors.DMError(dmerrors.ErrJsonUnMarshel, err)
	}
	val, ok := result[SecretParentKey]
	if !ok {
		be.logger.Err(err).Msg("Error reading secret data, secret key not found")
		return string(origVal.([]byte)), nil
	}
	_, ok = val.(string)
	if !ok {
		be.logger.Err(err).Msg("Error reading secret data, secret value not found")
		return "", dmerrors.DMError(ErrInvalidSecretData, nil)
	}
	return val.(string), nil
}

// DeleteSecret deletes the secret from the secret store
func (be *SecretStoreClient) DeleteSecret(ctx context.Context, secretId string) error {
	return be.SecretStoreOps.DeleteSecret(ctx, secretId)
}
