// Code generated by protoc-gen-wrapper. DO NOT EDIT.

package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
)

// WrapColorInfo	adapts a ColorInfoServer	and presents it as a ColorInfoClient
func WrapColorInfo(server ColorInfoServer) ColorInfoClient {
	return &colorInfoWrapper{server}
}

type colorInfoWrapper struct {
	server ColorInfoServer
}

// compile time check that we implement the interface we need
var _ ColorInfoClient = (*colorInfoWrapper)(nil)

// UnwrapServer returns the underlying server instance.
func (w *colorInfoWrapper) UnwrapServer() ColorInfoServer {
	return w.server
}

// Unwrap implements wrap.Unwrapper and returns the underlying server instance as an unknown type.
func (w *colorInfoWrapper) Unwrap() any {
	return w.UnwrapServer()
}

func (w *colorInfoWrapper) DescribeColor(ctx context.Context, req *DescribeColorRequest, _ ...grpc.CallOption) (*ColorSupport, error) {
	return w.server.DescribeColor(ctx, req)
}
