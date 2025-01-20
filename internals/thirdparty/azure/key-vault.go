package azure

import (
	"context"
	"encoding/json"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	azsecrets "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	utils "github.com/vapusdata-oss/aistudio/core/utils"
)

type AzureKVManager interface {
	WriteSecret(ctx context.Context, data any, name string) error
	ReadSecret(ctx context.Context, secretId string) (any, error)
	DeleteSecret(ctx context.Context, secretId string) error
}

type AzureKeyVault struct {
	cl                               *azsecrets.Client
	secretPrefix, secretNameTemplate string
}

func NewAzureKeyVault(ctx context.Context, opts *AzureConfig, valultURI string) (*AzureKeyVault, error) {
	credential, err := azidentity.NewClientSecretCredential(opts.TenantID, opts.ClientID, opts.ClientSecret, nil)
	if err != nil {
		return nil, dmerrors.DMError(ErrCreatingAzureCredential, err)
	}

	// Establish a connection to the Key Vault client
	client, err := azsecrets.NewClient(valultURI, credential, nil)
	if err != nil {
		return nil, dmerrors.DMError(ErrCreatingAzureKeyVaultClient, err)
	}
	return &AzureKeyVault{
		cl: client,
	}, nil
}

func (akv *AzureKeyVault) WriteSecret(ctx context.Context, data any, secretName string) error {
	// Convert the secret value to a byte array
	secretValue, err := json.Marshal(data)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrJsonMarshel, err)
	}
	_, err = akv.cl.SetSecret(ctx, secretName, azsecrets.SetSecretParameters{
		Value: utils.Str2Ptr(string(secretValue)),
	}, nil)
	if err != nil {
		return dmerrors.DMError(ErrCreatingAzureSecret, err)
	}

	return nil
}

func (akv *AzureKeyVault) ReadSecret(ctx context.Context, secretName string) (any, error) {
	resp, err := akv.cl.GetSecret(ctx, secretName, "", nil)
	// TO:DO check error for 404 or other using error.As
	if err != nil {
		return nil, dmerrors.DMError(ErrReadingAzureSecret, err)
	}

	return json.Marshal([]byte(*resp.Value))
}

func (akv *AzureKeyVault) DeleteSecret(ctx context.Context, secretName string) error {
	_, err := akv.cl.DeleteSecret(ctx, secretName, nil)
	if err != nil {
		return dmerrors.DMError(ErrDeletingAzureSecret, err)
	}
	return nil
}
