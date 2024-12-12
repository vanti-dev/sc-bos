// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.3
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
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	AccessApi_GetLastAccessAttempt_FullMethodName = "/smartcore.bos.AccessApi/GetLastAccessAttempt"
	AccessApi_PullAccessAttempts_FullMethodName   = "/smartcore.bos.AccessApi/PullAccessAttempts"
)

// AccessApiClient is the client API for AccessApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// AccessApi describes the capability to manage access to a resource.
// This could be a access card reader next to a door, or a barrier at a car park.
type AccessApiClient interface {
	GetLastAccessAttempt(ctx context.Context, in *GetLastAccessAttemptRequest, opts ...grpc.CallOption) (*AccessAttempt, error)
	PullAccessAttempts(ctx context.Context, in *PullAccessAttemptsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[PullAccessAttemptsResponse], error)
}

type accessApiClient struct {
	cc grpc.ClientConnInterface
}

func NewAccessApiClient(cc grpc.ClientConnInterface) AccessApiClient {
	return &accessApiClient{cc}
}

func (c *accessApiClient) GetLastAccessAttempt(ctx context.Context, in *GetLastAccessAttemptRequest, opts ...grpc.CallOption) (*AccessAttempt, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AccessAttempt)
	err := c.cc.Invoke(ctx, AccessApi_GetLastAccessAttempt_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessApiClient) PullAccessAttempts(ctx context.Context, in *PullAccessAttemptsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[PullAccessAttemptsResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &AccessApi_ServiceDesc.Streams[0], AccessApi_PullAccessAttempts_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[PullAccessAttemptsRequest, PullAccessAttemptsResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type AccessApi_PullAccessAttemptsClient = grpc.ServerStreamingClient[PullAccessAttemptsResponse]

// AccessApiServer is the server API for AccessApi service.
// All implementations must embed UnimplementedAccessApiServer
// for forward compatibility.
//
// AccessApi describes the capability to manage access to a resource.
// This could be a access card reader next to a door, or a barrier at a car park.
type AccessApiServer interface {
	GetLastAccessAttempt(context.Context, *GetLastAccessAttemptRequest) (*AccessAttempt, error)
	PullAccessAttempts(*PullAccessAttemptsRequest, grpc.ServerStreamingServer[PullAccessAttemptsResponse]) error
	mustEmbedUnimplementedAccessApiServer()
}

// UnimplementedAccessApiServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAccessApiServer struct{}

func (UnimplementedAccessApiServer) GetLastAccessAttempt(context.Context, *GetLastAccessAttemptRequest) (*AccessAttempt, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLastAccessAttempt not implemented")
}
func (UnimplementedAccessApiServer) PullAccessAttempts(*PullAccessAttemptsRequest, grpc.ServerStreamingServer[PullAccessAttemptsResponse]) error {
	return status.Errorf(codes.Unimplemented, "method PullAccessAttempts not implemented")
}
func (UnimplementedAccessApiServer) mustEmbedUnimplementedAccessApiServer() {}
func (UnimplementedAccessApiServer) testEmbeddedByValue()                   {}

// UnsafeAccessApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccessApiServer will
// result in compilation errors.
type UnsafeAccessApiServer interface {
	mustEmbedUnimplementedAccessApiServer()
}

func RegisterAccessApiServer(s grpc.ServiceRegistrar, srv AccessApiServer) {
	// If the following call pancis, it indicates UnimplementedAccessApiServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
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
		FullMethod: AccessApi_GetLastAccessAttempt_FullMethodName,
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
	return srv.(AccessApiServer).PullAccessAttempts(m, &grpc.GenericServerStream[PullAccessAttemptsRequest, PullAccessAttemptsResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type AccessApi_PullAccessAttemptsServer = grpc.ServerStreamingServer[PullAccessAttemptsResponse]

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
