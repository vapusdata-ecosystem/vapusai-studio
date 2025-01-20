package filemanager

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/vapusdata-oss/aistudio/core/models"
	gcp "github.com/vapusdata-oss/aistudio/core/thirdparty/gcp"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

type FileManager interface {
	Upload(ctx context.Context, opts *models.FileManageOpts) error
	ListFiles(ctx context.Context, opts *models.FileManageOpts) error
	DeleteFiles(ctx context.Context, opts *models.FileManageOpts) error
}

type FileManagerClient struct {
	service     string
	credentials *models.GenericCredentialModel
}

type FileManageFunc func(*FileManagerClient)

func WithService(service string) FileManageFunc {
	return func(f *FileManagerClient) {
		f.service = service
	}
}

func WithCredentials(credentials *models.GenericCredentialModel) FileManageFunc {
	return func(f *FileManagerClient) {
		f.credentials = credentials
	}
}

func NewFileManagerClient(ctx context.Context, logger zerolog.Logger, opts ...FileManageFunc) (FileManager, error) {
	fileManagerClient := &FileManagerClient{}
	for _, opt := range opts {
		opt(fileManagerClient)
	}
	switch fileManagerClient.service {
	case mpb.IntegrationPlugins_GOOGLE_DRIVE.String():
		return gcp.NewDriveAgent(ctx, &gcp.GcpConfig{
			ServiceAccountKey: []byte(fileManagerClient.credentials.GetGcpCreds().ServiceAccountKey),
			ProjectID:         fileManagerClient.credentials.GetGcpCreds().ProjectId,
			Region:            fileManagerClient.credentials.GetGcpCreds().Region,
			Zone:              fileManagerClient.credentials.GetGcpCreds().Zone,
		}, fileManagerClient.credentials.Username, logger)
	default:
		return nil, ErrInvalidFIleStoreService
	}
}
