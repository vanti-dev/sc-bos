// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: tenants.proto

package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	TenantApi_ListTenants_FullMethodName       = "/smartcore.bos.tenants.TenantApi/ListTenants"
	TenantApi_PullTenants_FullMethodName       = "/smartcore.bos.tenants.TenantApi/PullTenants"
	TenantApi_CreateTenant_FullMethodName      = "/smartcore.bos.tenants.TenantApi/CreateTenant"
	TenantApi_GetTenant_FullMethodName         = "/smartcore.bos.tenants.TenantApi/GetTenant"
	TenantApi_UpdateTenant_FullMethodName      = "/smartcore.bos.tenants.TenantApi/UpdateTenant"
	TenantApi_DeleteTenant_FullMethodName      = "/smartcore.bos.tenants.TenantApi/DeleteTenant"
	TenantApi_PullTenant_FullMethodName        = "/smartcore.bos.tenants.TenantApi/PullTenant"
	TenantApi_AddTenantZones_FullMethodName    = "/smartcore.bos.tenants.TenantApi/AddTenantZones"
	TenantApi_RemoveTenantZones_FullMethodName = "/smartcore.bos.tenants.TenantApi/RemoveTenantZones"
	TenantApi_ListSecrets_FullMethodName       = "/smartcore.bos.tenants.TenantApi/ListSecrets"
	TenantApi_PullSecrets_FullMethodName       = "/smartcore.bos.tenants.TenantApi/PullSecrets"
	TenantApi_CreateSecret_FullMethodName      = "/smartcore.bos.tenants.TenantApi/CreateSecret"
	TenantApi_VerifySecret_FullMethodName      = "/smartcore.bos.tenants.TenantApi/VerifySecret"
	TenantApi_GetSecret_FullMethodName         = "/smartcore.bos.tenants.TenantApi/GetSecret"
	TenantApi_UpdateSecret_FullMethodName      = "/smartcore.bos.tenants.TenantApi/UpdateSecret"
	TenantApi_DeleteSecret_FullMethodName      = "/smartcore.bos.tenants.TenantApi/DeleteSecret"
	TenantApi_PullSecret_FullMethodName        = "/smartcore.bos.tenants.TenantApi/PullSecret"
	TenantApi_RegenerateSecret_FullMethodName  = "/smartcore.bos.tenants.TenantApi/RegenerateSecret"
)

// TenantApiClient is the client API for TenantApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TenantApiClient interface {
	ListTenants(ctx context.Context, in *ListTenantsRequest, opts ...grpc.CallOption) (*ListTenantsResponse, error)
	PullTenants(ctx context.Context, in *PullTenantsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[PullTenantsResponse], error)
	CreateTenant(ctx context.Context, in *CreateTenantRequest, opts ...grpc.CallOption) (*Tenant, error)
	GetTenant(ctx context.Context, in *GetTenantRequest, opts ...grpc.CallOption) (*Tenant, error)
	UpdateTenant(ctx context.Context, in *UpdateTenantRequest, opts ...grpc.CallOption) (*Tenant, error)
	DeleteTenant(ctx context.Context, in *DeleteTenantRequest, opts ...grpc.CallOption) (*DeleteTenantResponse, error)
	PullTenant(ctx context.Context, in *PullTenantRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[PullTenantResponse], error)
	AddTenantZones(ctx context.Context, in *AddTenantZonesRequest, opts ...grpc.CallOption) (*Tenant, error)
	RemoveTenantZones(ctx context.Context, in *RemoveTenantZonesRequest, opts ...grpc.CallOption) (*Tenant, error)
	ListSecrets(ctx context.Context, in *ListSecretsRequest, opts ...grpc.CallOption) (*ListSecretsResponse, error)
	PullSecrets(ctx context.Context, in *PullSecretsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[PullSecretsResponse], error)
	CreateSecret(ctx context.Context, in *CreateSecretRequest, opts ...grpc.CallOption) (*Secret, error)
	// VerifySecret validates that the given tenant_id has a secret that matches the given secret, returning that secret.
	// An Unauthenticated error will be returned if the tenant+secret do not match or are not known.
	VerifySecret(ctx context.Context, in *VerifySecretRequest, opts ...grpc.CallOption) (*Secret, error)
	GetSecret(ctx context.Context, in *GetSecretRequest, opts ...grpc.CallOption) (*Secret, error)
	UpdateSecret(ctx context.Context, in *UpdateSecretRequest, opts ...grpc.CallOption) (*Secret, error)
	DeleteSecret(ctx context.Context, in *DeleteSecretRequest, opts ...grpc.CallOption) (*DeleteSecretResponse, error)
	PullSecret(ctx context.Context, in *PullSecretRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[PullSecretResponse], error)
	// Creates a new hash for the secret, returns that hash. Just like CreateSecret but using an existing secret
	RegenerateSecret(ctx context.Context, in *RegenerateSecretRequest, opts ...grpc.CallOption) (*Secret, error)
}

