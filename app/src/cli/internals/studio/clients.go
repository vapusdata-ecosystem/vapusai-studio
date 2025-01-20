package plclient

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/bufbuild/protoyaml-go"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"

	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
	grpcops "github.com/vapusdata-oss/aistudio/core/serviceops/grpcops"
	gwcl "github.com/vapusdata-oss/aistudio/core/serviceops/httpcls"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	aipb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
	dpb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

var AgentGoals = map[string][]interface{}{
	"account":               getAgentOps(pb.AccountAgentActions_name),
	"vapussearch":           getAgentOps(pb.VapusSearchType_name),
	"domains":               getAgentOps(dpb.DomainAgentActions_name),
	"user":                  getAgentOps(pb.UserAgentActions_name),
	"authz":                 getAgentOps(pb.AuthzAgentActions_name),
	"authorization":         getAgentOps(pb.AccessTokenAgentUtility_name),
	"utility":               getAgentOps(nil),
	pkg.AIModelNodeResource: getAgentOps(aipb.AIModelNodeConfiguratorActions_name),
}

type ActionHandlerOpts struct {
	ParentCmd   string
	Args        []string
	Action      string
	File        string
	La          bool
	AccessToken string
	SearchQ     string
	Params      map[string]string
	Resource    string
}

type VapusStudioClient struct {
	Host               string
	PlConn             pb.StudioServiceClient
	UserConn           pb.UserManagementServiceClient
	DomainConn         dpb.DomainServiceClient
	AIStudioConn       aipb.AIAgentStudioClient
	platformGrpcClient *grpcops.GrpcClient
	aiStudioGrpcClient *grpcops.GrpcClient
	CaCertFile         string
	ClientCertFile     string
	ClientKeyFile      string
	ValidTill          time.Time
	Error              error
	logger             zerolog.Logger
	ResourceActionMap  map[string][]interface{}
	inputFormat        string
	ActionHandler      ActionHandlerOpts
	protoyamlUnMarshal protoyaml.UnmarshalOptions
	protoyamlMarshal   protoyaml.MarshalOptions
	fileBytes          []byte
	GwClient           *gwcl.VapusHttpClient
	protojsonMarshal   protojson.MarshalOptions
}

func getAgentOps(enum_map map[int32]string) []interface{} {
	var ops []interface{}
	if enum_map == nil {
		return ops
	}
	for _, v := range enum_map {
		ops = append(ops, v)
	}
	return ops
}

func NewPlatFormClient(params map[string]string, logger zerolog.Logger) (*VapusStudioClient, error) {
	urlParsed, err := url.Parse(params["url"])
	if err != nil {
		return nil, err
	}
	// namespace, ok := params["namespace"]
	// if !ok {
	// 	return nil, errors.New("namespace is required, missing from the context")
	// }
	// port, ok := params["port"]
	// if !ok {
	// 	return nil, errors.New("port is required, missing from the context")
	// }
	// portI, err := strconv.Atoi(port)
	// if err != nil {
	// 	return nil, errors.Join(err, errors.New("port is not a valid integer"))
	// }
	dns := fmt.Sprintf("Connecting to vapusdata instance at - %s", urlParsed.String())
	// dns = "localhost:9013"

	// telnet, err := net.DialTimeout("tcp", dns, 1*time.Second)
	// if err != nil {
	// 	return nil, err
	// }
	// defer telnet.Close()

	// grpcClient := pbtools.NewGrpcClient(logger,
	// 	pbtools.ClientWithInsecure(true),
	// 	pbtools.ClientWithServiceAddress(dns))
	gwcls, err := gwcl.New(params["url"], logger)
	if err != nil {
		// return nil, err
		return nil, nil
	}

	cl := &VapusStudioClient{
		Host: dns,
		// PlConn:             pb.NewStudioServiceClient(grpcClient.Connection),
		// platformGrpcClient: grpcClient,
		// UserConn:           pb.NewUserManagementServiceClient(grpcClient.Connection),
		// DomainConn:         dpb.NewDomainServiceClient(grpcClient.Connection),
		logger:             logger,
		protoyamlUnMarshal: protoyaml.UnmarshalOptions{},
		protoyamlMarshal: protoyaml.MarshalOptions{
			EmitUnpopulated: true,
		},
		protojsonMarshal: protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
			UseEnumNumbers:  false,
		},
		ActionHandler: ActionHandlerOpts{},
		GwClient:      gwcls,
	}
	return cl, nil
}

func (x *VapusStudioClient) setAIStudioClient(url string) error {
	svcinfo, err := x.PlConn.StudioServicesInfo(context.Background(), &pb.StudioServicesRequest{})
	if err != nil {
		return err
	}
	if svcinfo.GetNetworkParams() != nil || len(svcinfo.GetNetworkParams()) < 1 {
		for _, svc := range svcinfo.GetNetworkParams() {
			if svc.SvcTag == mpb.VapusSvcs_AISTUDIO {
				telnet, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", url, svc.Port), 1*time.Second)
				if err != nil {
					return err
				}
				defer telnet.Close()
				x.aiStudioGrpcClient = grpcops.NewGrpcClient(x.logger,
					grpcops.ClientWithInsecure(true),
					grpcops.ClientWithServiceAddress(fmt.Sprintf("%s:%d", url, svc.Port)))
				x.AIStudioConn = aipb.NewAIAgentStudioClient(x.aiStudioGrpcClient.Connection)
				return nil
			}
		}
	}
	return errors.New("AI Studio service not found")
}

func (x *VapusStudioClient) Close() {
	x.PlConn = nil
	x.UserConn = nil
	x.DomainConn = nil
	if x.platformGrpcClient != nil {
		x.platformGrpcClient.Close()
	}
	if x.aiStudioGrpcClient != nil {
		x.aiStudioGrpcClient.Close()
	}
}

func (x *VapusStudioClient) PrintDescribe(data protoreflect.ProtoMessage, resource string) {
	bytes, err := x.protoyamlMarshal.Marshal(data)
	if err != nil {
		x.logger.Error().Msgf("Error in marshaling %v details", resource)
	}
	x.logger.Info().Msgf("\n%s", string(bytes))
}
