//*
// @fileoverview
// This file contains the protocol buffer definitions for the VapusDataStudio API.
// The VapusDataStudioService API allows users to create and manage data marketplacees, which are
// virtualized data environments that provide a unified view of data from
// multiple sources.
//
// @packageDocumentation

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: protos/vapusai-studio/v1alpha1/organization.proto

package v1alpha1

import (
	context "context"
	v1alpha1 "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	StudioService_StudioPublicInfo_FullMethodName    = "/vapusai.v1alpha1.StudioService/StudioPublicInfo"
	StudioService_OrganizationManager_FullMethodName = "/vapusai.v1alpha1.StudioService/OrganizationManager"
	StudioService_OrganizationGetter_FullMethodName  = "/vapusai.v1alpha1.StudioService/OrganizationGetter"
	StudioService_AccountGetter_FullMethodName       = "/vapusai.v1alpha1.StudioService/AccountGetter"
	StudioService_AccountManager_FullMethodName      = "/vapusai.v1alpha1.StudioService/AccountManager"
)

// StudioServiceClient is the client API for StudioService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// *
// The StudioService is the main service interface for managing vapusai services.
// It provides methods for creating accounts, retrieving account information,
// creating data marketplacees, storing data marketplace secrets, and retrieving data marketplace information.
type StudioServiceClient interface {
	StudioPublicInfo(ctx context.Context, in *v1alpha1.EmptyRequest, opts ...grpc.CallOption) (*StudioPublicInfoResponse, error)
	//*
	// Manages a org.
	// @param {OrganizationManagerRequest} request - The request object containing the org details.
	// @returns {OrganizationResponse} The response object containing the result of the org management actions.
	OrganizationManager(ctx context.Context, in *OrganizationManagerRequest, opts ...grpc.CallOption) (*OrganizationResponse, error)
	//*
	// Retrieves a org.
	// @param {OrganizationGetterRequest} request - The request object containing the org ID.
	// @returns {OrganizationResponse} The response object containing the retrieved org.
	OrganizationGetter(ctx context.Context, in *OrganizationGetterRequest, opts ...grpc.CallOption) (*OrganizationResponse, error)
	//*
	// Retrieves information about theaccount.
	// @param {AccountManager} request - The account creation request.
	// @returns {AccountResponse} - The account creation response.
	AccountGetter(ctx context.Context, in *v1alpha1.EmptyRequest, opts ...grpc.CallOption) (*AccountResponse, error)
	//*
	// Creates a new account with the specified name.
	// @param {AccountManager} request - The account creation request.
	// @returns {AccountResponse} - The account creation response.
	AccountManager(ctx context.Context, in *AccountManagerRequest, opts ...grpc.CallOption) (*AccountResponse, error)
}

type studioServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStudioServiceClient(cc grpc.ClientConnInterface) StudioServiceClient {
	return &studioServiceClient{cc}
}

func (c *studioServiceClient) StudioPublicInfo(ctx context.Context, in *v1alpha1.EmptyRequest, opts ...grpc.CallOption) (*StudioPublicInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StudioPublicInfoResponse)
	err := c.cc.Invoke(ctx, StudioService_StudioPublicInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *studioServiceClient) OrganizationManager(ctx context.Context, in *OrganizationManagerRequest, opts ...grpc.CallOption) (*OrganizationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrganizationResponse)
	err := c.cc.Invoke(ctx, StudioService_OrganizationManager_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *studioServiceClient) OrganizationGetter(ctx context.Context, in *OrganizationGetterRequest, opts ...grpc.CallOption) (*OrganizationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrganizationResponse)
	err := c.cc.Invoke(ctx, StudioService_OrganizationGetter_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *studioServiceClient) AccountGetter(ctx context.Context, in *v1alpha1.EmptyRequest, opts ...grpc.CallOption) (*AccountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AccountResponse)
	err := c.cc.Invoke(ctx, StudioService_AccountGetter_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *studioServiceClient) AccountManager(ctx context.Context, in *AccountManagerRequest, opts ...grpc.CallOption) (*AccountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AccountResponse)
	err := c.cc.Invoke(ctx, StudioService_AccountManager_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StudioServiceServer is the server API for StudioService service.
// All implementations must embed UnimplementedStudioServiceServer
// for forward compatibility
//
// *
// The StudioService is the main service interface for managing vapusai services.
// It provides methods for creating accounts, retrieving account information,
// creating data marketplacees, storing data marketplace secrets, and retrieving data marketplace information.
type StudioServiceServer interface {
	StudioPublicInfo(context.Context, *v1alpha1.EmptyRequest) (*StudioPublicInfoResponse, error)
	//*
	// Manages a org.
	// @param {OrganizationManagerRequest} request - The request object containing the org details.
	// @returns {OrganizationResponse} The response object containing the result of the org management actions.
	OrganizationManager(context.Context, *OrganizationManagerRequest) (*OrganizationResponse, error)
	//*
	// Retrieves a org.
	// @param {OrganizationGetterRequest} request - The request object containing the org ID.
	// @returns {OrganizationResponse} The response object containing the retrieved org.
	OrganizationGetter(context.Context, *OrganizationGetterRequest) (*OrganizationResponse, error)
	//*
	// Retrieves information about theaccount.
	// @param {AccountManager} request - The account creation request.
	// @returns {AccountResponse} - The account creation response.
	AccountGetter(context.Context, *v1alpha1.EmptyRequest) (*AccountResponse, error)
	//*
	// Creates a new account with the specified name.
	// @param {AccountManager} request - The account creation request.
	// @returns {AccountResponse} - The account creation response.
	AccountManager(context.Context, *AccountManagerRequest) (*AccountResponse, error)
	mustEmbedUnimplementedStudioServiceServer()
}

// UnimplementedStudioServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStudioServiceServer struct {
}

