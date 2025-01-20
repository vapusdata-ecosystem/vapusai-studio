package datasvc

import (
	"context"

	"github.com/rs/zerolog"
	vpostgres "github.com/vapusdata-oss/aistudio/core/dataservices/postgres"
	vredis "github.com/vapusdata-oss/aistudio/core/dataservices/redis"
	"github.com/vapusdata-oss/aistudio/core/models"

	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	logger "github.com/vapusdata-oss/aistudio/core/logger"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

type DataStoreClient struct {
	RedisClient     *vredis.RedisStore
	PostgresClient  *vpostgres.PostgresStore
	DataStoreParams *models.StoreParams
	Logger          zerolog.Logger
	Debug           bool
	InApp           bool
}

var (
	log = logger.CoreLogger
)

func (d *DataStoreClient) Close() {
	if d.RedisClient != nil {
		d.RedisClient.Close()
	}
	if d.PostgresClient != nil {
		d.PostgresClient.Close()
	}
}

type DataStoreOpts func(*DataStoreClient)

func WithDebug(debug bool) DataStoreOpts {
	return func(d *DataStoreClient) {
		d.Debug = debug
	}
}

func WithInApp(inApp bool) DataStoreOpts {
	return func(d *DataStoreClient) {
		d.InApp = inApp
	}
}

func WithLogger(log zerolog.Logger) DataStoreOpts {
	return func(d *DataStoreClient) {
		d.Logger = log
	}
}

func WithStoreParams(params *models.StoreParams) DataStoreOpts {
	return func(d *DataStoreClient) {
		d.DataStoreParams = params
	}
}

// NewDataConnClient creates a new DataConnClient
func New(ctx context.Context, opts ...DataStoreOpts) (*DataStoreClient, error) {
	var err error
	dsc := &DataStoreClient{}
	for _, opt := range opts {
		opt(dsc)
	}
	if dsc.DataStoreParams == nil {
		return nil, dmerrors.DMError(ErrDataStoreParams404, ErrDataStoreParams404)
	}
	log.Debug().Msgf("Creating new data source connection for %s", dsc.DataStoreParams.DataSourceEngine)
	switch dsc.DataStoreParams.DataSourceEngine {
	case mpb.StorageEngine_POSTGRES.String():
		var client *vpostgres.PostgresStore
		crds := &vpostgres.PostgresOpts{
			URL:      dsc.DataStoreParams.Creds.Address,
			Username: dsc.DataStoreParams.Creds.DsCreds.GenericCredentialModel.Username,
			Password: dsc.DataStoreParams.Creds.DsCreds.GenericCredentialModel.Password,
			Database: dsc.DataStoreParams.Creds.DsCreds.DB,
			Port:     int(dsc.DataStoreParams.Creds.Port),
		}
		if dsc.InApp {
			client, err = vpostgres.NewPostgresStoreLocal(crds, log)
		} else {
			client, err = vpostgres.NewPostgresStore(crds, log)
		}

		if err != nil {
			log.Err(err).Msg("Error connecting to postgres")
			return nil, err
		}
		return &DataStoreClient{PostgresClient: client, DataStoreParams: dsc.DataStoreParams, Logger: log}, nil
	default:
		return nil, dmerrors.DMError(ErrInvalidDataStoreEngine, ErrDataStoreConn)
	}
}

func connectRedis(ctx context.Context, opts *models.NetworkParams) (*vredis.RedisStore, error) {
	return vredis.NewRedisStore(ctx, &vredis.Redis{
		URL:      opts.Address,
		Port:     int(opts.Port),
		Password: opts.DsCreds.GenericCredentialModel.Password,
		Username: opts.DsCreds.GenericCredentialModel.Username,
	}, log)
}
