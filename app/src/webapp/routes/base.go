package routes

type NavMenuId string

type RouteBaseStruct struct {
	ItemName string
	Url      string
	ItemId   string
	Svg      string
	Children []RouteBaseStruct
}

var (
	HomeNav        NavMenuId = "vapushomeNavMenu"
	SettingsNav    NavMenuId = "settingsNavMenu"
	ManageAINav    NavMenuId = "manageAINavMenu"
	VapusStudioNav NavMenuId = "vapusStudioNavMenu"
	DevelopersNav  NavMenuId = "developersNavMenu"
)

func (p NavMenuId) String() string {
	return string(p)
}

type SidebarId string

var (
	OrganizationSettingsPage        SidebarId = "OrganizationSettings"
	OrganizationDataSourcesPage     SidebarId = "OrganizationDataSources"
	OrganizationsPage               SidebarId = "Organizations"
	OrganizationUsersPage           SidebarId = "OrganizationUsers"
	OrganizationOverviewPage        SidebarId = "OrganizationOverview"
	ManageAIOverviewPage            SidebarId = "aiStudioOverview"
	ManageAIModelNodesPage          SidebarId = "aiStudioModelNodes"
	AIPromptsPage                   SidebarId = "aiPrompts"
	AIModelDeploymentPage           SidebarId = "aiModelDeployments"
	AITrainerPage                   SidebarId = "aiTrainers"
	ManageAIAgentsPage              SidebarId = "aiStudioAgents"
	ManageAIModelInterfacePage      SidebarId = "aiStudioModelInterface"
	SettingsProfilePage             SidebarId = "settingsProfile"
	SettingsStudioPage              SidebarId = "settingsStudio"
	SettingsIntegrationsPage        SidebarId = "settingsIntegrations"
	SettingsTokensPage              SidebarId = "settingsTokens"
	SettingsUsersPage               SidebarId = "settingsUsers"
	AccessRequestsPage              SidebarId = "accessRequests"
	LogoutPage                      SidebarId = "logout"
	DevelopersResourcesPage         SidebarId = "settingsResources"
	DevelopersEnumsPage             SidebarId = "settingsEnums"
	SettingsPluginsPage             SidebarId = "settingsPlugins"
	AIStudioPage                    SidebarId = "aiStudio"
	AgentStudioPage                 SidebarId = "agentStudio"
	ManageAIGuardrailsPage          SidebarId = "aiGuardrails"
	SettingsStudioOrganizationsPage SidebarId = "settingsStudioOrganizations"
)

func (p SidebarId) String() string {
	return string(p)
}

var (
	NavMenuRoutes                []RouteBaseStruct
	DataMarketplaceSidebarRoutes []RouteBaseStruct
	OrganizationSidebarRoutes    []RouteBaseStruct
	ManageAISidebarRoutes        []RouteBaseStruct
)

