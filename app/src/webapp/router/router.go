package router

import (
	"context"
	"io"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vapusdata-oss/aistudio/webapp/pkgs"
	"github.com/vapusdata-oss/aistudio/webapp/routes"
	"github.com/vapusdata-oss/aistudio/webapp/services"
	aipb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func GetNewRouter() *echo.Echo {
	var err error
	app := echo.New()
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST},
	}))
	app.Static("/static", "static")
	pattern := filepath.Join("templates", "*.html")
	renderer := &TemplateRenderer{
		templates: template.Must(template.New("").Funcs(template.FuncMap{
			"limitWords":             limitWords,
			"epochConverter":         EpochConverter,
			"inSlice":                InSlice,
			"toJSON":                 toJSON,
			"stringCheck":            stringCheck,
			"limitletters":           limitletters,
			"sliceContains":          SliceContains,
			"escapeHTML":             escapeHTML,
			"addRand":                addRand,
			"randBool":               randBool,
			"marshalToYaml":          marshalToYaml,
			"epochConverterFull":     EpochConverterFull,
			"strContains":            strContains,
			"protoToJSON":            protoToJSON,
			"sliceLen":               sliceLen[any],
			"slugToTitle":            slugToTitle,
			"epochConverterTextDate": EpochConverterTextDate,
			"escapeJSON":             escapeJSON,
			"joinSlice":              joinSlice,
			"enumoTitle":             enumoTitle,
		}).ParseGlob(pattern)),
	}

	app.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	app.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/ui/datamarketplace")
	})
	services.InitWebappSvc()
	app.Renderer = renderer

	gwmux := runtime.NewServeMux(
		runtime.WithOutgoingHeaderMatcher(runtime.DefaultHeaderMatcher),
		runtime.WithIncomingHeaderMatcher(runtime.DefaultHeaderMatcher),
		runtime.WithMarshalerOption("application/protobuf", &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
				UseEnumNumbers:  false,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(100*1024*1024),
			grpc.MaxCallSendMsgSize(100*1024*1024),
		),
	}
	err = pb.RegisterStudioServiceHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.PlDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering platform service handler from endpoint")
	}
	err = pb.RegisterPluginServiceHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.PlDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering plugin service handler from endpoint")
	}
	err = pb.RegisterUserManagementServiceHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.PlDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering User service handler from endpoint")
	}

	err = pb.RegisterUtilityServiceHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.PlDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering UtilityService service handler from endpoint")
	}

	err = aipb.RegisterAIAgentsHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering AIAgents service handler from endpoint")
	}
	err = aipb.RegisterAIModelStudioHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering AIModelStudio service handler from endpoint")
	}
	err = aipb.RegisterAIModelsHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering AIModels service handler from endpoint")
	}
	err = aipb.RegisterAIPromptsHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering AIPrompts service handler from endpoint")
	}
	err = aipb.RegisterAIAgentStudioHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering AIPrompts service handler from endpoint")
	}
	err = aipb.RegisterAIGuardrailsHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering AIGuardrails service handler from endpoint")
	}
	app.Any("/api/*", echo.WrapHandler(gwmux))
	uigrp := app.Group(routes.UIRoute)
	homeRouter(uigrp)
	authnRouter(app)
	CommonRouter(uigrp)
	studioRouters(uigrp)
	exploreRouter(uigrp)
	manageAIRoutes(uigrp)
	settingsRouter(uigrp)
	developersRouters(uigrp)
	app.HTTPErrorHandler = func(err error, c echo.Context) {
		var (
			code    = http.StatusInternalServerError
			message interface{}
		)

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = he.Message
		} else {
			message = http.StatusText(code)
		}

		// If 404, render the custom 404 template
		if code == http.StatusNotFound {
			if err := c.Render(http.StatusNotFound, "404.html", nil); err != nil {
				c.Logger().Error(err)
			}
			return
		}
		if code == http.StatusInternalServerError {
			if err := c.Render(http.StatusNotFound, "400.html", nil); err != nil {
				c.Logger().Error(err)
			}
			return
		}
		if code == http.StatusUnauthorized {
			if err := c.Render(http.StatusNotFound, "403.html", nil); err != nil {
				c.Logger().Error(err)
			}
			return
		}

		// For other errors, send JSON response or default message
		if !c.Response().Committed {
			c.JSON(code, map[string]interface{}{
				"code":    code,
				"message": message,
			})
		}
	}

	return app
}
