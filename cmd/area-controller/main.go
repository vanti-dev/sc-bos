package main

import (
	"context"
	"flag"
	"github.com/vanti-dev/bsp-ew/internal/system"
	"github.com/vanti-dev/bsp-ew/internal/system/alerts"
	"os"

	"github.com/smart-core-os/sc-golang/pkg/trait/brightnesssensor"
	"github.com/vanti-dev/bsp-ew/internal/auto/export"
	"github.com/vanti-dev/bsp-ew/internal/driver/axiomxa"

	"github.com/smart-core-os/sc-golang/pkg/trait/airtemperature"
	"github.com/smart-core-os/sc-golang/pkg/trait/light"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/onoff"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"
	"github.com/vanti-dev/bsp-ew/internal/app"
	"github.com/vanti-dev/bsp-ew/internal/auto"
	"github.com/vanti-dev/bsp-ew/internal/auto/lights"
	"github.com/vanti-dev/bsp-ew/internal/driver"
	"github.com/vanti-dev/bsp-ew/internal/driver/bacnet"
	"github.com/vanti-dev/bsp-ew/internal/node"
	"github.com/vanti-dev/bsp-ew/internal/testapi"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.uber.org/zap"
)

var systemConfig app.SystemConfig

func init() {
	flag.StringVar(&systemConfig.ListenGRPC, "listen-grpc", ":23557", "address (host:port) to host a Smart Core gRPC server on")
	flag.StringVar(&systemConfig.ListenHTTPS, "listen-https", ":443", "address (host:port) to host a HTTPS server on")
	flag.StringVar(&systemConfig.DataDir, "data-dir", ".data/area-controller-01", "path to local data storage directory")
	flag.StringVar(&systemConfig.StaticDir, "static-dir", "", "(optional) path to directory to host static files over HTTP from")

	flag.BoolVar(&systemConfig.DisablePolicy, "insecure-disable-policy", false, "Insecure! Disable checking requests against the security policy. This option opens up the server to any request.")
	flag.BoolVar(&systemConfig.LocalOAuth, "local-auth", false, "Enable issuing password tokens based on credentials found in users.json")
	flag.BoolVar(&systemConfig.TenantOAuth, "tenant-auth", false, "enable issuing client tokens based on credentials found in tenants.json or verified via the enrollment manager node")
}

func main() {
	os.Exit(app.RunUntilInterrupt(run))
}

func run(ctx context.Context) error {
	flag.Parse()
	systemConfig.LocalConfigFileName = "area-controller.local.json"
	systemConfig.Logger = zap.NewDevelopmentConfig()
	systemConfig.DriverFactories = map[string]driver.Factory{
		axiomxa.DriverName: axiomxa.Factory,
		bacnet.DriverName:  bacnet.Factory,
	}
	systemConfig.AutoFactories = map[string]auto.Factory{
		lights.AutoType: lights.Factory,
		"export-mqtt":   export.MQTTFactory,
	}
	systemConfig.SystemFactories = map[string]system.Factory{
		"alerts": alerts.Factory,
	}

	controller, err := app.Bootstrap(ctx, systemConfig)
	if err != nil {
		return err
	}

	addNodeAPIs(controller.Node)

	gen.RegisterTestApiServer(controller.GRPC, testapi.NewAPI())

	return controller.Run(ctx)
}

func addNodeAPIs(supporter node.Supporter) {
	{
		r := airtemperature.NewApiRouter()
		c := airtemperature.WrapApi(r)
		supporter.Support(node.Routing(r), node.Clients(c))
	}
	{
		r := brightnesssensor.NewApiRouter()
		c := brightnesssensor.WrapApi(r)
		supporter.Support(node.Routing(r), node.Clients(c))
	}
	{
		r := light.NewApiRouter()
		c := light.WrapApi(r)
		supporter.Support(node.Routing(r), node.Clients(c))
	}
	{
		r := occupancysensor.NewApiRouter()
		c := occupancysensor.WrapApi(r)
		supporter.Support(node.Routing(r), node.Clients(c))
	}
	{
		r := onoff.NewApiRouter()
		c := onoff.WrapApi(r)
		supporter.Support(node.Routing(r), node.Clients(c))
	}
	{
		r := parent.NewApiRouter()
		c := parent.WrapApi(r)
		supporter.Support(node.Routing(r), node.Clients(c))
	}
}
