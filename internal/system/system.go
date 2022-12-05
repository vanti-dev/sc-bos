// Package system and sub packages add optional features to a controller.
package system

import (
	"github.com/vanti-dev/bsp-ew/internal/node"
	"github.com/vanti-dev/bsp-ew/internal/task"
	"go.uber.org/zap"
)

type Services struct {
	Logger *zap.Logger
	Node   *node.Node // for advertising devices
}

type Factory interface {
	New(services Services) task.Starter
}
