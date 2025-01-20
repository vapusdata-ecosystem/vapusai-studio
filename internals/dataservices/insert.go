package datasvc

import (
	"context"

	datasvcpkgs "github.com/vapusdata-oss/aistudio/core/dataservices/pkgs"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

func (svc *DataStoreClient) InsertBulkDataSet(ctx context.Context, request *datasvcpkgs.InsertDataRequest) (*datasvcpkgs.InsertDataResponse, error) {
	switch svc.DataStoreParams.DataSourceEngine {
	case mpb.StorageEngine_POSTGRES.String():
		return svc.PostgresClient.InsertInBulk(ctx, request, svc.Logger)
	default:
		return nil, ErrInvalidDataStorageEngine
	}
}

func (svc *DataStoreClient) InsertDataSet(ctx context.Context, request *datasvcpkgs.InsertDataRequest) (*datasvcpkgs.InsertDataResponse, error) {
	resp := &datasvcpkgs.InsertDataResponse{
		DataTable:       request.DataTable,
		RecordsInserted: 0,
	}
	var err error
	switch svc.DataStoreParams.DataSourceEngine {
	case mpb.StorageEngine_POSTGRES.String():
		err = svc.PostgresClient.Insert(ctx, request, svc.Logger)
	default:
		return nil, ErrInvalidDataStorageEngine
	}
	if err != nil {
		resp.RecordsFailed = 1
		return resp, err
	} else {
		resp.RecordsInserted = 1
		return resp, nil
	}
}
