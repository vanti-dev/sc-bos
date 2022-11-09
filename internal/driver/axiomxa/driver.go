package axiomxa

import (
	"context"
	"errors"
	"github.com/vanti-dev/bsp-ew/internal/driver"
	"github.com/vanti-dev/bsp-ew/internal/driver/axiomxa/config"
	"github.com/vanti-dev/bsp-ew/internal/driver/axiomxa/rpc"
	"github.com/vanti-dev/bsp-ew/internal/node"
	"github.com/vanti-dev/bsp-ew/internal/task"
	"go.uber.org/zap"
)

const DriverName = "axiomxa"

var Factory driver.Factory = factory{}

type factory struct{}

func (f factory) New(services driver.Services) task.Starter {
	d := &Driver{
		announcer: services.Node,
	}
	d.Lifecycle = task.NewLifecycle(d.applyConfig)
	d.Logger = services.Logger.Named(DriverName)
	d.ReadConfig = config.ReadBytes
	return d
}

func (f factory) AddSupport(supporter node.Supporter) {
	r := rpc.NewAxiomXaDriverServiceRouter()
	supporter.Support(node.Routing(r), node.Clients(rpc.WrapAxiomXaDriverService(r)))
}

type Driver struct {
	*task.Lifecycle[config.Root]
	announcer node.Announcer
}

func (d *Driver) applyConfig(_ context.Context, cfg config.Root) error {
	// todo: track announcements and undo them on config update - aka support more than one config update

	if cfg.HTTP == nil {
		return errors.New("http missing")
	}

	d.Logger.Debug("Setting up AxiomXa HTTP connector", zap.String("baseUrl", cfg.HTTP.BaseURL))
	httpImpl := &server{
		config: cfg,
		logger: d.Logger.Named("server"),
	}
	d.announcer.Announce(cfg.Name, node.HasClient(rpc.WrapAxiomXaDriverService(httpImpl)))

	return nil
}
