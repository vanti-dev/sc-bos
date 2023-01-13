package driver

import (
	"crypto/tls"
	"net/http"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task"
)

type Services struct {
	Logger          *zap.Logger
	Node            *node.Node // for advertising devices
	Tasks           *task.Group
	ClientTLSConfig *tls.Config // for connecting to other smartcore nodes
	HTTPMux         *http.ServeMux
}

type Driver interface {
}

type Factory interface {
	New(services Services) task.Starter
}
