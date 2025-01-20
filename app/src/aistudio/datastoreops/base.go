package dmstores

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	models "github.com/vapusdata-oss/aistudio/core/models"
	serviceopss "github.com/vapusdata-oss/aistudio/core/serviceops"
	vapusloader "github.com/vapusdata-oss/aistudio/core/serviceops"
)

type DMStore struct {
	// Inheriting VapusBESecretStorage for secret store
	*vapusloader.SecretStore

	// Backend store for data mesh data
	*BeDataStore
	BlobOps *vapusloader.BlobStore
	Error   error
}

// GLobal var for DM store, it can accessed across the service
var (
	DMStoreManager *DMStore
	logger         zerolog.Logger
)

// Constructor to create new object for DMStore struct
func newDMStore(conf *serviceopss.VapusAISvcConfig) *DMStore {
	dmSec, err := NewVapusBESecretStore(conf.GetSecretStoragePath())
	if err != nil {
		return &DMStore{
			Error: err,
		}
	}
	dmDb := NewVapusBEDataStore(conf, dmSec)
	if dmDb.Error != nil {
		return &DMStore{
			Error: err,
		}
	}
	blobOps, err := GetBlobStore(conf.GetFileStorePath(), dmSec)
	if err != nil {
		return &DMStore{
			Error: err,
		}
	}
	return &DMStore{
		SecretStore: dmSec,
		BeDataStore: dmDb,
		BlobOps:     blobOps,
	}
}

// Initializing DMStore struct with object and global var
func InitDMStore(conf *serviceopss.VapusAISvcConfig) {
	logger = pkgs.GetSubDMLogger(pkgs.DSTORES, "DBStore")
	if DMStoreManager == nil || DMStoreManager.SecretStore == nil || DMStoreManager.BeDataStore == nil {
		DMStoreManager = newDMStore(conf)
	}
}

func GetBlobStore(path string, dmSec *vapusloader.SecretStore) (*vapusloader.BlobStore, error) {
	logger = pkgs.GetSubDMLogger(pkgs.DSTORES, "BlobStore")
	creds := &models.StoreParams{}
	secretStr, err := dmSec.ReadSecret(context.TODO(), path)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while reading BlobStore store secret data")
	}
	err = json.Unmarshal([]byte(secretStr), creds)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while unmarshalling BlobStore store secret data")
	}
	return vapusloader.NewBlobStoreClient(context.Background(), creds, logger)
}
