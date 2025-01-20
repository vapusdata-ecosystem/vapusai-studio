package pkgs

import (
	"context"
	"path/filepath"

	"github.com/rs/zerolog"
	serviceops "github.com/vapusdata-oss/aistudio/core/serviceops"
	coreutils "github.com/vapusdata-oss/aistudio/core/utils"
	utils "github.com/vapusdata-oss/aistudio/webapp/utils"
)

var VapusSvcInternalClientManager *serviceops.VapusSvcInternalClients

func InitNetworkConfig(configRoot, fileName string) error {
	DmLogger.Info().Msgf("Reading network configuration with path - %v ", filepath.Join(configRoot, fileName))

	cf, err := coreutils.ReadBasicConfig(coreutils.GetConfFileType(fileName), filepath.Join(configRoot, fileName), &serviceops.NetworkConfig{})
	if err != nil {
		DmLogger.Panic().Err(err).Msg("error while loading service config")
		return err
	}

	svcnetConf, ok := cf.(*serviceops.NetworkConfig)
	if !ok {
		DmLogger.Panic().Msg("error while loading network config")
		return utils.ErrInvalidNetworkConfig
	} else {
		NetworkConfigManager = svcnetConf
	}
	return nil
}

func InitVapusSvcInternalClients(hostSvc string, logger zerolog.Logger) {
	// TODO: Handle error
	err := serviceops.SvcUpTimeCheck(context.Background(), NetworkConfigManager, "", logger, 0)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while checking service uptime.")
	} else {
		logger.Info().Msg("service is up and running.")
	}
	res, err := serviceops.SetupVapusSvcInternalClients(context.Background(), NetworkConfigManager, "", logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while initializing vapus svc internal clients.")
	}
	VapusSvcInternalClientManager = res
}
