package app

import (
	"context"
	"fmt"
	"path"

	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/smart-core-os/sc-golang/pkg/wrap"
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/system"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"github.com/vanti-dev/sc-bos/pkg/task/serviceapi"
	"github.com/vanti-dev/sc-bos/pkg/zone"
)

// addFactorySupport is used to register factories with a node to expose custom factory APIs.
// This checks each value in m and if that value has an API, via node.SelfSupporter, then it is registered with s.
func addFactorySupport[M ~map[K]F, K comparable, F any](s node.Supporter, m M) {
	for _, factory := range m {
		if api, ok := any(factory).(node.SelfSupporter); ok {
			api.AddSupport(s)
		}
	}
}

func (c *Controller) startDrivers(configs []driver.RawConfig) (*service.Map, error) {
	ctxServices := driver.Services{
		Logger:          c.Logger.Named("driver"),
		Node:            c.Node,
		ClientTLSConfig: c.ClientTLSConfig,
		HTTPMux:         c.Mux,
	}

	m := service.NewMap(func(id, kind string) (service.Lifecycle, error) {
		driverServices := ctxServices
		driverServices.Config = &serviceConfigStore{store: c.ControllerConfig.Drivers(), id: id}

		f, ok := c.SystemConfig.DriverFactories[kind]
		if !ok {
			return nil, fmt.Errorf("unsupported driver type %v", kind)
		}
		return f.New(ctxServices), nil
	}, service.IdIsRequired)

	var allErrs error
	for _, cfg := range configs {
		_, _, err := m.Create(cfg.Name, cfg.Type, service.State{Active: !cfg.Disabled, Config: cfg.Raw})
		allErrs = multierr.Append(allErrs, err)
	}
	return m, allErrs
}

func (c *Controller) startAutomations(configs []auto.RawConfig) (*service.Map, error) {
	ctxServices := auto.Services{
		Logger:          c.Logger.Named("auto"),
		Node:            c.Node,
		Database:        c.Database,
		GRPCServices:    c.GRPC,
		CohortManager:   c.ManagerConn,
		ClientTLSConfig: c.ClientTLSConfig,
	}

	m := service.NewMap(func(id, kind string) (service.Lifecycle, error) {
		autoServices := ctxServices
		autoServices.Config = &serviceConfigStore{store: c.ControllerConfig.Automations(), id: id}

		f, ok := c.SystemConfig.AutoFactories[kind]
		if !ok {
			return nil, fmt.Errorf("unsupported automation type %v", kind)
		}
		return f.New(ctxServices), nil
	}, service.IdIsRequired)

	var allErrs error
	for _, cfg := range configs {
		_, _, err := m.Create(cfg.Name, cfg.Type, service.State{Active: !cfg.Disabled, Config: cfg.Raw})
		allErrs = multierr.Append(allErrs, err)
	}
	return m, allErrs
}

func (c *Controller) startSystems() (*service.Map, error) {
	grpcEndpoint, err := c.grpcEndpoint()
	if err != nil {
		return nil, err
	}
	ctxServices := system.Services{
		ConfigDirs:       c.SystemConfig.ConfigDirs,
		DataDir:          c.SystemConfig.DataDir,
		Logger:           c.Logger.Named("system"),
		Node:             c.Node,
		GRPCEndpoint:     grpcEndpoint,
		Database:         c.Database,
		HTTPMux:          c.Mux,
		TokenValidators:  c.TokenValidators,
		ReflectionServer: c.ReflectionServer,
		GRPCCerts:        c.GRPCCerts,
		PrivateKey:       c.PrivateKey,
		CohortManager:    c.ManagerConn,
		ClientTLSConfig:  c.ClientTLSConfig,
	}
	m := service.NewMap(func(_, kind string) (service.Lifecycle, error) {
		f, ok := c.SystemConfig.SystemFactories[kind]
		if !ok {
			return nil, fmt.Errorf("unsupported system type %v", kind)
		}
		return f.New(ctxServices), nil
	}, service.IdIsKind)

	var allErrs error
	for kind, cfg := range c.SystemConfig.Systems {
		_, _, err := m.Create("", kind, service.State{Active: !cfg.Disabled, Config: cfg.Raw})
		allErrs = multierr.Append(allErrs, err)
	}
	return m, allErrs
}

