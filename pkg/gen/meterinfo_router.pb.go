// Code generated by protoc-gen-router. DO NOT EDIT.

package gen

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/smart-core-os/sc-golang/pkg/router"
)

// MeterInfoRouter is a gen.MeterInfoServer that allows routing named requests to specific gen.MeterInfoClient
type MeterInfoRouter struct {
	UnimplementedMeterInfoServer

	router.Router
}

// compile time check that we implement the interface we need
var _ MeterInfoServer = (*MeterInfoRouter)(nil)

func NewMeterInfoRouter(opts ...router.Option) *MeterInfoRouter {
	return &MeterInfoRouter{
		Router: router.NewRouter(opts...),
	}
}

// WithMeterInfoClientFactory instructs the router to create a new
// client the first time Get is called for that name.
func WithMeterInfoClientFactory(f func(name string) (MeterInfoClient, error)) router.Option {
	return router.WithFactory(func(name string) (any, error) {
		return f(name)
	})
}

func (r *MeterInfoRouter) Register(server *grpc.Server) {
	RegisterMeterInfoServer(server, r)
}

// Add extends Router.Add to panic if client is not of type gen.MeterInfoClient.
func (r *MeterInfoRouter) Add(name string, client any) any {
	if !r.HoldsType(client) {
		panic(fmt.Sprintf("not correct type: client of type %T is not a gen.MeterInfoClient", client))
	}
	return r.Router.Add(name, client)
}

func (r *MeterInfoRouter) HoldsType(client any) bool {
	_, ok := client.(MeterInfoClient)
	return ok
}

func (r *MeterInfoRouter) AddMeterInfoClient(name string, client MeterInfoClient) MeterInfoClient {
	res := r.Add(name, client)
	if res == nil {
		return nil
	}
	return res.(MeterInfoClient)
}

func (r *MeterInfoRouter) RemoveMeterInfoClient(name string) MeterInfoClient {
	res := r.Remove(name)
	if res == nil {
		return nil
	}
	return res.(MeterInfoClient)
}

func (r *MeterInfoRouter) GetMeterInfoClient(name string) (MeterInfoClient, error) {
	res, err := r.Get(name)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.(MeterInfoClient), nil
}

func (r *MeterInfoRouter) DescribeMeterReading(ctx context.Context, request *DescribeMeterReadingRequest) (*MeterReadingSupport, error) {
	child, err := r.GetMeterInfoClient(request.Name)
	if err != nil {
		return nil, err
	}

	return child.DescribeMeterReading(ctx, request)
}
