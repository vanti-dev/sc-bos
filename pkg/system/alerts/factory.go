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
	"github.com/vanti-dev/sc-bos/pkg/system/alerts/pgxalerts"
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
	s := &System{
		name:      services.Node.Name(),
		announcer: services.Node,
	}
	s.Service = service.New(service.MonoApply(s.applyConfig))
	return s
}

func Register(supporter node.Supporter) {
	alertApiRouter := gen.NewAlertApiRouter()
	alertAdminRouter := gen.NewAlertAdminApiRouter()
	supporter.Support(
		node.Routing(alertApiRouter), node.Clients(gen.WrapAlertApi(alertApiRouter)),
		node.Routing(alertAdminRouter), node.Clients(gen.WrapAlertAdminApi(alertAdminRouter)),
	)
}

type System struct {
	*service.Service[config.Root]

	name      string
	announcer node.Announcer
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
	default:
		return fmt.Errorf("unsuported storage type %s", cfg.Storage.Type)
	}

	return nil
}
