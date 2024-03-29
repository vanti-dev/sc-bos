// Code generated by protoc-gen-wrapper. DO NOT EDIT.

package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
)

// WrapAxiomXaDriverService	adapts a AxiomXaDriverServiceServer	and presents it as a AxiomXaDriverServiceClient
func WrapAxiomXaDriverService(server AxiomXaDriverServiceServer) AxiomXaDriverServiceClient {
	return &axiomXaDriverServiceWrapper{server}
}

type axiomXaDriverServiceWrapper struct {
	server AxiomXaDriverServiceServer
}

// compile time check that we implement the interface we need
var _ AxiomXaDriverServiceClient = (*axiomXaDriverServiceWrapper)(nil)

// UnwrapServer returns the underlying server instance.
func (w *axiomXaDriverServiceWrapper) UnwrapServer() AxiomXaDriverServiceServer {
	return w.server
}

// Unwrap implements wrap.Unwrapper and returns the underlying server instance as an unknown type.
func (w *axiomXaDriverServiceWrapper) Unwrap() any {
	return w.UnwrapServer()
}

func (w *axiomXaDriverServiceWrapper) SaveQRCredential(ctx context.Context, req *SaveQRCredentialRequest, _ ...grpc.CallOption) (*SaveQRCredentialResponse, error) {
	return w.server.SaveQRCredential(ctx, req)
}
