// Code generated by protoc-gen-wrapper. DO NOT EDIT.

package gen

import (
	wrap "github.com/smart-core-os/sc-golang/pkg/wrap"
	grpc "google.golang.org/grpc"
)

// WrapAccessHistory	adapts a AccessHistoryServer	and presents it as a AccessHistoryClient
func WrapAccessHistory(server AccessHistoryServer) *AccessHistoryWrapper {
	conn := wrap.ServerToClient(AccessHistory_ServiceDesc, server)
	client := NewAccessHistoryClient(conn)
	return &AccessHistoryWrapper{
		AccessHistoryClient: client,
		server:              server,
		conn:                conn,
		desc:                AccessHistory_ServiceDesc,
	}
}

type AccessHistoryWrapper struct {
	AccessHistoryClient

	server AccessHistoryServer
	conn   grpc.ClientConnInterface
	desc   grpc.ServiceDesc
}

// UnwrapServer returns the underlying server instance.
func (w *AccessHistoryWrapper) UnwrapServer() AccessHistoryServer {
	return w.server
}

// Unwrap implements wrap.Unwrapper and returns the underlying server instance as an unknown type.
func (w *AccessHistoryWrapper) Unwrap() any {
	return w.UnwrapServer()
}

func (w *AccessHistoryWrapper) UnwrapService() (grpc.ClientConnInterface, grpc.ServiceDesc) {
	return w.conn, w.desc
}
