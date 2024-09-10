// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: mqtt.proto

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
	MqttService_PullMessages_FullMethodName = "/smartcore.bos.MqttService/PullMessages"
)

// MqttServiceClient is the client API for MqttService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MqttServiceClient interface {
	PullMessages(ctx context.Context, in *PullMessagesRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[PullMessagesResponse], error)
}

type mqttServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMqttServiceClient(cc grpc.ClientConnInterface) MqttServiceClient {
	return &mqttServiceClient{cc}
}

func (c *mqttServiceClient) PullMessages(ctx context.Context, in *PullMessagesRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[PullMessagesResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &MqttService_ServiceDesc.Streams[0], MqttService_PullMessages_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[PullMessagesRequest, PullMessagesResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type MqttService_PullMessagesClient = grpc.ServerStreamingClient[PullMessagesResponse]

// MqttServiceServer is the server API for MqttService service.
// All implementations must embed UnimplementedMqttServiceServer
// for forward compatibility.
type MqttServiceServer interface {
	PullMessages(*PullMessagesRequest, grpc.ServerStreamingServer[PullMessagesResponse]) error
	mustEmbedUnimplementedMqttServiceServer()
}

// UnimplementedMqttServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMqttServiceServer struct{}

func (UnimplementedMqttServiceServer) PullMessages(*PullMessagesRequest, grpc.ServerStreamingServer[PullMessagesResponse]) error {
	return status.Errorf(codes.Unimplemented, "method PullMessages not implemented")
}
func (UnimplementedMqttServiceServer) mustEmbedUnimplementedMqttServiceServer() {}
func (UnimplementedMqttServiceServer) testEmbeddedByValue()                     {}

// UnsafeMqttServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MqttServiceServer will
// result in compilation errors.
type UnsafeMqttServiceServer interface {
	mustEmbedUnimplementedMqttServiceServer()
}

func RegisterMqttServiceServer(s grpc.ServiceRegistrar, srv MqttServiceServer) {
	// If the following call pancis, it indicates UnimplementedMqttServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MqttService_ServiceDesc, srv)
}

func _MqttService_PullMessages_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PullMessagesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MqttServiceServer).PullMessages(m, &grpc.GenericServerStream[PullMessagesRequest, PullMessagesResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type MqttService_PullMessagesServer = grpc.ServerStreamingServer[PullMessagesResponse]

// MqttService_ServiceDesc is the grpc.ServiceDesc for MqttService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MqttService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "smartcore.bos.MqttService",
	HandlerType: (*MqttServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PullMessages",
			Handler:       _MqttService_PullMessages_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "mqtt.proto",
}
