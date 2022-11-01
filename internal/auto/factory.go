package auto

import (
	"github.com/vanti-dev/bsp-ew/internal/node"
	"github.com/vanti-dev/bsp-ew/internal/task"
	"go.uber.org/zap"
)

type Services struct {
	Logger *zap.Logger
	Node   *node.Node // for advertising devices
}

// Factory constructs new automation instances.
type Factory interface {
	// note this is an interface, not a func type so that the controller can check for other interfaces, like GrpcApi.

	New(services Services) task.Starter
}

type FactoryFunc func(services Services) task.Starter

func (f FactoryFunc) New(services Services) task.Starter {
	return f(services)
}
