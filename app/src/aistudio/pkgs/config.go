package pkgs

import (
	utils "github.com/vapusdata-oss/aistudio/aistudio/utils"
	serviceopss "github.com/vapusdata-oss/aistudio/core/serviceops"
	dmutils "github.com/vapusdata-oss/aistudio/core/utils"
)

var ServiceConfigManager *serviceopss.VapusAISvcConfig
var NetworkConfigManager *serviceopss.NetworkConfig

func newServiceConfig(configRoot, path string) *serviceopss.VapusAISvcConfig {
	return LoadServiceConfig(configRoot, path)
}

func InitServiceConfig(configRoot, path string) {
	if ServiceConfigManager == nil {
		ServiceConfigManager = newServiceConfig(configRoot, path)
	}
}

func LoadServiceConfig(configRoot, path string) *serviceopss.VapusAISvcConfig {
	// Read the service configuration from the file
	DmLogger.Info().Msgf("Reading service configuration with path - %v ", path)

	cf, err := dmutils.ReadBasicConfig(dmutils.GetConfFileType(path), path, &serviceopss.VapusAISvcConfig{})
	if err != nil {
		DmLogger.Panic().Err(err).Msg("error while loading service config")
		return nil
	}

	svcConf := cf.(*serviceopss.VapusAISvcConfig)
	svcConf.Path = configRoot
	return svcConf
}

func InitNetworkConfig(configRoot, path string) error {
	DmLogger.Info().Msgf("Reading network configuration with path - %v ", path)

	cf, err := dmutils.ReadBasicConfig(dmutils.GetConfFileType(path), path, &serviceopss.NetworkConfig{})
	if err != nil {
		DmLogger.Panic().Err(err).Msg("error while loading service config")
		return err
	}

	svcnetConf, ok := cf.(*serviceopss.NetworkConfig)
	if !ok {
		DmLogger.Panic().Msg("error while loading network config")
		return utils.ErrInvalidNetworkConfig
	} else {
		NetworkConfigManager = svcnetConf
	}
	return nil
}
