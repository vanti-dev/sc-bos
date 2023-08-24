package main

import (
	"context"
	"os"
	"runtime"

	"github.com/vanti-dev/sc-bos/pkg/app"
	"github.com/vanti-dev/sc-bos/pkg/app/sysconf"
	"github.com/vanti-dev/sc-bos/pkg/auto/allautos"
	"github.com/vanti-dev/sc-bos/pkg/driver/alldrivers"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node/alltraits"
	"github.com/vanti-dev/sc-bos/pkg/system/allsystems"
	"github.com/vanti-dev/sc-bos/pkg/testapi"
	"github.com/vanti-dev/sc-bos/pkg/zone/allzones"
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

	if runtime.GOOS != "windows" {
		systemConfig.ConfigDirs = append(systemConfig.ConfigDirs, "/etc/sc-bos")
	}
	systemConfig.DriverFactories = alldrivers.Factories()
	systemConfig.AutoFactories = allautos.Factories()
	systemConfig.SystemFactories = allsystems.Factories()
	systemConfig.ZoneFactories = allzones.Factories()

	err := sysconf.Load(&systemConfig)
	return systemConfig, err
}
