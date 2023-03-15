// Code generated by protoc-gen-router. DO NOT EDIT.

package gen

import (
	context "context"
	fmt "fmt"
	router "github.com/smart-core-os/sc-golang/pkg/router"
	grpc "google.golang.org/grpc"
	io "io"
)

// ColorApiRouter is a ColorApiServer that allows routing named requests to specific ColorApiClient
type ColorApiRouter struct {
	UnimplementedColorApiServer

	router.Router
}

// compile time check that we implement the interface we need
var _ ColorApiServer = (*ColorApiRouter)(nil)

func NewColorApiRouter(opts ...router.Option) *ColorApiRouter {
	return &ColorApiRouter{
		Router: router.NewRouter(opts...),
	}
}

// WithColorApiClientFactory instructs the router to create a new
// client the first time Get is called for that name.
func WithColorApiClientFactory(f func(name string) (ColorApiClient, error)) router.Option {
	return router.WithFactory(func(name string) (any, error) {
		return f(name)
	})
}

func (r *ColorApiRouter) Register(server *grpc.Server) {
	RegisterColorApiServer(server, r)
}

// Add extends Router.Add to panic if client is not of type ColorApiClient.
func (r *ColorApiRouter) Add(name string, client any) any {
	if !r.HoldsType(client) {
		panic(fmt.Sprintf("not correct type: client of type %T is not a ColorApiClient", client))
	}
	return r.Router.Add(name, client)
}

func (r *ColorApiRouter) HoldsType(client any) bool {
	_, ok := client.(ColorApiClient)
	return ok
}

func (r *ColorApiRouter) AddColorApiClient(name string, client ColorApiClient) ColorApiClient {
	res := r.Add(name, client)
	if res == nil {
		return nil
	}
	return res.(ColorApiClient)
}

func (r *ColorApiRouter) RemoveColorApiClient(name string) ColorApiClient {
	res := r.Remove(name)
	if res == nil {
		return nil
	}
	return res.(ColorApiClient)
}

func (r *ColorApiRouter) GetColorApiClient(name string) (ColorApiClient, error) {
	res, err := r.Get(name)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.(ColorApiClient), nil
}

func (r *ColorApiRouter) GetColor(ctx context.Context, request *GetColorRequest) (*Color, error) {
	child, err := r.GetColorApiClient(request.Name)
	if err != nil {
		return nil, err
	}

	return child.GetColor(ctx, request)
}

func (r *ColorApiRouter) UpdateColor(ctx context.Context, request *UpdateColorRequest) (*Color, error) {
	child, err := r.GetColorApiClient(request.Name)
	if err != nil {
		return nil, err
	}

	return child.UpdateColor(ctx, request)
}

func (r *ColorApiRouter) PullColor(request *PullColorRequest, server ColorApi_PullColorServer) error {
	child, err := r.GetColorApiClient(request.Name)
	if err != nil {
		return err
	}

	// so we can cancel our forwarding request if we can't send responses to our caller
	reqCtx, reqDone := context.WithCancel(server.Context())
	// issue the request
	stream, err := child.PullColor(reqCtx, request)
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

		var msg *PullColorResponse
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
