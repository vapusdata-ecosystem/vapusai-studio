package services

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	pkgs "github.com/vapusdata-oss/aistudio/webapp/pkgs"
	routes "github.com/vapusdata-oss/aistudio/webapp/routes"
	"github.com/vapusdata-oss/aistudio/webapp/utils"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

func (x *WebappService) AuthOrganizationHandler(c echo.Context) error {
	at, err := utils.GetCookie(c, pkgs.ACCESS_TOKEN)
	if err != nil || at == "" {
		return c.Redirect(http.StatusSeeOther, routes.Login)

	}
	log.Println("Sdfsd", c.Request().URL.Path, "query", c.Request().URL.Query())
	redirectPath := ""
	val, ok := c.Request().URL.Query()["redirect"]
	if !ok {
		redirectPath = routes.UIRoute + routes.HomeGroup
	} else {
		redirectPath = val[0]
	}

	ctx := utils.GetBearerCtx(c.Request().Context(), at)
	Organization := c.Param("Organization")
	result, err := x.grpcClients.UserConn.AccessTokenInterface(ctx, &pb.AccessTokenInterfaceRequest{
		Organization: Organization,
		Utility:      pb.AccessTokenAgentUtility_ORG_LOGIN,
	})
	if err != nil {
		x.logger.Err(err).Msg("error while authenticating Organization")
		return err
	}
	atCookie := &http.Cookie{
		Name:    pkgs.ACCESS_TOKEN,
		Value:   result.Token.AccessToken,
		Expires: time.Unix(result.Token.ValidTill, 0),
		Path:    "/",
	}

	c.SetCookie(atCookie)
	time.Sleep(2 * time.Second)
	// Redirect to the home page
	return c.Redirect(http.StatusSeeOther, redirectPath)
}
