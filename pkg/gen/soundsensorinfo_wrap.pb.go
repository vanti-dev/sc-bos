// Code generated by protoc-gen-wrapper. DO NOT EDIT.

package gen

import (
	wrap "github.com/smart-core-os/sc-golang/pkg/wrap"
	grpc "google.golang.org/grpc"
)

// WrapSoundSensorInfo	adapts a SoundSensorInfoServer	and presents it as a SoundSensorInfoClient
func WrapSoundSensorInfo(server SoundSensorInfoServer) *SoundSensorInfoWrapper {
	conn := wrap.ServerToClient(SoundSensorInfo_ServiceDesc, server)
	client := NewSoundSensorInfoClient(conn)
	return &SoundSensorInfoWrapper{
		SoundSensorInfoClient: client,
		server:                server,
		conn:                  conn,
		desc:                  SoundSensorInfo_ServiceDesc,
	}
}

type SoundSensorInfoWrapper struct {
	SoundSensorInfoClient

	server SoundSensorInfoServer
	conn   grpc.ClientConnInterface
	desc   grpc.ServiceDesc
}

// UnwrapServer returns the underlying server instance.
func (w *SoundSensorInfoWrapper) UnwrapServer() SoundSensorInfoServer {
	return w.server
}

// Unwrap implements wrap.Unwrapper and returns the underlying server instance as an unknown type.
func (w *SoundSensorInfoWrapper) Unwrap() any {
	return w.UnwrapServer()
}

func (w *SoundSensorInfoWrapper) UnwrapService() (grpc.ClientConnInterface, grpc.ServiceDesc) {
	return w.conn, w.desc
}
