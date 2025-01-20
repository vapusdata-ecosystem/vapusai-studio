package datasvc

import (
	"context"
	"strings"

	"github.com/rs/zerolog"
	datasvcpkgs "github.com/vapusdata-oss/aistudio/core/dataservices/pkgs"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

func (dsc *DataStoreClient) prepareSqlQuery(qopts *datasvcpkgs.QueryOpts, logger zerolog.Logger) (*string, error) {
	query := "SELECT {field} FROM {table} where {condition}"
	query = strings.Replace(query, "{table}", qopts.DataCollection, -1)
	if qopts.QueryString != "" {
		query = strings.Replace(query, "{condition}", qopts.QueryString, -1)
	}
	if qopts.CountRecords {
		query = strings.Replace(query, "{field}", "count(*)", -1)
	} else {
		if len(qopts.IncludeFields) > 0 {
			query = strings.Replace(query, "{field}", strings.Join(qopts.IncludeFields, ","), -1)
		} else {
			query = strings.Replace(query, "{field}", "*", -1)
		}
	}
	logger.Info().Msgf("Query prepared by opts dataSVC - %v", query)
	return &query, nil
}

func (dsc *DataStoreClient) SelectWithFilter(ctx context.Context, qopts *datasvcpkgs.QueryOpts, resultObj interface{}, logger zerolog.Logger) ([]map[string]interface{}, error) {
	var err error
	result := make([]map[string]interface{}, 0)
	switch dsc.DataStoreParams.DataSourceEngine {
	case mpb.StorageEngine_POSTGRES.String():
		result, err = dsc.PostgresClient.SelectWithFilter(ctx, qopts)
		if err != nil {
			return result, err
		}
		return result, nil
	}
	return result, ErrInvalidDataStorageEngine
}

func (dsc *DataStoreClient) Select(ctx context.Context, qopts *datasvcpkgs.QueryOpts, dest interface{}, logger zerolog.Logger) error {
	// TO:DO - Add support for other data storage engines - https://github.com/Masterminds/squirrel
	switch dsc.DataStoreParams.DataSourceEngine {
	case mpb.StorageEngine_POSTGRES.String():
		var query *string
		if qopts.RawQuery != "" {
			query = &qopts.RawQuery
		} else {
			q, err := dsc.prepareSqlQuery(qopts, logger)
			if err != nil {
				return err
			}
			query = q
		}
		rows, err := dsc.PostgresClient.Select(ctx, query)
		if err != nil {
			return err
		}
		return ScanSql(rows, dest, logger)
	}
	return ErrInvalidDataStorageEngine
}
