// Package reflectionapi provides utilities for interacting with the gRPC reflection API.
package reflectionapi

import (
	"fmt"

	"google.golang.org/grpc/codes"
	reflectionpb "google.golang.org/grpc/reflection/grpc_reflection_v1"
	"google.golang.org/grpc/status"
)

// ListServices sends a ListServices message on stream and returns the list of services sent in the response.
func ListServices(stream reflectionpb.ServerReflection_ServerReflectionInfoClient) ([]*reflectionpb.ServiceResponse, error) {
	res, err := reflect[*reflectionpb.ServerReflectionResponse_ListServicesResponse](stream, &reflectionpb.ServerReflectionRequest{
		MessageRequest: &reflectionpb.ServerReflectionRequest_ListServices{},
	})
	if err != nil {
		return nil, err
	}
	return res.ListServicesResponse.GetService(), nil
}

// FileContainingSymbol sends a FileContainingSymbol message on stream and returns the files containing the symbol.
func FileContainingSymbol(stream reflectionpb.ServerReflection_ServerReflectionInfoClient, symbol string) ([][]byte, error) {
	res, err := reflect[*reflectionpb.ServerReflectionResponse_FileDescriptorResponse](stream, &reflectionpb.ServerReflectionRequest{
		MessageRequest: &reflectionpb.ServerReflectionRequest_FileContainingSymbol{FileContainingSymbol: symbol},
	})
	if err != nil {
		return nil, err
	}
	return res.FileDescriptorResponse.GetFileDescriptorProto(), nil
}

// FileByFilename sends a FileByFilename message on stream and returns the files at the path.
func FileByFilename(stream reflectionpb.ServerReflection_ServerReflectionInfoClient, path string) ([][]byte, error) {
	res, err := reflect[*reflectionpb.ServerReflectionResponse_FileDescriptorResponse](stream, &reflectionpb.ServerReflectionRequest{
		MessageRequest: &reflectionpb.ServerReflectionRequest_FileByFilename{FileByFilename: path},
	})
	if err != nil {
		return nil, err
	}
	return res.FileDescriptorResponse.GetFileDescriptorProto(), nil
}

func reflect[T any](client reflectionpb.ServerReflection_ServerReflectionInfoClient, req *reflectionpb.ServerReflectionRequest) (T, error) {
	var zero T
	err := client.Send(req)
	if err != nil {
		return zero, fmt.Errorf("failed to send %T for %[1]v: %[2]w", req.MessageRequest, err)
	}
	msg, err := client.Recv()
	if err != nil {
		return zero, fmt.Errorf("failed to receive response %T for %[1]v: %[2]w", req.MessageRequest, err)
	}

	errRes, ok := msg.MessageResponse.(*reflectionpb.ServerReflectionResponse_ErrorResponse)
	if ok {
		return zero, status.Error(codes.Code(errRes.ErrorResponse.GetErrorCode()), errRes.ErrorResponse.GetErrorMessage())
	}

	return msg.MessageResponse.(T), nil
}
