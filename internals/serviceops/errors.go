package svcops

import "errors"

var (
	ErrAuthenticatorInitFailed = errors.New("error while initializing authenticator for the service based on the provided configuration")
	ErrAuthenticatorParamsNil  = errors.New("error while initializing authenticator for the service, authn params are nil")
	ErrJwtParamsNil            = errors.New("error while initializing jwt authn for the service, jwt params are nil")
	ErrJwtAuthInitFailed       = errors.New("error while initializing jwt authn for the service based on the provided configuration")
	ErrValidatorInitFailed     = errors.New("error while initializing validator for the service")
	ErrPbacConfigPathEmpty     = errors.New("error while initializing pbac config for the service, pbac config path is empty")
	ErrPbacConfigInitFailed    = errors.New("error while initializing pbac config for the service based on the provided configuration")
)

var (
	ErrAIStudioConnNotInitialized          = errors.New("AIStudio connection not initialized")
	ErrUserConnNotInitialized              = errors.New("User connection not initialized")
	ErrDataProductServerConnNotInitialized = errors.New("DataProductServer connection not initialized")
	ErrPlatConnNotInitialized              = errors.New("Studio connection not initialized")
	Erruser404                             = errors.New("User not found")
	ErrNoDataFound                         = errors.New("No data found")
	ErrPluginActionFailed                  = errors.New("Plugin action failed, this plugin is not available for your account, enable it from the plugin store under settings section")
	ErrDataProductQueryCallFailed          = errors.New("Data product query call failed")
	ErrGeneratingContent                   = errors.New("Error generating content")
	ErrInvalidPluginActionRequest          = errors.New("Invalid plugin action request")
)
