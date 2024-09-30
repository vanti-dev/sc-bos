// Code generated by protoc-gen-wrapper. DO NOT EDIT.

package gen

import (
	wrap "github.com/smart-core-os/sc-golang/pkg/wrap"
	grpc "google.golang.org/grpc"
)

// WrapTemperatureApi	adapts a gen.TemperatureApiServer	and presents it as a gen.TemperatureApiClient
func WrapTemperatureApi(server TemperatureApiServer) *TemperatureApiWrapper {
	conn := wrap.ServerToClient(TemperatureApi_ServiceDesc, server)
	client := NewTemperatureApiClient(conn)
	return &TemperatureApiWrapper{
		TemperatureApiClient: client,
		server:               server,
		conn:                 conn,
		desc:                 TemperatureApi_ServiceDesc,
	}
}

type TemperatureApiWrapper struct {
	TemperatureApiClient

	server TemperatureApiServer
	conn   grpc.ClientConnInterface
	desc   grpc.ServiceDesc
}

// UnwrapServer returns the underlying server instance.
func (w *TemperatureApiWrapper) UnwrapServer() TemperatureApiServer {
	return w.server
}

// Unwrap implements wrap.Unwrapper and returns the underlying server instance as an unknown type.
func (w *TemperatureApiWrapper) Unwrap() any {
	return w.UnwrapServer()
}

func (w *TemperatureApiWrapper) UnwrapService() (grpc.ClientConnInterface, grpc.ServiceDesc) {
	return w.conn, w.desc
}
