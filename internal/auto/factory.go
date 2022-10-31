package auto

import (
	"github.com/vanti-dev/bsp-ew/internal/node"
	"go.uber.org/zap"
)

type Services struct {
	Logger *zap.Logger
	Node   *node.Node // for advertising devices
}

// Factory constructs new automation instances.
type Factory interface {
	// note this is an interface, not a func type so that the controller can check for other interfaces, like GrpcApi.

	New(services *Services) Starter
}

type FactoryFunc func(services *Services) Starter

func (f FactoryFunc) New(services *Services) Starter {
	return f(services)
}
