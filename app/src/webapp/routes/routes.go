package routes

var (
	LoginRedirect    = "/login/redirect"
	Login            string
	LoginCallBack    string
	Logout           string
	OrganizationAuth = "/auth/Organization/:Organization"
	UIRoute          = "/ui"
	ManagePrefix     = "/manage"
	Register         = "/register"
)

const (
	ManageAIGroup                   = "/ai/manage"
	ManageAIModelNodes              = "/model-nodes"
	ManageAIAgents                  = "/agents"
	ManageAIAgentsResource          = "/agents/:agentId"
	ManageAIPrompts                 = "/prompts"
	ManageAIPromptResource          = "/prompts/:promptId"
	ManageAIGuardrails              = "/guardrails"
	ManageAIGuardrailResource       = "/guardrails/:guardrailId"
	ManageAIManageModelNodeResource = "/model-nodes/:aiModelNodeId"
)

const (
	StudioGroup = "/studios"
	AIStudio    = "/ai"
	AgentStudio = "/agents"
)

const (
	SettingsGroup         = "/settings"
	SettingsOrganizations = "/Organization"
	SettingToken          = "/tokens"
	SettingsStudio        = "/platform"
	SettingsIntergation   = "/integrations"
	SettingsUsers         = "/users"
	SettingsUserResource  = "/users/:userId"

	SettingsPlugins             = "/plugins"
	SettingsPluginResource      = "/plugins/:pluginId"
	SettingsStudioOrganizations = "/platform-Organizations"
)

const (
	DevelopersGroup     = "/developers"
	DevelopersResources = "/resources"
	DevelopersEnums     = "/enums"
)

const (
	HomeGroup = "/"
)

const (
	ExploreGroup            = "/explore"
	ExpOrganizations        = "/Organizations"
	ExpOrganizationResource = "/Organizations/:OrganizationId"
	ExpOrganizationUsers    = "/Organizations/:OrganizationId/users"
)

const (
	SearchGroup = "/search"
)
