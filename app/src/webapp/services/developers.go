package services

import (
	"net/http"

	"github.com/labstack/echo/v4"
	serviceops "github.com/vapusdata-oss/aistudio/core/serviceops"
	"github.com/vapusdata-oss/aistudio/webapp/models"
	"github.com/vapusdata-oss/aistudio/webapp/routes"
)

func (x *WebappService) SettingsResources(c echo.Context) error {
	response := models.SettingsResponse{
		SpecMap:            serviceops.SpecMap,
		ResourceActionsMap: serviceops.ResourceActionsMap,
	}
	globalContext, err := x.getDeveloperSectionGlobals(c, routes.DevelopersResourcesPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting settings section globals")
		return HandleGLobalContextError(c, err)
	}
	response.CurrentOrganization = globalContext.CurrentOrganization
	return c.Render(http.StatusOK, "settings-resources.html", map[string]interface{}{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Resources & Specs",
	})
}

func (x *WebappService) SettingsEnums(c echo.Context) error {
	// enums := maps.Keys(vapussvc.EnumSpecs)
	response := models.SettingsResponse{
		Enums: serviceops.EnumSpecs,
	}
	globalContext, err := x.getDeveloperSectionGlobals(c, routes.DevelopersEnumsPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting settings section globals")
		return HandleGLobalContextError(c, err)
	}
	response.CurrentOrganization = globalContext.CurrentOrganization
	return c.Render(http.StatusOK, "settings-enums.html", map[string]interface{}{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Enums",
	})
}
