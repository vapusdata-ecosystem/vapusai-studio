package utils

import (
	"errors"
)

var (
	ErrNoCredentialFoundForDataSource = errors.New("no credential found for data node")
	ErrNoDataSourceFound              = errors.New("no data node found")
	ErrListingArtifacts               = errors.New("error while listing artifacts")
	ErrNoNameFoundForPackage          = errors.New("no name found for package")
	ErrDataSourceAttributesNotFound   = errors.New("error: no attributes found in data node")
	ErrDecodingCredential             = errors.New("error while decoding node credential")
	ErrNoPackagesFoundInECR           = errors.New("no packages found in ecr")
	ErrConnDataSource                 = errors.New("error while connecting to data node")

	ErrNoMappings = errors.New("no mappings found in ElasticSearch for this index")
	ErrNoProps    = errors.New("no properties found in ElasticSearch for this index")

	ErrInvalidDataSourceType = errors.New("invalid data source")

	ErrTokenExpired    = errors.New("token expired")
	ErrUnAuthenticated = errors.New("user is not authenticated")
)

var (
	ErrInvalidDataProductFormat = errors.New("invalid format for data product spec")
	ErrInvalidNabhikWorkflow    = errors.New("invalid workflow requested from nabhik engine")

	ErrInvalidNameSpaceOperation = errors.New("invalid action for namespace operations")
	ErrProductDeploymentFailed   = errors.New("data product deployment failed in k8s")

	ErrInvalidNabhikAgent     = errors.New("invalid configuration for nabhik agent")
	ErrInvalidNabhikOperation = errors.New("invalid action for nabhik agent")
	ErrNoDPContainers         = errors.New("error: no data workers are configured in this data product")

	ErrMissingOrganizationArtifactStore = errors.New("no artifact stores are present for current organization/platform")
	ErrImagePullSecretOperation         = errors.New("error while fetching or creating organization pull secrets")

	ErrInvalidDataProductPublisherAgent = errors.New("invalid configuration for data product publisher agent")
	ErrOrganizationJwtSecretFailed      = errors.New("error while fetching/creating organization jwt secret")
	ErrInvalidDataProductAgentOperation = errors.New("invalid action for data product agent")
	ErrBuildingDataProduct              = errors.New("error while building data product")

	ErrInvalidDataWorkerAgentOperation = errors.New("invalid action for data worker agent")
	ErrDataWorkerConfig404             = errors.New("data worker config not found")

	ErrInvalidK8SServiceType = errors.New("invalid k8s service type")

	ErrUnsupportedServiceProvider = errors.New("unsupported service provider for K8s Cluster")
	ErrK8SInfraParamsNil          = errors.New("K8S Infra Params is nil")

	ErrTrinoIsReadOnly = errors.New("trino is in read only mode")

	ErrHomeDirNotFound = errors.New("home directory not found")

	ErrKubeConfigNotFound = errors.New("kubeconfig not found")

	ErrK8SCouldNotConnect = errors.New("k8s client could not connect")

	ErrInvalidEmailerEngine      = errors.New("invalid emailer engine")
	ErrEmailerConn               = errors.New("error while connecting to emailer engine")
	ErrInvalidEmailerCredentials = errors.New("invalid emailer credentials")
)
