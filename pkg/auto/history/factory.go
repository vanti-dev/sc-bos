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
		clients:  services.Node,
		announce: services.Node,
		logger:   services.Logger.Named("history"),

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
	clients  node.Clienter
	announce node.Announcer
	logger   *zap.Logger

	db *bolthold.Store

	cohortManagerName string
	cohortManager     node.Remote
}

func (a *automation) applyConfig(ctx context.Context, cfg config.Root) error {
	a.logger.Info("applying config", zap.Any("storageType", cfg.Storage.Type), zap.Any("trait", cfg.Source.Trait))
	// work out where we're storing the history
	var store history.Store
	switch cfg.Storage.Type {
	case "postgres":
		pool, err := pgxutil.Connect(ctx, cfg.Storage.ConnectConfig)
		if err != nil {
			return err
		}
		store, err = pgxstore.SetupStoreFromPool(ctx, cfg.Source.SourceName(), pool)
		if err != nil {
			return err
		}
	case "memory":
		store = memstore.New()
	case "api":
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
		conn, err := a.cohortManager.Connect(ctx)
		if err != nil {
			return err
		}
		client := gen.NewHistoryAdminApiClient(conn)
		store = apistore.New(client, a.cohortManagerName, cfg.Source.SourceName())
	case "bolt":
		var err error
		store, err = boltstore.NewFromDb(a.db, cfg.Source.SourceName())
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported storage type %s", cfg.Storage.Type)
	}

	// work out where we're getting the records from
	var serverClient any
	payloads := make(chan []byte)
	switch cfg.Source.Trait {
	case trait.Electric:
		serverClient = gen.WrapElectricHistory(historypb.NewElectricServer(store))
		go a.collectElectricDemandChanges(ctx, *cfg.Source, payloads)
	case meter.TraitName:
		serverClient = gen.WrapMeterHistory(historypb.NewMeterServer(store))
		go a.collectMeterReadingChanges(ctx, *cfg.Source, payloads)
	case trait.OccupancySensor:
		serverClient = gen.WrapOccupancySensorHistory(historypb.NewOccupancySensorServer(store))
		go a.collectOccupancyChanges(ctx, *cfg.Source, payloads)
	case trait.AirQualitySensor:
		serverClient = gen.WrapAirQualitySensorHistory(historypb.NewAirQualitySensorServer(store))
		go a.collectAirQualityChanges(ctx, *cfg.Source, payloads)
	case statuspb.TraitName:
		serverClient = gen.WrapStatusHistory(historypb.NewStatusServer(store))
		go a.collectCurrentStatusChanges(ctx, *cfg.Source, payloads)
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

	announce := node.AnnounceContext(ctx, a.announce)
	// we could technically announce this as a trait client, but there's no real need for that
	announce.Announce(cfg.Source.Name, node.HasClient(serverClient))

	return nil
}
