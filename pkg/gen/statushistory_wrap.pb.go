// Code generated by protoc-gen-wrapper. DO NOT EDIT.

package gen

import (
	"google.golang.org/grpc"

	"github.com/smart-core-os/sc-golang/pkg/wrap"
)

// WrapStatusHistory	adapts a gen.StatusHistoryServer	and presents it as a gen.StatusHistoryClient
func WrapStatusHistory(server StatusHistoryServer) *StatusHistoryWrapper {
	conn := wrap.ServerToClient(StatusHistory_ServiceDesc, server)
	client := NewStatusHistoryClient(conn)
	return &StatusHistoryWrapper{
		StatusHistoryClient: client,
		server:              server,
		conn:                conn,
		desc:                StatusHistory_ServiceDesc,
	}
}

type StatusHistoryWrapper struct {
	StatusHistoryClient

	server StatusHistoryServer
	conn   grpc.ClientConnInterface
	desc   grpc.ServiceDesc
}

// UnwrapServer returns the underlying server instance.
func (w *StatusHistoryWrapper) UnwrapServer() StatusHistoryServer {
	return w.server
}

// Unwrap implements wrap.Unwrapper and returns the underlying server instance as an unknown type.
func (w *StatusHistoryWrapper) Unwrap() any {
	return w.UnwrapServer()
}

func (w *StatusHistoryWrapper) UnwrapService() (grpc.ClientConnInterface, grpc.ServiceDesc) {
	return w.conn, w.desc
}
