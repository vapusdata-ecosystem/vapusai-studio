package services

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vapusdata-oss/aistudio/webapp/models"
	pkgs "github.com/vapusdata-oss/aistudio/webapp/pkgs"
	"github.com/vapusdata-oss/aistudio/webapp/routes"
)

var Logger = pkgs.DmLogger

func (x *WebappService) HomePageHandler(c echo.Context) error {
	response := x.grpcClients.GetDashboard(c)
	globalContext, err := x.getHomeSectionGlobals(c)
	if err != nil {
		Logger.Err(err).Msg("error while getting home section globals")
		return c.Render(http.StatusBadRequest, "403.html", map[string]interface{}{
			"GlobalContext": &models.GlobalContexts{
				LoginUrl: routes.Login,
			},
		})
	}
	return c.Render(http.StatusOK, "home.html", map[string]interface{}{
		"GlobalContext": globalContext,
		"Response":      response,
		"Section":       "Dashboard",
		"SectionHeader": "Dashboard",
	})
}
