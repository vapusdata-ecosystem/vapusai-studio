package datasvc

import (
	"context"
	"encoding/base64"

	"github.com/rs/zerolog"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	models "github.com/vapusdata-oss/aistudio/core/models"
	tpgcp "github.com/vapusdata-oss/aistudio/core/thirdparty/gcp"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

type BlobStore interface {
	Close()
	CreateBucket(ctx context.Context, bucketName string) error
	DeleteBucket(ctx context.Context, bucketName string) error
	ListBuckets(ctx context.Context) ([]string, error)
	GetBucket(ctx context.Context, bucketName string) (string, error)
	ListObjects(ctx context.Context, bucketName string) ([]string, error)
	UploadObject(ctx context.Context, bucketName, objectName string, data []byte) error
	DeleteObject(ctx context.Context, bucketName, objectName string) error
	DownloadObject(ctx context.Context, bucketName, objectName string) ([]byte, error)
}

type BlobStoreClient struct {
	BlobStore
	logger zerolog.Logger
	creds  *models.StoreParams
}

func (d *BlobStoreClient) Close() {
	if d.BlobStore != nil {
		d.Close()
	}
}

func NewBlobStoreClient(ctx context.Context, opts *models.StoreParams, logger zerolog.Logger) (BlobStore, error) {
	log.Debug().Msg("Creating new BlobStoreClient client")
	resultCl := &BlobStoreClient{
		logger: logger,
	}
	if opts == nil || opts.DataSourceEngine == "" || opts.Creds == nil {
		return nil, dmerrors.DMError(ErrInvalidDataStoreEngine, ErrDataStoreConn)
	}
	switch opts.DataSourceService {
	case mpb.StorageService_GCP_CLOUD_STORAGE.String():
		decodeData, err := base64.StdEncoding.DecodeString(opts.Creds.DsCreds.GcpCreds.ServiceAccountKey)
		if err != nil {
			log.Err(err).Msg("Error decoding gcp service account key")
			return nil, err
		}
		client, err := tpgcp.NewBucketAgent(ctx, &tpgcp.GcpConfig{
			ServiceAccountKey: []byte(decodeData),
			ProjectID:         opts.Creds.DsCreds.GcpCreds.ProjectId,
			Region:            opts.Creds.DsCreds.GcpCreds.Region,
		}, logger)
		if err != nil {
			log.Err(err).Msg("Error creating gcp secret manager client")
			return nil, err
		}
		resultCl.BlobStore = client
		return resultCl, nil
	default:
		return nil, dmerrors.DMError(ErrInvalidBlobStoreService, ErrDataStoreConn)
	}
}
