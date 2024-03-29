// Code generated by protoc-gen-wrapper. DO NOT EDIT.

package gen

import (
	context "context"
	wrap "github.com/smart-core-os/sc-golang/pkg/wrap"
	grpc "google.golang.org/grpc"
)

// WrapUdmiService	adapts a UdmiServiceServer	and presents it as a UdmiServiceClient
func WrapUdmiService(server UdmiServiceServer) UdmiServiceClient {
	return &udmiServiceWrapper{server}
}

type udmiServiceWrapper struct {
	server UdmiServiceServer
}

// compile time check that we implement the interface we need
var _ UdmiServiceClient = (*udmiServiceWrapper)(nil)

// UnwrapServer returns the underlying server instance.
func (w *udmiServiceWrapper) UnwrapServer() UdmiServiceServer {
	return w.server
}

// Unwrap implements wrap.Unwrapper and returns the underlying server instance as an unknown type.
func (w *udmiServiceWrapper) Unwrap() any {
	return w.UnwrapServer()
}

func (w *udmiServiceWrapper) PullControlTopics(ctx context.Context, in *PullControlTopicsRequest, opts ...grpc.CallOption) (UdmiService_PullControlTopicsClient, error) {
	stream := wrap.NewClientServerStream(ctx)
	server := &pullControlTopicsUdmiServiceServerWrapper{stream.Server()}
	client := &pullControlTopicsUdmiServiceClientWrapper{stream.Client()}
	go func() {
		err := w.server.PullControlTopics(in, server)
		stream.Close(err)
	}()
	return client, nil
}

type pullControlTopicsUdmiServiceClientWrapper struct {
	grpc.ClientStream
}

func (w *pullControlTopicsUdmiServiceClientWrapper) Recv() (*PullControlTopicsResponse, error) {
	m := new(PullControlTopicsResponse)
	if err := w.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

type pullControlTopicsUdmiServiceServerWrapper struct {
	grpc.ServerStream
}

func (s *pullControlTopicsUdmiServiceServerWrapper) Send(response *PullControlTopicsResponse) error {
	return s.ServerStream.SendMsg(response)
}

func (w *udmiServiceWrapper) OnMessage(ctx context.Context, req *OnMessageRequest, _ ...grpc.CallOption) (*OnMessageResponse, error) {
	return w.server.OnMessage(ctx, req)
}

func (w *udmiServiceWrapper) PullExportMessages(ctx context.Context, in *PullExportMessagesRequest, opts ...grpc.CallOption) (UdmiService_PullExportMessagesClient, error) {
	stream := wrap.NewClientServerStream(ctx)
	server := &pullExportMessagesUdmiServiceServerWrapper{stream.Server()}
	client := &pullExportMessagesUdmiServiceClientWrapper{stream.Client()}
	go func() {
		err := w.server.PullExportMessages(in, server)
		stream.Close(err)
	}()
	return client, nil
}

type pullExportMessagesUdmiServiceClientWrapper struct {
	grpc.ClientStream
}

func (w *pullExportMessagesUdmiServiceClientWrapper) Recv() (*PullExportMessagesResponse, error) {
	m := new(PullExportMessagesResponse)
	if err := w.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

type pullExportMessagesUdmiServiceServerWrapper struct {
	grpc.ServerStream
}

func (s *pullExportMessagesUdmiServiceServerWrapper) Send(response *PullExportMessagesResponse) error {
	return s.ServerStream.SendMsg(response)
}

func (w *udmiServiceWrapper) GetExportMessage(ctx context.Context, req *GetExportMessageRequest, _ ...grpc.CallOption) (*MqttMessage, error) {
	return w.server.GetExportMessage(ctx, req)
}
