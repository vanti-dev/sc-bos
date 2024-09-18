// Code generated by protoc-gen-wrapper. DO NOT EDIT.

package gen

import (
	"google.golang.org/grpc"

	"github.com/smart-core-os/sc-golang/pkg/wrap"
)

// WrapMeterHistory	adapts a gen.MeterHistoryServer	and presents it as a gen.MeterHistoryClient
func WrapMeterHistory(server MeterHistoryServer) *MeterHistoryWrapper {
	conn := wrap.ServerToClient(MeterHistory_ServiceDesc, server)
	client := NewMeterHistoryClient(conn)
	return &MeterHistoryWrapper{
		MeterHistoryClient: client,
		server:             server,
		conn:               conn,
		desc:               MeterHistory_ServiceDesc,
	}
}

type MeterHistoryWrapper struct {
	MeterHistoryClient

	server MeterHistoryServer
	conn   grpc.ClientConnInterface
	desc   grpc.ServiceDesc
}

// UnwrapServer returns the underlying server instance.
func (w *MeterHistoryWrapper) UnwrapServer() MeterHistoryServer {
	return w.server
}

// Unwrap implements wrap.Unwrapper and returns the underlying server instance as an unknown type.
func (w *MeterHistoryWrapper) Unwrap() any {
	return w.UnwrapServer()
}

func (w *MeterHistoryWrapper) UnwrapService() (grpc.ClientConnInterface, grpc.ServiceDesc) {
	return w.conn, w.desc
}
