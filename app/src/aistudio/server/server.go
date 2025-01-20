package server

import (
	"context"
	"flag"
	"os"

	"github.com/rs/zerolog"
	dmcontrollers "github.com/vapusdata-oss/aistudio/aistudio/controllers"
	middlewares "github.com/vapusdata-oss/aistudio/aistudio/middlewares"
	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	"github.com/vapusdata-oss/aistudio/core/globals"
	grpcops "github.com/vapusdata-oss/aistudio/core/serviceops/grpcops"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"

	interceptors "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	rpcauth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	selector "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	grpc "google.golang.org/grpc"
	health "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	debugLogFlag bool
	flagconfPath string
	configName   = "config/aistudio-service-config.yaml"
	networkDir   = "network"
	logger       zerolog.Logger
)

func init() {
	flag.StringVar(&flagconfPath, "conf", "", "config path, eg: --conf=/data/vapusdata")
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
	packagesInit()
}

func initServer(grpcServer *grpcops.GRPCServer) {

	// Setup auth matcher.
	allButTheez := func(ctx context.Context, callMeta interceptors.CallMeta) bool {
		return healthpb.Health_ServiceDesc.ServiceName != callMeta.Service
	}

	// Add unary and stream interceptors for prometheus
	var opts = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			grpc_prometheus.UnaryServerInterceptor,
			selector.UnaryServerInterceptor(rpcauth.UnaryServerInterceptor(middlewares.AuthnMiddleware), selector.MatchFunc(allButTheez)),
			selector.UnaryServerInterceptor(middlewares.UnaryRequestValidator(), selector.MatchFunc(allButTheez)),
		),
		grpc.ChainStreamInterceptor(
			grpc_prometheus.StreamServerInterceptor,
			selector.StreamServerInterceptor(rpcauth.StreamServerInterceptor(middlewares.AuthnMiddleware), selector.MatchFunc(allButTheez)),
		),
	}

	// Create a new GRPC server
	//First step is to configure the vapusStudio server
	grpcServer.ServerPort = pkgs.NetworkConfigManager.AIStudioSvc.Port

	// Initialize the server
	logger.Info().Msg("Configuring VapusData Studio Server")

	// Initialize the grpc server ops and net listner for the server
	grpcServer.ConfigureGrpcServer(opts, debugLogFlag)
	logger.Info().Msg("VapusData Studio Server configured successfully.")
	if grpcServer.GrpcServ == nil {
		logger.Info().Msg("Failed to initialize VapusData Studio Server")
	}
	healthcheck := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer.GrpcServ, healthcheck)
	// Register the VapusData Studio and Node controller
	reflection.Register(grpcServer.GrpcServ)
	pb.RegisterAIPromptsServer(grpcServer.GrpcServ, dmcontrollers.NewAIPrompts())
	pb.RegisterAIAgentsServer(grpcServer.GrpcServ, dmcontrollers.NewVapusAIAgents())
	pb.RegisterAIModelsServer(grpcServer.GrpcServ, dmcontrollers.NewAIModels())
	pb.RegisterAIModelStudioServer(grpcServer.GrpcServ, dmcontrollers.NewAIModelStudio())
	pb.RegisterAIAgentStudioServer(grpcServer.GrpcServ, dmcontrollers.NewVapusAIAgentStudio())
	pb.RegisterAIGuardrailsServer(grpcServer.GrpcServ, dmcontrollers.NewVapusAIGuardrails())
	// if grpcServer.SetUpHttp {
	// 	grpcServer.HttpGwPort = int32(pkgs.NetworkConfigManager.AIStudioSvc.HttpGwPort)
	// 	grpcServer.MuxInstance = runtime.NewServeMux(
	// 		runtime.WithMetadata(middlewares.HttpAuthnMiddleware),
	// 	)
	// 	if err := pb.RegisterVapusAIStudioHandlerServer(context.Background(), grpcServer.MuxInstance, dmcontrollers.NewVapusAIStudio()); err != nil {
	// 		logger.Fatal().Err(err).Msg("Failed to register the Vapus Data Product Server handler")
	// 	}
	// 	logger.Info().Msgf("Http Server configured at - %v", grpcServer.Http1Address)
	// }
	logger.Info().Msgf("Grpc Server configured at - %v", grpcServer.ServerAddress)
}

func GrpcServer() *grpcops.GRPCServer {
	var grpcServer *grpcops.GRPCServer
	var serverOpts []grpcops.GrpcOptions
	// Initialize the service configuration
	if pkgs.ServiceConfigManager.ServerCerts.Mtls {
		logger.Info().Msg("Configuring VapusData AIStudio Server with MTLS connection")
		serverOpts = append(serverOpts, grpcops.WithMtls(pkgs.ServiceConfigManager.GetMtlsCerts()))
	} else if pkgs.ServiceConfigManager.ServerCerts.PlainTls {
		logger.Info().Msg("Configuring VapusData AIStudio Server with PlainTLS connection")
		serverOpts = append(serverOpts, grpcops.WithPlainTls(pkgs.ServiceConfigManager.GetPlainTlsCerts()))
	} else {
		logger.Info().Msg("Configuring VapusData AIStudio Server with insecure connection")
		serverOpts = append(serverOpts, grpcops.WithInsecure(true))
	}
	grpcServer = grpcops.NewGRPCServer(serverOpts...)
	grpcServer.Logger = logger
	grpcServer.SetUpHttp = false
	initServer(grpcServer)
	return grpcServer
}
