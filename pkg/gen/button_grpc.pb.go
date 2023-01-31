// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: button.proto

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

// ButtonApiClient is the client API for ButtonApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ButtonApiClient interface {
	// Gets the current state of the button.
	// Contain the most recent gesture, so clients using polling can still detect and respond to gestures.
	// However, to reduce latency, PullButtonState is recommended for this use case.
	GetButtonState(ctx context.Context, in *GetButtonStateRequest, opts ...grpc.CallOption) (*ButtonState, error)
	// Fetches changes to button Press state and gestures, and optionally the initial state.
	PullButtonState(ctx context.Context, in *PullButtonStateRequest, opts ...grpc.CallOption) (ButtonApi_PullButtonStateClient, error)
}

type buttonApiClient struct {
	cc grpc.ClientConnInterface
}

func NewButtonApiClient(cc grpc.ClientConnInterface) ButtonApiClient {
	return &buttonApiClient{cc}
}

func (c *buttonApiClient) GetButtonState(ctx context.Context, in *GetButtonStateRequest, opts ...grpc.CallOption) (*ButtonState, error) {
	out := new(ButtonState)
	err := c.cc.Invoke(ctx, "/smartcore.bos.ButtonApi/GetButtonState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *buttonApiClient) PullButtonState(ctx context.Context, in *PullButtonStateRequest, opts ...grpc.CallOption) (ButtonApi_PullButtonStateClient, error) {
	stream, err := c.cc.NewStream(ctx, &ButtonApi_ServiceDesc.Streams[0], "/smartcore.bos.ButtonApi/PullButtonState", opts...)
	if err != nil {
		return nil, err
	}
	x := &buttonApiPullButtonStateClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ButtonApi_PullButtonStateClient interface {
	Recv() (*PullButtonStateResponse, error)
	grpc.ClientStream
}

type buttonApiPullButtonStateClient struct {
	grpc.ClientStream
}

func (x *buttonApiPullButtonStateClient) Recv() (*PullButtonStateResponse, error) {
	m := new(PullButtonStateResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ButtonApiServer is the server API for ButtonApi service.
// All implementations must embed UnimplementedButtonApiServer
// for forward compatibility
type ButtonApiServer interface {
	// Gets the current state of the button.
	// Contain the most recent gesture, so clients using polling can still detect and respond to gestures.
	// However, to reduce latency, PullButtonState is recommended for this use case.
	GetButtonState(context.Context, *GetButtonStateRequest) (*ButtonState, error)
	// Fetches changes to button Press state and gestures, and optionally the initial state.
	PullButtonState(*PullButtonStateRequest, ButtonApi_PullButtonStateServer) error
	mustEmbedUnimplementedButtonApiServer()
}

// UnimplementedButtonApiServer must be embedded to have forward compatible implementations.
type UnimplementedButtonApiServer struct {
}

func (UnimplementedButtonApiServer) GetButtonState(context.Context, *GetButtonStateRequest) (*ButtonState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetButtonState not implemented")
}
func (UnimplementedButtonApiServer) PullButtonState(*PullButtonStateRequest, ButtonApi_PullButtonStateServer) error {
	return status.Errorf(codes.Unimplemented, "method PullButtonState not implemented")
}
func (UnimplementedButtonApiServer) mustEmbedUnimplementedButtonApiServer() {}

// UnsafeButtonApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ButtonApiServer will
// result in compilation errors.
type UnsafeButtonApiServer interface {
	mustEmbedUnimplementedButtonApiServer()
}

func RegisterButtonApiServer(s grpc.ServiceRegistrar, srv ButtonApiServer) {
	s.RegisterService(&ButtonApi_ServiceDesc, srv)
}

func _ButtonApi_GetButtonState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetButtonStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ButtonApiServer).GetButtonState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smartcore.bos.ButtonApi/GetButtonState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ButtonApiServer).GetButtonState(ctx, req.(*GetButtonStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ButtonApi_PullButtonState_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PullButtonStateRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ButtonApiServer).PullButtonState(m, &buttonApiPullButtonStateServer{stream})
}

type ButtonApi_PullButtonStateServer interface {
	Send(*PullButtonStateResponse) error
	grpc.ServerStream
}

type buttonApiPullButtonStateServer struct {
	grpc.ServerStream
}

func (x *buttonApiPullButtonStateServer) Send(m *PullButtonStateResponse) error {
	return x.ServerStream.SendMsg(m)
}

// ButtonApi_ServiceDesc is the grpc.ServiceDesc for ButtonApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ButtonApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "smartcore.bos.ButtonApi",
	HandlerType: (*ButtonApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetButtonState",
			Handler:    _ButtonApi_GetButtonState_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PullButtonState",
			Handler:       _ButtonApi_PullButtonState_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "button.proto",
}
