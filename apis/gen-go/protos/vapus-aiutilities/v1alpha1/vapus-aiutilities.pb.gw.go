// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: protos/vapus-aiutilities/v1alpha1/vapus-aiutilities.proto

/*
Package v1alpha1 is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package v1alpha1

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = metadata.Join

func request_AIUtility_GenerateEmbedding_0(ctx context.Context, marshaler runtime.Marshaler, client AIUtilityClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq GenerateEmbeddingRequest
	var metadata runtime.ServerMetadata

	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.GenerateEmbedding(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_AIUtility_GenerateEmbedding_0(ctx context.Context, marshaler runtime.Marshaler, server AIUtilityServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq GenerateEmbeddingRequest
	var metadata runtime.ServerMetadata

	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := server.GenerateEmbedding(ctx, &protoReq)
	return msg, metadata, err

}

func request_AIUtility_SensitivityAnalyzer_0(ctx context.Context, marshaler runtime.Marshaler, client AIUtilityClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq SensitivityAnalyzerRequest
	var metadata runtime.ServerMetadata

	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.SensitivityAnalyzer(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_AIUtility_SensitivityAnalyzer_0(ctx context.Context, marshaler runtime.Marshaler, server AIUtilityServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq SensitivityAnalyzerRequest
	var metadata runtime.ServerMetadata

	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := server.SensitivityAnalyzer(ctx, &protoReq)
	return msg, metadata, err

}

// RegisterAIUtilityHandlerServer registers the http handlers for service AIUtility to "mux".
// UnaryRPC     :call AIUtilityServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterAIUtilityHandlerFromEndpoint instead.
func RegisterAIUtilityHandlerServer(ctx context.Context, mux *runtime.ServeMux, server AIUtilityServer) error {

	mux.Handle("POST", pattern_AIUtility_GenerateEmbedding_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/vapusdata.ai_studio.v1alpha1.AIUtility/GenerateEmbedding", runtime.WithHTTPPathPattern("/vapusdata.ai_studio.v1alpha1.AIUtility/GenerateEmbedding"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_AIUtility_GenerateEmbedding_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_AIUtility_GenerateEmbedding_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("POST", pattern_AIUtility_SensitivityAnalyzer_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/vapusdata.ai_studio.v1alpha1.AIUtility/SensitivityAnalyzer", runtime.WithHTTPPathPattern("/vapusdata.ai_studio.v1alpha1.AIUtility/SensitivityAnalyzer"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_AIUtility_SensitivityAnalyzer_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_AIUtility_SensitivityAnalyzer_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

// RegisterAIUtilityHandlerFromEndpoint is same as RegisterAIUtilityHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterAIUtilityHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.NewClient(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Errorf("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Errorf("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterAIUtilityHandler(ctx, mux, conn)
}

// RegisterAIUtilityHandler registers the http handlers for service AIUtility to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterAIUtilityHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterAIUtilityHandlerClient(ctx, mux, NewAIUtilityClient(conn))
}

// RegisterAIUtilityHandlerClient registers the http handlers for service AIUtility
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "AIUtilityClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "AIUtilityClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "AIUtilityClient" to call the correct interceptors.
func RegisterAIUtilityHandlerClient(ctx context.Context, mux *runtime.ServeMux, client AIUtilityClient) error {

	mux.Handle("POST", pattern_AIUtility_GenerateEmbedding_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/vapusdata.ai_studio.v1alpha1.AIUtility/GenerateEmbedding", runtime.WithHTTPPathPattern("/vapusdata.ai_studio.v1alpha1.AIUtility/GenerateEmbedding"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_AIUtility_GenerateEmbedding_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_AIUtility_GenerateEmbedding_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("POST", pattern_AIUtility_SensitivityAnalyzer_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/vapusdata.ai_studio.v1alpha1.AIUtility/SensitivityAnalyzer", runtime.WithHTTPPathPattern("/vapusdata.ai_studio.v1alpha1.AIUtility/SensitivityAnalyzer"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_AIUtility_SensitivityAnalyzer_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_AIUtility_SensitivityAnalyzer_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_AIUtility_GenerateEmbedding_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"vapusdata.ai_studio.v1alpha1.AIUtility", "GenerateEmbedding"}, ""))

	pattern_AIUtility_SensitivityAnalyzer_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"vapusdata.ai_studio.v1alpha1.AIUtility", "SensitivityAnalyzer"}, ""))
)

var (
	forward_AIUtility_GenerateEmbedding_0 = runtime.ForwardResponseMessage

	forward_AIUtility_SensitivityAnalyzer_0 = runtime.ForwardResponseMessage
)
