package services

import (
	"github.com/rs/zerolog"
	dmstores "github.com/vapusdata-oss/aistudio/aistudio/datastoreops"
	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
)

type StudioServices struct {
	DMStore *dmstores.DMStore
	logger  zerolog.Logger
}

var StudioServicesManager *StudioServices
var helperLogger zerolog.Logger

func newStudioServices(dmstore *dmstores.DMStore) *StudioServices {
	return &StudioServices{
		DMStore: dmstore,
	}
}

func InitStudioServices(dmstore *dmstores.DMStore) {
	InitStudioServices(dmstore)
	if StudioServicesManager == nil {
		StudioServicesManager = newStudioServices(dmstore)
		StudioServicesManager.logger = pkgs.GetSubDMLogger(pkgs.SVCS, "StudioServices")
	}
}
