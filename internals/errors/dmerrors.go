package dmerrors

import "errors"

var (
	// Error constants for JSON operations
	ErrJsonMarshel   = errors.New("error while marshalling JSON")
	ErrJsonUnMarshel = errors.New("error while unmarshalling JSON")
	ErrStruct2Json   = errors.New("failed to convert struct to json")

	// Error constants for viper operations
	ErrViperConfigRead = errors.New("error while reading configuration file")
	ErrViperConfigSet  = errors.New("error while setting configuration file")

	ErrUserOrganization404            = errors.New("error- invalid organization requested, user is not attached to requested organization")
	ErrWriteYAMLFile                  = errors.New("error while writing to yaml file")
	ErrInvalidArgs                    = errors.New("invalid arguments provided for the command, please provide the required arguments")
	ErrNoCredentialFoundForDataSource = errors.New("no credentials found for the data source")
	ErrInvalidComplianceAgentParams   = errors.New("invalid compliance agent parameters")
	ErrInvalidComplianceAction        = errors.New("invalid compliance action")
)
