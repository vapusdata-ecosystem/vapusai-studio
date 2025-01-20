package router

import (
	"github.com/labstack/echo/v4"
	"github.com/vapusdata-oss/aistudio/webapp/pkgs"
	routes "github.com/vapusdata-oss/aistudio/webapp/routes"
	services "github.com/vapusdata-oss/aistudio/webapp/services"
)

func authnRouter(e *echo.Echo) {
	services.InitAuthnService(pkgs.WebAppConfigManager.Path)
	e.GET(routes.Login, services.AuthnServiceManager.LoginHandler)
	e.GET(routes.LoginCallBack, services.AuthnServiceManager.CallbackHandler)
	e.GET(routes.LoginRedirect, services.AuthnServiceManager.LoginRedirect)
	e.GET(routes.Logout, services.AuthnServiceManager.Logout)
}
