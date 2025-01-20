package svcops

import (
	"context"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	datasvc "github.com/vapusdata-oss/aistudio/core/dataservices"
	dmlogger "github.com/vapusdata-oss/aistudio/core/logger"
	models "github.com/vapusdata-oss/aistudio/core/models"
	"gopkg.in/yaml.v2"
)

type SecretStore struct {
	Path  string
	creds *models.StoreParams
	*datasvc.SecretStoreClient
}

type secretOpts func(*SecretStore)

func WithSecretStorePath(r string) secretOpts {
	return func(s *SecretStore) {
		s.Path = r
	}
}

func NewSecretStoreCreds(filePath string, log zerolog.Logger) (*models.StoreParams, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal().Msgf("Error reading file: %v", err)
		return nil, err
	}
	conf := &models.StoreParams{}
	log.Info().Msgf("File Content after reading : %s", string(bytes))
	err = yaml.Unmarshal(bytes, conf)
	if err != nil {
		log.Fatal().Msgf("Error unmarshalling file: %v", err)
		return nil, err
	}
	err = validator.New().Struct(conf)

	if err != nil {
		return nil, err
	}
	return conf, nil
}

func NewSecretStoreClient(ctx context.Context, log zerolog.Logger, opts ...secretOpts) (*SecretStore, error) {
	s := &SecretStore{}
	for _, opt := range opts {
		opt(s)
	}
	creds, err := NewSecretStoreCreds(s.Path, log)
	if err != nil {
		return nil, err
	}
	client, err := datasvc.NewSecretStoreClient(ctx, creds, dmlogger.GetSubDMLogger(log, "datasvc", "secret_store_main"))
	if err != nil {
		return nil, err
	}
	s.creds = creds
	s.SecretStoreClient = client
	return s, nil
}

func (s *SecretStore) GetCreds() *models.NetworkParams {
	return s.GetCreds()
}
