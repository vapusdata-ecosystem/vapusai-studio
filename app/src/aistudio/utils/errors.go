package utils

import "errors"

// Error constants
var (
	ErrInvalidAction                   = errors.New("invalid action for this request")
	ErrDataSourceCredsNotFound         = errors.New("error while getting the credentials for the data source")
	ErrDataSourceCredsSecretGet        = errors.New("error while reading the secret from the secret store")
	ErrInvalidAIModelNodeRequestSpec   = errors.New("invalid AI model node request spec")
	ErrInvalidAIModelNodeRequested     = errors.New("invalid AI model node requested")
	ErrSavingAIModelNode               = errors.New("error while saving AI model node metadata in elasticsearch")
	ErrInvalidNetworkConfig            = errors.New("invalid network configuration")
	ErrSettingAIModelNodeCredentials   = errors.New("error while setting AI model node credentials")
	ErrCrawlingAIModels                = errors.New("error while crawling AI models")
	ErrGetCredsSecret                  = errors.New("error while getting the credentials from secret store")
	ErrAIModelConn                     = errors.New("error while creating AI model connection")
	ErrAIModelNode404                  = errors.New("AI model node not found")
	ErrAIModelManagerAction404         = errors.New("AI model manager action not found")
	ErrAIModelManagerPromptTemplate404 = errors.New("AI model manager prompt template not found")
	ErrGetAIModelNetParams             = errors.New("error while getting network params for AI model")
	ErrAIPrompt404                     = errors.New("AI prompt not found")
	ErrSavingAIModelPrompt             = errors.New("error while saving AI model prompt to datastore")
	ErrInvalidAIModelPromptRequestSpec = errors.New("invalid AI model prompt request spec")
	ErrAIPromptNotEditable             = errors.New("AI prompt is not editable")
	ErrAIModelNode403                  = errors.New("AI model node not available for org")
	ErrUpdatingAIModelNode             = errors.New("error while updating AI model node metadata in elasticsearch")
	ErrPrompt403                       = errors.New("error while validating user access")
	ErrAIAgent404                      = errors.New("AI agent not found")
	ErrAIAgentCreate500                = errors.New("error while creating AI agent")
	ErrAIAgentNotEditable              = errors.New("AI agent is not editable")
	ErrAIAgentPatch400                 = errors.New("error while patching AI agent")
	InvalidDataProductParam            = errors.New("invalid data product params")
	InvalidContentFormatParam          = errors.New("invalid content format params")
	ErrInvalidAIModels                 = errors.New("error while getting AI model list from datastore")
	ErrInvalidAgentChainParams         = errors.New("invalid agent chain parameters")
	ErrInvalidAgentType                = errors.New("invalid agent type")
	ErrFileDataNotFound                = errors.New("error while getting file data")
	ErrInvalidDataFormat               = errors.New("invalid data format")
	ErrAIGuardrailNotEditable          = errors.New("AI guardrail is not editable")
	ErrAIGuardrail404                  = errors.New("AI guardrail not found")
	ErrAIGuardrailPatch400             = errors.New("error while patching AI guardrail")
	ErrAIGuardrailCreate400            = errors.New("error while creating AI guardrail")
	ErrAIGuardrail403                  = errors.New("error while validating user access to AI guardrail")
	ErrOrganization404                 = errors.New("organization not found")
	ErrMissingUploadResourceName       = errors.New("resource is empty for upload request")
	Err401                             = errors.New("unauthorized access")
	ErrUser404                         = errors.New("user not found")
	ErrInvalidAccessTokenAgentUtility  = errors.New("invalid access token agent utility")
	ErrRefreshToken404                 = errors.New("refresh token not found")
	ErrRefreshTokenExpired             = errors.New("refresh token expired")
	ErrRefreshTokenInactive            = errors.New("refresh token inactive")
	ErrUserInviteCreateFailed          = errors.New("error while creating user invite")
	ErrUserInviteExists                = errors.New("user invite already exists")
	ErrInvalidUserManagerAction        = errors.New("invalid user manager action")
	ErrSavingOrganizationAuthJwt       = errors.New("error while saving organization auth jwt")
	ErrCreateOrganization              = errors.New("error while creating organization")
	ErrPlugin404                       = errors.New("plugin not found")
	ErrPlugin403                       = errors.New("plugin not in organization")
	ErrPluginNotEditable               = errors.New("plugin not editable")
	ErrPluginPatch400                  = errors.New("error while patching plugin")
	ErrPluginCreate400                 = errors.New("error while creating plugin")
	ErrPluginOrganizationScope403      = errors.New("plugin organization scope 403")
	ErrPluginScope403                  = errors.New("plugin scope 403")
	ErrUnSupportedPluginType           = errors.New("unsupported plugin type")
	ErrInvalidResourceScope            = errors.New("invalid resource scope")
	ErrPluginServiceScopeExists        = errors.New("plugin service with scope already exists in the organization")
	ErrFileStorePlugin400              = errors.New("error while creating file store plugin")
	ErrFileStorePlugin404              = errors.New("file store plugin not found")
	ErrInvalidPluginTypeForAction      = errors.New("invalid plugin type for action")
	ErrInvalidOrganizationRequested    = errors.New("invalid organization requested")
	ErrInvalidUserRequested            = errors.New("invalid user requested")
	ErrUserOrganizationMapping         = errors.New("user organization mapping not found")
	ErrInvalidAddOrganizationRequest   = errors.New("invalid add organization request")
	ErrInvalidOrganizationAction       = errors.New("invalid organization action")
	ErrGetAccount                      = errors.New("error while getting account")
	ErrInvalidAccountAction            = errors.New("invalid account action")
	ErrAccountOps403                   = errors.New("account operation 403")
	ErrConfigureAIStudioModel          = errors.New("error while configuring AI studio model")
	ErrUserOrganization404             = errors.New("user organization not found")
	ErrOrganizationInitialization      = errors.New("error while initializing organization")
	ErrListingAccount                  = errors.New("error while listing account")
	ErrNotServiceOrganization          = errors.New("organization is not a service organization")
	ErrLoginFailed                     = errors.New("login failed")
	ErrAuthenticatorInitFailed         = errors.New("authenticator initialization failed")
	ErrCannotCreateServiceOrg          = errors.New("cannot create service organization")
)
