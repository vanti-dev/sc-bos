// Code generated by protoc-gen-wrapper. DO NOT EDIT.

package gen

import (
	"github.com/smart-core-os/sc-golang/pkg/wrap"
)

// WrapMeterApi	adapts a gen.MeterApiServer	and presents it as a gen.MeterApiClient
func WrapMeterApi(server MeterApiServer) MeterApiClient {
	conn := wrap.ServerToClient(MeterApi_ServiceDesc, server)
	client := NewMeterApiClient(conn)
	return &meterApiWrapper{
		MeterApiClient: client,
		server:         server,
	}
}

type meterApiWrapper struct {
	MeterApiClient

	server MeterApiServer
}

// compile time check that we implement the interface we need
var _ MeterApiClient = (*meterApiWrapper)(nil)

// UnwrapServer returns the underlying server instance.
func (w *meterApiWrapper) UnwrapServer() MeterApiServer {
	return w.server
}

// Unwrap implements wrap.Unwrapper and returns the underlying server instance as an unknown type.
func (w *meterApiWrapper) Unwrap() any {
	return w.UnwrapServer()
}
