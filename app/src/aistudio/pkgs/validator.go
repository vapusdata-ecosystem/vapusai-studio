package pkgs

import dmutils "github.com/vapusdata-oss/aistudio/core/utils"

var platformRequestValidator *dmutils.DMValidator

func initStudioRequestValidator() {
	validator, err := dmutils.NewDMValidator()
	if err != nil {
		pkgLogger.Panic().Err(err).Msg("Error while loading validator")
	}
	platformRequestValidator = validator
}

func GetStudioRequestValidator() *dmutils.DMValidator {
	if platformRequestValidator == nil {
		initStudioRequestValidator()
	}
	return platformRequestValidator
}
