package driver

import (
	"crypto/tls"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task"
	"go.uber.org/zap"
)

type Services struct {
	Logger          *zap.Logger
	Node            *node.Node // for advertising devices
	Tasks           *task.Group
	ClientTLSConfig *tls.Config // for connecting to other smartcore nodes
}

type Driver interface {
}

type Factory interface {
	New(services Services) task.Starter
}
