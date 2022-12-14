package auto

import (
	"github.com/timshannon/bolthold"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Services struct {
	Logger       *zap.Logger
	Node         *node.Node // for advertising devices
	Database     *bolthold.Store
	GRPCServices grpc.ServiceRegistrar // for registering non-routed services
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
