package pkgs

import (
	"path/filepath"

	serviceops "github.com/vapusdata-oss/aistudio/core/serviceops"
	utils "github.com/vapusdata-oss/aistudio/core/utils"
)

var WebAppConfigManager *serviceops.WebAppConfig
var NetworkConfigManager *serviceops.NetworkConfig

func InitServiceConfig(configRoot, path string) {
	DmLogger.Info().Msgf("WebAppConfigManager current value - %v", WebAppConfigManager)
	if WebAppConfigManager == nil {
		DmLogger.Info().Msg("Initializing the WebAppConfigManager")
		WebAppConfigManager = loadServiceConfig(configRoot, path)
	}
}

func loadServiceConfig(configRoot, fileName string) *serviceops.WebAppConfig {
	// Read the service configuration from the file
	DmLogger.Info().Msgf("Reading service configuration with path - %v", filepath.Join(configRoot, fileName))
	cf, err := utils.ReadBasicConfig(utils.GetConfFileType(fileName), filepath.Join(configRoot, fileName), &serviceops.WebAppConfig{})
	if err != nil {
		DmLogger.Panic().Err(err).Msg("error while loading webapp config")
	}

	svcConf := cf.(*serviceops.WebAppConfig)
	svcConf.Path = configRoot
	return svcConf
}
