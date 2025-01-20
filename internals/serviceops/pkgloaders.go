package svcops

import (
	"log"

	validator "github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/vapusdata-oss/aistudio/core/authn"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	pbac "github.com/vapusdata-oss/aistudio/core/pbac"
	utils "github.com/vapusdata-oss/aistudio/core/utils"
)

type VapusSvcPackageParams struct {
	JwtParams      *encryption.JWTAuthn
	AuthnParams    *authn.AuthnSecrets
	PbacConfigPath string
}

type VapusSvcPackages struct {
	VapusJwtAuth         *encryption.VapusDataJwtAuthn
	AuthnManager         *authn.Authenticator
	StudioRBACManager    *pbac.PbacConfig
	GrpcRequestValidator *utils.DMValidator
	ModelValidator       *validator.Validate
}

type VapusSvcPkgOpts func(*VapusSvcPackageParams)

func WithJwtParams(params *encryption.JWTAuthn) VapusSvcPkgOpts {
	return func(p *VapusSvcPackageParams) {
		p.JwtParams = params
	}
}

func WithAuthnParams(params *authn.AuthnSecrets) VapusSvcPkgOpts {
	return func(p *VapusSvcPackageParams) {
		p.AuthnParams = params
	}
}

func WithPbacConfigPath(path string) VapusSvcPkgOpts {
	return func(p *VapusSvcPackageParams) {
		p.PbacConfigPath = path
	}
}

func InitSvcPackages(params *VapusSvcPackageParams, VapusSvcPackageManager *VapusSvcPackages, logger zerolog.Logger, opts ...VapusSvcPkgOpts) (*VapusSvcPackageParams, *VapusSvcPackages, error) {
	var err error
	if params == nil {
		params = &VapusSvcPackageParams{}
		for _, opt := range opts {
			opt(params)
		}
	}
	if VapusSvcPackageManager == nil {
		VapusSvcPackageManager = &VapusSvcPackages{}
	}
	VapusSvcPackageManager.ModelValidator = validator.New()
	if VapusSvcPackageManager.GrpcRequestValidator == nil {
		VapusSvcPackageManager.GrpcRequestValidator, err = utils.NewDMValidator()
		if err != nil {
			logger.Err(err).Msg("Error while loading validator")
			return nil, nil, ErrValidatorInitFailed
		}
		logger.Info().Msg("GRPC request Validator initialized successfully")
	}
	if VapusSvcPackageManager.AuthnManager == nil {
		if params.AuthnParams == nil {
			logger.Err(ErrAuthenticatorParamsNil).Msg("Error while loading authn config")
			return nil, nil, ErrAuthenticatorParamsNil
		}
		log.Println("params.AuthnParams - ", params.AuthnParams)
		VapusSvcPackageManager.AuthnManager, err = authn.New(params.AuthnParams.OIDCSecrets, params.AuthnParams.AuthnMethod)
		if err != nil {
			logger.Err(err).Msg("Error while initializing authenticator")
			return nil, nil, ErrAuthenticatorInitFailed
		}
		logger.Info().Msg("Authenticator initialized successfully")
	}
	if VapusSvcPackageManager.VapusJwtAuth == nil {
		if params.JwtParams == nil {
			logger.Err(ErrJwtParamsNil).Msg("Error while loading jwt config")
			return nil, nil, ErrJwtParamsNil
		}
		log.Println("params.JwtParams - ", params.JwtParams)
		VapusSvcPackageManager.VapusJwtAuth, err = encryption.NewVapusDataJwtAuthn(params.JwtParams)
		if err != nil {
			logger.Err(err).Msg("Error while initializing jwt authn")
			return nil, nil, ErrJwtAuthInitFailed
		}
		if err := VapusSvcPackageManager.ModelValidator.Struct(VapusSvcPackageManager.VapusJwtAuth.Opts); err != nil {
			logger.Err(err).Msg("Error while validating jwt config")
			return nil, nil, err
		}
		logger.Info().Msg("JWT Authn initialized successfully")
	}
	if VapusSvcPackageManager.StudioRBACManager == nil {
		log.Println(params.PbacConfigPath)
		if params.PbacConfigPath == "" {
			logger.Err(ErrPbacConfigPathEmpty).Msg("Error while loading pbac config")
			return nil, nil, ErrPbacConfigPathEmpty
		}
		VapusSvcPackageManager.StudioRBACManager, err = pbac.LoadPbacConfig(params.PbacConfigPath)
		if err != nil {
			logger.Err(err).Msg("Error while initializing pbac config")
			return params, VapusSvcPackageManager, ErrPbacConfigInitFailed
		}
		logger.Info().Msg("PBAC config initialized successfully")
	}
	logger.Info().Msg("Service packages initialized successfully")
	log.Println(VapusSvcPackageManager)
	return params, VapusSvcPackageManager, nil
}
