// Code generated by protoc-gen-wrapper. DO NOT EDIT.

package gen

import (
	wrap "github.com/smart-core-os/sc-golang/pkg/wrap"
)

// WrapElectricHistory	adapts a ElectricHistoryServer	and presents it as a ElectricHistoryClient
func WrapElectricHistory(server ElectricHistoryServer) ElectricHistoryClient {
	conn := wrap.ServerToClient(ElectricHistory_ServiceDesc, server)
	client := NewElectricHistoryClient(conn)
	return &electricHistoryWrapper{
		ElectricHistoryClient: client,
		server:                server,
	}
}

type electricHistoryWrapper struct {
	ElectricHistoryClient

	server ElectricHistoryServer
}

// compile time check that we implement the interface we need
var _ ElectricHistoryClient = (*electricHistoryWrapper)(nil)

// UnwrapServer returns the underlying server instance.
func (w *electricHistoryWrapper) UnwrapServer() ElectricHistoryServer {
	return w.server
}

// Unwrap implements wrap.Unwrapper and returns the underlying server instance as an unknown type.
func (w *electricHistoryWrapper) Unwrap() any {
	return w.UnwrapServer()
}
