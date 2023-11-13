// Package history provides an implementation of HistoryAdminApi backed by a history.Store.
// Enabling this system on a controller will allow history automations to store their records with us if configured to do so.
package history

import (
	"context"
	"errors"
	"fmt"

	"github.com/timshannon/bolthold"
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/internal/util/pgxutil"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/history"
	"github.com/vanti-dev/sc-bos/pkg/history/boltstore"
	"github.com/vanti-dev/sc-bos/pkg/history/pgxstore"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/system"
	"github.com/vanti-dev/sc-bos/pkg/system/history/config"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

var Factory factory

type factory struct{}

func (_ factory) New(services system.Services) service.Lifecycle {
	return NewSystem(services)
}

func (_ factory) AddSupport(supporter node.Supporter) {
	Register(supporter)
}

func NewSystem(services system.Services) *System {
	logger := services.Logger.Named("history")
	s := &System{
		name:      services.Node.Name(),
		announcer: services.Node,
		db:        services.Database,
		logger:    logger,
	}
	s.Service = service.New(
		service.MonoApply(s.applyConfig),
		service.WithRetry[config.Root](service.RetryWithLogger(func(logContext service.RetryContext) {
			logContext.LogTo("applyConfig", logger)
		})),
	)
	return s
}
func Register(supporter node.Supporter) {
	historyAdminApiRouter := gen.NewHistoryAdminApiRouter()
	supporter.Support(
		node.Routing(historyAdminApiRouter), node.Clients(gen.WrapHistoryAdminApi(historyAdminApiRouter)),
	)
}

type System struct {
	*service.Service[config.Root]
	name      string
	announcer node.Announcer
	db        *bolthold.Store

	logger *zap.Logger
}

func (s *System) applyConfig(ctx context.Context, cfg config.Root) error {
	// using AnnounceContext only makes when using MonoApply, which we are in NewSystem
	announcer := node.AnnounceContext(ctx, s.announcer)

	if cfg.Storage == nil {
		return errors.New("no storage")
	}

	var store func(string) history.Store

	switch cfg.Storage.Type {
	case config.StorageTypePostgres:
		pool, err := pgxutil.Connect(ctx, cfg.Storage.ConnectConfig)
		if err != nil {
			return fmt.Errorf("connect: %w", err)
		}

		if err := pgxstore.SetupDB(ctx, pool); err != nil {
			return fmt.Errorf("setup: %w", err)
		}

		store = func(source string) history.Store {
			return pgxstore.NewStoreFromPool(source, pool)
		}
	case config.StorageTypeBolt:
		store = func(source string) history.Store {
			st, err := boltstore.NewFromDb(s.db, source, s.logger)
			if err != nil {
				s.logger.Error("failed to create bolt store", zap.Error(err))
			}
			return st
		}
	default:
		return fmt.Errorf("unsuported storage type %s", cfg.Storage.Type)
	}

	server := &storeServer{store: store}
	announcer.Announce(s.name, node.HasClient(gen.WrapHistoryAdminApi(server)))

	return nil
}
