// Code generated by protoc-gen-wrapper. DO NOT EDIT.

package gen

import (
	"github.com/smart-core-os/sc-golang/pkg/wrap"
)

// WrapHubApi	adapts a gen.HubApiServer	and presents it as a gen.HubApiClient
func WrapHubApi(server HubApiServer) HubApiClient {
	conn := wrap.ServerToClient(HubApi_ServiceDesc, server)
	client := NewHubApiClient(conn)
	return &hubApiWrapper{
		HubApiClient: client,
		server:       server,
	}
}

type hubApiWrapper struct {
	HubApiClient

	server HubApiServer
}

// compile time check that we implement the interface we need
var _ HubApiClient = (*hubApiWrapper)(nil)

// UnwrapServer returns the underlying server instance.
func (w *hubApiWrapper) UnwrapServer() HubApiServer {
	return w.server
}

// Unwrap implements wrap.Unwrapper and returns the underlying server instance as an unknown type.
func (w *hubApiWrapper) Unwrap() any {
	return w.UnwrapServer()
}