type tenantApiClient struct {
	cc grpc.ClientConnInterface
}

func NewTenantApiClient(cc grpc.ClientConnInterface) TenantApiClient {
	return &tenantApiClient{cc}
}

func (c *tenantApiClient) ListTenants(ctx context.Context, in *ListTenantsRequest, opts ...grpc.CallOption) (*ListTenantsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListTenantsResponse)
	err := c.cc.Invoke(ctx, TenantApi_ListTenants_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantApiClient) PullTenants(ctx context.Context, in *PullTenantsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[PullTenantsResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &TenantApi_ServiceDesc.Streams[0], TenantApi_PullTenants_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[PullTenantsRequest, PullTenantsResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type TenantApi_PullTenantsClient = grpc.ServerStreamingClient[PullTenantsResponse]

func (c *tenantApiClient) CreateTenant(ctx context.Context, in *CreateTenantRequest, opts ...grpc.CallOption) (*Tenant, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Tenant)
	err := c.cc.Invoke(ctx, TenantApi_CreateTenant_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantApiClient) GetTenant(ctx context.Context, in *GetTenantRequest, opts ...grpc.CallOption) (*Tenant, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Tenant)
	err := c.cc.Invoke(ctx, TenantApi_GetTenant_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantApiClient) UpdateTenant(ctx context.Context, in *UpdateTenantRequest, opts ...grpc.CallOption) (*Tenant, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Tenant)
	err := c.cc.Invoke(ctx, TenantApi_UpdateTenant_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantApiClient) DeleteTenant(ctx context.Context, in *DeleteTenantRequest, opts ...grpc.CallOption) (*DeleteTenantResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteTenantResponse)
	err := c.cc.Invoke(ctx, TenantApi_DeleteTenant_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantApiClient) PullTenant(ctx context.Context, in *PullTenantRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[PullTenantResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &TenantApi_ServiceDesc.Streams[1], TenantApi_PullTenant_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[PullTenantRequest, PullTenantResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type TenantApi_PullTenantClient = grpc.ServerStreamingClient[PullTenantResponse]

func (c *tenantApiClient) AddTenantZones(ctx context.Context, in *AddTenantZonesRequest, opts ...grpc.CallOption) (*Tenant, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Tenant)
	err := c.cc.Invoke(ctx, TenantApi_AddTenantZones_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantApiClient) RemoveTenantZones(ctx context.Context, in *RemoveTenantZonesRequest, opts ...grpc.CallOption) (*Tenant, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Tenant)
	err := c.cc.Invoke(ctx, TenantApi_RemoveTenantZones_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantApiClient) ListSecrets(ctx context.Context, in *ListSecretsRequest, opts ...grpc.CallOption) (*ListSecretsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListSecretsResponse)
	err := c.cc.Invoke(ctx, TenantApi_ListSecrets_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantApiClient) PullSecrets(ctx context.Context, in *PullSecretsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[PullSecretsResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &TenantApi_ServiceDesc.Streams[2], TenantApi_PullSecrets_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[PullSecretsRequest, PullSecretsResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type TenantApi_PullSecretsClient = grpc.ServerStreamingClient[PullSecretsResponse]

func (c *tenantApiClient) CreateSecret(ctx context.Context, in *CreateSecretRequest, opts ...grpc.CallOption) (*Secret, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Secret)
	err := c.cc.Invoke(ctx, TenantApi_CreateSecret_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantApiClient) VerifySecret(ctx context.Context, in *VerifySecretRequest, opts ...grpc.CallOption) (*Secret, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Secret)
	err := c.cc.Invoke(ctx, TenantApi_VerifySecret_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantApiClient) GetSecret(ctx context.Context, in *GetSecretRequest, opts ...grpc.CallOption) (*Secret, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Secret)
	err := c.cc.Invoke(ctx, TenantApi_GetSecret_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantApiClient) UpdateSecret(ctx context.Context, in *UpdateSecretRequest, opts ...grpc.CallOption) (*Secret, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Secret)
	err := c.cc.Invoke(ctx, TenantApi_UpdateSecret_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantApiClient) DeleteSecret(ctx context.Context, in *DeleteSecretRequest, opts ...grpc.CallOption) (*DeleteSecretResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteSecretResponse)
	err := c.cc.Invoke(ctx, TenantApi_DeleteSecret_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantApiClient) PullSecret(ctx context.Context, in *PullSecretRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[PullSecretResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &TenantApi_ServiceDesc.Streams[3], TenantApi_PullSecret_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[PullSecretRequest, PullSecretResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type TenantApi_PullSecretClient = grpc.ServerStreamingClient[PullSecretResponse]

func (c *tenantApiClient) RegenerateSecret(ctx context.Context, in *RegenerateSecretRequest, opts ...grpc.CallOption) (*Secret, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Secret)
	err := c.cc.Invoke(ctx, TenantApi_RegenerateSecret_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TenantApiServer is the server API for TenantApi service.
// All implementations must embed UnimplementedTenantApiServer
// for forward compatibility.
type TenantApiServer interface {
	ListTenants(context.Context, *ListTenantsRequest) (*ListTenantsResponse, error)
	PullTenants(*PullTenantsRequest, grpc.ServerStreamingServer[PullTenantsResponse]) error
	CreateTenant(context.Context, *CreateTenantRequest) (*Tenant, error)
	GetTenant(context.Context, *GetTenantRequest) (*Tenant, error)
	UpdateTenant(context.Context, *UpdateTenantRequest) (*Tenant, error)
	DeleteTenant(context.Context, *DeleteTenantRequest) (*DeleteTenantResponse, error)
	PullTenant(*PullTenantRequest, grpc.ServerStreamingServer[PullTenantResponse]) error
	AddTenantZones(context.Context, *AddTenantZonesRequest) (*Tenant, error)
	RemoveTenantZones(context.Context, *RemoveTenantZonesRequest) (*Tenant, error)
	ListSecrets(context.Context, *ListSecretsRequest) (*ListSecretsResponse, error)
	PullSecrets(*PullSecretsRequest, grpc.ServerStreamingServer[PullSecretsResponse]) error
	CreateSecret(context.Context, *CreateSecretRequest) (*Secret, error)
	// VerifySecret validates that the given tenant_id has a secret that matches the given secret, returning that secret.
	// An Unauthenticated error will be returned if the tenant+secret do not match or are not known.
	VerifySecret(context.Context, *VerifySecretRequest) (*Secret, error)
	GetSecret(context.Context, *GetSecretRequest) (*Secret, error)
	UpdateSecret(context.Context, *UpdateSecretRequest) (*Secret, error)
	DeleteSecret(context.Context, *DeleteSecretRequest) (*DeleteSecretResponse, error)
	PullSecret(*PullSecretRequest, grpc.ServerStreamingServer[PullSecretResponse]) error
	// Creates a new hash for the secret, returns that hash. Just like CreateSecret but using an existing secret
	RegenerateSecret(context.Context, *RegenerateSecretRequest) (*Secret, error)
	mustEmbedUnimplementedTenantApiServer()
}

// UnimplementedTenantApiServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTenantApiServer struct{}

func (UnimplementedTenantApiServer) ListTenants(context.Context, *ListTenantsRequest) (*ListTenantsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTenants not implemented")
}
func (UnimplementedTenantApiServer) PullTenants(*PullTenantsRequest, grpc.ServerStreamingServer[PullTenantsResponse]) error {
	return status.Errorf(codes.Unimplemented, "method PullTenants not implemented")
}
func (UnimplementedTenantApiServer) CreateTenant(context.Context, *CreateTenantRequest) (*Tenant, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTenant not implemented")
}
func (UnimplementedTenantApiServer) GetTenant(context.Context, *GetTenantRequest) (*Tenant, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTenant not implemented")
}
func (UnimplementedTenantApiServer) UpdateTenant(context.Context, *UpdateTenantRequest) (*Tenant, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTenant not implemented")
}
func (UnimplementedTenantApiServer) DeleteTenant(context.Context, *DeleteTenantRequest) (*DeleteTenantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTenant not implemented")
}
func (UnimplementedTenantApiServer) PullTenant(*PullTenantRequest, grpc.ServerStreamingServer[PullTenantResponse]) error {
	return status.Errorf(codes.Unimplemented, "method PullTenant not implemented")
}
func (UnimplementedTenantApiServer) AddTenantZones(context.Context, *AddTenantZonesRequest) (*Tenant, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddTenantZones not implemented")
}
func (UnimplementedTenantApiServer) RemoveTenantZones(context.Context, *RemoveTenantZonesRequest) (*Tenant, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveTenantZones not implemented")
}
func (UnimplementedTenantApiServer) ListSecrets(context.Context, *ListSecretsRequest) (*ListSecretsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSecrets not implemented")
}
func (UnimplementedTenantApiServer) PullSecrets(*PullSecretsRequest, grpc.ServerStreamingServer[PullSecretsResponse]) error {
	return status.Errorf(codes.Unimplemented, "method PullSecrets not implemented")
}
func (UnimplementedTenantApiServer) CreateSecret(context.Context, *CreateSecretRequest) (*Secret, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSecret not implemented")
}
func (UnimplementedTenantApiServer) VerifySecret(context.Context, *VerifySecretRequest) (*Secret, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifySecret not implemented")
}
func (UnimplementedTenantApiServer) GetSecret(context.Context, *GetSecretRequest) (*Secret, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSecret not implemented")
}
func (UnimplementedTenantApiServer) UpdateSecret(context.Context, *UpdateSecretRequest) (*Secret, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSecret not implemented")
}
func (UnimplementedTenantApiServer) DeleteSecret(context.Context, *DeleteSecretRequest) (*DeleteSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSecret not implemented")
}
func (UnimplementedTenantApiServer) PullSecret(*PullSecretRequest, grpc.ServerStreamingServer[PullSecretResponse]) error {
	return status.Errorf(codes.Unimplemented, "method PullSecret not implemented")
}
func (UnimplementedTenantApiServer) RegenerateSecret(context.Context, *RegenerateSecretRequest) (*Secret, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegenerateSecret not implemented")
}
func (UnimplementedTenantApiServer) mustEmbedUnimplementedTenantApiServer() {}
func (UnimplementedTenantApiServer) testEmbeddedByValue()                   {}

// UnsafeTenantApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TenantApiServer will
// result in compilation errors.
type UnsafeTenantApiServer interface {
	mustEmbedUnimplementedTenantApiServer()
}

func RegisterTenantApiServer(s grpc.ServiceRegistrar, srv TenantApiServer) {
	// If the following call pancis, it indicates UnimplementedTenantApiServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&TenantApi_ServiceDesc, srv)
}

func _TenantApi_ListTenants_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListTenantsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantApiServer).ListTenants(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TenantApi_ListTenants_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantApiServer).ListTenants(ctx, req.(*ListTenantsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TenantApi_PullTenants_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PullTenantsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TenantApiServer).PullTenants(m, &grpc.GenericServerStream[PullTenantsRequest, PullTenantsResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type TenantApi_PullTenantsServer = grpc.ServerStreamingServer[PullTenantsResponse]

func _TenantApi_CreateTenant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTenantRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantApiServer).CreateTenant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TenantApi_CreateTenant_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantApiServer).CreateTenant(ctx, req.(*CreateTenantRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TenantApi_GetTenant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTenantRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantApiServer).GetTenant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TenantApi_GetTenant_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantApiServer).GetTenant(ctx, req.(*GetTenantRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TenantApi_UpdateTenant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTenantRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantApiServer).UpdateTenant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TenantApi_UpdateTenant_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantApiServer).UpdateTenant(ctx, req.(*UpdateTenantRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TenantApi_DeleteTenant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTenantRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantApiServer).DeleteTenant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TenantApi_DeleteTenant_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantApiServer).DeleteTenant(ctx, req.(*DeleteTenantRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TenantApi_PullTenant_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PullTenantRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TenantApiServer).PullTenant(m, &grpc.GenericServerStream[PullTenantRequest, PullTenantResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type TenantApi_PullTenantServer = grpc.ServerStreamingServer[PullTenantResponse]

func _TenantApi_AddTenantZones_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddTenantZonesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantApiServer).AddTenantZones(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TenantApi_AddTenantZones_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantApiServer).AddTenantZones(ctx, req.(*AddTenantZonesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TenantApi_RemoveTenantZones_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveTenantZonesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantApiServer).RemoveTenantZones(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TenantApi_RemoveTenantZones_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantApiServer).RemoveTenantZones(ctx, req.(*RemoveTenantZonesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TenantApi_ListSecrets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSecretsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantApiServer).ListSecrets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TenantApi_ListSecrets_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantApiServer).ListSecrets(ctx, req.(*ListSecretsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TenantApi_PullSecrets_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PullSecretsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TenantApiServer).PullSecrets(m, &grpc.GenericServerStream[PullSecretsRequest, PullSecretsResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type TenantApi_PullSecretsServer = grpc.ServerStreamingServer[PullSecretsResponse]

func _TenantApi_CreateSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantApiServer).CreateSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TenantApi_CreateSecret_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantApiServer).CreateSecret(ctx, req.(*CreateSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TenantApi_VerifySecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifySecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantApiServer).VerifySecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TenantApi_VerifySecret_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantApiServer).VerifySecret(ctx, req.(*VerifySecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TenantApi_GetSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantApiServer).GetSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TenantApi_GetSecret_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantApiServer).GetSecret(ctx, req.(*GetSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TenantApi_UpdateSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantApiServer).UpdateSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TenantApi_UpdateSecret_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantApiServer).UpdateSecret(ctx, req.(*UpdateSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TenantApi_DeleteSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantApiServer).DeleteSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TenantApi_DeleteSecret_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantApiServer).DeleteSecret(ctx, req.(*DeleteSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TenantApi_PullSecret_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PullSecretRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TenantApiServer).PullSecret(m, &grpc.GenericServerStream[PullSecretRequest, PullSecretResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type TenantApi_PullSecretServer = grpc.ServerStreamingServer[PullSecretResponse]

func _TenantApi_RegenerateSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegenerateSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantApiServer).RegenerateSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TenantApi_RegenerateSecret_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantApiServer).RegenerateSecret(ctx, req.(*RegenerateSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TenantApi_ServiceDesc is the grpc.ServiceDesc for TenantApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TenantApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "smartcore.bos.tenants.TenantApi",
	HandlerType: (*TenantApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListTenants",
			Handler:    _TenantApi_ListTenants_Handler,
		},
		{
			MethodName: "CreateTenant",
			Handler:    _TenantApi_CreateTenant_Handler,
		},
		{
			MethodName: "GetTenant",
			Handler:    _TenantApi_GetTenant_Handler,
		},
		{
			MethodName: "UpdateTenant",
			Handler:    _TenantApi_UpdateTenant_Handler,
		},
		{
			MethodName: "DeleteTenant",
			Handler:    _TenantApi_DeleteTenant_Handler,
		},
		{
			MethodName: "AddTenantZones",
			Handler:    _TenantApi_AddTenantZones_Handler,
		},
		{
			MethodName: "RemoveTenantZones",
			Handler:    _TenantApi_RemoveTenantZones_Handler,
		},
		{
			MethodName: "ListSecrets",
			Handler:    _TenantApi_ListSecrets_Handler,
		},
		{
			MethodName: "CreateSecret",
			Handler:    _TenantApi_CreateSecret_Handler,
		},
		{
			MethodName: "VerifySecret",
			Handler:    _TenantApi_VerifySecret_Handler,
		},
		{
			MethodName: "GetSecret",
			Handler:    _TenantApi_GetSecret_Handler,
		},
		{
			MethodName: "UpdateSecret",
			Handler:    _TenantApi_UpdateSecret_Handler,
		},
		{
			MethodName: "DeleteSecret",
			Handler:    _TenantApi_DeleteSecret_Handler,
		},
		{
			MethodName: "RegenerateSecret",
			Handler:    _TenantApi_RegenerateSecret_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PullTenants",
			Handler:       _TenantApi_PullTenants_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "PullTenant",
			Handler:       _TenantApi_PullTenant_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "PullSecrets",
			Handler:       _TenantApi_PullSecrets_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "PullSecret",
			Handler:       _TenantApi_PullSecret_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "tenants.proto",
}
