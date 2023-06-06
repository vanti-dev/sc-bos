package area

import (
	"context"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"github.com/vanti-dev/sc-bos/pkg/task/serviceapi"
	"github.com/vanti-dev/sc-bos/pkg/zone"
	"github.com/vanti-dev/sc-bos/pkg/zone/area/config"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/electric"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/hvac"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/lighting"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/meter"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/mode"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/occupancy"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/status"
)

// DefaultFeatures lists all the default features for an area.
var DefaultFeatures = []zone.Factory{
	electric.Feature,
	hvac.Feature,
	lighting.Feature,
	meter.Feature,
	mode.Feature,
	occupancy.Feature,
	status.Feature,
}

// Factory builds a generic area using DefaultFeatures.
var Factory = FactoryWithFeatures(DefaultFeatures...)

// FactoryWithFeatures returns an area with the given features.
func FactoryWithFeatures(features ...zone.Factory) zone.Factory {
	return zone.FactoryFunc(func(services zone.Services) service.Lifecycle {
		f := &Area{
			services: services,
			features: features,
		}
		f.Service = service.New(service.MonoApply(f.applyConfig))
		return f
	})
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
	services.Devices = &zone.Devices{}

	featureImpls := make([]service.Lifecycle, len(a.features))
	for i, feature := range a.features {
		featureImpls[i] = feature.New(services)
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
	for _, impl := range featureImpls {
		_, err := impl.Configure(cfg.Raw)
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
