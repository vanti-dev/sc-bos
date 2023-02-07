package main

import (
	"context"
	"os"

	"github.com/vanti-dev/sc-bos/pkg/app"
	"github.com/vanti-dev/sc-bos/pkg/app/sysconf"
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/auto/export"
	"github.com/vanti-dev/sc-bos/pkg/auto/lights"
	"github.com/vanti-dev/sc-bos/pkg/auto/udmi"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/driver/axiomxa"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet"
	"github.com/vanti-dev/sc-bos/pkg/driver/mock"
	"github.com/vanti-dev/sc-bos/pkg/driver/xovis"
	"github.com/vanti-dev/sc-bos/pkg/node/alltraits"
	"github.com/vanti-dev/sc-bos/pkg/system"
	"github.com/vanti-dev/sc-bos/pkg/system/alerts"
	"github.com/vanti-dev/sc-bos/pkg/system/authn"
	"github.com/vanti-dev/sc-bos/pkg/system/hub"
	"github.com/vanti-dev/sc-bos/pkg/system/publications"
	"github.com/vanti-dev/sc-bos/pkg/system/tenants"
	"github.com/vanti-dev/sc-bos/pkg/testapi"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

func main() {
	os.Exit(app.RunUntilInterrupt(run))
}

func run(ctx context.Context) error {
	systemConfig, err := loadSystemConfig()
	if err != nil {
		return err
	}

	controller, err := app.Bootstrap(ctx, systemConfig)
	if err != nil {
		return err
	}

	alltraits.AddSupport(controller.Node)

	gen.RegisterTestApiServer(controller.GRPC, testapi.NewAPI())

	return controller.Run(ctx)
}

func loadSystemConfig() (sysconf.Config, error) {
	systemConfig := sysconf.Default()
	systemConfig.DataDir = ".data/area-controller-01"
	systemConfig.AppConfigFile = "area-controller.local.json"

	systemConfig.DriverFactories = map[string]driver.Factory{
		axiomxa.DriverName: axiomxa.Factory,
		bacnet.DriverName:  bacnet.Factory,
		mock.DriverName:    mock.Factory,
		xovis.DriverName:   xovis.Factory,
	}
	systemConfig.AutoFactories = map[string]auto.Factory{
		lights.AutoType: lights.Factory,
		"export-mqtt":   export.MQTTFactory,
		udmi.AutoType:   udmi.Factory,
	}
	systemConfig.SystemFactories = map[string]system.Factory{
		"alerts":       alerts.Factory,
		"authn":        authn.Factory(),
		"hub":          hub.Factory(),
		"tenants":      tenants.Factory,
		"publications": publications.Factory,
	}

	err := sysconf.Load(&systemConfig)
	return systemConfig, err
}
