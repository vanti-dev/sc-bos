package meter

import (
	"context"
	"path"
	"time"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-bos/pkg/gentrait/meter"
	"github.com/smart-core-os/sc-bos/pkg/node"
	"github.com/smart-core-os/sc-bos/pkg/task/service"
	"github.com/smart-core-os/sc-bos/pkg/zone"
	"github.com/smart-core-os/sc-bos/pkg/zone/feature/meter/config"
)

var Feature = zone.FactoryFunc(func(services zone.Services) service.Lifecycle {
	services.Logger = services.Logger.Named("meter")
	f := &feature{
		announcer: node.NewReplaceAnnouncer(services.Node),
		devices:   services.Devices,
		clients:   services.Node,
		logger:    services.Logger,
	}
	f.Service = service.New(service.MonoApply(f.applyConfig), service.WithParser(config.ParseConfig))
	return f
})

type feature struct {
	*service.Service[config.Root]
	announcer *node.ReplaceAnnouncer
	devices   *zone.Devices
	clients   node.ClientConner
	logger    *zap.Logger
}

func (f *feature) applyConfig(ctx context.Context, cfg config.Root) error {
	announce := f.announcer.Replace(ctx)
	logger := f.logger

	conn := f.clients.ClientConn()
	apiClient := gen.NewMeterApiClient(conn)
	infoClient := gen.NewMeterInfoClient(conn)
	historyClient := gen.NewMeterHistoryClient(conn)
	announceGroup := func(name string, devices []string) {
		if len(devices) == 0 {
			return
		}

		group := &Group{
			apiClient:        apiClient,
			infoClient:       infoClient,
			historyApiClient: historyClient,
			names:            devices,
			logger:           logger,

			now: time.Now,

			historyBackupConf: cfg.HistoryBackup,
		}
		f.devices.Add(devices...)
		announce.Announce(name, node.HasTrait(meter.TraitName, node.WithClients(gen.WrapMeterApi(group), gen.WrapMeterInfo(group))))
	}

	announceGroup(cfg.Name, cfg.Meters)
	for name, meters := range cfg.MeterGroups {
		announceGroup(path.Join(cfg.Name, name), meters)
	}

	return nil
}