func (UnimplementedStudioServiceServer) StudioPublicInfo(context.Context, *v1alpha1.EmptyRequest) (*StudioPublicInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StudioPublicInfo not implemented")
}
func (UnimplementedStudioServiceServer) OrganizationManager(context.Context, *OrganizationManagerRequest) (*OrganizationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrganizationManager not implemented")
}
func (UnimplementedStudioServiceServer) OrganizationGetter(context.Context, *OrganizationGetterRequest) (*OrganizationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrganizationGetter not implemented")
}
func (UnimplementedStudioServiceServer) AccountGetter(context.Context, *v1alpha1.EmptyRequest) (*AccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AccountGetter not implemented")
}
func (UnimplementedStudioServiceServer) AccountManager(context.Context, *AccountManagerRequest) (*AccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AccountManager not implemented")
}
func (UnimplementedStudioServiceServer) mustEmbedUnimplementedStudioServiceServer() {}

// UnsafeStudioServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StudioServiceServer will
// result in compilation errors.
type UnsafeStudioServiceServer interface {
	mustEmbedUnimplementedStudioServiceServer()
}

func RegisterStudioServiceServer(s grpc.ServiceRegistrar, srv StudioServiceServer) {
	s.RegisterService(&StudioService_ServiceDesc, srv)
}

