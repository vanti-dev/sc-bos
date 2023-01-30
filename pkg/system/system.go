// Package system and sub packages add optional features to a controller.
package system

import (
	"github.com/timshannon/bolthold"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/node"
)

type Services struct {
	Logger   *zap.Logger
	Node     *node.Node // for advertising devices
	Database *bolthold.Store
}

type Factory interface {
	New(services Services) service.Lifecycle
}
