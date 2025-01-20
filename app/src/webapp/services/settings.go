package services

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/vapusdata-oss/aistudio/core/globals"
	serviceops "github.com/vapusdata-oss/aistudio/core/serviceops"
	"github.com/vapusdata-oss/aistudio/core/serviceops/grpcops"
	"github.com/vapusdata-oss/aistudio/webapp/models"
	pkgs "github.com/vapusdata-oss/aistudio/webapp/pkgs"
	"github.com/vapusdata-oss/aistudio/webapp/routes"
	"github.com/vapusdata-oss/aistudio/webapp/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	dpb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

func (x *WebappService) SettingsProfile(c echo.Context) error {
	var updateSpec string
	response := models.SettingsResponse{
		ActionParams: &models.ResourceManagerParams{},
	}
	globalContext, err := x.getSettingsSectionGlobals(c, routes.SettingsProfilePage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting home section globals")
		return HandleGLobalContextError(c, err)
	}
	if globalContext.UserInfo != nil {
		obj := &dpb.UserManagerRequest{
			Spec:   globalContext.UserInfo,
			Action: dpb.UserAgentActions_PATCH_USER,
		}
		bytess, err := grpcops.ProtoYamlMarshal(obj)
		if err != nil {
			Logger.Err(err).Msg("error while marshaling user spec")
			return x.HandleError(c, err, http.StatusBadRequest, globalContext)
		}
		updateSpec = string(bytess)
	}

	response.ActionParams.API = fmt.Sprintf("%s/api/v1alpha1/platform/users", pkgs.NetworkConfigManager.ExternalURL)
	response.ActionParams.ActionMap = map[string]string{
		utils.UPDATE: updateSpec,
	}
	return c.Render(http.StatusOK, "settings-profile.html", map[string]interface{}{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Your Profile",
	})
}

func (x *WebappService) SettingsVapusStudio(c echo.Context) error {
	var updateSpec string
	response := models.SettingsResponse{
		ActionParams: &models.ResourceManagerParams{},
	}
	globalContext, err := x.getSettingsSectionGlobals(c, routes.SettingsStudioPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting home section globals")
		return HandleGLobalContextError(c, err)
	}
	if globalContext.Account != nil {
		obj := &dpb.AccountManagerRequest{
			Spec:    globalContext.Account,
			Actions: dpb.AccountAgentActions_UPDATE_PROFILE,
		}
		bytess, err := grpcops.ProtoYamlMarshal(obj)
		if err != nil {
			Logger.Err(err).Msg("error while marshaling platform account spec")
			return x.HandleError(c, err, http.StatusBadRequest, globalContext)
		}
		updateSpec = string(bytess)
	}
	response.CurrentOrganization = globalContext.CurrentOrganization
	response.ActionParams.API = fmt.Sprintf("%s/api/v1alpha1/platform", pkgs.NetworkConfigManager.ExternalURL)
	response.ActionParams.ActionMap = map[string]string{
		utils.UPDATE: updateSpec,
	}
	return c.Render(http.StatusOK, "settings-platform.html", map[string]interface{}{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Studio Settings",
	})
}

func (x *WebappService) OrganizationSettings(c echo.Context) error {
	var updateSpec string
	var addUsersSpec string
	var err error
	response := models.OrganizationSvcResponse{
		ActionParams: &models.ResourceManagerParams{},
	}
	globalContext, err := x.getSettingsSectionGlobals(c, routes.OrganizationSettingsPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting home section globals")
		return HandleGLobalContextError(c, err)
		// return HandleGLobalContextError(c, err)
	}
	response.CurrentOrganization = globalContext.CurrentOrganization
	if response.CurrentOrganization != nil {
		obj := &dpb.OrganizationManagerRequest{
			Spec:    response.CurrentOrganization,
			Actions: dpb.OrganizationAgentActions_PATCH_ORG,
		}
		bytess, err := grpcops.ProtoYamlMarshal(obj)
		if err != nil {
			Logger.Err(err).Msg("error while marshaling Organization spec")
			return x.HandleError(c, err, http.StatusNotFound, globalContext)
		}
		updateSpec = string(bytess)
		addUserObj := &dpb.OrganizationManagerRequest{
			Users: []*mpb.OrganizationUserOps{{
				UserId:           "",
				Role:             []string{},
				ValidTill:        0,
				InviteIfNotFound: true,
			}},
			Actions: dpb.OrganizationAgentActions_ADD_ORG_USER,
		}
		uBytess, err := grpcops.ProtoYamlMarshal(addUserObj)
		if err != nil {
			Logger.Err(err).Msg("error while marshaling Organization add user spec")
			return x.HandleError(c, err, http.StatusNotFound, globalContext)
		}
		addUsersSpec = string(uBytess)
	} else {
		Logger.Err(err).Msg("error while getting Organization")
		return x.HandleError(c, err, http.StatusNotFound, globalContext)
	}

	response.ActionParams.API = fmt.Sprintf("%s/api/v1alpha1/platform/Organizations", pkgs.NetworkConfigManager.ExternalURL)
	response.ActionParams.ActionMap = map[string]string{
		utils.UPDATE:    updateSpec,
		utils.ADD_USERS: addUsersSpec,
	}
	return c.Render(http.StatusOK, "settings-Organization.html", map[string]interface{}{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": response.CurrentOrganization.Name,
		"ResourceBase":  response.CurrentOrganization.ResourceBase,
	})
}

