package datasvc

import (
	"context"

	"github.com/rs/zerolog"
	datasvcpkgs "github.com/vapusdata-oss/aistudio/core/dataservices/pkgs"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

func (dsc *DataStoreClient) CountRows(ctx context.Context, qopts *datasvcpkgs.QueryOpts, logger zerolog.Logger) (int64, error) {
	switch dsc.DataStoreParams.DataSourceEngine {
	case mpb.StorageEngine_POSTGRES.String():
		return dsc.PostgresClient.Count(ctx, qopts)
	default:
		logger.Error().Msg("Invalid data storage engine")
		return 0, ErrInvalidDataStorageEngine
	}
}
