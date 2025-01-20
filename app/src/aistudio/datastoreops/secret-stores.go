package dmstores

import (
	"context"
	"encoding/json"
	"log"

	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	serviceopss "github.com/vapusdata-oss/aistudio/core/serviceops"
)

// NewVapusBESecretStore creates a new BeSecretStore object with driver client of different secret backend store based on the configuration provided
func NewVapusBESecretStore(path string) (*serviceopss.SecretStore, error) {
	logger.Debug().Msgf("Creating secret store client with path: %s", path)
	ctx := context.Background()
	client, err := serviceopss.NewSecretStoreClient(ctx, logger, serviceopss.WithSecretStorePath(path))
	if err != nil {
		logger.Info().Msgf("Error while creating secret store client: %v", err)
		return nil, err
	}
	return client, nil
}

func InitStoreDependencies(ctx context.Context, conf *serviceopss.VapusAISvcConfig) {
	if pkgs.JwtAuthnParams == nil {
		bootJwtAuthn(ctx, conf.GetJwtAuthSecretPath())
	}
}

func bootJwtAuthn(ctx context.Context, secName string) {
	logger.Info().Msgf("Boot Jwt Authn with secret path: %s", secName)
	secretStr, err := DMStoreManager.ReadSecret(ctx, secName)
	if err != nil {
		logger.Fatal().Err(err).Msgf("error while reading Jwt secret %s", secName)
	}
	tmp := &encryption.JWTAuthn{}
	log.Println("Secret data:  ++++++++++++++++++++++++++++  ", secretStr)
	err = json.Unmarshal([]byte(secretStr), tmp)
	if err != nil {
		logger.Fatal().Err(err).Msgf("error while unmarshalling Jwt secret %s", secName)
	}
	pkgs.JwtAuthnParams = tmp
}
