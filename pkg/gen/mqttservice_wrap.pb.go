// Code generated by protoc-gen-wrapper. DO NOT EDIT.

package gen

import (
	"github.com/smart-core-os/sc-golang/pkg/wrap"
)

// WrapMqttService	adapts a gen.MqttServiceServer	and presents it as a gen.MqttServiceClient
func WrapMqttService(server MqttServiceServer) MqttServiceClient {
	conn := wrap.ServerToClient(MqttService_ServiceDesc, server)
	client := NewMqttServiceClient(conn)
	return &mqttServiceWrapper{
		MqttServiceClient: client,
		server:            server,
	}
}

type mqttServiceWrapper struct {
	MqttServiceClient

	server MqttServiceServer
}

// compile time check that we implement the interface we need
var _ MqttServiceClient = (*mqttServiceWrapper)(nil)

// UnwrapServer returns the underlying server instance.
func (w *mqttServiceWrapper) UnwrapServer() MqttServiceServer {
	return w.server
}

// Unwrap implements wrap.Unwrapper and returns the underlying server instance as an unknown type.
func (w *mqttServiceWrapper) Unwrap() any {
	return w.UnwrapServer()
}
