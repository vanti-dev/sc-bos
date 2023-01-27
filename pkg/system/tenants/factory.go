package tenants

import (
	"context"
	"errors"
	"fmt"

	"github.com/vanti-dev/sc-bos/internal/util/pgxutil"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/system"
	"github.com/vanti-dev/sc-bos/pkg/system/tenants/config"
	"github.com/vanti-dev/sc-bos/pkg/system/tenants/hold"
	"github.com/vanti-dev/sc-bos/pkg/system/tenants/pgxtenants"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"go.uber.org/zap"
)

var Factory = factory{
	server: &hold.Server{}, // shared by all grpc servers
}

type factory struct {
	server *hold.Server
}

func (f factory) New(services system.Services) service.Lifecycle {
	return NewSystem(services)
}

func (f factory) AddSupport(supporter node.Supporter) {
	supporter.Support(node.Api(f.server))
}

func NewSystem(services system.Services) *System {
	s := &System{
		logger: services.Logger.Named("tenants"),
	}
	s.Service = service.New(s.applyConfig)
	return s
}

type System struct {
	*service.Service[config.Root]
	logger *zap.Logger
}

func (s *System) applyConfig(ctx context.Context, cfg config.Root) error {
	if cfg.Storage == nil {
		return errors.New("no storage")
	}
	if cfg.Storage.Type != "postgres" {
		return fmt.Errorf("unsuported storage type %s, want one of [postgres]", cfg.Storage.Type)
	}

	pool, err := pgxutil.Connect(ctx, cfg.Storage.ConnectConfig)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	server, err := pgxtenants.NewServerFromPool(ctx, pool, pgxtenants.WithLogger(s.logger))
	if err != nil {
		return fmt.Errorf("init: %w", err)
	}

	// There's only one tenant api, each time we run we make sure to take over control of it.
	Factory.server.Fill(gen.WrapTenantApi(server))

	return nil
}
