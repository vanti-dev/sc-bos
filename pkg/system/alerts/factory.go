package alerts

import (
	"context"
	"errors"
	"fmt"

	"github.com/vanti-dev/sc-bos/internal/util/pgxutil"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/system"
	"github.com/vanti-dev/sc-bos/pkg/system/alerts/config"
	"github.com/vanti-dev/sc-bos/pkg/system/alerts/hubalerts"
	"github.com/vanti-dev/sc-bos/pkg/system/alerts/pgxalerts"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

var Factory factory

type factory struct{}

func (_ factory) New(services system.Services) service.Lifecycle {
	return NewSystem(services)
}

func NewSystem(services system.Services) *System {
	logger := services.Logger.Named("alerts")
	s := &System{
		name:      services.Node.Name(),
		announcer: services.Node,

		cohortManagerName: "", // use the default
		cohortManager:     services.CohortManager,
	}
	s.Service = service.New(
		service.MonoApply(s.applyConfig),
		service.WithRetry[config.Root](service.RetryWithLogger(func(logContext service.RetryContext) {
			logContext.LogTo("applyConfig", logger)
		})),
	)
	return s
}

type System struct {
	*service.Service[config.Root]

	name      string
	announcer node.Announcer

	cohortManagerName string
	cohortManager     node.Remote
}

func (s *System) applyConfig(ctx context.Context, cfg config.Root) error {
	// using AnnounceContext only makes when using MonoApply, which we are in NewSystem
	announcer := node.AnnounceContext(ctx, s.announcer)

	if cfg.Storage == nil {
		return errors.New("no storage")
	}
	switch cfg.Storage.Type {
	case config.StorageTypePostgres:
		pool, err := pgxutil.Connect(ctx, cfg.Storage.ConnectConfig)
		if err != nil {
			return fmt.Errorf("connect: %w", err)
		}

		server, err := pgxalerts.NewServerFromPool(ctx, pool)
		if err != nil {
			return fmt.Errorf("init: %w", err)
		}

		announcer.Announce(s.name, node.HasClient(
			gen.WrapAlertApi(server),
			gen.WrapAlertAdminApi(server),
		))
	case config.StorageTypeHub:
		server := hubalerts.NewServer("", s.name, s.cohortManager)
		announcer.Announce(s.name, node.HasClient(
			gen.WrapAlertApi(server),
			gen.WrapAlertAdminApi(server),
		))
	default:
		return fmt.Errorf("unsuported storage type %s", cfg.Storage.Type)
	}

	return nil
}
