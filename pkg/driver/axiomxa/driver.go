package axiomxa

import (
	"context"
	"errors"

	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/driver/axiomxa/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/axiomxa/rpc"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

const DriverName = "axiomxa"

var Factory driver.Factory = factory{}

type factory struct{}

func (f factory) New(services driver.Services) service.Lifecycle {
	d := &Driver{
		announcer: services.Node,
	}
	d.Service = service.New(d.applyConfig, service.WithParser(config.ReadBytes))
	d.logger = services.Logger.Named(DriverName)
	return d
}

func (f factory) AddSupport(supporter node.Supporter) {
	r := rpc.NewAxiomXaDriverServiceRouter()
	supporter.Support(node.Routing(r), node.Clients(rpc.WrapAxiomXaDriverService(r)))
}

type Driver struct {
	*service.Service[config.Root]
	announcer node.Announcer
	logger    *zap.Logger
}

func (d *Driver) applyConfig(ctx context.Context, cfg config.Root) error {
	announcer := node.AnnounceContext(ctx, d.announcer)

	if cfg.HTTP == nil {
		return errors.New("http missing")
	}

	d.logger.Debug("Setting up AxiomXa HTTP connector", zap.String("baseUrl", cfg.HTTP.BaseURL))
	httpImpl := &server{
		config: cfg,
		logger: d.logger.Named("server"),
	}
	announcer.Announce(cfg.Name, node.HasClient(rpc.WrapAxiomXaDriverService(httpImpl)))

	return nil
}
