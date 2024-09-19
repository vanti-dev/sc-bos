package tenants

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/internal/util/pgxutil"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/system"
	"github.com/vanti-dev/sc-bos/pkg/system/tenants/config"
	"github.com/vanti-dev/sc-bos/pkg/system/tenants/hold"
	"github.com/vanti-dev/sc-bos/pkg/system/tenants/pgxtenants"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
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

func NewSystem(services system.Services) *System {
	s := &System{
		hubNode: services.CohortManager,
		logger:  services.Logger.Named("tenants"),
	}
	s.Service = service.New(
		s.applyConfig,
		service.WithRetry[config.Root](service.RetryWithLogger(func(logContext service.RetryContext) {
			logContext.LogTo("applyConfig", s.logger)
		})),
	)
	return s
}

type System struct {
	*service.Service[config.Root]
	hubNode node.Remote
	logger  *zap.Logger
}

func (s *System) applyConfig(ctx context.Context, cfg config.Root) error {
	if cfg.Storage == nil {
		return errors.New("no storage")
	}
	switch cfg.Storage.Type {
	case config.StorageTypeProxy:
		conn, err := s.hubNode.Connect(ctx)
		if err != nil {
			return err
		}
		Factory.server.Fill(gen.NewTenantApiClient(conn))
	case config.StorageTypePostgres:
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
	default:
		return fmt.Errorf("unsuported storage type %s", cfg.Storage.Type)
	}

	return nil
}
