package driver

import (
	"github.com/vanti-dev/bsp-ew/internal/node"
	"github.com/vanti-dev/bsp-ew/internal/task"
	"go.uber.org/zap"
)

type Services struct {
	Logger *zap.Logger
	Node   *node.Node // for advertising devices
	Tasks  *task.Group
}

type Driver interface {
}

type Factory interface {
	New(services Services) task.Starter
}
