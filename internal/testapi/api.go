package testapi

import (
	"context"
	"sync"

	"github.com/vanti-dev/bsp-ew/pkg/gen"
)

type API struct {
	gen.UnimplementedTestApiServer

	m    sync.RWMutex
	data string
}

func (api *API) GetTest(ctx context.Context, request *gen.GetTestRequest) (*gen.Test, error) {
	api.m.RLock()
	defer api.m.RUnlock()

	data := api.data
	return &gen.Test{Data: data}, nil
}

func (api *API) UpdateTest(ctx context.Context, request *gen.UpdateTestRequest) (*gen.Test, error) {
	api.m.Lock()
	defer api.m.Unlock()

	data := request.GetTest().GetData()
	api.data = data
	return &gen.Test{Data: data}, nil
}

func NewAPI() *API {
	api := &API{
		data: "default data",
	}

	return api
}
