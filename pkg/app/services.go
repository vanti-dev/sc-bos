package app

import (
	"path/filepath"

	"github.com/vanti-dev/sc-bos/pkg/app/files"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"github.com/vanti-dev/sc-bos/pkg/task/serviceapi"
)

func announceServices[M ~map[string]T, T any](c *Controller, name string, services *service.Map, factories M) node.Undo {
	return c.Node.Announce(name, node.HasClient(gen.WrapServicesApi(serviceapi.NewApi(services,
		serviceapi.WithKnownTypesFromMapKeys(factories),
		serviceapi.WithLogger(c.Logger.Named("serviceapi")),
		// results in .data/config/user/{name}/my-service.json
		serviceapi.WithStore(serviceapi.StoreDir(files.Path(c.SystemConfig.DataDir, filepath.Join("config/user", name)))),
		serviceapi.WithMarshaller(serviceapi.MarshalArrayConfig(name)),
	))))
}

func announceAutoServices[M ~map[string]T, T any](c *Controller, services *service.Map, factories M) node.Undo {
	// special because the config name isn't the name we announce as
	return c.Node.Announce("automations", node.HasClient(gen.WrapServicesApi(serviceapi.NewApi(services,
		serviceapi.WithKnownTypesFromMapKeys(factories),
		serviceapi.WithLogger(c.Logger.Named("serviceapi")),
		// results in .data/config/user/{name}/my-service.json
		serviceapi.WithStore(serviceapi.StoreDir(files.Path(c.SystemConfig.DataDir, filepath.Join("config/user", "automations")))),
		serviceapi.WithMarshaller(serviceapi.MarshalArrayConfig("automation")),
	))))
}

func announceSystemServices[M ~map[string]T, T any](c *Controller, services *service.Map, factories M) node.Undo {
	// special because we don't support writing this config, yet
	// todo: support writing system config
	return c.Node.Announce("automations", node.HasClient(gen.WrapServicesApi(serviceapi.NewApi(services,
		serviceapi.WithKnownTypesFromMapKeys(factories),
		serviceapi.WithLogger(c.Logger.Named("serviceapi")),
	))))
}
