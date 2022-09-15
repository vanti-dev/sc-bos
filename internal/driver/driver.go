package driver

import (
	"context"
	"encoding/json"

	"github.com/vanti-dev/bsp-ew/internal/node"
	"go.uber.org/zap"
)

type Services struct {
	Logger *zap.Logger
	Node   *node.Node // for advertising devices
}

type Driver interface {
	Name() string
}

type Factory func(ctx context.Context, services Services, config json.RawMessage) (Driver, error)
