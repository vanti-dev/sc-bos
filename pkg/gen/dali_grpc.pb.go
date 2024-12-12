// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.3
// source: dali.proto

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
	DaliApi_AddToGroup_FullMethodName           = "/smartcore.bos.driver.dali.DaliApi/AddToGroup"
	DaliApi_RemoveFromGroup_FullMethodName      = "/smartcore.bos.driver.dali.DaliApi/RemoveFromGroup"
	DaliApi_GetGroupMembership_FullMethodName   = "/smartcore.bos.driver.dali.DaliApi/GetGroupMembership"
	DaliApi_GetControlGearStatus_FullMethodName = "/smartcore.bos.driver.dali.DaliApi/GetControlGearStatus"
	DaliApi_GetEmergencyStatus_FullMethodName   = "/smartcore.bos.driver.dali.DaliApi/GetEmergencyStatus"
	DaliApi_Identify_FullMethodName             = "/smartcore.bos.driver.dali.DaliApi/Identify"
	DaliApi_StartTest_FullMethodName            = "/smartcore.bos.driver.dali.DaliApi/StartTest"
	DaliApi_StopTest_FullMethodName             = "/smartcore.bos.driver.dali.DaliApi/StopTest"
	DaliApi_GetTestResult_FullMethodName        = "/smartcore.bos.driver.dali.DaliApi/GetTestResult"
	DaliApi_DeleteTestResult_FullMethodName     = "/smartcore.bos.driver.dali.DaliApi/DeleteTestResult"
)

// DaliApiClient is the client API for DaliApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DaliApiClient interface {
	// Group commands
	AddToGroup(ctx context.Context, in *AddToGroupRequest, opts ...grpc.CallOption) (*AddToGroupResponse, error)
	RemoveFromGroup(ctx context.Context, in *RemoveFromGroupRequest, opts ...grpc.CallOption) (*RemoveFromGroupResponse, error)
	GetGroupMembership(ctx context.Context, in *GetGroupMembershipRequest, opts ...grpc.CallOption) (*GetGroupMembershipResponse, error)
	// Control Gear Commands
	GetControlGearStatus(ctx context.Context, in *GetControlGearStatusRequest, opts ...grpc.CallOption) (*ControlGearStatus, error)
	// Emergency Light commands
	GetEmergencyStatus(ctx context.Context, in *GetEmergencyStatusRequest, opts ...grpc.CallOption) (*EmergencyStatus, error)
	// Start identification for the light; typically this will flash an indicator LED for a few seconds.
	Identify(ctx context.Context, in *IdentifyRequest, opts ...grpc.CallOption) (*IdentifyResponse, error)
	// Attempt to start a function or duration test.
	StartTest(ctx context.Context, in *StartTestRequest, opts ...grpc.CallOption) (*StartTestResponse, error)
	// Stop any test that is in progress.
	StopTest(ctx context.Context, in *StopTestRequest, opts ...grpc.CallOption) (*StopTestResponse, error)
	// Retrieve the results (pass/fail) of the most recent function or duration test to be performed.
	GetTestResult(ctx context.Context, in *GetTestResultRequest, opts ...grpc.CallOption) (*TestResult, error)
	// Can be used to clear a test pass from the light's internal memory. Only passes can be deleted in this way -
	// a failure will stick until it's replaced with a pass.
	//
	// Useful to make sure you don't record the same test multiple times.
	DeleteTestResult(ctx context.Context, in *DeleteTestResultRequest, opts ...grpc.CallOption) (*TestResult, error)
}

type daliApiClient struct {
	cc grpc.ClientConnInterface
}

func NewDaliApiClient(cc grpc.ClientConnInterface) DaliApiClient {
	return &daliApiClient{cc}
}

