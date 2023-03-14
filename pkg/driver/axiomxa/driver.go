package axiomxa

import (
	"context"
	"net/http"

	"github.com/olebedev/emitter"
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/driver/axiomxa/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/axiomxa/jsonapi"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
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
	r := gen.NewAxiomXaDriverServiceRouter()
	supporter.Support(node.Routing(r), node.Clients(gen.WrapAxiomXaDriverService(r)))
}

type Driver struct {
	*service.Service[config.Root]
	announcer node.Announcer
	logger    *zap.Logger
}

func (d *Driver) applyConfig(ctx context.Context, cfg config.Root) error {
	announcer := node.AnnounceContext(ctx, d.announcer)

	if cfg.HTTP != nil {
		d.logger.Debug("Setting up AxiomXa HTTP connector", zap.String("baseUrl", cfg.HTTP.BaseURL))
		username, password, err := cfg.HTTP.Credentials()
		if err != nil {
			return err
		}
		tlsConfig, err := cfg.HTTP.TLS.TLSConfig()
		if err != nil {
			return err
		}
		client := jsonapi.NewClient(cfg.HTTP.BaseURL, username, password)
		if tlsConfig != nil {
			client.HTTPClient.Transport = &http.Transport{
				TLSClientConfig: tlsConfig,
			}
		}

		qrServerImpl := &qrServer{
			config: cfg,
			logger: d.logger.Named("server"),
		}
		announcer.Announce(cfg.Name, node.HasClient(gen.WrapAxiomXaDriverService(qrServerImpl)))
	}

	// bus topics will be one of the Keys in bsp-ew.go, like KeyAccessGranted ("AG")
	// The event argument is always of type mps.Fields.
	bus := emitter.New(0)
	if err := d.setupMessagePortServer(ctx, cfg, bus); err != nil {
		return err
	}

	// devices maps axiom devices (and controllers) to smart core names
	devices := devicesFromConfig(cfg.Devices)

	// announce traits that expose this functionality over smart core
	if err := d.announceTraits(ctx, cfg, announcer, bus, devices); err != nil {
		return err
	}

	return nil
}
