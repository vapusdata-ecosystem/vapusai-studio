package server

import (
	"context"
	"path/filepath"

	dmstores "github.com/vapusdata-oss/aistudio/aistudio/datastoreops"
	"github.com/vapusdata-oss/aistudio/aistudio/nabhiksvc"
	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	services "github.com/vapusdata-oss/aistudio/aistudio/services"
	serviceopss "github.com/vapusdata-oss/aistudio/core/serviceops"
)

func packagesInit() {
	//Initialize the logger
	pkgs.InitWAPLogger(debugLogFlag)

	logger = pkgs.GetSubDMLogger(pkgs.IDEN, "AIstudio server init")

	logger.Info().Msg("Loading service config for VapusData server")
	// Load the service configuration, secrets inton the memory of the service. These information will be used by the service to connect to the database, vault etc connections
	pkgs.InitServiceConfig(flagconfPath, filepath.Join(flagconfPath, configName))

	pkgs.InitNetworkConfig(flagconfPath, filepath.Join(flagconfPath, pkgs.ServiceConfigManager.NetworkConfigFile))

	logger.Info().Msg("Service config loaded successfully")

	ctx := context.Background()

	bootStores(ctx, pkgs.ServiceConfigManager)
	logger.Info().Msg("Service data stores loaded successfully")

	dmstores.InitStoreDependencies(ctx, pkgs.ServiceConfigManager)
	logger.Info().Msg("Service store dependencies loaded successfully")

	logger.Info().Msg("Service config loaded successfully")
	// Initialize the jwt authn validator
	logger.Info().Msgf("Loading JWT authn with secret path: %s", pkgs.ServiceConfigManager.GetJwtAuthSecretPath())

	// Initialize the NewVapusAuth
	pkgs.InitAIStudioAuth(pkgs.JwtAuthnParams, pkgs.DmValidator)

	pkgs.InitserviceopsInternalClients()
	bootConnectionPool()
	defer ctx.Done()
}

func bootStores(ctx context.Context, conf *serviceopss.VapusAISvcConfig) {
	//Boot the stores
	logger.Info().Msg("Booting the data stores")
	dmstores.InitDMStore(conf)
	if dmstores.DMStoreManager.Error != nil {
		logger.Fatal().Err(dmstores.DMStoreManager.Error).Msg("error while initializing data stores.")
	}
	services.InitStudioServices(dmstores.DMStoreManager)
}

func bootConnectionPool() {
	// Boot the connection pool
	logger.Info().Msg("Booting the connection pool")
	resp := nabhiksvc.InitAIModelNodeConnectionPool(nabhiksvc.WithLogger(logger),
		nabhiksvc.WithDMStore(dmstores.DMStoreManager))
	if resp != nil {
		nabhiksvc.AIModelNodeConnectionPoolManager = resp
	}
	logger.Info().Msg("Connection pool booted successfully")
	nabhiksvc.InitGuardrailPool()
	logger.Info().Msg("Guardrail pool booted successfully")
}
