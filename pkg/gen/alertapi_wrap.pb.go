// Code generated by protoc-gen-wrapper. DO NOT EDIT.

package gen

import (
	wrap "github.com/smart-core-os/sc-golang/pkg/wrap"
)

// WrapAlertApi	adapts a AlertApiServer	and presents it as a AlertApiClient
func WrapAlertApi(server AlertApiServer) AlertApiClient {
	conn := wrap.ServerToClient(AlertApi_ServiceDesc, server)
	client := NewAlertApiClient(conn)
	return &alertApiWrapper{
		AlertApiClient: client,
		server:         server,
	}
}

type alertApiWrapper struct {
	AlertApiClient

	server AlertApiServer
}

// compile time check that we implement the interface we need
var _ AlertApiClient = (*alertApiWrapper)(nil)

// UnwrapServer returns the underlying server instance.
func (w *alertApiWrapper) UnwrapServer() AlertApiServer {
	return w.server
}

// Unwrap implements wrap.Unwrapper and returns the underlying server instance as an unknown type.
func (w *alertApiWrapper) Unwrap() any {
	return w.UnwrapServer()
}
