package area

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/app/sysconf"
	"github.com/vanti-dev/sc-bos/pkg/block"
	"github.com/vanti-dev/sc-bos/pkg/block/mdblock"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"github.com/vanti-dev/sc-bos/pkg/task/serviceapi"
	"github.com/vanti-dev/sc-bos/pkg/zone"
	"github.com/vanti-dev/sc-bos/pkg/zone/area/config"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/airquality"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/electric"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/enterleave"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/hvac"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/lighting"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/meter"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/mode"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/occupancy"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/openclose"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/status"
)

// DefaultFeatures lists all the default features for an area.
var DefaultFeatures = []zone.Factory{
	electric.Feature,
	enterleave.Feature,
	hvac.Feature,
	lighting.Feature,
	meter.Feature,
	mode.Feature,
	occupancy.Feature,
	openclose.Feature,
	status.Feature,
	airquality.Feature,
}

// Factory builds a generic area using DefaultFeatures.
var Factory = FactoryWithFeatures(DefaultFeatures...)

// FactoryWithFeatures returns an area with the given features.
func FactoryWithFeatures(features ...zone.Factory) zone.Factory {
	return factory{features: features}
}

type factory struct {
	features []zone.Factory
}

func (f factory) New(services zone.Services) service.Lifecycle {
	services.Logger = services.Logger.Named("area")
	a := &Area{
		services: services,
		features: f.features,
	}
	a.Service = service.New(service.MonoApply(a.applyConfig))
	return a
}

// ConfigBlocks implements sysconf.BlockSource2 supporting blocks from the nested services the zone is hosting.
func (f factory) ConfigBlocks(cfg *sysconf.Config) []block.Block {
	// todo: a lot of this logic is shared with sysconf, figure out if there's a way to share it

	defaultBlocks := []block.Block{
		{Path: []string{"disabled"}},
	}
	blocks := []block.Block{
		{Path: []string{"metadata"}, Blocks: mdblock.Categories},
		{
			Path:         []string{"drivers"},
			Key:          "name",
			TypeKey:      "type",
			BlocksByType: cfg.DriverConfigBlocks(),
			Blocks:       defaultBlocks,
		},
	}

	for _, feature := range f.features {
		switch source := any(feature).(type) {
		case sysconf.BlockSource:
			blocks = append(blocks, source.ConfigBlocks()...)
		case sysconf.BlockSource2:
			blocks = append(blocks, source.ConfigBlocks(cfg)...)
		}
	}

	return blocks
}

type Area struct {
	*service.Service[config.Root]
	services zone.Services
	features []zone.Factory
}

func (a *Area) applyConfig(ctx context.Context, cfg config.Root) error {
	announce := node.AnnounceContext(ctx, a.services.Node)
	if cfg.Metadata != nil {
		announce.Announce(cfg.Name, node.HasMetadata(cfg.Metadata))
	}

	services := a.services
	services.Logger = a.services.Logger.With(zap.String("zone", cfg.Name))
	services.Devices = &zone.Devices{}

	type serviceConfig struct {
		service.Lifecycle
		cfg []byte
	}
	serviceConfigs := make([]serviceConfig, 0, len(a.features)+len(cfg.Drivers))
	featureImpls := make([]service.Lifecycle, 0, len(a.features)+len(cfg.Drivers))
	for _, feature := range a.features {
		impl := feature.New(services)
		serviceConfigs = append(serviceConfigs, serviceConfig{Lifecycle: impl, cfg: cfg.Raw})
		featureImpls = append(featureImpls, impl)
	}

	driverServices := driver.Services{
		Logger:          services.Logger.Named("driver"),
		Node:            services.Node,
		ClientTLSConfig: services.ClientTLSConfig,
		HTTPMux:         services.HTTPMux,
	}
	for _, d := range cfg.Drivers {
		f, ok := a.services.DriverFactories[d.Type]
		if !ok {
			return fmt.Errorf("unsupported driver type %v", d.Type)
		}
		impl := f.New(driverServices)
		serviceConfigs = append(serviceConfigs, serviceConfig{Lifecycle: impl, cfg: d.Raw})
		featureImpls = append(featureImpls, impl)
	}

	// make the zone area implement the ServicesApi
	m := service.NewMapOf(featureImpls)
	api := serviceapi.NewApi(m)
	announce.Announce(cfg.Name, node.HasClient(gen.WrapServicesApi(api)))

	// stop all features if the area is stopped
	go func() {
		<-ctx.Done()
		for _, impl := range featureImpls {
			_, _ = impl.Stop()
		}
	}()

	// configure and start all the features
	// might want to split the configure and start steps to pick up on any config errors early?
	for _, impl := range serviceConfigs {
		_, err := impl.Configure(impl.cfg)
		if err != nil {
			// change this if we ever want to support incomplete area deployments
			return err
		}
		_, err = impl.Start()
		if err != nil {
			// change this if we ever want to support incomplete area deployments
			return err
		}
	}

	for _, impl := range featureImpls {
		a.waitUntilLoaded(ctx, impl)
	}
	services.Devices.Freeze()

	return nil
}

func (a *Area) waitUntilLoaded(ctx context.Context, impl service.Lifecycle) {
	ctx, stop := context.WithCancel(ctx)
	defer stop()

	settled := func(state service.State) bool {
		return state.Active && !state.Loading ||
			!state.Active && state.Err != nil
	}

	state, changes := impl.StateAndChanges(ctx)
	if settled(state) {
		return
	}
	for state := range changes {
		if settled(state) {
			return
		}
	}
}