func (x *WebappService) StudioUsersList(c echo.Context) error {
	response := models.SettingsResponse{
		Users:           x.grpcClients.GetStudioUsers(c),
		ActionParams:    &models.ResourceManagerParams{},
		BackListingLink: c.Request().URL.String(),
	}
	globalContext, err := x.getSettingsSectionGlobals(c, routes.SettingsUsersPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting settings section globals")
		return HandleGLobalContextError(c, err)
	}
	response.CurrentOrganization = globalContext.CurrentOrganization
	return c.Render(http.StatusOK, "settings-users.html", map[string]interface{}{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Studio Users",
	})
}

func (x *WebappService) OrganizationUsersList(c echo.Context) error {
	globalContext, err := x.getSettingsSectionGlobals(c, routes.SettingsUsersPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting home section globals")
		return HandleGLobalContextError(c, err)
	}
	var addUsersSpec string
	response := models.OrganizationSvcResponse{
		BackListingLink: c.Request().URL.String(),
		ActionParams:    &models.ResourceManagerParams{},
	}
	if globalContext.CurrentOrganization.OrgType == mpb.OrganizationTypes_SERVICE_ORG {
		response.Users = x.grpcClients.GetStudioUsers(c)
	} else {
		response.Users = x.grpcClients.GetMyOrganizationUsers(c)
	}

	response.CurrentOrganization = globalContext.CurrentOrganization
	response.CurrentOrganization = globalContext.CurrentOrganization
	if response.CurrentOrganization != nil {
		addUObj := &dpb.OrganizationManagerRequest{
			Actions: dpb.OrganizationAgentActions_PATCH_ORG,
			Users: []*mpb.OrganizationUserOps{{
				UserId:           "",
				Role:             []string{},
				ValidTill:        0,
				InviteIfNotFound: false,
			}},
		}
		bytess, err := grpcops.ProtoYamlMarshal(addUObj)
		if err != nil {
			Logger.Err(err).Msg("error while marshaling Organization spec for adding users")
			return x.HandleError(c, err, http.StatusNotFound, globalContext)
		}
		addUsersSpec = string(bytess)
	} else {
		Logger.Err(err).Msg("error while getting Organization")
		return x.HandleError(c, err, http.StatusNotFound, globalContext)
	}
	response.ActionParams.API = fmt.Sprintf("%s/api/v1alpha1/platform/Organizations", pkgs.NetworkConfigManager.ExternalURL)
	response.ActionParams.ActionMap = map[string]string{
		utils.ADD_USERS: addUsersSpec,
	}
	return c.Render(http.StatusOK, "settings-Organization-users.html", map[string]interface{}{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Organization Users",
	})
}

func (x *WebappService) UserDetails(c echo.Context) error {
	globalContext, err := x.getSettingsSectionGlobals(c, routes.SettingsUsersPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting settings section globals")
		return HandleGLobalContextError(c, err)
	}
	userId, err := GetUrlParams(c, globals.UserId)
	user, _, err := x.grpcClients.GetUserInfo(c, userId)
	if err != nil || user == nil {
		Logger.Err(err).Msg("error while getting user info")
		return x.HandleError(c, err, http.StatusNotFound, globalContext)
	}
	response := models.SettingsResponse{
		User:         user,
		ActionParams: &models.ResourceManagerParams{},
	}

	cRoles := GetUserCurrentOrganizationRole(globalContext.UserInfo, globalContext.CurrentOrganization.OrgId)
	if slices.Contains(cRoles, mpb.StudioUserRoles_ORG_OWNERS.String()) {
		var updateSpec = ""
		bbtes, err := grpcops.ProtoYamlMarshal(&dpb.UserManagerRequest{
			Spec:   user,
			Action: dpb.UserAgentActions_PATCH_USER,
		})
		if err != nil {
			Logger.Err(err).Msg("error while marshaling user spec")
			return x.HandleError(c, err, http.StatusBadRequest, globalContext)
		}
		updateSpec = string(bbtes)
		response.ActionParams.API = fmt.Sprintf("%s/api/v1alpha1/platform/users", pkgs.NetworkConfigManager.ExternalURL)
		response.ActionParams.ActionMap = map[string]string{
			utils.UPDATE: updateSpec,
		}
	}
	response.CurrentOrganization = globalContext.CurrentOrganization
	return c.Render(http.StatusOK, "settings-user-details.html", map[string]interface{}{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": user.UserId,
	})
}
func (x *WebappService) ManagePluginsHandler(c echo.Context) error {
	response := &models.SettingsResponse{
		Plugins:         x.grpcClients.ListPlugins(c),
		ActionParams:    &models.ResourceManagerParams{},
		BackListingLink: c.Request().URL.String(),
	}
	globalContext, err := x.getSettingsSectionGlobals(c, routes.SettingsPluginsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting plugin section globals")
		return HandleGLobalContextError(c, err)
	}

	Specs := make(map[string]string)
	for _, plugin := range response.Plugins {
		plugin.NetworkParams.Credentials = &mpb.GenericCredentialObj{
			AwsCreds:   &mpb.AWSCreds{},
			GcpCreds:   &mpb.GCPCreds{},
			AzureCreds: &mpb.AzureCreds{},
		}
		obj := &dpb.PluginManagerRequest{
			Spec:   plugin,
			Action: dpb.PluginAgentAction_PATCH_PLUGIN,
		}
		bytess, err := grpcops.ProtoYamlMarshal(obj)
		if err != nil {
			x.logger.Err(err).Msg("error while marshaling plugin spec")
			Specs[plugin.PluginId] = ""
		} else {
			Specs[plugin.PluginId] = string(bytess)
		}
	}

	response.ActionParams.API = fmt.Sprintf("%s/api/v1alpha1/plugins", pkgs.NetworkConfigManager.ExternalURL)
	return c.Render(http.StatusOK, "settings-plugins.html", map[string]interface{}{
		"GlobalContext":  globalContext,
		"Response":       response,
		"SectionHeader":  "Plugins",
		"Specs":          Specs,
		"CreateTemplate": GetProtoYamlString(serviceops.PluginManagerRequest),
	})
}