func (c *daliApiClient) AddToGroup(ctx context.Context, in *AddToGroupRequest, opts ...grpc.CallOption) (*AddToGroupResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddToGroupResponse)
	err := c.cc.Invoke(ctx, DaliApi_AddToGroup_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daliApiClient) RemoveFromGroup(ctx context.Context, in *RemoveFromGroupRequest, opts ...grpc.CallOption) (*RemoveFromGroupResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RemoveFromGroupResponse)
	err := c.cc.Invoke(ctx, DaliApi_RemoveFromGroup_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daliApiClient) GetGroupMembership(ctx context.Context, in *GetGroupMembershipRequest, opts ...grpc.CallOption) (*GetGroupMembershipResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetGroupMembershipResponse)
	err := c.cc.Invoke(ctx, DaliApi_GetGroupMembership_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daliApiClient) GetControlGearStatus(ctx context.Context, in *GetControlGearStatusRequest, opts ...grpc.CallOption) (*ControlGearStatus, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ControlGearStatus)
	err := c.cc.Invoke(ctx, DaliApi_GetControlGearStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daliApiClient) GetEmergencyStatus(ctx context.Context, in *GetEmergencyStatusRequest, opts ...grpc.CallOption) (*EmergencyStatus, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EmergencyStatus)
	err := c.cc.Invoke(ctx, DaliApi_GetEmergencyStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daliApiClient) Identify(ctx context.Context, in *IdentifyRequest, opts ...grpc.CallOption) (*IdentifyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(IdentifyResponse)
	err := c.cc.Invoke(ctx, DaliApi_Identify_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daliApiClient) StartTest(ctx context.Context, in *StartTestRequest, opts ...grpc.CallOption) (*StartTestResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StartTestResponse)
	err := c.cc.Invoke(ctx, DaliApi_StartTest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daliApiClient) StopTest(ctx context.Context, in *StopTestRequest, opts ...grpc.CallOption) (*StopTestResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StopTestResponse)
	err := c.cc.Invoke(ctx, DaliApi_StopTest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daliApiClient) GetTestResult(ctx context.Context, in *GetTestResultRequest, opts ...grpc.CallOption) (*TestResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TestResult)
	err := c.cc.Invoke(ctx, DaliApi_GetTestResult_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daliApiClient) DeleteTestResult(ctx context.Context, in *DeleteTestResultRequest, opts ...grpc.CallOption) (*TestResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TestResult)
	err := c.cc.Invoke(ctx, DaliApi_DeleteTestResult_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DaliApiServer is the server API for DaliApi service.
// All implementations must embed UnimplementedDaliApiServer
// for forward compatibility.
type DaliApiServer interface {
	// Group commands
	AddToGroup(context.Context, *AddToGroupRequest) (*AddToGroupResponse, error)
	RemoveFromGroup(context.Context, *RemoveFromGroupRequest) (*RemoveFromGroupResponse, error)
	GetGroupMembership(context.Context, *GetGroupMembershipRequest) (*GetGroupMembershipResponse, error)
	// Control Gear Commands
	GetControlGearStatus(context.Context, *GetControlGearStatusRequest) (*ControlGearStatus, error)
	// Emergency Light commands
	GetEmergencyStatus(context.Context, *GetEmergencyStatusRequest) (*EmergencyStatus, error)
	// Start identification for the light; typically this will flash an indicator LED for a few seconds.
	Identify(context.Context, *IdentifyRequest) (*IdentifyResponse, error)
	// Attempt to start a function or duration test.
	StartTest(context.Context, *StartTestRequest) (*StartTestResponse, error)
	// Stop any test that is in progress.
	StopTest(context.Context, *StopTestRequest) (*StopTestResponse, error)
	// Retrieve the results (pass/fail) of the most recent function or duration test to be performed.
	GetTestResult(context.Context, *GetTestResultRequest) (*TestResult, error)
	// Can be used to clear a test pass from the light's internal memory. Only passes can be deleted in this way -
	// a failure will stick until it's replaced with a pass.
	//
	// Useful to make sure you don't record the same test multiple times.
	DeleteTestResult(context.Context, *DeleteTestResultRequest) (*TestResult, error)
	mustEmbedUnimplementedDaliApiServer()
}

// UnimplementedDaliApiServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDaliApiServer struct{}

func (UnimplementedDaliApiServer) AddToGroup(context.Context, *AddToGroupRequest) (*AddToGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddToGroup not implemented")
}
func (UnimplementedDaliApiServer) RemoveFromGroup(context.Context, *RemoveFromGroupRequest) (*RemoveFromGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveFromGroup not implemented")
}
func (UnimplementedDaliApiServer) GetGroupMembership(context.Context, *GetGroupMembershipRequest) (*GetGroupMembershipResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroupMembership not implemented")
}
func (UnimplementedDaliApiServer) GetControlGearStatus(context.Context, *GetControlGearStatusRequest) (*ControlGearStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetControlGearStatus not implemented")
}
func (UnimplementedDaliApiServer) GetEmergencyStatus(context.Context, *GetEmergencyStatusRequest) (*EmergencyStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEmergencyStatus not implemented")
}
func (UnimplementedDaliApiServer) Identify(context.Context, *IdentifyRequest) (*IdentifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Identify not implemented")
}
func (UnimplementedDaliApiServer) StartTest(context.Context, *StartTestRequest) (*StartTestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartTest not implemented")
}
func (UnimplementedDaliApiServer) StopTest(context.Context, *StopTestRequest) (*StopTestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopTest not implemented")
}
func (UnimplementedDaliApiServer) GetTestResult(context.Context, *GetTestResultRequest) (*TestResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTestResult not implemented")
}
func (UnimplementedDaliApiServer) DeleteTestResult(context.Context, *DeleteTestResultRequest) (*TestResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTestResult not implemented")
}
func (UnimplementedDaliApiServer) mustEmbedUnimplementedDaliApiServer() {}
func (UnimplementedDaliApiServer) testEmbeddedByValue()                 {}

// UnsafeDaliApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DaliApiServer will
// result in compilation errors.
type UnsafeDaliApiServer interface {
	mustEmbedUnimplementedDaliApiServer()
}

func RegisterDaliApiServer(s grpc.ServiceRegistrar, srv DaliApiServer) {
	// If the following call pancis, it indicates UnimplementedDaliApiServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DaliApi_ServiceDesc, srv)
}

func _DaliApi_AddToGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddToGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaliApiServer).AddToGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaliApi_AddToGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaliApiServer).AddToGroup(ctx, req.(*AddToGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaliApi_RemoveFromGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveFromGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaliApiServer).RemoveFromGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaliApi_RemoveFromGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaliApiServer).RemoveFromGroup(ctx, req.(*RemoveFromGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaliApi_GetGroupMembership_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupMembershipRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaliApiServer).GetGroupMembership(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaliApi_GetGroupMembership_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaliApiServer).GetGroupMembership(ctx, req.(*GetGroupMembershipRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaliApi_GetControlGearStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetControlGearStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaliApiServer).GetControlGearStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaliApi_GetControlGearStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaliApiServer).GetControlGearStatus(ctx, req.(*GetControlGearStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaliApi_GetEmergencyStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEmergencyStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaliApiServer).GetEmergencyStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaliApi_GetEmergencyStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaliApiServer).GetEmergencyStatus(ctx, req.(*GetEmergencyStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaliApi_Identify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdentifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaliApiServer).Identify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaliApi_Identify_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaliApiServer).Identify(ctx, req.(*IdentifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaliApi_StartTest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartTestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaliApiServer).StartTest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaliApi_StartTest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaliApiServer).StartTest(ctx, req.(*StartTestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaliApi_StopTest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopTestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaliApiServer).StopTest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaliApi_StopTest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaliApiServer).StopTest(ctx, req.(*StopTestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaliApi_GetTestResult_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTestResultRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaliApiServer).GetTestResult(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaliApi_GetTestResult_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaliApiServer).GetTestResult(ctx, req.(*GetTestResultRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaliApi_DeleteTestResult_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTestResultRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaliApiServer).DeleteTestResult(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DaliApi_DeleteTestResult_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaliApiServer).DeleteTestResult(ctx, req.(*DeleteTestResultRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DaliApi_ServiceDesc is the grpc.ServiceDesc for DaliApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DaliApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "smartcore.bos.driver.dali.DaliApi",
	HandlerType: (*DaliApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddToGroup",
			Handler:    _DaliApi_AddToGroup_Handler,
		},
		{
			MethodName: "RemoveFromGroup",
			Handler:    _DaliApi_RemoveFromGroup_Handler,
		},
		{
			MethodName: "GetGroupMembership",
			Handler:    _DaliApi_GetGroupMembership_Handler,
		},
		{
			MethodName: "GetControlGearStatus",
			Handler:    _DaliApi_GetControlGearStatus_Handler,
		},
		{
			MethodName: "GetEmergencyStatus",
			Handler:    _DaliApi_GetEmergencyStatus_Handler,
		},
		{
			MethodName: "Identify",
			Handler:    _DaliApi_Identify_Handler,
		},
		{
			MethodName: "StartTest",
			Handler:    _DaliApi_StartTest_Handler,
		},
		{
			MethodName: "StopTest",
			Handler:    _DaliApi_StopTest_Handler,
		},
		{
			MethodName: "GetTestResult",
			Handler:    _DaliApi_GetTestResult_Handler,
		},
		{
			MethodName: "DeleteTestResult",
			Handler:    _DaliApi_DeleteTestResult_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dali.proto",
}
