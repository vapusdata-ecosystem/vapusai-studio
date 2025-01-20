package pkgs

import (
	"log"

	"github.com/rs/zerolog"
	serviceops "github.com/vapusdata-oss/aistudio/core/serviceops"
)

var SvcPackageManager *serviceops.VapusSvcPackages
var SvcPackageParams *serviceops.VapusSvcPackageParams

func InitSvcPackageParams() {
	SvcPackageParams = &serviceops.VapusSvcPackageParams{}
}

func InitStudioSvcPackages(logger zerolog.Logger, opts ...serviceops.VapusSvcPkgOpts) error {
	if SvcPackageParams == nil {
		SvcPackageParams = &serviceops.VapusSvcPackageParams{}
	}
	for _, opt := range opts {
		opt(SvcPackageParams)
	}
	var err error
	SvcPackageParams, SvcPackageManager, err = serviceops.InitSvcPackages(SvcPackageParams, SvcPackageManager, logger, opts...)
	if err != nil {
		return err
	}
	log.Println("==================================start")
	log.Println(SvcPackageManager)
	log.Println(SvcPackageParams)
	log.Println("==================================end")
	return nil
}
