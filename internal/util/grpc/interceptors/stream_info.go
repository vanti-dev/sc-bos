package interceptors

import (
	"google.golang.org/grpc"
)

type StreamServerInfoSource interface {
	StreamServerInfo(method string) (grpc.StreamServerInfo, bool)
}

// CorrectStreamInfo will modify the StreamServerInfo of the incoming request to match the one looked up from the source.
//
// This is useful for UnknownServiceHandler implementations, as they will always report as bidirectional streams
// even if the actual service definition says otherwise.
// Use this interceptor as the first in the chain to ensure that the StreamServerInfo is correct for other
// interceptors.
func CorrectStreamInfo(source StreamServerInfoSource) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		newInfo, ok := source.StreamServerInfo(info.FullMethod)
		if ok {
			info.IsServerStream = newInfo.IsServerStream
			info.IsClientStream = newInfo.IsClientStream
		}

		return handler(srv, ss)
	}
}
