// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.27.1
// source: color.proto

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

// ColorApiClient is the client API for ColorApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ColorApiClient interface {
	// get the current color
	GetColor(ctx context.Context, in *GetColorRequest, opts ...grpc.CallOption) (*Color, error)
	// request that the color be changed
	UpdateColor(ctx context.Context, in *UpdateColorRequest, opts ...grpc.CallOption) (*Color, error)
	// request updates to changes to the color value
	PullColor(ctx context.Context, in *PullColorRequest, opts ...grpc.CallOption) (ColorApi_PullColorClient, error)
}

type colorApiClient struct {
	cc grpc.ClientConnInterface
}

func NewColorApiClient(cc grpc.ClientConnInterface) ColorApiClient {
	return &colorApiClient{cc}
}

func (c *colorApiClient) GetColor(ctx context.Context, in *GetColorRequest, opts ...grpc.CallOption) (*Color, error) {
	out := new(Color)
	err := c.cc.Invoke(ctx, "/smartcore.bos.ColorApi/GetColor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *colorApiClient) UpdateColor(ctx context.Context, in *UpdateColorRequest, opts ...grpc.CallOption) (*Color, error) {
	out := new(Color)
	err := c.cc.Invoke(ctx, "/smartcore.bos.ColorApi/UpdateColor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *colorApiClient) PullColor(ctx context.Context, in *PullColorRequest, opts ...grpc.CallOption) (ColorApi_PullColorClient, error) {
	stream, err := c.cc.NewStream(ctx, &ColorApi_ServiceDesc.Streams[0], "/smartcore.bos.ColorApi/PullColor", opts...)
	if err != nil {
		return nil, err
	}
	x := &colorApiPullColorClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ColorApi_PullColorClient interface {
	Recv() (*PullColorResponse, error)
	grpc.ClientStream
}

type colorApiPullColorClient struct {
	grpc.ClientStream
}

func (x *colorApiPullColorClient) Recv() (*PullColorResponse, error) {
	m := new(PullColorResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ColorApiServer is the server API for ColorApi service.
// All implementations must embed UnimplementedColorApiServer
// for forward compatibility
type ColorApiServer interface {
	// get the current color
	GetColor(context.Context, *GetColorRequest) (*Color, error)
	// request that the color be changed
	UpdateColor(context.Context, *UpdateColorRequest) (*Color, error)
	// request updates to changes to the color value
	PullColor(*PullColorRequest, ColorApi_PullColorServer) error
	mustEmbedUnimplementedColorApiServer()
}

// UnimplementedColorApiServer must be embedded to have forward compatible implementations.
type UnimplementedColorApiServer struct {
}

func (UnimplementedColorApiServer) GetColor(context.Context, *GetColorRequest) (*Color, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetColor not implemented")
}
func (UnimplementedColorApiServer) UpdateColor(context.Context, *UpdateColorRequest) (*Color, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateColor not implemented")
}
func (UnimplementedColorApiServer) PullColor(*PullColorRequest, ColorApi_PullColorServer) error {
	return status.Errorf(codes.Unimplemented, "method PullColor not implemented")
}
func (UnimplementedColorApiServer) mustEmbedUnimplementedColorApiServer() {}

// UnsafeColorApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ColorApiServer will
// result in compilation errors.
type UnsafeColorApiServer interface {
	mustEmbedUnimplementedColorApiServer()
}

func RegisterColorApiServer(s grpc.ServiceRegistrar, srv ColorApiServer) {
	s.RegisterService(&ColorApi_ServiceDesc, srv)
}

func _ColorApi_GetColor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetColorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ColorApiServer).GetColor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smartcore.bos.ColorApi/GetColor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ColorApiServer).GetColor(ctx, req.(*GetColorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ColorApi_UpdateColor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateColorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ColorApiServer).UpdateColor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smartcore.bos.ColorApi/UpdateColor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ColorApiServer).UpdateColor(ctx, req.(*UpdateColorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ColorApi_PullColor_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PullColorRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ColorApiServer).PullColor(m, &colorApiPullColorServer{stream})
}

type ColorApi_PullColorServer interface {
	Send(*PullColorResponse) error
	grpc.ServerStream
}

type colorApiPullColorServer struct {
	grpc.ServerStream
}

func (x *colorApiPullColorServer) Send(m *PullColorResponse) error {
	return x.ServerStream.SendMsg(m)
}

// ColorApi_ServiceDesc is the grpc.ServiceDesc for ColorApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ColorApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "smartcore.bos.ColorApi",
	HandlerType: (*ColorApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetColor",
			Handler:    _ColorApi_GetColor_Handler,
		},
		{
			MethodName: "UpdateColor",
			Handler:    _ColorApi_UpdateColor_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PullColor",
			Handler:       _ColorApi_PullColor_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "color.proto",
}

// ColorInfoClient is the client API for ColorInfo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ColorInfoClient interface {
	// Get information about how a named device implements Color features
	DescribeColor(ctx context.Context, in *DescribeColorRequest, opts ...grpc.CallOption) (*ColorSupport, error)
}

type colorInfoClient struct {
	cc grpc.ClientConnInterface
}

func NewColorInfoClient(cc grpc.ClientConnInterface) ColorInfoClient {
	return &colorInfoClient{cc}
}

func (c *colorInfoClient) DescribeColor(ctx context.Context, in *DescribeColorRequest, opts ...grpc.CallOption) (*ColorSupport, error) {
	out := new(ColorSupport)
	err := c.cc.Invoke(ctx, "/smartcore.bos.ColorInfo/DescribeColor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ColorInfoServer is the server API for ColorInfo service.
// All implementations must embed UnimplementedColorInfoServer
// for forward compatibility
type ColorInfoServer interface {
	// Get information about how a named device implements Color features
	DescribeColor(context.Context, *DescribeColorRequest) (*ColorSupport, error)
	mustEmbedUnimplementedColorInfoServer()
}

// UnimplementedColorInfoServer must be embedded to have forward compatible implementations.
type UnimplementedColorInfoServer struct {
}

func (UnimplementedColorInfoServer) DescribeColor(context.Context, *DescribeColorRequest) (*ColorSupport, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeColor not implemented")
}
func (UnimplementedColorInfoServer) mustEmbedUnimplementedColorInfoServer() {}

// UnsafeColorInfoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ColorInfoServer will
// result in compilation errors.
type UnsafeColorInfoServer interface {
	mustEmbedUnimplementedColorInfoServer()
}

func RegisterColorInfoServer(s grpc.ServiceRegistrar, srv ColorInfoServer) {
	s.RegisterService(&ColorInfo_ServiceDesc, srv)
}

func _ColorInfo_DescribeColor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DescribeColorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ColorInfoServer).DescribeColor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smartcore.bos.ColorInfo/DescribeColor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ColorInfoServer).DescribeColor(ctx, req.(*DescribeColorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ColorInfo_ServiceDesc is the grpc.ServiceDesc for ColorInfo service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ColorInfo_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "smartcore.bos.ColorInfo",
	HandlerType: (*ColorInfoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DescribeColor",
			Handler:    _ColorInfo_DescribeColor_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "color.proto",
}
