package router

import (
	"github.com/labstack/echo/v4"
	"github.com/vapusdata-oss/aistudio/webapp/routes"
	services "github.com/vapusdata-oss/aistudio/webapp/services"
)

func settingsRouter(e *echo.Group) {
	settingsGroup := e.Group(routes.SettingsGroup)
	settingsGroup.GET("", services.WebappServiceManager.SettingsProfile)
	settingsGroup.GET(routes.SettingsStudio, services.WebappServiceManager.SettingsVapusStudio)
	settingsGroup.GET(routes.SettingsUsers, services.WebappServiceManager.OrganizationUsersList)
	settingsGroup.GET(routes.SettingsUserResource, services.WebappServiceManager.UserDetails)

	settingsGroup.GET(routes.SettingsPlugins, services.WebappServiceManager.ManagePluginsHandler)
	settingsGroup.GET(routes.SettingsPluginResource, services.WebappServiceManager.ManagePluginDetailHandler)
	settingsGroup.GET(routes.SettingsOrganizations, services.WebappServiceManager.OrganizationSettings)
	settingsGroup.GET(routes.SettingsStudioOrganizations, services.WebappServiceManager.StudioOrganizationsList)
}
