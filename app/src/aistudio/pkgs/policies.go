package pkgs

import (
	pbac "github.com/vapusdata-oss/aistudio/core/pbac"
)

var StudioRBACManager *pbac.PbacConfig

func InitPolicyLib(configPath string) {
	if StudioRBACManager == nil {
		val, err := pbac.LoadPbacConfig(configPath)
		if err != nil {
			pkgLogger.Err(err).Msg("error initializing and loading of pbac config")
			panic(err)
		}
		StudioRBACManager = val
	}
}
