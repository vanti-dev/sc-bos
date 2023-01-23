package auto

import (
	"github.com/timshannon/bolthold"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/vanti-dev/sc-bos/pkg/node"
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

	New(services Services) service.Lifecycle
}

type FactoryFunc func(services Services) service.Lifecycle

func (f FactoryFunc) New(services Services) service.Lifecycle {
	return f(services)
}
