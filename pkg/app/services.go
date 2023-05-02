package app

import (
	"context"
	"fmt"
	"path/filepath"

	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/app/files"
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

func (c *Controller) startDrivers() (*service.Map, error) {
	ctxServices := driver.Services{
		Logger:          c.Logger.Named("driver"),
		Node:            c.Node,
		ClientTLSConfig: c.ClientTLSConfig,
		HTTPMux:         c.Mux,
	}

	m := service.NewMap(func(kind string) (service.Lifecycle, error) {
		f, ok := c.SystemConfig.DriverFactories[kind]
		if !ok {
			return nil, fmt.Errorf("unsupported driver type %v", kind)
		}
		return f.New(ctxServices), nil
	}, service.IdIsRequired)

	var allErrs error
	for _, cfg := range c.ControllerConfig.Drivers {
		_, _, err := m.Create(cfg.Name, cfg.Type, service.State{Active: !cfg.Disabled, Config: cfg.Raw})
		allErrs = multierr.Append(allErrs, err)
	}
	return m, allErrs
}

func (c *Controller) startAutomations() (*service.Map, error) {
	ctxServices := auto.Services{
		Logger:       c.Logger.Named("auto"),
		Node:         c.Node,
		Database:     c.Database,
		GRPCServices: c.GRPC,
	}

	m := service.NewMap(func(kind string) (service.Lifecycle, error) {
		f, ok := c.SystemConfig.AutoFactories[kind]
		if !ok {
			return nil, fmt.Errorf("unsupported automation type %v", kind)
		}
		return f.New(ctxServices), nil
	}, service.IdIsRequired)

	var allErrs error
	for _, cfg := range c.ControllerConfig.Automation {
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
		DataDir:         c.SystemConfig.DataDir,
		Logger:          c.Logger.Named("system"),
		Node:            c.Node,
		GRPCEndpoint:    grpcEndpoint,
		Database:        c.Database,
		HTTPMux:         c.Mux,
		TokenValidators: c.TokenValidators,
		GRPCCerts:       c.GRPCCerts,
		PrivateKey:      c.PrivateKey,
		CohortManager:   c.ManagerConn,
		ClientTLSConfig: c.ClientTLSConfig,
	}
	m := service.NewMap(func(kind string) (service.Lifecycle, error) {
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

func (c *Controller) startZones() (*service.Map, error) {
	ctxServices := zone.Services{
		Logger: c.Logger.Named("auto"),
		Node:   c.Node,
	}

	m := service.NewMap(func(kind string) (service.Lifecycle, error) {
		f, ok := c.SystemConfig.ZoneFactories[kind]
		if !ok {
			return nil, fmt.Errorf("unsupported zone type %v", kind)
		}
		return f.New(ctxServices), nil
	}, service.IdIsRequired)

	var allErrs error
	for _, cfg := range c.ControllerConfig.Zones {
		_, _, err := m.Create(cfg.Name, cfg.Type, service.State{Active: !cfg.Disabled, Config: cfg.Raw})
		allErrs = multierr.Append(allErrs, err)
	}
	return m, allErrs
}

func logServiceMapChanges(ctx context.Context, logger *zap.Logger, m *service.Map) {
	known := map[string]func(){}
	changes := m.Listen(ctx)
	for _, record := range m.Values() {
		ctx, stop := context.WithCancel(ctx)
		known[record.Id] = stop
		record := record
		go logServiceRecordChanges(ctx, logger, record)
	}
	for change := range changes {
		if change.OldValue == nil && change.NewValue != nil {
			// add
			if _, ok := known[change.NewValue.Id]; ok {
				continue // deal with potential race between Listen and Values
			}
			ctx, stop := context.WithCancel(ctx)
			known[change.NewValue.Id] = stop
			go logServiceRecordChanges(ctx, logger, change.NewValue)
		} else if change.OldValue != nil && change.NewValue == nil {
			// remove
			stop, ok := known[change.OldValue.Id]
			if !ok {
				continue
			}
			delete(known, change.OldValue.Id)
			stop()
		}
	}
}

func logServiceRecordChanges(ctx context.Context, logger *zap.Logger, r *service.Record) {
	logger = logger.With(zap.String("id", r.Id), zap.String("kind", r.Kind))
	state, changes := r.Service.StateAndChanges(ctx)
	lastMode := ""
	logMode := func(change service.State) {
		mode := ""
		switch {
		case !change.Active && change.Err != nil:
			mode = "error"
		case !change.Active:
			mode = "Stopped"
		case change.Loading:
			mode = "Loading"
		case change.Active:
			mode = "Running"
		}
		if mode == lastMode {
			return
		}
		switch mode {
		case "error":
			logger.Warn("Failed to load", zap.Error(change.Err))
		case "":
			return
		case "Stopped":
			if lastMode == "" {
				logger.Debug("Created")
			} else {
				logger.Debug(mode)
			}
		default:
			logger.Debug(mode)
		}
		lastMode = mode
	}
	logMode(state)
	for change := range changes {
		logMode(change)
	}
}

func announceServices[M ~map[string]T, T any](c *Controller, name string, services *service.Map, factories M) node.Undo {
	client := gen.WrapServicesApi(serviceapi.NewApi(services,
		serviceapi.WithKnownTypesFromMapKeys(factories),
		serviceapi.WithLogger(c.Logger.Named("serviceapi")),
		// results in .data/config/user/{name}/my-service.json
		serviceapi.WithStore(serviceapi.StoreDir(files.Path(c.SystemConfig.DataDir, filepath.Join("config/user", name)))),
		serviceapi.WithMarshaller(serviceapi.MarshalArrayConfig(name)),
	))
	return node.UndoAll(
		c.Node.Announce(name, node.HasClient(client)),
		c.Node.Announce(filepath.Join(c.Node.Name(), name), node.HasClient(client)),
	)
}

func announceAutoServices[M ~map[string]T, T any](c *Controller, services *service.Map, factories M) node.Undo {
	// special because the config name isn't the name we announce as
	client := gen.WrapServicesApi(serviceapi.NewApi(services,
		serviceapi.WithKnownTypesFromMapKeys(factories),
		serviceapi.WithLogger(c.Logger.Named("serviceapi")),
		// results in .data/config/user/{name}/my-service.json
		serviceapi.WithStore(serviceapi.StoreDir(files.Path(c.SystemConfig.DataDir, filepath.Join("config/user", "automations")))),
		serviceapi.WithMarshaller(serviceapi.MarshalArrayConfig("automation")),
	))
	return node.UndoAll(
		c.Node.Announce("automations", node.HasClient(client)),
		c.Node.Announce(filepath.Join(c.Node.Name(), "automations"), node.HasClient(client)),
	)
}

func announceSystemServices[M ~map[string]T, T any](c *Controller, services *service.Map, factories M) node.Undo {
	// special because we don't support writing this config, yet
	// todo: support writing system config
	client := gen.WrapServicesApi(serviceapi.NewApi(services,
		serviceapi.WithKnownTypesFromMapKeys(factories),
		serviceapi.WithLogger(c.Logger.Named("serviceapi")),
	))
	return node.UndoAll(
		c.Node.Announce("systems", node.HasClient(client)),
		c.Node.Announce(filepath.Join(c.Node.Name(), "systems"), node.HasClient(client)),
	)
}
