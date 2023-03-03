package driver

import (
	"crypto/tls"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/task/service"

	"github.com/vanti-dev/sc-bos/pkg/node"
)

type Services struct {
	Logger          *zap.Logger
	Node            *node.Node  // for advertising devices
	ClientTLSConfig *tls.Config // for connecting to other smartcore nodes
	HTTPMux         *mux.Router
}

type Driver interface {
}

type Factory interface {
	New(services Services) service.Lifecycle
}
