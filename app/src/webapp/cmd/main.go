package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rs/zerolog"
	"github.com/vapusdata-oss/aistudio/core/globals"
	utils "github.com/vapusdata-oss/aistudio/core/utils"
	middlewares "github.com/vapusdata-oss/aistudio/webapp/middlewares"
	pkgs "github.com/vapusdata-oss/aistudio/webapp/pkgs"
	router "github.com/vapusdata-oss/aistudio/webapp/router"
	// "github.com/lucas-clemente/quic-go/http3"
)

var debugLogFlag bool
var configName = "config/webapp-config.yaml"
var logger zerolog.Logger
var flagconfPath string

func loadConfPath() {
	flag.StringVar(&flagconfPath, "conf", "", "config path, eg: -conf /data/vapusdata")
	flag.BoolVar(&debugLogFlag, "debug", false, "debug loggin, set it to true to enable the debug logs")
	flag.Parse()
	if flagconfPath == "" {
		var ok bool
		flagconfPath, ok = os.LookupEnv(globals.SVC_MOUNT_PATH)
		if !ok {
			logger.Fatal().Msgf("SVC_MOUNT_PATH env not found, please set env variable '%v' with dataproduct config to run the product service", globals.SVC_MOUNT_PATH)
		}
	}
	logger.Info().Msgf("Config root Path: %s", flagconfPath)
}

// Initialize the echo server for webapp
func init() {
	// INitialize the logger
	pkgs.InitWAPLogger(debugLogFlag)

	logger = pkgs.GetSubDMLogger(pkgs.IDEN, "VapusData platform server init")

	logger.Info().Msg("Logger middleware Initialized Successfully")

	loadConfPath()
	// Initialize the webapp configuration
	pkgs.InitServiceConfig(flagconfPath, configName)

	pkgs.InitNetworkConfig(flagconfPath, pkgs.WebAppConfigManager.NetworkConfigFile)

	// CHeck if the webapp configuration is initialized
	if pkgs.WebAppConfigManager == nil {
		logger.Fatal().Msg("Failed to initialize the WebAppConfigManager")
	}

	log.Println("WebAppConfigManager Initialized Successfully", pkgs.WebAppConfigManager.URIs.Login)
	// CHeck if the webapp authentication configuration is initialized
	if pkgs.WebAppConfigManager.GetAuthnSecretPath() == utils.EMPTYSTR {
		logger.Fatal().Msg("Failed to initialize the AuthnSecrets")
	}

	// authsrv.InitAuthnService(pkgs.WebAppConfigManager.GetAuthnSecretPath())
	// logger.Info().Msgf("Authentication service Initialized Successfully")

	// pkgs.InitAuthn(pkgs.WebAppConfigManager.GetJwtAuthSecretPath(), validator)
	// logger.Info().Msgf("Jwt Authentication package Initialized Successfully")

	pkgs.InitVapusSvcInternalClients(pkgs.NetworkConfigManager.WebAppSvc.ServiceName, logger)

}

func main() {

	// Initialize the fibr server for webapp
	router := router.GetNewRouter()

	router.Use(middlewares.LoggingMiddleware)

	// server := http3.Server{
	//     Addr:      ":443",
	//     Handler:   e,
	//     TLSConfig: tlsConfig,
	// }

	// server := http.Server{
	// 	Addr:    ":9011",
	// 	Handler: route,
	// }

	// log.Fatal(server.ListenAndServe())
	router.Logger.Fatal(router.Start(fmt.Sprintf(":%d", pkgs.NetworkConfigManager.WebAppSvc.Port)))

}
