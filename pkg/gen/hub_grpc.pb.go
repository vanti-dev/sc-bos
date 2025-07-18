// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: hub.proto

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
	HubApi_GetHubNode_FullMethodName     = "/smartcore.bos.HubApi/GetHubNode"
	HubApi_ListHubNodes_FullMethodName   = "/smartcore.bos.HubApi/ListHubNodes"
	HubApi_PullHubNodes_FullMethodName   = "/smartcore.bos.HubApi/PullHubNodes"
	HubApi_InspectHubNode_FullMethodName = "/smartcore.bos.HubApi/InspectHubNode"
	HubApi_EnrollHubNode_FullMethodName  = "/smartcore.bos.HubApi/EnrollHubNode"
	HubApi_RenewHubNode_FullMethodName   = "/smartcore.bos.HubApi/RenewHubNode"
	HubApi_TestHubNode_FullMethodName    = "/smartcore.bos.HubApi/TestHubNode"
	HubApi_ForgetHubNode_FullMethodName  = "/smartcore.bos.HubApi/ForgetHubNode"
)

// HubApiClient is the client API for HubApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HubApiClient interface {
	GetHubNode(ctx context.Context, in *GetHubNodeRequest, opts ...grpc.CallOption) (*HubNode, error)
	ListHubNodes(ctx context.Context, in *ListHubNodesRequest, opts ...grpc.CallOption) (*ListHubNodesResponse, error)
	PullHubNodes(ctx context.Context, in *PullHubNodesRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[PullHubNodesResponse], error)
	// Query the hub node for information that can be used to knowledgeably enroll that node with this hub.
	// This request will return both the node metadata and public certificates presented by the node.
	InspectHubNode(ctx context.Context, in *InspectHubNodeRequest, opts ...grpc.CallOption) (*HubNodeInspection, error)
	// Enroll the node with this hub.
	// Enrollment involves the hub signing the nodes public key and issuing that cert to the node.
	// A node can only be enrolled with one hub, the first to enroll the node wins.
	// Use RenewHubNode to refresh the certificate issued to the node.
	EnrollHubNode(ctx context.Context, in *EnrollHubNodeRequest, opts ...grpc.CallOption) (*HubNode, error)
	// Re-sign and re-issue a certificate to the node.
	// Fails if the node isn't already enrolled.
	RenewHubNode(ctx context.Context, in *RenewHubNodeRequest, opts ...grpc.CallOption) (*HubNode, error)
	// Test that communications with an enrolled node is working.
	// This checks communication and the TLS stack, only returning success if the node presents a public certificate signed
	// by this hub.
	TestHubNode(ctx context.Context, in *TestHubNodeRequest, opts ...grpc.CallOption) (*TestHubNodeResponse, error)
	// Forget a node that was previously enrolled with this hub.
	ForgetHubNode(ctx context.Context, in *ForgetHubNodeRequest, opts ...grpc.CallOption) (*ForgetHubNodeResponse, error)
}

type hubApiClient struct {
	cc grpc.ClientConnInterface
}

func NewHubApiClient(cc grpc.ClientConnInterface) HubApiClient {
	return &hubApiClient{cc}
}

