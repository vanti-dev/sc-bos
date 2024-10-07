// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: lighting_test.proto

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
	LightingTestApi_GetLightHealth_FullMethodName  = "/smartcore.bos.LightingTestApi/GetLightHealth"
	LightingTestApi_ListLightHealth_FullMethodName = "/smartcore.bos.LightingTestApi/ListLightHealth"
	LightingTestApi_ListLightEvents_FullMethodName = "/smartcore.bos.LightingTestApi/ListLightEvents"
	LightingTestApi_GetReportCSV_FullMethodName    = "/smartcore.bos.LightingTestApi/GetReportCSV"
)

// LightingTestApiClient is the client API for LightingTestApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LightingTestApiClient interface {
	GetLightHealth(ctx context.Context, in *GetLightHealthRequest, opts ...grpc.CallOption) (*LightHealth, error)
	ListLightHealth(ctx context.Context, in *ListLightHealthRequest, opts ...grpc.CallOption) (*ListLightHealthResponse, error)
	ListLightEvents(ctx context.Context, in *ListLightEventsRequest, opts ...grpc.CallOption) (*ListLightEventsResponse, error)
	GetReportCSV(ctx context.Context, in *GetReportCSVRequest, opts ...grpc.CallOption) (*ReportCSV, error)
}

type lightingTestApiClient struct {
	cc grpc.ClientConnInterface
}

func NewLightingTestApiClient(cc grpc.ClientConnInterface) LightingTestApiClient {
	return &lightingTestApiClient{cc}
}

func (c *lightingTestApiClient) GetLightHealth(ctx context.Context, in *GetLightHealthRequest, opts ...grpc.CallOption) (*LightHealth, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LightHealth)
	err := c.cc.Invoke(ctx, LightingTestApi_GetLightHealth_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lightingTestApiClient) ListLightHealth(ctx context.Context, in *ListLightHealthRequest, opts ...grpc.CallOption) (*ListLightHealthResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListLightHealthResponse)
	err := c.cc.Invoke(ctx, LightingTestApi_ListLightHealth_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lightingTestApiClient) ListLightEvents(ctx context.Context, in *ListLightEventsRequest, opts ...grpc.CallOption) (*ListLightEventsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListLightEventsResponse)
	err := c.cc.Invoke(ctx, LightingTestApi_ListLightEvents_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lightingTestApiClient) GetReportCSV(ctx context.Context, in *GetReportCSVRequest, opts ...grpc.CallOption) (*ReportCSV, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReportCSV)
	err := c.cc.Invoke(ctx, LightingTestApi_GetReportCSV_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LightingTestApiServer is the server API for LightingTestApi service.
// All implementations must embed UnimplementedLightingTestApiServer
// for forward compatibility.
type LightingTestApiServer interface {
	GetLightHealth(context.Context, *GetLightHealthRequest) (*LightHealth, error)
	ListLightHealth(context.Context, *ListLightHealthRequest) (*ListLightHealthResponse, error)
	ListLightEvents(context.Context, *ListLightEventsRequest) (*ListLightEventsResponse, error)
	GetReportCSV(context.Context, *GetReportCSVRequest) (*ReportCSV, error)
	mustEmbedUnimplementedLightingTestApiServer()
}

// UnimplementedLightingTestApiServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedLightingTestApiServer struct{}

func (UnimplementedLightingTestApiServer) GetLightHealth(context.Context, *GetLightHealthRequest) (*LightHealth, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLightHealth not implemented")
}
func (UnimplementedLightingTestApiServer) ListLightHealth(context.Context, *ListLightHealthRequest) (*ListLightHealthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLightHealth not implemented")
}
func (UnimplementedLightingTestApiServer) ListLightEvents(context.Context, *ListLightEventsRequest) (*ListLightEventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLightEvents not implemented")
}
func (UnimplementedLightingTestApiServer) GetReportCSV(context.Context, *GetReportCSVRequest) (*ReportCSV, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReportCSV not implemented")
}
func (UnimplementedLightingTestApiServer) mustEmbedUnimplementedLightingTestApiServer() {}
func (UnimplementedLightingTestApiServer) testEmbeddedByValue()                         {}

// UnsafeLightingTestApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LightingTestApiServer will
// result in compilation errors.
type UnsafeLightingTestApiServer interface {
	mustEmbedUnimplementedLightingTestApiServer()
}

func RegisterLightingTestApiServer(s grpc.ServiceRegistrar, srv LightingTestApiServer) {
	// If the following call pancis, it indicates UnimplementedLightingTestApiServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&LightingTestApi_ServiceDesc, srv)
}

func _LightingTestApi_GetLightHealth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLightHealthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LightingTestApiServer).GetLightHealth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LightingTestApi_GetLightHealth_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LightingTestApiServer).GetLightHealth(ctx, req.(*GetLightHealthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LightingTestApi_ListLightHealth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListLightHealthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LightingTestApiServer).ListLightHealth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LightingTestApi_ListLightHealth_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LightingTestApiServer).ListLightHealth(ctx, req.(*ListLightHealthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LightingTestApi_ListLightEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListLightEventsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LightingTestApiServer).ListLightEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LightingTestApi_ListLightEvents_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LightingTestApiServer).ListLightEvents(ctx, req.(*ListLightEventsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LightingTestApi_GetReportCSV_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReportCSVRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LightingTestApiServer).GetReportCSV(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LightingTestApi_GetReportCSV_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LightingTestApiServer).GetReportCSV(ctx, req.(*GetReportCSVRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LightingTestApi_ServiceDesc is the grpc.ServiceDesc for LightingTestApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LightingTestApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "smartcore.bos.LightingTestApi",
	HandlerType: (*LightingTestApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLightHealth",
			Handler:    _LightingTestApi_GetLightHealth_Handler,
		},
		{
			MethodName: "ListLightHealth",
			Handler:    _LightingTestApi_ListLightHealth_Handler,
		},
		{
			MethodName: "ListLightEvents",
			Handler:    _LightingTestApi_ListLightEvents_Handler,
		},
		{
			MethodName: "GetReportCSV",
			Handler:    _LightingTestApi_GetReportCSV_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lighting_test.proto",
}
