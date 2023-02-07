package main

import (
	"context"
	"embed"
	"os"

	"github.com/vanti-dev/sc-bos/pkg/app"
	"github.com/vanti-dev/sc-bos/pkg/app/sysconf"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/driver/proxy"
	"github.com/vanti-dev/sc-bos/pkg/node/alltraits"
	"github.com/vanti-dev/sc-bos/pkg/system"
	"github.com/vanti-dev/sc-bos/pkg/system/authn"
	"github.com/vanti-dev/sc-bos/pkg/testapi"

	"github.com/vanti-dev/sc-bos/pkg/auth/policy"

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
	systemConfig.ConfigDirs = []string{".data/gateway-01"}
	systemConfig.AppConfigFile = "gateway.local.json"

	pol, err := policy.FromFS(policyFS)
	if err != nil {
		return systemConfig, err
	}
	systemConfig.Policy = pol

	systemConfig.DriverFactories = map[string]driver.Factory{
		proxy.DriverName: proxy.Factory,
	}
	systemConfig.SystemFactories = map[string]system.Factory{
		"authn": authn.Factory(),
	}

	err = sysconf.Load(&systemConfig)
	return systemConfig, err
}

//go:embed policy
var policyFS embed.FS