func (x *WebappService) ManagePluginDetailHandler(c echo.Context) error {
	globalContext, err := x.getSettingsSectionGlobals(c, routes.SettingsPluginsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting plugin section globals")
		return HandleGLobalContextError(c, err)
	}
	var updateSpec, yamlSpec string
	PluginId, err := GetUrlParams(c, globals.PluginId)
	if err != nil {
		return c.Render(http.StatusBadRequest, "400.html", map[string]interface{}{
			"GlobalContext": globalContext,
		})
	}
	response := &models.SettingsResponse{
		Plugin:       x.grpcClients.GetPlugin(c, PluginId),
		ActionParams: &models.ResourceManagerParams{},
	}
	if response.Plugin != nil {
		upObj := &dpb.PluginManagerRequest{
			Spec:   response.Plugin,
			Action: dpb.PluginAgentAction_PATCH_PLUGIN,
		}
		bytess, err := grpcops.ProtoYamlMarshal(upObj)
		if err != nil {
			x.logger.Err(err).Msg("error while marshaling plugin spec")
			return x.HandleError(c, err, http.StatusBadRequest, globalContext)
		}
		updateSpec = string(bytess)
		yBytes, err := grpcops.ProtoYamlMarshal(response.Plugin)
		if err != nil {
			x.logger.Err(err).Msg("error while marshaling plugin spec")
		}
		yamlSpec = string(yBytes)
	} else {
		x.logger.Err(err).Msg("error while getting plugin details")
		return x.HandleError(c, err, http.StatusNotFound, globalContext)
	}
	response.ActionParams.API = fmt.Sprintf("%s/api/v1alpha1/plugins", pkgs.NetworkConfigManager.ExternalURL)
	if response.Plugin.Editable && response.Plugin.Org == globalContext.CurrentOrganization.OrgId {
		response.ActionParams.ActionMap = map[string]string{
			utils.UPDATE: updateSpec,
		}
	}
	return c.Render(http.StatusOK, "settings-plugin-detail.html", map[string]interface{}{
		"GlobalContext": globalContext,
		"Response":      response,
		"YamlSpec":      yamlSpec,
		"SectionHeader": response.Plugin.PluginId,
	})
}

func (x *WebappService) StudioOrganizationsList(c echo.Context) error {
	response := models.SettingsResponse{
		Organizations: x.grpcClients.GetOrganizations(c),
		ActionParams:  &models.ResourceManagerParams{},
	}
	globalContext, err := x.getSettingsSectionGlobals(c, routes.SettingsStudioOrganizationsPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting explore page global context")
		return HandleGLobalContextError(c, err)
	}
	for _, Organization := range globalContext.UserInfo.OrganizationRoles {
		response.YourOrganizations = append(response.YourOrganizations, Organization.OrganizationId)
	}
	// x.LoadOrganizationMap(c)
	response.ActionParams.API = fmt.Sprintf("%s/api/v1alpha1/platform/Organizations", pkgs.NetworkConfigManager.ExternalURL)
	return c.Render(http.StatusOK, "settings-platform-Organizations.html", map[string]interface{}{
		"GlobalContext":  globalContext,
		"Response":       response,
		"SectionHeader":  "Organizations",
		"CreateTemplate": GetProtoYamlString(serviceops.OrganizationManagerRequest),
	})
}
