// Code generated by protoc-gen-router. DO NOT EDIT.

package gen

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"

	"github.com/smart-core-os/sc-golang/pkg/router"
)

// StatusApiRouter is a gen.StatusApiServer that allows routing named requests to specific gen.StatusApiClient
type StatusApiRouter struct {
	UnimplementedStatusApiServer

	router.Router
}

// compile time check that we implement the interface we need
var _ StatusApiServer = (*StatusApiRouter)(nil)

func NewStatusApiRouter(opts ...router.Option) *StatusApiRouter {
	return &StatusApiRouter{
		Router: router.NewRouter(opts...),
	}
}

// WithStatusApiClientFactory instructs the router to create a new
// client the first time Get is called for that name.
func WithStatusApiClientFactory(f func(name string) (StatusApiClient, error)) router.Option {
	return router.WithFactory(func(name string) (any, error) {
		return f(name)
	})
}

func (r *StatusApiRouter) Register(server *grpc.Server) {
	RegisterStatusApiServer(server, r)
}

// Add extends Router.Add to panic if client is not of type gen.StatusApiClient.
func (r *StatusApiRouter) Add(name string, client any) any {
	if !r.HoldsType(client) {
		panic(fmt.Sprintf("not correct type: client of type %T is not a gen.StatusApiClient", client))
	}
	return r.Router.Add(name, client)
}

func (r *StatusApiRouter) HoldsType(client any) bool {
	_, ok := client.(StatusApiClient)
	return ok
}

func (r *StatusApiRouter) AddStatusApiClient(name string, client StatusApiClient) StatusApiClient {
	res := r.Add(name, client)
	if res == nil {
		return nil
	}
	return res.(StatusApiClient)
}

func (r *StatusApiRouter) RemoveStatusApiClient(name string) StatusApiClient {
	res := r.Remove(name)
	if res == nil {
		return nil
	}
	return res.(StatusApiClient)
}

func (r *StatusApiRouter) GetStatusApiClient(name string) (StatusApiClient, error) {
	res, err := r.Get(name)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.(StatusApiClient), nil
}

func (r *StatusApiRouter) GetCurrentStatus(ctx context.Context, request *GetCurrentStatusRequest) (*StatusLog, error) {
	child, err := r.GetStatusApiClient(request.Name)
	if err != nil {
		return nil, err
	}

	return child.GetCurrentStatus(ctx, request)
}

func (r *StatusApiRouter) PullCurrentStatus(request *PullCurrentStatusRequest, server StatusApi_PullCurrentStatusServer) error {
	child, err := r.GetStatusApiClient(request.Name)
	if err != nil {
		return err
	}

	// so we can cancel our forwarding request if we can't send responses to our caller
	reqCtx, reqDone := context.WithCancel(server.Context())
	// issue the request
	stream, err := child.PullCurrentStatus(reqCtx, request)
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

		var msg *PullCurrentStatusResponse
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
