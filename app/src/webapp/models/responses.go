package models

import (
	"github.com/vapusdata-oss/aistudio/webapp/routes"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

type ResourceManagerParams struct {
	API              string
	SupportedActions []string
	YamlSpec         string
	ActionMap        map[string]string
	StreamAPI        string
}

type GlobalContexts struct {
	NavMenuMap          []routes.RouteBaseStruct
	SidebarMap          []routes.RouteBaseStruct
	CurrentNav          string
	CurrentSideBar      string
	UserInfo            *mpb.User
	Account             *mpb.Account
	CurrentOrganization *mpb.Organization
	LoginUrl            string
	AccessTokenKey      string
	Manager             bool
	OrganizationMap     map[string]string
	CurrentUrl          string
}

type SettingsResponse struct {
	ActionParams        *ResourceManagerParams
	CurrentOrganization *mpb.Organization
	Users               []*mpb.User
	User                *mpb.User
	BackListingLink     string
	SpecMap             map[mpb.RequestObjects]string
	ResourceActionsMap  map[mpb.RequestObjects][]string
	Enums               map[string]map[string]int32
	Plugins             []*mpb.Plugin
	Plugin              *mpb.Plugin
	Organizations       []*mpb.Organization
	YourOrganizations   []string
}

type SearchResponse struct {
	ActionParams        *ResourceManagerParams
	CurrentOrganization *mpb.Organization
}

type OrganizationSvcResponse struct {
	CurrentOrganization *mpb.Organization
	Users               []*mpb.User
	User                *mpb.User
	BackListingLink     string
	ActionParams        *ResourceManagerParams
}

type AIStudioResponse struct {
	AIModelNodes    []*mpb.AIModelNode
	AIModelNode     *mpb.AIModelNode
	ActionParams    *ResourceManagerParams
	BackListingLink string
	AIPrompts       []*mpb.AIModelPrompt
	AIPrompt        *mpb.AIModelPrompt
	AIAgents        []*mpb.VapusAIAgent
	AIAgent         *mpb.VapusAIAgent
	AIGuardrail     *mpb.AIGuardrails
	AIGuardrails    []*mpb.AIGuardrails
}

type AgentStudioResponse struct {
	AIModelNodes    []*mpb.AIModelNode
	AIModelNode     *mpb.AIModelNode
	ActionParams    *ResourceManagerParams
	BackListingLink string
	AIPrompts       []*mpb.AIModelPrompt
	AIPrompt        *mpb.AIModelPrompt
	AIAgents        []*mpb.VapusAIAgent
	AIAgent         *mpb.VapusAIAgent
}