func (c *hubApiClient) GetHubNode(ctx context.Context, in *GetHubNodeRequest, opts ...grpc.CallOption) (*HubNode, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HubNode)
	err := c.cc.Invoke(ctx, HubApi_GetHubNode_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hubApiClient) ListHubNodes(ctx context.Context, in *ListHubNodesRequest, opts ...grpc.CallOption) (*ListHubNodesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListHubNodesResponse)
	err := c.cc.Invoke(ctx, HubApi_ListHubNodes_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hubApiClient) PullHubNodes(ctx context.Context, in *PullHubNodesRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[PullHubNodesResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &HubApi_ServiceDesc.Streams[0], HubApi_PullHubNodes_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[PullHubNodesRequest, PullHubNodesResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type HubApi_PullHubNodesClient = grpc.ServerStreamingClient[PullHubNodesResponse]

func (c *hubApiClient) InspectHubNode(ctx context.Context, in *InspectHubNodeRequest, opts ...grpc.CallOption) (*HubNodeInspection, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HubNodeInspection)
	err := c.cc.Invoke(ctx, HubApi_InspectHubNode_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hubApiClient) EnrollHubNode(ctx context.Context, in *EnrollHubNodeRequest, opts ...grpc.CallOption) (*HubNode, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HubNode)
	err := c.cc.Invoke(ctx, HubApi_EnrollHubNode_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hubApiClient) RenewHubNode(ctx context.Context, in *RenewHubNodeRequest, opts ...grpc.CallOption) (*HubNode, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HubNode)
	err := c.cc.Invoke(ctx, HubApi_RenewHubNode_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hubApiClient) TestHubNode(ctx context.Context, in *TestHubNodeRequest, opts ...grpc.CallOption) (*TestHubNodeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TestHubNodeResponse)
	err := c.cc.Invoke(ctx, HubApi_TestHubNode_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hubApiClient) ForgetHubNode(ctx context.Context, in *ForgetHubNodeRequest, opts ...grpc.CallOption) (*ForgetHubNodeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ForgetHubNodeResponse)
	err := c.cc.Invoke(ctx, HubApi_ForgetHubNode_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HubApiServer is the server API for HubApi service.
// All implementations must embed UnimplementedHubApiServer
// for forward compatibility.
type HubApiServer interface {
	GetHubNode(context.Context, *GetHubNodeRequest) (*HubNode, error)
	ListHubNodes(context.Context, *ListHubNodesRequest) (*ListHubNodesResponse, error)
	PullHubNodes(*PullHubNodesRequest, grpc.ServerStreamingServer[PullHubNodesResponse]) error
	// Query the hub node for information that can be used to knowledgeably enroll that node with this hub.
	// This request will return both the node metadata and public certificates presented by the node.
	InspectHubNode(context.Context, *InspectHubNodeRequest) (*HubNodeInspection, error)
	// Enroll the node with this hub.
	// Enrollment involves the hub signing the nodes public key and issuing that cert to the node.
	// A node can only be enrolled with one hub, the first to enroll the node wins.
	// Use RenewHubNode to refresh the certificate issued to the node.
	EnrollHubNode(context.Context, *EnrollHubNodeRequest) (*HubNode, error)
	// Re-sign and re-issue a certificate to the node.
	// Fails if the node isn't already enrolled.
	RenewHubNode(context.Context, *RenewHubNodeRequest) (*HubNode, error)
	// Test that communications with an enrolled node is working.
	// This checks communication and the TLS stack, only returning success if the node presents a public certificate signed
	// by this hub.
	TestHubNode(context.Context, *TestHubNodeRequest) (*TestHubNodeResponse, error)
	// Forget a node that was previously enrolled with this hub.
	ForgetHubNode(context.Context, *ForgetHubNodeRequest) (*ForgetHubNodeResponse, error)
	mustEmbedUnimplementedHubApiServer()
}

// UnimplementedHubApiServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedHubApiServer struct{}

func (UnimplementedHubApiServer) GetHubNode(context.Context, *GetHubNodeRequest) (*HubNode, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHubNode not implemented")
}
func (UnimplementedHubApiServer) ListHubNodes(context.Context, *ListHubNodesRequest) (*ListHubNodesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListHubNodes not implemented")
}
func (UnimplementedHubApiServer) PullHubNodes(*PullHubNodesRequest, grpc.ServerStreamingServer[PullHubNodesResponse]) error {
	return status.Errorf(codes.Unimplemented, "method PullHubNodes not implemented")
}
func (UnimplementedHubApiServer) InspectHubNode(context.Context, *InspectHubNodeRequest) (*HubNodeInspection, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InspectHubNode not implemented")
}
func (UnimplementedHubApiServer) EnrollHubNode(context.Context, *EnrollHubNodeRequest) (*HubNode, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnrollHubNode not implemented")
}
func (UnimplementedHubApiServer) RenewHubNode(context.Context, *RenewHubNodeRequest) (*HubNode, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RenewHubNode not implemented")
}
func (UnimplementedHubApiServer) TestHubNode(context.Context, *TestHubNodeRequest) (*TestHubNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestHubNode not implemented")
}
func (UnimplementedHubApiServer) ForgetHubNode(context.Context, *ForgetHubNodeRequest) (*ForgetHubNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ForgetHubNode not implemented")
}
func (UnimplementedHubApiServer) mustEmbedUnimplementedHubApiServer() {}
func (UnimplementedHubApiServer) testEmbeddedByValue()                {}

// UnsafeHubApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HubApiServer will
// result in compilation errors.
type UnsafeHubApiServer interface {
	mustEmbedUnimplementedHubApiServer()
}

func RegisterHubApiServer(s grpc.ServiceRegistrar, srv HubApiServer) {
	// If the following call pancis, it indicates UnimplementedHubApiServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&HubApi_ServiceDesc, srv)
}

func _HubApi_GetHubNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetHubNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HubApiServer).GetHubNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HubApi_GetHubNode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HubApiServer).GetHubNode(ctx, req.(*GetHubNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HubApi_ListHubNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListHubNodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HubApiServer).ListHubNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HubApi_ListHubNodes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HubApiServer).ListHubNodes(ctx, req.(*ListHubNodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HubApi_PullHubNodes_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PullHubNodesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(HubApiServer).PullHubNodes(m, &grpc.GenericServerStream[PullHubNodesRequest, PullHubNodesResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type HubApi_PullHubNodesServer = grpc.ServerStreamingServer[PullHubNodesResponse]

func _HubApi_InspectHubNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InspectHubNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HubApiServer).InspectHubNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HubApi_InspectHubNode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HubApiServer).InspectHubNode(ctx, req.(*InspectHubNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HubApi_EnrollHubNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnrollHubNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HubApiServer).EnrollHubNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HubApi_EnrollHubNode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HubApiServer).EnrollHubNode(ctx, req.(*EnrollHubNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HubApi_RenewHubNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RenewHubNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HubApiServer).RenewHubNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HubApi_RenewHubNode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HubApiServer).RenewHubNode(ctx, req.(*RenewHubNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HubApi_TestHubNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestHubNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HubApiServer).TestHubNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HubApi_TestHubNode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HubApiServer).TestHubNode(ctx, req.(*TestHubNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HubApi_ForgetHubNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ForgetHubNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HubApiServer).ForgetHubNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HubApi_ForgetHubNode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HubApiServer).ForgetHubNode(ctx, req.(*ForgetHubNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// HubApi_ServiceDesc is the grpc.ServiceDesc for HubApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var HubApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "smartcore.bos.HubApi",
	HandlerType: (*HubApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetHubNode",
			Handler:    _HubApi_GetHubNode_Handler,
		},
		{
			MethodName: "ListHubNodes",
			Handler:    _HubApi_ListHubNodes_Handler,
		},
		{
			MethodName: "InspectHubNode",
			Handler:    _HubApi_InspectHubNode_Handler,
		},
		{
			MethodName: "EnrollHubNode",
			Handler:    _HubApi_EnrollHubNode_Handler,
		},
		{
			MethodName: "RenewHubNode",
			Handler:    _HubApi_RenewHubNode_Handler,
		},
		{
			MethodName: "TestHubNode",
			Handler:    _HubApi_TestHubNode_Handler,
		},
		{
			MethodName: "ForgetHubNode",
			Handler:    _HubApi_ForgetHubNode_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PullHubNodes",
			Handler:       _HubApi_PullHubNodes_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "hub.proto",
}
