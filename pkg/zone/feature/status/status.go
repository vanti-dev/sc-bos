package status

import (
	"context"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"github.com/vanti-dev/sc-bos/pkg/zone"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/status/config"
)

var Feature = zone.FactoryFunc(func(services zone.Services) service.Lifecycle {
	f := &feature{
		announce: services.Node,
		devices:  services.Devices,
		clients:  services.Node,
		logger:   services.Logger,
	}
	f.Service = service.New(service.MonoApply(f.applyConfig))
	return f
})

type feature struct {
	*service.Service[config.Root]
	announce node.Announcer
	devices  *zone.Devices
	clients  node.Clienter
	logger   *zap.Logger
}

func (f *feature) applyConfig(ctx context.Context, cfg config.Root) error {
	announce := node.AnnounceContext(ctx, f.announce)
	logger := f.logger.With(zap.String("zone", cfg.Name))

	if len(cfg.StatusLogs) > 0 || cfg.StatusLogAll {
		var client gen.StatusApiClient
		if err := f.clients.Client(&client); err != nil {
			return err
		}

		f.devices.Add(cfg.StatusLogs...)
		if cfg.StatusLogAll {
			go func() {
				select {
				case <-ctx.Done():
					return
				case <-f.devices.Frozen():
					names := f.namesThatImplementTrait(ctx, statuspb.TraitName, f.devices.Names()...)
					group := &Group{
						client: client,
						names:  names,
						logger: logger,
					}
					announce.Announce(cfg.Name, node.HasTrait(statuspb.TraitName, node.WithClients(gen.WrapStatusApi(group))))
				}
			}()
		} else {
			group := &Group{
				client: client,
				names:  cfg.StatusLogs,
				logger: logger,
			}
			announce.Announce(cfg.Name, node.HasTrait(statuspb.TraitName, node.WithClients(gen.WrapStatusApi(group))))
		}
	}

	return nil
}

func (f *feature) namesThatImplementTrait(ctx context.Context, tn trait.Name, names ...string) []string {
	var mdClient traits.MetadataApiClient
	if err := f.clients.Client(&mdClient); err != nil {
		f.logger.Warn("cannot discover status devices, metadata api client not supported")
		return nil
	}

	var res []string
	for _, name := range names {
		md, err := mdClient.GetMetadata(ctx, &traits.GetMetadataRequest{Name: name})
		if err != nil {
			continue
		}
		for _, tmd := range md.Traits {
			if tmd.Name == string(tn) {
				res = append(res, name)
				break
			}
		}
	}
	return res
}
