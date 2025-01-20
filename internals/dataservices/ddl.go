package datasvc

import (
	"context"

	pkgs "github.com/vapusdata-oss/aistudio/core/dataservices/pkgs"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

func (svc *DataStoreClient) RunDDLs(ctx context.Context, request *string) error {
	switch svc.DataStoreParams.DataSourceEngine {
	case mpb.StorageEngine_POSTGRES.String():
		return svc.PostgresClient.RunDDL(ctx, request)
	default:
		return ErrInvalidDataStorageEngine
	}
}

func (svc *DataStoreClient) CreateDataTables(ctx context.Context, opts *pkgs.DataTablesOpts) error {
	switch svc.DataStoreParams.DataSourceEngine {
	case mpb.StorageEngine_POSTGRES.String():
		return svc.PostgresClient.CreateDataTables(ctx, opts)
	default:
		return ErrInvalidDataStorageEngine
	}
}
