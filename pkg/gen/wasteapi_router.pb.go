// Code generated by protoc-gen-router. DO NOT EDIT.

package gen

import (
	context "context"
	fmt "fmt"
	router "github.com/smart-core-os/sc-golang/pkg/router"
	grpc "google.golang.org/grpc"
	io "io"
)

// WasteApiRouter is a WasteApiServer that allows routing named requests to specific WasteApiClient
type WasteApiRouter struct {
	UnimplementedWasteApiServer

	router.Router
}

// compile time check that we implement the interface we need
var _ WasteApiServer = (*WasteApiRouter)(nil)

func NewWasteApiRouter(opts ...router.Option) *WasteApiRouter {
	return &WasteApiRouter{
		Router: router.NewRouter(opts...),
	}
}

// WithWasteApiClientFactory instructs the router to create a new
// client the first time Get is called for that name.
func WithWasteApiClientFactory(f func(name string) (WasteApiClient, error)) router.Option {
	return router.WithFactory(func(name string) (any, error) {
		return f(name)
	})
}

func (r *WasteApiRouter) Register(server *grpc.Server) {
	RegisterWasteApiServer(server, r)
}

// Add extends Router.Add to panic if client is not of type WasteApiClient.
func (r *WasteApiRouter) Add(name string, client any) any {
	if !r.HoldsType(client) {
		panic(fmt.Sprintf("not correct type: client of type %T is not a WasteApiClient", client))
	}
	return r.Router.Add(name, client)
}

func (r *WasteApiRouter) HoldsType(client any) bool {
	_, ok := client.(WasteApiClient)
	return ok
}

func (r *WasteApiRouter) AddWasteApiClient(name string, client WasteApiClient) WasteApiClient {
	res := r.Add(name, client)
	if res == nil {
		return nil
	}
	return res.(WasteApiClient)
}

func (r *WasteApiRouter) RemoveWasteApiClient(name string) WasteApiClient {
	res := r.Remove(name)
	if res == nil {
		return nil
	}
	return res.(WasteApiClient)
}

func (r *WasteApiRouter) GetWasteApiClient(name string) (WasteApiClient, error) {
	res, err := r.Get(name)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.(WasteApiClient), nil
}

func (r *WasteApiRouter) ListWasteRecords(ctx context.Context, request *ListWasteRecordsRequest) (*ListWasteRecordsResponse, error) {
	child, err := r.GetWasteApiClient(request.Name)
	if err != nil {
		return nil, err
	}

	return child.ListWasteRecords(ctx, request)
}

func (r *WasteApiRouter) PullWasteRecords(request *PullWasteRecordsRequest, server WasteApi_PullWasteRecordsServer) error {
	child, err := r.GetWasteApiClient(request.Name)
	if err != nil {
		return err
	}

	// so we can cancel our forwarding request if we can't send responses to our caller
	reqCtx, reqDone := context.WithCancel(server.Context())
	// issue the request
	stream, err := child.PullWasteRecords(reqCtx, request)
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

		var msg *PullWasteRecordsResponse
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