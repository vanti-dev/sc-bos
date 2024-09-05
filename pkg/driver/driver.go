package driver

import (
	"context"
	"crypto/tls"
	"net/http"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

type Services struct {
	Logger          *zap.Logger
	Node            *node.Node  // for advertising devices
	ClientTLSConfig *tls.Config // for connecting to other smartcore nodes
	HTTPMux         *http.ServeMux
	ConfigUpdater   ConfigUpdater
}

type Factory interface {
	New(services Services) service.Lifecycle
}

type ConfigUpdater interface {
	UpdateConfig(ctx context.Context, config []byte) error
}
