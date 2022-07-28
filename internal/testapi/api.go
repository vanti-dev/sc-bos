package testapi

import (
	"context"
	"sync"

	"github.com/vanti-dev/bsp-ew/internal/testgen"
)

type API struct {
	testgen.UnimplementedTestApiServer

	m    sync.RWMutex
	data string
}

func (api *API) GetTest(ctx context.Context, request *testgen.GetTestRequest) (*testgen.Test, error) {
	api.m.RLock()
	defer api.m.RUnlock()

	data := api.data
	return &testgen.Test{Data: data}, nil
}

func (api *API) UpdateTest(ctx context.Context, request *testgen.UpdateTestRequest) (*testgen.Test, error) {
	api.m.Lock()
	defer api.m.Unlock()

	data := request.GetTest().GetData()
	api.data = data
	return &testgen.Test{Data: data}, nil
}

func NewAPI() *API {
	api := &API{
		data: "default data",
	}

	return api
}
