// Code generated by protoc-gen-wrapper. DO NOT EDIT.

package gen

import (
	wrap "github.com/smart-core-os/sc-golang/pkg/wrap"
)

// WrapLightingTestApi	adapts a LightingTestApiServer	and presents it as a LightingTestApiClient
func WrapLightingTestApi(server LightingTestApiServer) LightingTestApiClient {
	conn := wrap.ServerToClient(LightingTestApi_ServiceDesc, server)
	client := NewLightingTestApiClient(conn)
	return &lightingTestApiWrapper{
		LightingTestApiClient: client,
		server:                server,
	}
}

type lightingTestApiWrapper struct {
	LightingTestApiClient

	server LightingTestApiServer
}

// compile time check that we implement the interface we need
var _ LightingTestApiClient = (*lightingTestApiWrapper)(nil)

// UnwrapServer returns the underlying server instance.
func (w *lightingTestApiWrapper) UnwrapServer() LightingTestApiServer {
	return w.server
}

// Unwrap implements wrap.Unwrapper and returns the underlying server instance as an unknown type.
func (w *lightingTestApiWrapper) Unwrap() any {
	return w.UnwrapServer()
}
