// Package system and sub packages add optional features to a controller.
package system

import (
	"net/http"

	"github.com/timshannon/bolthold"
	"github.com/vanti-dev/sc-bos/pkg/auth/token"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/node"
)

type Services struct {
	DataDir         string
	Logger          *zap.Logger
	Node            *node.Node // for advertising devices
	CohortManager   node.Remote
	Database        *bolthold.Store
	HTTPMux         *http.ServeMux      // to allow systems to serve http requests
	TokenValidators *token.ValidatorSet // to allow systems to contribute towards client validation
}

type Factory interface {
	New(services Services) service.Lifecycle
}
