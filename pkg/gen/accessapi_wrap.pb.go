// Code generated by protoc-gen-wrapper. DO NOT EDIT.

package gen

import (
	wrap "github.com/smart-core-os/sc-golang/pkg/wrap"
	grpc "google.golang.org/grpc"
)

// WrapAccessApi	adapts a AccessApiServer	and presents it as a AccessApiClient
func WrapAccessApi(server AccessApiServer) *AccessApiWrapper {
	conn := wrap.ServerToClient(AccessApi_ServiceDesc, server)
	client := NewAccessApiClient(conn)
	return &AccessApiWrapper{
		AccessApiClient: client,
		server:          server,
		conn:            conn,
		desc:            AccessApi_ServiceDesc,
	}
}

type AccessApiWrapper struct {
	AccessApiClient

	server AccessApiServer
	conn   grpc.ClientConnInterface
	desc   grpc.ServiceDesc
}

// UnwrapServer returns the underlying server instance.
func (w *AccessApiWrapper) UnwrapServer() AccessApiServer {
	return w.server
}

// Unwrap implements wrap.Unwrapper and returns the underlying server instance as an unknown type.
func (w *AccessApiWrapper) Unwrap() any {
	return w.UnwrapServer()
}

func (w *AccessApiWrapper) UnwrapService() (grpc.ClientConnInterface, grpc.ServiceDesc) {
	return w.conn, w.desc
}
