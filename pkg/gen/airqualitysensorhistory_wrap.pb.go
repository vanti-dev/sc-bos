// Code generated by protoc-gen-wrapper. DO NOT EDIT.

package gen

import (
	wrap "github.com/smart-core-os/sc-golang/pkg/wrap"
)

// WrapAirQualitySensorHistory	adapts a AirQualitySensorHistoryServer	and presents it as a AirQualitySensorHistoryClient
func WrapAirQualitySensorHistory(server AirQualitySensorHistoryServer) AirQualitySensorHistoryClient {
	conn := wrap.ServerToClient(AirQualitySensorHistory_ServiceDesc, server)
	client := NewAirQualitySensorHistoryClient(conn)
	return &airQualitySensorHistoryWrapper{
		AirQualitySensorHistoryClient: client,
		server:                        server,
	}
}

type airQualitySensorHistoryWrapper struct {
	AirQualitySensorHistoryClient

	server AirQualitySensorHistoryServer
}

// compile time check that we implement the interface we need
var _ AirQualitySensorHistoryClient = (*airQualitySensorHistoryWrapper)(nil)

// UnwrapServer returns the underlying server instance.
func (w *airQualitySensorHistoryWrapper) UnwrapServer() AirQualitySensorHistoryServer {
	return w.server
}

// Unwrap implements wrap.Unwrapper and returns the underlying server instance as an unknown type.
func (w *airQualitySensorHistoryWrapper) Unwrap() any {
	return w.UnwrapServer()
}
