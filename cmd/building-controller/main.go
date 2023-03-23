package main

import (
	"context"
	"os"

	"github.com/vanti-dev/sc-bos/pkg/app"
	"github.com/vanti-dev/sc-bos/pkg/app/sysconf"
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/node/alltraits"
	"github.com/vanti-dev/sc-bos/pkg/system/allsystems"
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
	systemConfig.DataDir = ".data/building-controller"
	systemConfig.AppConfigFile = "building-controller.local.json"

	systemConfig.DriverFactories = map[string]driver.Factory{}
	systemConfig.AutoFactories = map[string]auto.Factory{}
	systemConfig.SystemFactories = allsystems.Factories()

	err := sysconf.Load(&systemConfig)
	return systemConfig, err
}
