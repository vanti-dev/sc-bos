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

type Status string

const (
	StatusInactive Status = "inactive" // Stopped driver. State after calling Stop
	StatusLoading  Status = "loading"  // The driver is loading configuration. State while Configure is running
	StatusActive   Status = "active"   // The driver has valid config and is serving requests. State after calling Start
	StatusError    Status = "error"    // The driver failed to start. If start fails, the driver will be in this state.
)
