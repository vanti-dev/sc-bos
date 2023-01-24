// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.3
// source: services.proto

package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ServicesApiClient is the client API for ServicesApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServicesApiClient interface {
	GetService(ctx context.Context, in *GetServiceRequest, opts ...grpc.CallOption) (*Service, error)
	PullService(ctx context.Context, in *PullServiceRequest, opts ...grpc.CallOption) (ServicesApi_PullServiceClient, error)
	CreateService(ctx context.Context, in *CreateServiceRequest, opts ...grpc.CallOption) (*Service, error)
	DeleteService(ctx context.Context, in *DeleteServiceRequest, opts ...grpc.CallOption) (*Service, error)
	ListServices(ctx context.Context, in *ListServicesRequest, opts ...grpc.CallOption) (*ListServicesResponse, error)
	PullServices(ctx context.Context, in *PullServicesRequest, opts ...grpc.CallOption) (ServicesApi_PullServicesClient, error)
	StartService(ctx context.Context, in *StartServiceRequest, opts ...grpc.CallOption) (*Service, error)
	ConfigureService(ctx context.Context, in *ConfigureServiceRequest, opts ...grpc.CallOption) (*Service, error)
	StopService(ctx context.Context, in *StopServiceRequest, opts ...grpc.CallOption) (*Service, error)
}

type servicesApiClient struct {
	cc grpc.ClientConnInterface
}

func NewServicesApiClient(cc grpc.ClientConnInterface) ServicesApiClient {
	return &servicesApiClient{cc}
}

