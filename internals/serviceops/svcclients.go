package svcops

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/rs/zerolog"
	grpcops "github.com/vapusdata-oss/aistudio/core/serviceops/grpcops"
	utils "github.com/vapusdata-oss/aistudio/core/utils"
	aipb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

type VapusSvcInternalClients struct {
	Host                    string
	UserConn                pb.UserManagementServiceClient
	SvcConn                 pb.StudioServiceClient
	AIStudioConn            aipb.AIModelStudioClient
	AIAgentClient           aipb.AIAgentsClient
	AIPromptClient          aipb.AIPromptsClient
	AIModelClient           aipb.AIModelsClient
	AIGurdrailsClient       aipb.AIGuardrailsClient
	platformGrpcClient      *grpcops.GrpcClient
	aiStudioGrpcClient      *grpcops.GrpcClient
	PluginServiceClient     pb.PluginServiceClient
	NetworkConfig           *NetworkConfig
	PlDns                   string
	AIStudioDns             string
	NabhikServerDns         string
	dataProductServerClient *grpcops.GrpcClient
}

func SvcUpTimeCheck(ctx context.Context, networkConfig *NetworkConfig, self string, logger zerolog.Logger, counter int64) error {
	if counter > 1 {
		time.Sleep(15 * time.Second)
	}
	counter++
	logger.Info().Msg("Checking if all services are up........")
	platformSvcDns := fmt.Sprintf("%s:%d", networkConfig.StudioSvc.ServiceName, networkConfig.StudioSvc.ServicePort)
	aiStudioSvcDns := fmt.Sprintf("%s:%d", networkConfig.AIStudioSvc.ServiceName, networkConfig.AIStudioSvc.ServicePort)
	nabhikServerDns := fmt.Sprintf("%s:%d", networkConfig.NabhikServer.ServiceName, networkConfig.NabhikServer.ServicePort)
	err := utils.Telnet("tcp", platformSvcDns)
	if err != nil {
		logger.Error().Err(err).Msg("Studio service is not up yet")
		if counter > 6 {
			return err
		} else {
			return SvcUpTimeCheck(ctx, networkConfig, self, logger, counter)
		}
	}
	err = utils.Telnet("tcp", aiStudioSvcDns)
	if err != nil {
		logger.Error().Err(err).Msg("AI Studio service is not up yet")
		if counter > 6 {
			return err
		} else {
			return SvcUpTimeCheck(ctx, networkConfig, self, logger, counter)
		}
	}
	err = utils.Telnet("tcp", nabhikServerDns)
	if err != nil {
		logger.Error().Err(err).Msg("Data Product Server service is not up yet")
		if counter > 6 {
			return err
		} else {
			return SvcUpTimeCheck(ctx, networkConfig, self, logger, counter)
		}
	}
	return nil
}

