package router

import (
	"github.com/labstack/echo/v4"
	"github.com/vapusdata-oss/aistudio/webapp/routes"
	services "github.com/vapusdata-oss/aistudio/webapp/services"
)

func homeRouter(e *echo.Group) {
	e.GET("/", services.WebappServiceManager.HomePageHandler)
}

func CommonRouter(e *echo.Group) {
	e.GET(routes.OrganizationAuth, services.WebappServiceManager.AuthOrganizationHandler)
}
