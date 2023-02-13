package zone

import (
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"go.uber.org/zap"
)

type Services struct {
	Logger *zap.Logger
	Node   *node.Node
}

type Factory interface {
	New(Services) service.Lifecycle
}

type FactoryFunc func(services Services) service.Lifecycle

func (f FactoryFunc) New(services Services) service.Lifecycle {
	return f(services)
}
