package services

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vapusdata-oss/aistudio/webapp/models"
	pkgs "github.com/vapusdata-oss/aistudio/webapp/pkgs"
	routes "github.com/vapusdata-oss/aistudio/webapp/routes"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func (x *WebappService) HandleError(c echo.Context, err error, httpCode uint, globalContext *models.GlobalContexts) error {
	if httpCode != 0 {
		switch httpCode {
		case http.StatusUnauthorized:
			return c.Render(http.StatusUnauthorized, "401.html", map[string]interface{}{
				"GlobalContext": globalContext,
			})
		case http.StatusForbidden:
			return c.Render(http.StatusForbidden, "403.html", map[string]interface{}{
				"GlobalContext": globalContext,
			})
		case http.StatusNotFound:
			return c.Render(http.StatusNotFound, "404.html", map[string]interface{}{
				"GlobalContext": globalContext,
			})
		case http.StatusInternalServerError:
			return c.Render(http.StatusInternalServerError, "500.html", map[string]interface{}{
				"GlobalContext": globalContext,
			})
		case http.StatusBadRequest:
			return c.Render(http.StatusBadRequest, "400.html", map[string]interface{}{
				"GlobalContext": globalContext,
			})
		default:
			return c.Render(http.StatusBadRequest, "400.html", map[string]interface{}{
				"GlobalContext": globalContext,
			})
		}
	}
	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.Unauthenticated:
			return c.Render(http.StatusUnauthorized, "401.html", map[string]interface{}{
				"GlobalContext": globalContext,
			})
		case codes.NotFound:
			return c.Render(http.StatusNotFound, "404.html", map[string]interface{}{
				"GlobalContext": globalContext,
			})
		case codes.Aborted:
			return c.Render(http.StatusBadRequest, "400.html", map[string]interface{}{
				"GlobalContext": globalContext,
			})
		case codes.PermissionDenied:
			return c.Render(http.StatusForbidden, "403.html", map[string]interface{}{
				"GlobalContext": globalContext,
			})
		case codes.InvalidArgument:
			return c.Render(http.StatusBadRequest, "400.html", map[string]interface{}{
				"GlobalContext": globalContext,
			})
		case codes.Unavailable:
			return c.Render(http.StatusServiceUnavailable, "503.html", map[string]interface{}{
				"GlobalContext": globalContext,
			})
		default:
			return c.Render(http.StatusBadRequest, "400.html", map[string]interface{}{
				"GlobalContext": globalContext,
			})
		}
	}
	return c.Render(http.StatusBadRequest, "400.html", map[string]interface{}{
		"GlobalContext": globalContext,
	})
}

func HandleGLobalContextError(c echo.Context, err error) error {
	if err != nil {
		if status.Code(err).String() == codes.Unauthenticated.String() {
			return c.Redirect(http.StatusSeeOther, routes.Login)
		}
		publicInfo, _ := pkgs.VapusSvcInternalClientManager.SvcConn.StudioPublicInfo(c.Request().Context(), &mpb.EmptyRequest{})

		return c.Render(http.StatusBadRequest, "400.html", map[string]interface{}{
			"publicInfo": publicInfo,
		})
	}
	return nil
}
