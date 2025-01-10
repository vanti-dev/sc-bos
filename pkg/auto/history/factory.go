// Package history provides an automation that pulls data from a trait and inserts them into store.
// The automation announces a history api to allow API retrieval of these records filtered by time period.
package history

import (
	"context"
	"errors"
	"fmt"

	"github.com/timshannon/bolthold"
	"go.uber.org/zap"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/wrap"
	"github.com/vanti-dev/sc-bos/internal/util/pgxutil"
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/auto/history/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/historypb"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/meter"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/history"
	"github.com/vanti-dev/sc-bos/pkg/history/apistore"
	"github.com/vanti-dev/sc-bos/pkg/history/boltstore"
	"github.com/vanti-dev/sc-bos/pkg/history/memstore"
	"github.com/vanti-dev/sc-bos/pkg/history/pgxstore"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

var Factory = auto.FactoryFunc(NewAutomation)

func NewAutomation(services auto.Services) service.Lifecycle {
	a := &automation{
		clients:   services.Node,
		announcer: node.NewReplaceAnnouncer(services.Node),
		logger:    services.Logger.Named("history"),

		db: services.Database,

		cohortManagerName: "", // use the default
		cohortManager:     services.CohortManager,
	}
	a.Service = service.New(
		service.MonoApply(a.applyConfig),
		service.WithRetry[config.Root](service.RetryWithLogger(func(logContext service.RetryContext) {
			logContext.LogTo("applyConfig", a.logger)
		})),
	)
	return a
}

type automation struct {
	*service.Service[config.Root]
	clients   node.Clienter
	announcer *node.ReplaceAnnouncer
	logger    *zap.Logger

	db *bolthold.Store

	cohortManagerName string
	cohortManager     node.Remote
}

func (a *automation) applyConfig(ctx context.Context, cfg config.Root) error {
	a.logger.Debug("applying config", zap.Any("storageType", cfg.Storage.Type), zap.Any("trait", cfg.Source.Trait))
	// work out where we're storing the history
	var store history.Store
	switch cfg.Storage.Type {
	case "postgres":
		pool, err := pgxutil.Connect(ctx, cfg.Storage.ConnectConfig)
		if err != nil {
			return err
		}
		opts := []pgxstore.Option{
			pgxstore.WithLogger(a.logger),
		}
		if ttl := cfg.Storage.TTL; ttl != nil {
			if ttl.MaxAge.Duration > 0 {
				opts = append(opts, pgxstore.WithMaxAge(ttl.MaxAge.Duration))
			}
			if ttl.MaxCount > 0 {
				opts = append(opts, pgxstore.WithMaxCount(ttl.MaxCount))
			}
		}
		store, err = pgxstore.SetupStoreFromPool(ctx, cfg.Source.SourceName(), pool, opts...)
		if err != nil {
			return err
		}
	case "memory":
		var opts []memstore.Option
		if ttl := cfg.Storage.TTL; ttl != nil {
			if ttl.MaxAge.Duration > 0 {
				opts = append(opts, memstore.WithMaxAge(ttl.MaxAge.Duration))
			}
			if ttl.MaxCount > 0 {
				opts = append(opts, memstore.WithMaxCount(ttl.MaxCount))
			}
		}
		store = memstore.New(opts...)
	case "api":
		if cfg.Storage.TTL != nil {
			a.logger.Warn("storage.ttl ignored when storage.type is \"api\"")
		}
		name := cfg.Storage.Name
		if name == "" {
			return errors.New("storage.name missing, must exist when storage.type is \"api\"")
		}
		var client gen.HistoryAdminApiClient
		if err := a.clients.Client(&client); err != nil {
			return err
		}
		store = apistore.New(client, name, cfg.Source.SourceName())
	case "hub":
		if cfg.Storage.TTL != nil {
			a.logger.Warn("storage.ttl ignored when storage.type is \"hub\"")
		}
		conn, err := a.cohortManager.Connect(ctx)
		if err != nil {
			return err
		}
		client := gen.NewHistoryAdminApiClient(conn)
		store = apistore.New(client, a.cohortManagerName, cfg.Source.SourceName())
	case "bolt":
		var err error
		opts := []boltstore.Option{
			boltstore.WithLogger(a.logger),
		}
		if ttl := cfg.Storage.TTL; ttl != nil {
			if ttl.MaxAge.Duration > 0 {
				opts = append(opts, boltstore.WithMaxAge(ttl.MaxAge.Duration))
			}
			if ttl.MaxCount > 0 {
				opts = append(opts, boltstore.WithMaxCount(ttl.MaxCount))
			}
		}
		store, err = boltstore.NewFromDb(ctx, a.db, cfg.Source.SourceName(), opts...)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported storage type %s", cfg.Storage.Type)
	}

	// work out where we're getting the records from
	var serverClient wrap.ServiceUnwrapper
	payloads := make(chan []byte)
	var collect func(context.Context, config.Source, chan<- []byte)
	switch cfg.Source.Trait {
	case trait.AirQualitySensor:
		serverClient = gen.WrapAirQualitySensorHistory(historypb.NewAirQualitySensorServer(store))
		collect = a.collectAirQualityChanges
	case trait.AirTemperature:
		serverClient = gen.WrapAirTemperatureHistory(historypb.NewAirTemperatureServer(store))
		collect = a.collectAirTemperatureChanges
	case trait.Electric:
		serverClient = gen.WrapElectricHistory(historypb.NewElectricServer(store))
		collect = a.collectElectricDemandChanges
	case meter.TraitName:
		serverClient = gen.WrapMeterHistory(historypb.NewMeterServer(store))
		collect = a.collectMeterReadingChanges
	case trait.OccupancySensor:
		serverClient = gen.WrapOccupancySensorHistory(historypb.NewOccupancySensorServer(store))
		collect = a.collectOccupancyChanges
	case statuspb.TraitName:
		serverClient = gen.WrapStatusHistory(historypb.NewStatusServer(store))
		collect = a.collectCurrentStatusChanges
	default:
		return fmt.Errorf("unsupported trait %s", cfg.Source.Trait)
	}

	// each time the source emits, we append it to the store
	go func() {
		defer close(payloads)
		for {
			select {
			case <-ctx.Done():
				return
			case payload := <-payloads:
				_, err := store.Append(ctx, payload)
				if err != nil {
					a.logger.Warn("storage failed", zap.Error(err))
				}
			}
		}
	}()

	announce := a.announcer.Replace(ctx)
	// announce the trait too to ensure its services get added to the router before the collect routine starts
	announce.Announce(cfg.Source.Name, node.HasClient(serverClient), node.HasTrait(cfg.Source.Trait))

	go collect(ctx, *cfg.Source, payloads)

	return nil
}