var NavMenuList = []RouteBaseStruct{
	{
		ItemName: "Dashboard",
		ItemId:   HomeNav.String(),
		Url:      UIRoute + HomeGroup,
		Svg: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="h-6 w-6 m-2" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M3 9l9-6 9 6v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
  <path d="M9 22V12h6v10"></path>
</svg>
`,
	},
	{
		ItemName: "Playground",
		ItemId:   VapusStudioNav.String(),
		Url:      UIRoute + StudioGroup + AIStudio,
		Svg: `<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 m-2"  viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <rect x="7" y="9" width="10" height="10" rx="2"></rect>
  <path d="M9 3h6v3H9z"></path>
  <path d="M5 10h1v4H5z"></path>
  <path d="M18 10h1v4h-1z"></path>
  <circle cx="9" cy="13" r="0.5"></circle>
  <circle cx="15" cy="13" r="0.5"></circle>
</svg>
`,
		Children: StudioNavList,
	},

	{
		ItemName: "AI Center",
		ItemId:   ManageAINav.String(),
		Url:      UIRoute + ManageAIGroup + ManageAIModelNodes,
		Svg: `<div>
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="h-8 w-8 m-2" fill="currentColor">
        <circle cx="12" cy="12" r="10" fill="none" />
        <path d="M15.5 8.5h-7c-.55 0-1 .45-1 1v5c0 .55.45 1 1 1h7c.55 0 1-.45 1-1v-5c0-.55-.45-1-1-1zm-7 6v-5h7v5h-7z" fill="#ffffff"/>
        <path d="M12 14.5c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zm0-2.5a.5.5 0 110 1 .5.5 0 010-1z" fill="#ffffff"/>
        <rect x="11.25" y="5.5" width="1.5" height="2.5" fill="#ffffff"/>
        <rect x="11.25" y="16" width="1.5" height="2.5" fill="#ffffff"/>
        <rect x="16" y="11.25" width="2.5" height="1.5" fill="#ffffff"/>
        <rect x="5.5" y="11.25" width="2.5" height="1.5" fill="#ffffff"/>
    </svg>
</div>
`,
		Children: ManageAINavSideList,
	},
	{
		ItemName: "Settings",
		ItemId:   SettingsNav.String(),
		Url:      UIRoute + SettingsGroup,
		Svg: `<svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-settings h-6 w-6 m-2" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
  <path stroke="none" d="M0 0h24v24H0z" fill="none"/>
  <path d="M10.325 4.317a1.724 1.724 0 0 1 3.35 0l.333 1.35a7.04 7.04 0 0 1 1.817 .621l1.308 -.478a1.724 1.724 0 0 1 2.156 2.156l-.478 1.308a7.034 7.034 0 0 1 .621 1.817l1.35 .333a1.724 1.724 0 0 1 0 3.35l-1.35 .333a7.034 7.034 0 0 1 -.621 1.817l.478 1.308a1.724 1.724 0 0 1 -2.156 2.156l-1.308 -.478a7.04 7.04 0 0 1 -1.817 .621l-.333 1.35a1.724 1.724 0 0 1 -3.35 0l-.333 -1.35a7.04 7.04 0 0 1 -1.817 -.621l-1.308 .478a1.724 1.724 0 0 1 -2.156 -2.156l.478 -1.308a7.034 7.034 0 0 1 -.621 -1.817l-1.35 -.333a1.724 1.724 0 0 1 0 -3.35l1.35 -.333a7.034 7.034 0 0 1 .621 -1.817l-.478 -1.308a1.724 1.724 0 0 1 2.156 -2.156l1.308 .478a7.04 7.04 0 0 1 1.817 -.621z" />
  <circle cx="12" cy="12" r="3" />
</svg>
`,
		Children: SettingsSideList,
	},
	{
		ItemName: "Developer Resources",
		ItemId:   DevelopersNav.String(),
		Url:      UIRoute + DevelopersGroup + DevelopersResources,
		Svg: `<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 m-2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <polyline points="16 18 22 12 16 6"></polyline>
  <polyline points="8 6 2 12 8 18"></polyline>
</svg>
`,
		Children: DevelopersSideList,
	},
}

var ManageAINavSideList = []RouteBaseStruct{
	{ItemName: "Models Registry", ItemId: ManageAIModelNodesPage.String(), Url: UIRoute + ManageAIGroup + ManageAIModelNodes},
	{ItemName: "Prompts", ItemId: AIPromptsPage.String(), Url: UIRoute + ManageAIGroup + ManageAIPrompts},
	{ItemName: "Agents", ItemId: ManageAIAgentsPage.String(), Url: UIRoute + ManageAIGroup + ManageAIAgents},
	{ItemName: "Guardrails", ItemId: ManageAIGuardrailsPage.String(), Url: UIRoute + ManageAIGroup + ManageAIGuardrails},
}

var StudioNavList = []RouteBaseStruct{
	{ItemName: "AI Studio", ItemId: AIStudioPage.String(), Url: UIRoute + StudioGroup + AIStudio},
	{ItemName: "Agent Studio", ItemId: AgentStudioPage.String(), Url: UIRoute + StudioGroup + AgentStudio},
}

var SettingsSideList = []RouteBaseStruct{
	{ItemName: "Profile", ItemId: SettingsProfilePage.String(), Url: UIRoute + SettingsGroup},
	{ItemName: "Organization", ItemId: OrganizationSettingsPage.String(), Url: UIRoute + SettingsGroup + SettingsOrganizations},
	{ItemName: "Studio", ItemId: SettingsStudioPage.String(), Url: UIRoute + SettingsGroup + SettingsStudio},
	// {ItemName: "Integrations", ItemId: SettingsIntergationsPage.String(), Url: UIRoute + SettingsGroup + SettingsIntergation},
	{ItemName: "Users", ItemId: SettingsUsersPage.String(), Url: UIRoute + SettingsGroup + SettingsUsers},
	{ItemName: "Plugins", ItemId: SettingsPluginsPage.String(), Url: UIRoute + SettingsGroup + SettingsPlugins},
	{ItemName: "Studio Organizations", ItemId: SettingsStudioOrganizationsPage.String(), Url: UIRoute + SettingsGroup + SettingsStudioOrganizations},
}

var DevelopersSideList = []RouteBaseStruct{
	{ItemName: "Resources", ItemId: DevelopersResourcesPage.String(), Url: UIRoute + DevelopersGroup + DevelopersResources},
	{ItemName: "Enums", ItemId: DevelopersEnumsPage.String(), Url: UIRoute + DevelopersGroup + DevelopersEnums},
	// {ItemName: "Tokens", ItemId: SettingTokenPage.String(), Url: UIRoute + SettingsGroup + SettingToken},
}

var (
	UIHome = UIRoute + HomeGroup

	DataServerHome = UIRoute + StudioGroup + AIStudio

	ManageAIHome = UIRoute + ManageAIGroup + ManageAIModelNodes

	AIStudioHome = UIRoute + StudioGroup + AIStudio

	AgentStudioHome = UIRoute + StudioGroup + AgentStudio
)

// var NavMenuMap = map[string]map[string]string{
// 	HomeNav.String():        {HomeGroup: "Home"},
// 	SettingsNav.String():    {SettingsGroup: "Settings"},
// 	MyOrganizationNav.String():    {MyOrganizationGroup: "My Organizations"},
// 	DataMarketplaceNav.String():    {DataMarketplaceGroup: "Data Marketplace"},
// 	ManageAINav.String():    {ManageAIGroup: "AI Studio"},
// 	VapusInterfaceNav.String(): {VapusInterfaceGroup: "Query Ground"},
// }

// var NavMenuList = []string{
// 	HomeNav.String(),
// 	DataMarketplaceNav.String(),
// 	MyOrganizationNav.String(),
// 	ManageAINav.String(),
// 	VapusInterfaceNav.String(),
// 	SettingsNav.String(),
// }

// var DatamarketplaceSideMap = map[string]map[string]string{
// 	DataMarketplaceOverviewPage.String():      {DataMarketplaceGroup: "Data Marketplace Overview"},
// 	DataProductsPage.String():          {DataMarketplaceDataProduct: "Data Products"},
// 	MyDataProductsPage.String():        {MyDataProducts: "My Data Products"},
// 	DataCatalogPage.String():           {DataCatalogs: "Data Catalogs"},
// 	DataProductDiscoverPage.String():   {DataProductDiscover: "Discover"},
// 	RequestNewDataProductPage.String(): {RequestNewDataProduct: "Request New Data Product"},
// 	// UsersPage.String():                 map[string]string{:"Users"},
// 	DataMarketplacePage.String(): {DataMarketplaceDataProductResource: "Data Marketplace"},
// 	OrganizationsPage.String():  {Organizations: "Organizations"},
// }

// var OrganizationSideMap = map[string]map[string]string{
// 	OrganizationOverviewPage.String():     {MyOrganizationGroup: "Organization Overview"},
// 	DataProductsPage.String():       {DataProducts: "Data Products"},
// 	SettingsPage.String():           {OrganizationSettings: "Settings"},
// 	DataSourcesPage.String():        {DataSources: "Data Sources"},
// 	DataWorkersPage.String():        {DataWorkers: "Data Workers"},
// 	VdcDeploymentsPage.String():     {VdcDeployments: "VDC Deployments"},
// 	UsersPage.String():              {OrganizationUsers: "Users"},
// 	WorkersDeploymentsPage.String(): {WorkersDeployments: "Workers Deployments"},
// }

// var ManageAINavSideMap = map[string]map[string]string{
// 	ManageAIOverviewPage.String():   {ManageAIGroup: "AI Studio Overview"},
// 	ManageAIModelNodesPage.String(): {ManageAIModelNodes: "AI Studio Model Nodes"},
// 	AIPromptsPage.String():          {ManageAIPrompts: "AI Prompts"},
// 	ManageAIAskPage.String():        {ManageAIAsk: "AI Studio Ask"},
// }
