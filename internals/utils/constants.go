package utils

const (
	WINDOWS = "windows"
	MACOS   = "macos"
	LINUX   = "linux"

	EMPTYSTR  = ""
	LOCALHOST = "localhost"

	DEFAULTSERVERKEYTLS  = "x509/server_key.pem"
	DEFAULTSERVERCERTTLS = "x509/server_cert.pem"

	DEFAULT_CONFIG_TYPE = "json"
	JSON                = "json"
	TOML                = "toml"
	DOT                 = "."

	SERVICE_CONFIG_READ_ERROR = "failed to read the service configuration file"

	HASHICORPVAULT    = "hashicorpVault"
	AWS_SECRETMANAGER = "awsSecretManager"

	CUSTOM_MESSAGE = "customMessage"

	ERROR_CTX   = "error"
	WARNING_CTX = "warning"
	INFO_CTX    = "info"
)

const (
	REDACT_MASK                       = "X"
	REDACT_MASK_RUNE                  = 'X'
	REDACT_SPACE_MASK_RUNE            = '_'
	REDACT_SYMBOL_MASK_RUNE           = '+'
	REDACT_SPACE_MASK                 = "_"
	REDACT_SYMBOL_MASK                = "+"
	DATAPRODUCT_ENV                   = "dataProductId"
	DATACONTAINER_ENV                 = "dataWorkerId"
	DATAPRODUCT_STATUS_ENV            = "dataProductStatus"
	DATAPRODUCT_DEPLOYMENT_STATUS_ENV = "dataProductDeploymentStatus"
	DATAPRODUCT_CONFIG_DIGEST         = "dataProductSpecDigest"
	VDC_IDEN                          = "vdc"
	DATAWORKER_IDEN                   = "dataworker"
)

const (
	JWT_SECRET_SUFFIX               = "_organization_jwt"
	DATASOURCE_SECRET_SUFFIX        = "_secrets"
	DEPLOYMENT_CONFIG_SUFFIX        = "_deployment_config"
	VDC_PLATFORM_BRIDGE_FILE_SUFFIX = "vdc_platform_bridge"
	DEFAULT_VDC_MOUNT               = "/mnt/vapuscontainers"
	VDC_AUTHN_FILE_NAME             = "vdc_authn.yaml"
	VDC_SECRETS_FILE_NAME           = "secrets/container-secrets.yaml"
	VDC_CONFIG_FILE_NAME            = "config/container-config.yaml"
	VDC_K8S_MOUNT_PATH              = "/data/vapusdata/dataproduct"
	MOUNT_CONFIG_PATH               = "mountConfigPath"
)
