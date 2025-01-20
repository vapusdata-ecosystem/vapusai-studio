package dmstores

import (
	"context"
	"encoding/json"
	"fmt"

	datasvc "github.com/vapusdata-oss/aistudio/core/dataservices"
	dmlogger "github.com/vapusdata-oss/aistudio/core/logger"
	models "github.com/vapusdata-oss/aistudio/core/models"
	serviceopss "github.com/vapusdata-oss/aistudio/core/serviceops"
	vapusloader "github.com/vapusdata-oss/aistudio/core/serviceops"
)

var Account *models.Account

// BeDbStore is the dbstores for Data Mesh data, containing the Redis client currently
type BeDataStore struct {
	Db                  *datasvc.DataStoreClient
	PubSub              *datasvc.DataStoreClient
	Cacher              *datasvc.DataStoreClient
	IsRedisStackEnabled bool
	Error               error
}

// NewVapusBEDataStore creates a new BeDbStore object with driver client of different db backend store
func NewVapusBEDataStore(conf *serviceopss.VapusAISvcConfig, secretClient *vapusloader.SecretStore) *BeDataStore {
	ctx := context.Background()
	bds := &BeDataStore{}
	dbClient, err := initDbStores(ctx, conf.GetDBStoragePath(), secretClient)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while initializing db data store")
	}
	bds.Db = dbClient
	cacheClient, err := initDbStores(ctx, conf.GetCachStoragePath(), secretClient)
	if err != nil {
		logger.Error().Err(err).Msg("error while initializing cache data store")
	}
	bds.Cacher = cacheClient

	if Account == nil {
		bootAccount(bds)
	}
	return bds
}

func initDbStores(ctx context.Context, secName string, secretClient *vapusloader.SecretStore) (*datasvc.DataStoreClient, error) {
	secretStr, err := secretClient.ReadSecret(ctx, secName)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while reading secret data for data store")
		return nil, err
	}
	creds := &models.StoreParams{}
	err = json.Unmarshal([]byte(secretStr), creds)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while unmarshalling secret data")
		return nil, err
	}
	return datasvc.New(ctx, datasvc.WithInApp(true), datasvc.WithStoreParams(creds), datasvc.WithLogger(dmlogger.GetSubDMLogger(logger, "Aistudio Dbstore", "DataStoreClient")))
}

func (ds *DMStore) GetDbStoreParams() *models.StoreParams {
	return ds.BeDataStore.Db.DataStoreParams
}

func bootAccount(dbStore *BeDataStore) {
	var result []*models.Account
	query := fmt.Sprintf("SELECT * FROM %v", serviceopss.AccountsTable)
	err := dbStore.Db.PostgresClient.SelectInApp(context.TODO(), &query, &result)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while fetching accounts info from datastore")
	}
	if len(result) == 0 {
		logger.Fatal().Msg("No account found in the system")
	}
	Account = result[0]
}
