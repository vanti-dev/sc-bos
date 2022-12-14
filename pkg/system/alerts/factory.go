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
	"github.com/vanti-dev/sc-bos/pkg/task"
	"sync"
)

var Factory factory

type factory struct{}

func (_ factory) New(services system.Services) task.Starter {
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
	s.Lifecycle = task.NewLifecycle(s.applyConfig)
	s.Logger = services.Logger.Named("alerts")
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
	*task.Lifecycle[config.Root]

	name      string
	announcer node.Announcer

	mu   sync.Mutex
	undo node.Undo
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

	server, err := pgxalerts.NewServerFromPool(ctx, pool)
	if err != nil {
		return fmt.Errorf("init: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.undo != nil {
		s.undo()
	}
	s.undo = s.announcer.Announce(s.name, node.HasClient(
		gen.WrapAlertApi(server),
		gen.WrapAlertAdminApi(server),
	))

	return nil
}