func (c *Controller) startZones(configs []zone.RawConfig) (*service.Map, error) {
	ctxServices := zone.Services{
		Logger:          c.Logger.Named("zone"),
		Node:            c.Node,
		ClientTLSConfig: c.ClientTLSConfig,
		HTTPMux:         c.Mux,
		DriverFactories: c.SystemConfig.DriverFactories,
	}

	m := service.NewMap(func(id, kind string) (service.Lifecycle, error) {
		zoneServices := ctxServices
		zoneServices.Config = &serviceConfigStore{store: c.ControllerConfig.Zones(), id: id}

		f, ok := c.SystemConfig.ZoneFactories[kind]
		if !ok {
			return nil, fmt.Errorf("unsupported zone type %v", kind)
		}
		return f.New(ctxServices), nil
	}, service.IdIsRequired)

	var allErrs error
	for _, cfg := range configs {
		_, _, err := m.Create(cfg.Name, cfg.Type, service.State{Active: !cfg.Disabled, Config: cfg.Raw})
		allErrs = multierr.Append(allErrs, err)
	}
	return m, allErrs
}

func logServiceMapChanges(ctx context.Context, logger *zap.Logger, m *service.Map) {
	now, changes := m.GetAndListenState(ctx)
	for _, record := range now {
		logServiceRecordChange(logger, nil, record)
	}
	for change := range changes {
		logServiceRecordChange(logger, change.OldValue, change.NewValue)
	}
}

func logServiceRecordChange(logger *zap.Logger, oldVal, newVal *service.StateRecord) {
	switch {
	case newVal != nil:
		logger = logger.With(zap.String("id", newVal.Id), zap.String("kind", newVal.Kind))
	case oldVal != nil:
		logger = logger.With(zap.String("id", oldVal.Id), zap.String("kind", oldVal.Kind))
	}
	switch {
	case oldVal == nil && newVal != nil: // do nothing
	case newVal == nil: // removed
		logger.Debug("Removed")
	case oldVal == nil: // created
		logger.Debug("Created", zap.Bool("active", newVal.State.Active), zap.Bool("loading", newVal.State.Loading), zap.Error(newVal.State.Err))
	case !newVal.State.Active && newVal.State.Err != nil && oldVal.State.Err == nil: // error
		logger.Warn("Failed to load", zap.Error(newVal.State.Err))
	case oldVal.State.Active && !newVal.State.Active: // stopped
		logger.Debug("Stopped", zap.Error(newVal.State.Err))
	case newVal.State.Active && newVal.State.Loading && !newVal.State.NextAttemptTime.IsZero(): // retrying
		// rely on the service itself to log any issues that caused the retry
	case newVal.State.Active && newVal.State.Loading: // loading
		logger.Debug("Loading")
	case !oldVal.State.Active && newVal.State.Active || oldVal.State.Loading && !newVal.State.Loading: // started
		logger.Debug("Started")
	default:
		type state struct {
			Active, Loading bool
			Error           error
		}
		oldState := state{Active: oldVal.State.Active, Loading: oldVal.State.Loading, Error: oldVal.State.Err}
		newState := state{Active: newVal.State.Active, Loading: newVal.State.Loading, Error: newVal.State.Err}
		logger.Debug("Updated", zap.Any("old", oldState), zap.Any("new", newState))
	}
}

func announceServices[M ~map[string]T, T any](c *Controller, name string, services *service.Map, factories M, store serviceapi.Store) node.Undo {
	client := gen.WrapServicesApi(serviceapi.NewApi(services,
		serviceapi.WithKnownTypesFromMapKeys(factories),
		serviceapi.WithLogger(c.Logger.Named("serviceapi")),
		serviceapi.WithStore(store),
	))
	return announceNodeClient(c.Node, name, client)
}

func announceAutoServices[M ~map[string]T, T any](c *Controller, services *service.Map, factories M) node.Undo {
	// special because the config name isn't the name we announce as
	client := gen.WrapServicesApi(serviceapi.NewApi(services,
		serviceapi.WithKnownTypesFromMapKeys(factories),
		serviceapi.WithLogger(c.Logger.Named("serviceapi")),
		serviceapi.WithStore(c.ControllerConfig.Automations()),
	))
	return announceNodeClient(c.Node, "automations", client)
}

func announceSystemServices[M ~map[string]T, T any](c *Controller, services *service.Map, factories M) node.Undo {
	// special because we don't support writing this config, yet
	// todo: support writing system config
	client := gen.WrapServicesApi(serviceapi.NewApi(services,
		serviceapi.WithKnownTypesFromMapKeys(factories),
		serviceapi.WithLogger(c.Logger.Named("serviceapi")),
	))
	return announceNodeClient(c.Node, "systems", client)
}

func announceNodeClient(n *node.Node, base string, client wrap.ServiceUnwrapper) node.Undo {
	var undos []node.Undo
	undos = append(undos, n.Announce(base, node.HasClient(client)))
	if n.Name() != "" {
		undos = append(undos, n.Announce(path.Join(n.Name(), base), node.HasClient(client)))
	}
	return node.UndoAll(undos...)
}

type serviceConfigStore struct {
	store serviceapi.Store
	id    string
}

func (s *serviceConfigStore) UpdateConfig(ctx context.Context, data []byte) error {
	return s.store.SaveConfig(ctx, s.id, "", data)
}
