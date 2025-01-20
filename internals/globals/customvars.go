package globals

import "time"

var (
	DataProductDescriptionVectorDimensions int = 1536
	DataSourceMetadataVectorDimensions     int = 1536
	DEFAULT_ACCESS_TOKEN_VALIDITY          int = 3600
	IdSeparator                                = "_"
	BASE_TABLE_TYPE                            = "BASE TABLE"
	VIEW_TYPE                                  = "VIEW"
	DEFAULT_PAGE_SIZE                          = 10
	DEFAULT_AT_VALIDITY                        = time.Duration(3600)
)

type SearchKey string

const (
	DataProuductSK    SearchKey = "dataProduct"
	DataWorkerSK      SearchKey = "dataWorker"
	DataSourceSK      SearchKey = "dataSource"
	OrganizationSK    SearchKey = "organization"
	DataCatalogSK     SearchKey = "dataCatalog"
	UserSK            SearchKey = "user"
	VdcDeploymentSK   SearchKey = "vdcDeploymentId"
	AIModelNodeSK     SearchKey = "aiModelNode"
	WorkerDeplomentSK SearchKey = "workerDeployment"
	DataStore         SearchKey = "datastore"
)

func (x SearchKey) String() string {
	return string(x)
}

type AgentType string

func (at AgentType) String() string {
	return string(at)
}

const (
	NABHIK_AGENT                 AgentType = "nabhik"
	VAPUSDATACONTAINERAGENT      AgentType = "vapusDataContainerAgent"
	DATAWORKERAGENT              AgentType = "dataWorkerAgent"
	VDCDEPLOYMENTAGENT           AgentType = "vdcDeploymentAgent"
	DATASOURCEAGENT              AgentType = "dataSourceAgent"
	DATAPRODUCTAGENT             AgentType = "dataProductAgent"
	DATAMARKETPLACEAGENT         AgentType = "dataMarketplaceAgent"
	AISTUDIONODE                 AgentType = "aiStudioNode"
	ACCOUNTAGENT                 AgentType = "accountAgent"
	DATASOURCEMETADATAAGENT      AgentType = "dataSourceMetadataAgent"
	VAPUSDATAPLATFORMSEARCHAGENT AgentType = "vapusDataStudioSearchAgent"
	DOMAINAGENT                  AgentType = "organizationAgent"
	USERMANAGERAGENT             AgentType = "userManagerAgent"
	AUTHZMANAGERAGENT            AgentType = "authzManagerAgent"
	AIPROMPTAGENT                AgentType = "aiPromptAgent"
	VAPUSAIAGENT                 AgentType = "vapusAIAgent"
	FABRICGETTERAGENT            AgentType = "fabricGetterAgent"
)

const (
	SENDEMAIL    = "sendEmail"
	SECRETENGINE = "secretEngine"
)

// REDIS action and keys that are constant
const (
	LIST              = "list"
	ADD               = "Add"
	EXISTS            = "exists"
	COUNT             = "count"
	DEL               = "del"
	MADD              = "add-mulitple"
	ACCOUNT_KEY       = "account"
	ACCOUNT_DM_CF_KEY = "account:datamarketplace"
	DM_IDENTIFIER     = "datamarketplace"
	MARKETPLACEID     = "marketplaceId"
	DMNID             = "dmnId"
	DNID              = "dnId"
)

var FileContentTypes = map[string]string{
	"csv":   "text/csv",
	"json":  "application/json",
	"xml":   "application/xml",
	"yaml":  "application/yaml",
	"yml":   "application/yaml",
	"txt":   "text/plain",
	"tsv":   "text/tab-separated-values",
	"tab":   "text/tab-separated-values",
	"xls":   "application/vnd.ms-excel",
	"xlsx":  "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"doc":   "application/msword",
	"docx":  "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"pdf":   "application/pdf",
	"ppt":   "application/vnd.ms-powerpoint",
	"pptx":  "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"zip":   "application/zip",
	"tar":   "application/x-tar",
	"gz":    "application/gzip",
	"bz2":   "application/x-bzip2",
	"rar":   "application/x-rar-compressed",
	".gif":  "image/gif",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".svg":  "image/svg+xml",
	".bmp":  "image/bmp",
}

const (
	MarketplaceId      = "marketplaceId"
	OrganizationId     = "organizationId"
	DataProductId      = "dataproductId"
	CatalogId          = "catalogId"
	PromptId           = "promptId"
	WorkerDeploymentId = "workerDeploymentId"
	VdcDeploymentId    = "vdcDeploymentId"
	UserId             = "userId"
	DataWorkerId       = "dataworkerId"
	DataSourceId       = "dataSourceId"
	AIModelNodeId      = "aiModelNodeId"
	AgentId            = "agentId"
	PluginId           = "pluginId"
	GuardrailId        = "guardrailId"
	FabricChatId       = "fabricChatId"
)

var UrlResouceIdMap = map[string]bool{
	MarketplaceId:      true,
	OrganizationId:     true,
	DataProductId:      true,
	CatalogId:          true,
	PromptId:           true,
	WorkerDeploymentId: true,
	VdcDeploymentId:    true,
	UserId:             true,
	DataWorkerId:       true,
	DataSourceId:       true,
	AIModelNodeId:      true,
	AgentId:            true,
	PluginId:           true,
	GuardrailId:        true,
	FabricChatId:       true,
}

const (
	FabricChatFileKey = "fabricChat"
)
