// Code generated by protoc-gen-router. DO NOT EDIT.

package gen

import (
	context "context"
	fmt "fmt"
	router "github.com/smart-core-os/sc-golang/pkg/router"
	grpc "google.golang.org/grpc"
)

// PriorityApiRouter is a PriorityApiServer that allows routing named requests to specific PriorityApiClient
type PriorityApiRouter struct {
	UnimplementedPriorityApiServer

	router.Router
}

// compile time check that we implement the interface we need
var _ PriorityApiServer = (*PriorityApiRouter)(nil)

func NewPriorityApiRouter(opts ...router.Option) *PriorityApiRouter {
	return &PriorityApiRouter{
		Router: router.NewRouter(opts...),
	}
}

// WithPriorityApiClientFactory instructs the router to create a new
// client the first time Get is called for that name.
func WithPriorityApiClientFactory(f func(name string) (PriorityApiClient, error)) router.Option {
	return router.WithFactory(func(name string) (any, error) {
		return f(name)
	})
}

func (r *PriorityApiRouter) Register(server *grpc.Server) {
	RegisterPriorityApiServer(server, r)
}

// Add extends Router.Add to panic if client is not of type PriorityApiClient.
func (r *PriorityApiRouter) Add(name string, client any) any {
	if !r.HoldsType(client) {
		panic(fmt.Sprintf("not correct type: client of type %T is not a PriorityApiClient", client))
	}
	return r.Router.Add(name, client)
}

func (r *PriorityApiRouter) HoldsType(client any) bool {
	_, ok := client.(PriorityApiClient)
	return ok
}

func (r *PriorityApiRouter) AddPriorityApiClient(name string, client PriorityApiClient) PriorityApiClient {
	res := r.Add(name, client)
	if res == nil {
		return nil
	}
	return res.(PriorityApiClient)
}

func (r *PriorityApiRouter) RemovePriorityApiClient(name string) PriorityApiClient {
	res := r.Remove(name)
	if res == nil {
		return nil
	}
	return res.(PriorityApiClient)
}

func (r *PriorityApiRouter) GetPriorityApiClient(name string) (PriorityApiClient, error) {
	res, err := r.Get(name)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.(PriorityApiClient), nil
}

func (r *PriorityApiRouter) ClearPriorityEntry(ctx context.Context, request *ClearPriorityValueRequest) (*ClearPriorityValueResponse, error) {
	child, err := r.GetPriorityApiClient(request.Name)
	if err != nil {
		return nil, err
	}

	return child.ClearPriorityEntry(ctx, request)
}