func _StudioService_StudioPublicInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1alpha1.EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StudioServiceServer).StudioPublicInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StudioService_StudioPublicInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StudioServiceServer).StudioPublicInfo(ctx, req.(*v1alpha1.EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StudioService_OrganizationManager_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrganizationManagerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StudioServiceServer).OrganizationManager(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StudioService_OrganizationManager_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StudioServiceServer).OrganizationManager(ctx, req.(*OrganizationManagerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StudioService_OrganizationGetter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrganizationGetterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StudioServiceServer).OrganizationGetter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StudioService_OrganizationGetter_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StudioServiceServer).OrganizationGetter(ctx, req.(*OrganizationGetterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StudioService_AccountGetter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1alpha1.EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StudioServiceServer).AccountGetter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StudioService_AccountGetter_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StudioServiceServer).AccountGetter(ctx, req.(*v1alpha1.EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StudioService_AccountManager_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccountManagerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StudioServiceServer).AccountManager(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StudioService_AccountManager_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StudioServiceServer).AccountManager(ctx, req.(*AccountManagerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// StudioService_ServiceDesc is the grpc.ServiceDesc for StudioService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StudioService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "vapusai.v1alpha1.StudioService",
	HandlerType: (*StudioServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StudioPublicInfo",
			Handler:    _StudioService_StudioPublicInfo_Handler,
		},
		{
			MethodName: "OrganizationManager",
			Handler:    _StudioService_OrganizationManager_Handler,
		},
		{
			MethodName: "OrganizationGetter",
			Handler:    _StudioService_OrganizationGetter_Handler,
		},
		{
			MethodName: "AccountGetter",
			Handler:    _StudioService_AccountGetter_Handler,
		},
		{
			MethodName: "AccountManager",
			Handler:    _StudioService_AccountManager_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/vapusai-studio/v1alpha1/organization.proto",
}

const (
	UtilityService_StoreDMSecrets_FullMethodName = "/vapusai.v1alpha1.UtilityService/StoreDMSecrets"
	UtilityService_Upload_FullMethodName         = "/vapusai.v1alpha1.UtilityService/Upload"
	UtilityService_UploadStream_FullMethodName   = "/vapusai.v1alpha1.UtilityService/UploadStream"
)

// UtilityServiceClient is the client API for UtilityService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UtilityServiceClient interface {
	//*
	// Stores the specified data marketplace secrets.
	// @param {StoreDMSecretsRequest} request - The data marketplace secrets storage request.
	// @returns {StoreDMSecretsResponse} - The data marketplace secrets storage response.
	StoreDMSecrets(ctx context.Context, in *StoreDMSecretsRequest, opts ...grpc.CallOption) (*StoreDMSecretsResponse, error)
	Upload(ctx context.Context, in *UploadRequest, opts ...grpc.CallOption) (*UploadResponse, error)
	UploadStream(ctx context.Context, opts ...grpc.CallOption) (UtilityService_UploadStreamClient, error)
}

type utilityServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUtilityServiceClient(cc grpc.ClientConnInterface) UtilityServiceClient {
	return &utilityServiceClient{cc}
}

func (c *utilityServiceClient) StoreDMSecrets(ctx context.Context, in *StoreDMSecretsRequest, opts ...grpc.CallOption) (*StoreDMSecretsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StoreDMSecretsResponse)
	err := c.cc.Invoke(ctx, UtilityService_StoreDMSecrets_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *utilityServiceClient) Upload(ctx context.Context, in *UploadRequest, opts ...grpc.CallOption) (*UploadResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UploadResponse)
	err := c.cc.Invoke(ctx, UtilityService_Upload_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *utilityServiceClient) UploadStream(ctx context.Context, opts ...grpc.CallOption) (UtilityService_UploadStreamClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &UtilityService_ServiceDesc.Streams[0], UtilityService_UploadStream_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &utilityServiceUploadStreamClient{ClientStream: stream}
	return x, nil
}

type UtilityService_UploadStreamClient interface {
	Send(*UploadRequest) error
	CloseAndRecv() (*UploadResponse, error)
	grpc.ClientStream
}

type utilityServiceUploadStreamClient struct {
	grpc.ClientStream
}

func (x *utilityServiceUploadStreamClient) Send(m *UploadRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *utilityServiceUploadStreamClient) CloseAndRecv() (*UploadResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// UtilityServiceServer is the server API for UtilityService service.
// All implementations must embed UnimplementedUtilityServiceServer
// for forward compatibility
type UtilityServiceServer interface {
	//*
	// Stores the specified data marketplace secrets.
	// @param {StoreDMSecretsRequest} request - The data marketplace secrets storage request.
	// @returns {StoreDMSecretsResponse} - The data marketplace secrets storage response.
	StoreDMSecrets(context.Context, *StoreDMSecretsRequest) (*StoreDMSecretsResponse, error)
	Upload(context.Context, *UploadRequest) (*UploadResponse, error)
	UploadStream(UtilityService_UploadStreamServer) error
	mustEmbedUnimplementedUtilityServiceServer()
}

// UnimplementedUtilityServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUtilityServiceServer struct {
}

func (UnimplementedUtilityServiceServer) StoreDMSecrets(context.Context, *StoreDMSecretsRequest) (*StoreDMSecretsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreDMSecrets not implemented")
}
func (UnimplementedUtilityServiceServer) Upload(context.Context, *UploadRequest) (*UploadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (UnimplementedUtilityServiceServer) UploadStream(UtilityService_UploadStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadStream not implemented")
}
func (UnimplementedUtilityServiceServer) mustEmbedUnimplementedUtilityServiceServer() {}

// UnsafeUtilityServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UtilityServiceServer will
// result in compilation errors.
type UnsafeUtilityServiceServer interface {
	mustEmbedUnimplementedUtilityServiceServer()
}

func RegisterUtilityServiceServer(s grpc.ServiceRegistrar, srv UtilityServiceServer) {
	s.RegisterService(&UtilityService_ServiceDesc, srv)
}

func _UtilityService_StoreDMSecrets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreDMSecretsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UtilityServiceServer).StoreDMSecrets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UtilityService_StoreDMSecrets_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UtilityServiceServer).StoreDMSecrets(ctx, req.(*StoreDMSecretsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UtilityService_Upload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UtilityServiceServer).Upload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UtilityService_Upload_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UtilityServiceServer).Upload(ctx, req.(*UploadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UtilityService_UploadStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(UtilityServiceServer).UploadStream(&utilityServiceUploadStreamServer{ServerStream: stream})
}

type UtilityService_UploadStreamServer interface {
	SendAndClose(*UploadResponse) error
	Recv() (*UploadRequest, error)
	grpc.ServerStream
}

type utilityServiceUploadStreamServer struct {
	grpc.ServerStream
}

func (x *utilityServiceUploadStreamServer) SendAndClose(m *UploadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *utilityServiceUploadStreamServer) Recv() (*UploadRequest, error) {
	m := new(UploadRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// UtilityService_ServiceDesc is the grpc.ServiceDesc for UtilityService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UtilityService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "vapusai.v1alpha1.UtilityService",
	HandlerType: (*UtilityServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StoreDMSecrets",
			Handler:    _UtilityService_StoreDMSecrets_Handler,
		},
		{
			MethodName: "Upload",
			Handler:    _UtilityService_Upload_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadStream",
			Handler:       _UtilityService_UploadStream_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "protos/vapusai-studio/v1alpha1/organization.proto",
}

const (
	PluginService_PluginManager_FullMethodName = "/vapusai.v1alpha1.PluginService/PluginManager"
	PluginService_PluginGetter_FullMethodName  = "/vapusai.v1alpha1.PluginService/PluginGetter"
	PluginService_PluginAction_FullMethodName  = "/vapusai.v1alpha1.PluginService/PluginAction"
)

// PluginServiceClient is the client API for PluginService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PluginServiceClient interface {
	PluginManager(ctx context.Context, in *PluginManagerRequest, opts ...grpc.CallOption) (*PluginResponse, error)
	PluginGetter(ctx context.Context, in *PluginGetterRequest, opts ...grpc.CallOption) (*PluginResponse, error)
	PluginAction(ctx context.Context, in *PluginActionRequest, opts ...grpc.CallOption) (*PluginActionResponse, error)
}

type pluginServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPluginServiceClient(cc grpc.ClientConnInterface) PluginServiceClient {
	return &pluginServiceClient{cc}
}

func (c *pluginServiceClient) PluginManager(ctx context.Context, in *PluginManagerRequest, opts ...grpc.CallOption) (*PluginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PluginResponse)
	err := c.cc.Invoke(ctx, PluginService_PluginManager_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginServiceClient) PluginGetter(ctx context.Context, in *PluginGetterRequest, opts ...grpc.CallOption) (*PluginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PluginResponse)
	err := c.cc.Invoke(ctx, PluginService_PluginGetter_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginServiceClient) PluginAction(ctx context.Context, in *PluginActionRequest, opts ...grpc.CallOption) (*PluginActionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PluginActionResponse)
	err := c.cc.Invoke(ctx, PluginService_PluginAction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PluginServiceServer is the server API for PluginService service.
// All implementations must embed UnimplementedPluginServiceServer
// for forward compatibility
type PluginServiceServer interface {
	PluginManager(context.Context, *PluginManagerRequest) (*PluginResponse, error)
	PluginGetter(context.Context, *PluginGetterRequest) (*PluginResponse, error)
	PluginAction(context.Context, *PluginActionRequest) (*PluginActionResponse, error)
	mustEmbedUnimplementedPluginServiceServer()
}

// UnimplementedPluginServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPluginServiceServer struct {
}

func (UnimplementedPluginServiceServer) PluginManager(context.Context, *PluginManagerRequest) (*PluginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PluginManager not implemented")
}
func (UnimplementedPluginServiceServer) PluginGetter(context.Context, *PluginGetterRequest) (*PluginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PluginGetter not implemented")
}
func (UnimplementedPluginServiceServer) PluginAction(context.Context, *PluginActionRequest) (*PluginActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PluginAction not implemented")
}
func (UnimplementedPluginServiceServer) mustEmbedUnimplementedPluginServiceServer() {}

// UnsafePluginServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PluginServiceServer will
// result in compilation errors.
type UnsafePluginServiceServer interface {
	mustEmbedUnimplementedPluginServiceServer()
}

func RegisterPluginServiceServer(s grpc.ServiceRegistrar, srv PluginServiceServer) {
	s.RegisterService(&PluginService_ServiceDesc, srv)
}

func _PluginService_PluginManager_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PluginManagerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServiceServer).PluginManager(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PluginService_PluginManager_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServiceServer).PluginManager(ctx, req.(*PluginManagerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PluginService_PluginGetter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PluginGetterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServiceServer).PluginGetter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PluginService_PluginGetter_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServiceServer).PluginGetter(ctx, req.(*PluginGetterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PluginService_PluginAction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PluginActionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServiceServer).PluginAction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PluginService_PluginAction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServiceServer).PluginAction(ctx, req.(*PluginActionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PluginService_ServiceDesc is the grpc.ServiceDesc for PluginService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PluginService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "vapusai.v1alpha1.PluginService",
	HandlerType: (*PluginServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PluginManager",
			Handler:    _PluginService_PluginManager_Handler,
		},
		{
			MethodName: "PluginGetter",
			Handler:    _PluginService_PluginGetter_Handler,
		},
		{
			MethodName: "PluginAction",
			Handler:    _PluginService_PluginAction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/vapusai-studio/v1alpha1/organization.proto",
}