func (c *servicesApiClient) GetService(ctx context.Context, in *GetServiceRequest, opts ...grpc.CallOption) (*Service, error) {
	out := new(Service)
	err := c.cc.Invoke(ctx, "/smartcore.bos.ServicesApi/GetService", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesApiClient) PullService(ctx context.Context, in *PullServiceRequest, opts ...grpc.CallOption) (ServicesApi_PullServiceClient, error) {
	stream, err := c.cc.NewStream(ctx, &ServicesApi_ServiceDesc.Streams[0], "/smartcore.bos.ServicesApi/PullService", opts...)
	if err != nil {
		return nil, err
	}
	x := &servicesApiPullServiceClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ServicesApi_PullServiceClient interface {
	Recv() (*PullServiceResponse, error)
	grpc.ClientStream
}

type servicesApiPullServiceClient struct {
	grpc.ClientStream
}

func (x *servicesApiPullServiceClient) Recv() (*PullServiceResponse, error) {
	m := new(PullServiceResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *servicesApiClient) CreateService(ctx context.Context, in *CreateServiceRequest, opts ...grpc.CallOption) (*Service, error) {
	out := new(Service)
	err := c.cc.Invoke(ctx, "/smartcore.bos.ServicesApi/CreateService", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesApiClient) DeleteService(ctx context.Context, in *DeleteServiceRequest, opts ...grpc.CallOption) (*Service, error) {
	out := new(Service)
	err := c.cc.Invoke(ctx, "/smartcore.bos.ServicesApi/DeleteService", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesApiClient) ListServices(ctx context.Context, in *ListServicesRequest, opts ...grpc.CallOption) (*ListServicesResponse, error) {
	out := new(ListServicesResponse)
	err := c.cc.Invoke(ctx, "/smartcore.bos.ServicesApi/ListServices", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesApiClient) PullServices(ctx context.Context, in *PullServicesRequest, opts ...grpc.CallOption) (ServicesApi_PullServicesClient, error) {
	stream, err := c.cc.NewStream(ctx, &ServicesApi_ServiceDesc.Streams[1], "/smartcore.bos.ServicesApi/PullServices", opts...)
	if err != nil {
		return nil, err
	}
	x := &servicesApiPullServicesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ServicesApi_PullServicesClient interface {
	Recv() (*PullServicesResponse, error)
	grpc.ClientStream
}

type servicesApiPullServicesClient struct {
	grpc.ClientStream
}

func (x *servicesApiPullServicesClient) Recv() (*PullServicesResponse, error) {
	m := new(PullServicesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *servicesApiClient) StartService(ctx context.Context, in *StartServiceRequest, opts ...grpc.CallOption) (*Service, error) {
	out := new(Service)
	err := c.cc.Invoke(ctx, "/smartcore.bos.ServicesApi/StartService", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesApiClient) ConfigureService(ctx context.Context, in *ConfigureServiceRequest, opts ...grpc.CallOption) (*Service, error) {
	out := new(Service)
	err := c.cc.Invoke(ctx, "/smartcore.bos.ServicesApi/ConfigureService", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesApiClient) StopService(ctx context.Context, in *StopServiceRequest, opts ...grpc.CallOption) (*Service, error) {
	out := new(Service)
	err := c.cc.Invoke(ctx, "/smartcore.bos.ServicesApi/StopService", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServicesApiServer is the server API for ServicesApi service.
// All implementations must embed UnimplementedServicesApiServer
// for forward compatibility
type ServicesApiServer interface {
	GetService(context.Context, *GetServiceRequest) (*Service, error)
	PullService(*PullServiceRequest, ServicesApi_PullServiceServer) error
	CreateService(context.Context, *CreateServiceRequest) (*Service, error)
	DeleteService(context.Context, *DeleteServiceRequest) (*Service, error)
	ListServices(context.Context, *ListServicesRequest) (*ListServicesResponse, error)
	PullServices(*PullServicesRequest, ServicesApi_PullServicesServer) error
	StartService(context.Context, *StartServiceRequest) (*Service, error)
	ConfigureService(context.Context, *ConfigureServiceRequest) (*Service, error)
	StopService(context.Context, *StopServiceRequest) (*Service, error)
	mustEmbedUnimplementedServicesApiServer()
}

// UnimplementedServicesApiServer must be embedded to have forward compatible implementations.
type UnimplementedServicesApiServer struct {
}

func (UnimplementedServicesApiServer) GetService(context.Context, *GetServiceRequest) (*Service, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetService not implemented")
}
func (UnimplementedServicesApiServer) PullService(*PullServiceRequest, ServicesApi_PullServiceServer) error {
	return status.Errorf(codes.Unimplemented, "method PullService not implemented")
}
func (UnimplementedServicesApiServer) CreateService(context.Context, *CreateServiceRequest) (*Service, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateService not implemented")
}
func (UnimplementedServicesApiServer) DeleteService(context.Context, *DeleteServiceRequest) (*Service, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteService not implemented")
}
func (UnimplementedServicesApiServer) ListServices(context.Context, *ListServicesRequest) (*ListServicesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListServices not implemented")
}
func (UnimplementedServicesApiServer) PullServices(*PullServicesRequest, ServicesApi_PullServicesServer) error {
	return status.Errorf(codes.Unimplemented, "method PullServices not implemented")
}
func (UnimplementedServicesApiServer) StartService(context.Context, *StartServiceRequest) (*Service, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartService not implemented")
}
func (UnimplementedServicesApiServer) ConfigureService(context.Context, *ConfigureServiceRequest) (*Service, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfigureService not implemented")
}
func (UnimplementedServicesApiServer) StopService(context.Context, *StopServiceRequest) (*Service, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopService not implemented")
}
func (UnimplementedServicesApiServer) mustEmbedUnimplementedServicesApiServer() {}

// UnsafeServicesApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServicesApiServer will
// result in compilation errors.
type UnsafeServicesApiServer interface {
	mustEmbedUnimplementedServicesApiServer()
}

func RegisterServicesApiServer(s grpc.ServiceRegistrar, srv ServicesApiServer) {
	s.RegisterService(&ServicesApi_ServiceDesc, srv)
}

func _ServicesApi_GetService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesApiServer).GetService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smartcore.bos.ServicesApi/GetService",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesApiServer).GetService(ctx, req.(*GetServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServicesApi_PullService_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PullServiceRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ServicesApiServer).PullService(m, &servicesApiPullServiceServer{stream})
}

type ServicesApi_PullServiceServer interface {
	Send(*PullServiceResponse) error
	grpc.ServerStream
}

type servicesApiPullServiceServer struct {
	grpc.ServerStream
}

func (x *servicesApiPullServiceServer) Send(m *PullServiceResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _ServicesApi_CreateService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesApiServer).CreateService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smartcore.bos.ServicesApi/CreateService",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesApiServer).CreateService(ctx, req.(*CreateServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServicesApi_DeleteService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesApiServer).DeleteService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smartcore.bos.ServicesApi/DeleteService",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesApiServer).DeleteService(ctx, req.(*DeleteServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServicesApi_ListServices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListServicesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesApiServer).ListServices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smartcore.bos.ServicesApi/ListServices",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesApiServer).ListServices(ctx, req.(*ListServicesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServicesApi_PullServices_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PullServicesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ServicesApiServer).PullServices(m, &servicesApiPullServicesServer{stream})
}

type ServicesApi_PullServicesServer interface {
	Send(*PullServicesResponse) error
	grpc.ServerStream
}

type servicesApiPullServicesServer struct {
	grpc.ServerStream
}

func (x *servicesApiPullServicesServer) Send(m *PullServicesResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _ServicesApi_StartService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesApiServer).StartService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smartcore.bos.ServicesApi/StartService",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesApiServer).StartService(ctx, req.(*StartServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServicesApi_ConfigureService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigureServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesApiServer).ConfigureService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smartcore.bos.ServicesApi/ConfigureService",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesApiServer).ConfigureService(ctx, req.(*ConfigureServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServicesApi_StopService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesApiServer).StopService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smartcore.bos.ServicesApi/StopService",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesApiServer).StopService(ctx, req.(*StopServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ServicesApi_ServiceDesc is the grpc.ServiceDesc for ServicesApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ServicesApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "smartcore.bos.ServicesApi",
	HandlerType: (*ServicesApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetService",
			Handler:    _ServicesApi_GetService_Handler,
		},
		{
			MethodName: "CreateService",
			Handler:    _ServicesApi_CreateService_Handler,
		},
		{
			MethodName: "DeleteService",
			Handler:    _ServicesApi_DeleteService_Handler,
		},
		{
			MethodName: "ListServices",
			Handler:    _ServicesApi_ListServices_Handler,
		},
		{
			MethodName: "StartService",
			Handler:    _ServicesApi_StartService_Handler,
		},
		{
			MethodName: "ConfigureService",
			Handler:    _ServicesApi_ConfigureService_Handler,
		},
		{
			MethodName: "StopService",
			Handler:    _ServicesApi_StopService_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PullService",
			Handler:       _ServicesApi_PullService_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "PullServices",
			Handler:       _ServicesApi_PullServices_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "services.proto",
}
