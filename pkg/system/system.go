// Package system and sub packages add optional features to a controller.
package system

import (
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task"
)

type Services struct {
	Logger *zap.Logger
	Node   *node.Node // for advertising devices
}

type Factory interface {
	New(services Services) task.Starter
}
