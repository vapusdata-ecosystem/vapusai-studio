package grpcops

import (
	"context"
	"fmt"
	"log"

	"buf.build/go/protoyaml"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	utils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	rpcauth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	grpcstatus "google.golang.org/grpc/status"
)

var ErrorFinal = "errorFinal"

var LocalAuthzTokenCtx = "authorization:Bearer"

// HandleGrpcError is a utility function to handle grpc errors
func HandleGrpcError(err error, code grpccodes.Code) error {
	e, ok := err.(dmerrors.Error)
	if !ok {
		return grpcstatus.Error(code, err.Error())
	}
	return grpcstatus.Error(code, e.Error())
}

func HandleCtxCustomMessage(ctx context.Context, msgType string, msg ...string) context.Context {
	if ctx.Value(utils.CUSTOM_MESSAGE) == nil {
		cm := map[string]interface{}{msgType: msg}
		return context.WithValue(ctx, utils.CUSTOM_MESSAGE, cm)
	}
	cm := ctx.Value(utils.CUSTOM_MESSAGE).(map[string]interface{})
	cm[msgType] = append(cm[msgType].([]string), msg...)
	return context.WithValue(ctx, utils.CUSTOM_MESSAGE, cm)
}

func GetRpcAuthFromCtx(ctx context.Context) (string, context.Context, error) {
	token, err := rpcauth.AuthFromMD(ctx, "bearer")
	if err != nil || token == "" {
		return "", ctx, grpcstatus.Error(grpccodes.Unauthenticated, "Authentication bearer token not found in request metadata")
	}
	return token, context.WithValue(ctx, LocalAuthzTokenCtx, token), nil
}

// HandleResponse is a utility function to handle the base response
func HandleDMResponse(ctx context.Context, opts ...string) *mpb.DMResponse {
	if ctx.Value(utils.CUSTOM_MESSAGE) == nil {
		return nil
	}
	cm := ctx.Value(utils.CUSTOM_MESSAGE).(map[string]interface{})
	return &mpb.DMResponse{
		Message:      opts[0],
		DmStatusCode: opts[1],
		CustomMessage: func(cms map[string]interface{}) []*mpb.MapList {
			var cm []*mpb.MapList
			for k, v := range cms {
				cm = append(cm, &mpb.MapList{
					Key:    k,
					Values: v.([]string),
				})
			}
			return cm
		}(cm),
	}
}

func GetSvcDns(svcName string, namespace string, port int64) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local:%d", svcName, namespace, port)
}

var ProtoJsonMarshaller = protojson.MarshalOptions{}
var ProtoJsonUnMarshaller = protojson.UnmarshalOptions{}

// var ProtoYamlMarshaller = protojson.MarshalOptions{}
// var ProtoYamlUnMarshaller = protojson.MarshalOptions{}

func SwapNewContextWithAuthToken(ctx context.Context) context.Context {
	token, err := rpcauth.AuthFromMD(ctx, "bearer")
	if err != nil {
		token = ""
	}
	return metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer "+token))
}

func NewContextWithCopyAuthToken(ctx context.Context) context.Context {
	token, err := rpcauth.AuthFromMD(ctx, "bearer")
	if err != nil {
		token = ""
	}
	log.Println("ctx: <><>??<><>??", token)
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println("Outgoing Metadata ctx :", md)
	} else {
		fmt.Println("Failed to attach metadata to the new context ctx")
	}

	newCtx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))
	token, err = rpcauth.AuthFromMD(newCtx, "bearer")
	if err != nil {
		token = ""
	}
	log.Println("newCtx: <><>??<><>??", token)
	md, ok = metadata.FromOutgoingContext(newCtx)
	if ok {
		fmt.Println("Outgoing Metadata newCtx :", md)
	} else {
		fmt.Println("Failed to attach metadata to the new context newCtx")
	}

	backgroundCtx := context.WithValue(context.Background(), "authorization", token)
	newCtx1 := metadata.NewOutgoingContext(backgroundCtx, metadata.Pairs("authorization", "Bearer "+token))
	token, err = rpcauth.AuthFromMD(newCtx1, "bearer")
	if err != nil {
		token = ""
	}
	log.Println("newCtx1: <><>??<><>??", token)
	md, ok = metadata.FromOutgoingContext(newCtx1)
	if ok {
		fmt.Println("Outgoing Metadata newCtx1:", md)
	} else {
		fmt.Println("Failed to attach metadata to the new context newCtx1")
	}

	return newCtx
}

func GetPbAnyToGoAny(pbAny *anypb.Any) any {
	var result any
	if pbAny == nil {
		return ""
	}
	listValue := &structpb.ListValue{}

	if err := pbAny.UnmarshalTo(listValue); err == nil {
		sliceResult := []any{}
		for _, value := range listValue.Values {
			switch val := value.Kind.(type) {
			case *structpb.Value_StringValue:
				sliceResult = append(sliceResult, val.StringValue)
			case *structpb.Value_NumberValue:
				sliceResult = append(sliceResult, val.NumberValue)
			default:
				log.Printf("Unknown type in list value: %v", value)
			}
		}
		return sliceResult
	}

	// Handle scalar types
	stringValue := &wrapperspb.StringValue{}
	intValue := &wrapperspb.Int64Value{}

	if err := pbAny.UnmarshalTo(stringValue); err == nil {
		result = stringValue.Value
		return result
	}

	if err := pbAny.UnmarshalTo(intValue); err == nil {
		result = intValue.Value
		return result
	}
	return result
}

func ProtoYamlMarshal(message protoreflect.ProtoMessage) ([]byte, error) {
	mars := protoyaml.MarshalOptions{
		Indent:          2,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
	}
	return mars.Marshal(message)
}
