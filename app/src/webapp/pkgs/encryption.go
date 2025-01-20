package pkgs

import (
	validator "github.com/go-playground/validator/v10"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
)

var DmAuth *encryption.VapusDataJwtAuthn

func NewDMAuth(path string, validator *validator.Validate) (*encryption.VapusDataJwtAuthn, error) {
	obj, err := encryption.NewVapusDataJwtAuthnWithConfig(path)

	if err != nil {
		pkgLogger.Panic().Err(err).Msg("error while loading jwt authn for dm access")
		return nil, err
	}

	if err := validator.Struct(obj.Opts); err != nil {
		pkgLogger.Panic().Err(err).Msg("error while JWT authn object")
		return nil, err
	}

	return obj, nil
}

func InitAuthn(path string, validator *validator.Validate) {
	var err error
	if DmAuth == nil {
		NewDMAuth(path, validator)
		DmAuth, err = NewDMAuth(path, validator)
		if err != nil {
			pkgLogger.Panic().Err(err).Msg("error while while init JWT authn")
			panic(err)
		}
	}
}

func ValidateAccessToken(tokenString string) (*encryption.VapusDataStudioAccessClaims, error) {
	claim, err := DmAuth.ParseAndValidateAT(tokenString)
	if err != nil {
		return nil, err
	}
	return claim, nil
}
