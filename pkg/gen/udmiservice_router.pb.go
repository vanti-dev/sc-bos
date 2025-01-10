// Code generated by protoc-gen-router. DO NOT EDIT.

package gen

import (
	context "context"
	fmt "fmt"
	router "github.com/smart-core-os/sc-golang/pkg/router"
	grpc "google.golang.org/grpc"
	io "io"
)

// UdmiServiceRouter is a UdmiServiceServer that allows routing named requests to specific UdmiServiceClient
type UdmiServiceRouter struct {
	UnimplementedUdmiServiceServer

	router.Router
}

// compile time check that we implement the interface we need
var _ UdmiServiceServer = (*UdmiServiceRouter)(nil)

func NewUdmiServiceRouter(opts ...router.Option) *UdmiServiceRouter {
	return &UdmiServiceRouter{
		Router: router.NewRouter(opts...),
	}
}

// WithUdmiServiceClientFactory instructs the router to create a new
// client the first time Get is called for that name.
func WithUdmiServiceClientFactory(f func(name string) (UdmiServiceClient, error)) router.Option {
	return router.WithFactory(func(name string) (any, error) {
		return f(name)
	})
}

func (r *UdmiServiceRouter) Register(server grpc.ServiceRegistrar) {
	RegisterUdmiServiceServer(server, r)
}

// Add extends Router.Add to panic if client is not of type UdmiServiceClient.
func (r *UdmiServiceRouter) Add(name string, client any) any {
	if !r.HoldsType(client) {
		panic(fmt.Sprintf("not correct type: client of type %T is not a UdmiServiceClient", client))
	}
	return r.Router.Add(name, client)
}

func (r *UdmiServiceRouter) HoldsType(client any) bool {
	_, ok := client.(UdmiServiceClient)
	return ok
}

func (r *UdmiServiceRouter) AddUdmiServiceClient(name string, client UdmiServiceClient) UdmiServiceClient {
	res := r.Add(name, client)
	if res == nil {
		return nil
	}
	return res.(UdmiServiceClient)
}

func (r *UdmiServiceRouter) RemoveUdmiServiceClient(name string) UdmiServiceClient {
	res := r.Remove(name)
	if res == nil {
		return nil
	}
	return res.(UdmiServiceClient)
}

func (r *UdmiServiceRouter) GetUdmiServiceClient(name string) (UdmiServiceClient, error) {
	res, err := r.Get(name)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.(UdmiServiceClient), nil
}

func (r *UdmiServiceRouter) PullControlTopics(request *PullControlTopicsRequest, server UdmiService_PullControlTopicsServer) error {
	child, err := r.GetUdmiServiceClient(request.Name)
	if err != nil {
		return err
	}

	// so we can cancel our forwarding request if we can't send responses to our caller
	reqCtx, reqDone := context.WithCancel(server.Context())
	// issue the request
	stream, err := child.PullControlTopics(reqCtx, request)
	if err != nil {
		return err
	}

	// send the stream header
	header, err := stream.Header()
	if err != nil {
		return err
	}
	if err = server.SendHeader(header); err != nil {
		return err
	}

	// send all the messages
	// false means the error is from the child, true means the error is from the caller
	var callerError bool
	for {
		// Impl note: we could improve throughput here by issuing the Recv and Send in different goroutines, but we're doing
		// it synchronously until we have a need to change the behaviour

		var msg *PullControlTopicsResponse
		msg, err = stream.Recv()
		if err != nil {
			break
		}

		err = server.Send(msg)
		if err != nil {
			callerError = true
			break
		}
	}

	// err is guaranteed to be non-nil as it's the only way to exit the loop
	if callerError {
		// cancel the request
		reqDone()
		return err
	} else {
		if trailer := stream.Trailer(); trailer != nil {
			server.SetTrailer(trailer)
		}
		if err == io.EOF {
			return nil
		}
		return err
	}
}

func (r *UdmiServiceRouter) OnMessage(ctx context.Context, request *OnMessageRequest) (*OnMessageResponse, error) {
	child, err := r.GetUdmiServiceClient(request.Name)
	if err != nil {
		return nil, err
	}

	return child.OnMessage(ctx, request)
}

func (r *UdmiServiceRouter) PullExportMessages(request *PullExportMessagesRequest, server UdmiService_PullExportMessagesServer) error {
	child, err := r.GetUdmiServiceClient(request.Name)
	if err != nil {
		return err
	}

	// so we can cancel our forwarding request if we can't send responses to our caller
	reqCtx, reqDone := context.WithCancel(server.Context())
	// issue the request
	stream, err := child.PullExportMessages(reqCtx, request)
	if err != nil {
		return err
	}

	// send the stream header
	header, err := stream.Header()
	if err != nil {
		return err
	}
	if err = server.SendHeader(header); err != nil {
		return err
	}

	// send all the messages
	// false means the error is from the child, true means the error is from the caller
	var callerError bool
	for {
		// Impl note: we could improve throughput here by issuing the Recv and Send in different goroutines, but we're doing
		// it synchronously until we have a need to change the behaviour

		var msg *PullExportMessagesResponse
		msg, err = stream.Recv()
		if err != nil {
			break
		}

		err = server.Send(msg)
		if err != nil {
			callerError = true
			break
		}
	}

	// err is guaranteed to be non-nil as it's the only way to exit the loop
	if callerError {
		// cancel the request
		reqDone()
		return err
	} else {
		if trailer := stream.Trailer(); trailer != nil {
			server.SetTrailer(trailer)
		}
		if err == io.EOF {
			return nil
		}
		return err
	}
}

func (r *UdmiServiceRouter) GetExportMessage(ctx context.Context, request *GetExportMessageRequest) (*MqttMessage, error) {
	child, err := r.GetUdmiServiceClient(request.Name)
	if err != nil {
		return nil, err
	}

	return child.GetExportMessage(ctx, request)
}
