// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.28.2
// source: access.proto

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

// AccessApiClient is the client API for AccessApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccessApiClient interface {
	GetLastAccessAttempt(ctx context.Context, in *GetLastAccessAttemptRequest, opts ...grpc.CallOption) (*AccessAttempt, error)
	PullAccessAttempts(ctx context.Context, in *PullAccessAttemptsRequest, opts ...grpc.CallOption) (AccessApi_PullAccessAttemptsClient, error)
}

type accessApiClient struct {
	cc grpc.ClientConnInterface
}

func NewAccessApiClient(cc grpc.ClientConnInterface) AccessApiClient {
	return &accessApiClient{cc}
}

func (c *accessApiClient) GetLastAccessAttempt(ctx context.Context, in *GetLastAccessAttemptRequest, opts ...grpc.CallOption) (*AccessAttempt, error) {
	out := new(AccessAttempt)
	err := c.cc.Invoke(ctx, "/smartcore.bos.AccessApi/GetLastAccessAttempt", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessApiClient) PullAccessAttempts(ctx context.Context, in *PullAccessAttemptsRequest, opts ...grpc.CallOption) (AccessApi_PullAccessAttemptsClient, error) {
	stream, err := c.cc.NewStream(ctx, &AccessApi_ServiceDesc.Streams[0], "/smartcore.bos.AccessApi/PullAccessAttempts", opts...)
	if err != nil {
		return nil, err
	}
	x := &accessApiPullAccessAttemptsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type AccessApi_PullAccessAttemptsClient interface {
	Recv() (*PullAccessAttemptsResponse, error)
	grpc.ClientStream
}

type accessApiPullAccessAttemptsClient struct {
	grpc.ClientStream
}

func (x *accessApiPullAccessAttemptsClient) Recv() (*PullAccessAttemptsResponse, error) {
	m := new(PullAccessAttemptsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AccessApiServer is the server API for AccessApi service.
// All implementations must embed UnimplementedAccessApiServer
// for forward compatibility
type AccessApiServer interface {
	GetLastAccessAttempt(context.Context, *GetLastAccessAttemptRequest) (*AccessAttempt, error)
	PullAccessAttempts(*PullAccessAttemptsRequest, AccessApi_PullAccessAttemptsServer) error
	mustEmbedUnimplementedAccessApiServer()
}

// UnimplementedAccessApiServer must be embedded to have forward compatible implementations.
type UnimplementedAccessApiServer struct {
}

func (UnimplementedAccessApiServer) GetLastAccessAttempt(context.Context, *GetLastAccessAttemptRequest) (*AccessAttempt, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLastAccessAttempt not implemented")
}
func (UnimplementedAccessApiServer) PullAccessAttempts(*PullAccessAttemptsRequest, AccessApi_PullAccessAttemptsServer) error {
	return status.Errorf(codes.Unimplemented, "method PullAccessAttempts not implemented")
}
func (UnimplementedAccessApiServer) mustEmbedUnimplementedAccessApiServer() {}

// UnsafeAccessApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccessApiServer will
// result in compilation errors.
type UnsafeAccessApiServer interface {
	mustEmbedUnimplementedAccessApiServer()
}

func RegisterAccessApiServer(s grpc.ServiceRegistrar, srv AccessApiServer) {
	s.RegisterService(&AccessApi_ServiceDesc, srv)
}

func _AccessApi_GetLastAccessAttempt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLastAccessAttemptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessApiServer).GetLastAccessAttempt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smartcore.bos.AccessApi/GetLastAccessAttempt",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessApiServer).GetLastAccessAttempt(ctx, req.(*GetLastAccessAttemptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessApi_PullAccessAttempts_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PullAccessAttemptsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AccessApiServer).PullAccessAttempts(m, &accessApiPullAccessAttemptsServer{stream})
}

type AccessApi_PullAccessAttemptsServer interface {
	Send(*PullAccessAttemptsResponse) error
	grpc.ServerStream
}

type accessApiPullAccessAttemptsServer struct {
	grpc.ServerStream
}

func (x *accessApiPullAccessAttemptsServer) Send(m *PullAccessAttemptsResponse) error {
	return x.ServerStream.SendMsg(m)
}

// AccessApi_ServiceDesc is the grpc.ServiceDesc for AccessApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccessApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "smartcore.bos.AccessApi",
	HandlerType: (*AccessApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLastAccessAttempt",
			Handler:    _AccessApi_GetLastAccessAttempt_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PullAccessAttempts",
			Handler:       _AccessApi_PullAccessAttempts_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "access.proto",
}
