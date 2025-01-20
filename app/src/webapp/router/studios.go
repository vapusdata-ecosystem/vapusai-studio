package router

import (
	"github.com/labstack/echo/v4"
	"github.com/vapusdata-oss/aistudio/webapp/routes"
	services "github.com/vapusdata-oss/aistudio/webapp/services"
)

func studioRouters(e *echo.Group) {
	studioGroup := e.Group(routes.StudioGroup)
	studioGroup.GET(routes.AIStudio, services.WebappServiceManager.AIStudioHandler)
	studioGroup.GET(routes.AgentStudio, services.WebappServiceManager.AgentStudioHandler)
}
