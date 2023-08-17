package azureiot

import (
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

const DriverType = "azureiot"

var Factory auto.Factory = factory{}

type factory struct{}

func (f factory) New(services auto.Services) service.Lifecycle {
	panic("unimplemented")
}