func SetupVapusSvcInternalClients(ctx context.Context, networkConfig *NetworkConfig, self string, logger zerolog.Logger) (*VapusSvcInternalClients, error) {
	var err error
	client := &VapusSvcInternalClients{
		AIStudioDns:     fmt.Sprintf("%s:%d", networkConfig.AIStudioSvc.ServiceName, networkConfig.AIStudioSvc.ServicePort),
		PlDns:           fmt.Sprintf("%s:%d", networkConfig.StudioSvc.ServiceName, networkConfig.StudioSvc.ServicePort),
		NabhikServerDns: fmt.Sprintf("%s:%d", networkConfig.NabhikServer.ServiceName, networkConfig.NabhikServer.ServicePort),
	}
	logger.Info().Msg("Setting up VapusSvcInternalClients........")
	logger.Info().Msgf("StudioSvcDns: %s", client.PlDns)
	logger.Info().Msgf("AIStudioSvcDns: %s", client.AIStudioDns)
	logger.Info().Msgf("DataproductServerDns: %s", client.NabhikServerDns)
	err = utils.Telnet("tcp", client.PlDns)
	if err != nil && self != "" && self != networkConfig.StudioSvc.ServiceName {
		logger.Error().Err(err).Msg("Studio service is not up yet")
		client.platformGrpcClient = nil
	}
	logger.Info().Msg("Setting up VapusSvcInternalClients........")
	log.Println("client.platformGrpcClient: ", client.platformGrpcClient)
	if client.platformGrpcClient == nil {
		client.platformGrpcClient = grpcops.NewGrpcClient(logger,
			grpcops.ClientWithInsecure(true),
			grpcops.ClientWithServiceAddress(client.PlDns))
		client.SvcConn = pb.NewStudioServiceClient(client.platformGrpcClient.Connection)
		client.UserConn = pb.NewUserManagementServiceClient(client.platformGrpcClient.Connection)
		client.PluginServiceClient = pb.NewPluginServiceClient(client.platformGrpcClient.Connection)
	}

	err = utils.Telnet("tcp", client.AIStudioDns)
	if err != nil && self != "" && self != networkConfig.AIStudioSvc.ServiceName {
		logger.Error().Err(err).Msg("AI Studio service is not up yet")
		client.aiStudioGrpcClient = nil
	}
	log.Println("client.aiStudioGrpcClient: ", client.aiStudioGrpcClient)
	if client.aiStudioGrpcClient == nil {
		client.aiStudioGrpcClient = grpcops.NewGrpcClient(logger,
			grpcops.ClientWithInsecure(true),
			grpcops.ClientWithServiceAddress(client.AIStudioDns))
		client.AIStudioConn = aipb.NewAIModelStudioClient(client.aiStudioGrpcClient.Connection)
		client.AIAgentClient = aipb.NewAIAgentsClient(client.aiStudioGrpcClient.Connection)
		client.AIPromptClient = aipb.NewAIPromptsClient(client.aiStudioGrpcClient.Connection)
		client.AIModelClient = aipb.NewAIModelsClient(client.aiStudioGrpcClient.Connection)
		client.AIGurdrailsClient = aipb.NewAIGuardrailsClient(client.aiStudioGrpcClient.Connection)
	}
	client.NetworkConfig = networkConfig

	err = utils.Telnet("tcp", client.NabhikServerDns)
	if err != nil && self != "" && self != networkConfig.NabhikServer.ServiceName {
		logger.Error().Err(err).Msg("NabhikServer service is not up yet")
		client.dataProductServerClient = nil
	}
	log.Println("client.dataProductServerClient: ", client.dataProductServerClient)
	if client.dataProductServerClient == nil {
		client.dataProductServerClient = grpcops.NewGrpcClient(logger,
			grpcops.ClientWithInsecure(true),
			grpcops.ClientWithServiceAddress(client.NabhikServerDns))
	}
	client.NetworkConfig = networkConfig
	log.Println("platformGrpcClient: ", client.platformGrpcClient)
	log.Println("aiStudioGrpcClient: ", client.aiStudioGrpcClient)
	log.Println("NabhikServer: ", client.dataProductServerClient)
	return client, nil
}

func SetupAIStudioClient(ctx context.Context, dns string, self string, logger zerolog.Logger) (aipb.AIModelStudioClient, error) {
	// dns = "localhost:9013"
	telnet2, err := net.DialTimeout("tcp", dns, 1*time.Second)
	if err != nil {
		if self != "ai-studio" {
			return nil, err
		}
	}
	if telnet2 != nil {
		defer telnet2.Close()
	}
	grpcClient := grpcops.NewGrpcClient(logger,
		grpcops.ClientWithInsecure(true),
		grpcops.ClientWithServiceAddress(dns))
	return aipb.NewAIModelStudioClient(grpcClient.Connection), nil
}

func (x *VapusSvcInternalClients) Close() {
	x.UserConn = nil
	x.SvcConn = nil
	x.AIStudioConn = nil
	x.aiStudioGrpcClient = nil
	x.platformGrpcClient = nil
	x.dataProductServerClient = nil
}

func (x *VapusSvcInternalClients) PingTestAndReconnect(ctx context.Context, dns string, logger zerolog.Logger) error {
	err := utils.Telnet("tcp", dns)
	if err == nil {
		w, err := SetupVapusSvcInternalClients(ctx, x.NetworkConfig, "", logger)
		if err != nil {
			return err
		}
		x = w
		// x.PlConn = w.PlConn
		// x.UserConn = w.UserConn
		// x.OrganizationConn = w.OrganizationConn
		// x.aiStudioGrpcClient = w.aiStudioGrpcClient
		// x.platformGrpcClient = w.platformGrpcClient
		// x.AIStudioConn = w.AIStudioConn
		// x.DataProductServerConn = w.DataProductServerConn
		return nil
	}
	return err
}
