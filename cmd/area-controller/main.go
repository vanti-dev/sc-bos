package main

import (
	"context"
	"flag"
	"os"

	"github.com/smart-core-os/sc-golang/pkg/trait/light"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"
	"github.com/vanti-dev/bsp-ew/internal/app"
	"github.com/vanti-dev/bsp-ew/internal/driver"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali"
	"github.com/vanti-dev/bsp-ew/internal/testapi"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.uber.org/zap"
)

var (
	flagStaticDir string
)

var systemConfig app.SystemConfig

func init() {
	flag.StringVar(&systemConfig.ListenGRPC, "listen-grpc", ":23557", "address (host:port) to host a Smart Core gRPC server on")
	flag.StringVar(&systemConfig.ListenHTTPS, "listen-https", ":443", "address (host:port) to host a HTTPS server on")
	flag.StringVar(&systemConfig.DataDir, "data-dir", ".data/area-controller-01", "path to local data storage directory")
	flag.StringVar(&flagStaticDir, "static-dir", "ui/dist", "path for HTTP static resources")

	flag.BoolVar(&systemConfig.DisablePolicy, "insecure-disable-policy", false, "Insecure! Disable checking requests against the security policy. This option opens up the server to any request.")
	flag.BoolVar(&systemConfig.LocalOAuth, "local-auth", false, "Enable issuing password tokens based on credentials found in users.json")
	flag.BoolVar(&systemConfig.TenantOAuth, "tenant-auth", false, "enable issuing client tokens based on credentials found in tenants.json or verified via the enrollment manager node")
}

func main() {
	os.Exit(app.RunUntilInterrupt(run))
}

func run(ctx context.Context) error {
	flag.Parse()
	systemConfig.Logger = zap.NewDevelopmentConfig()
	systemConfig.DriverFactories = map[string]driver.Factory{
		tc3dali.DriverName: tc3dali.Factory,
	}
	systemConfig.Routers = []app.Router{
		light.NewApiRouter(),
		occupancysensor.NewApiRouter(),
		parent.NewApiRouter(),
	}

	controller, err := app.Bootstrap(ctx, systemConfig)
	if err != nil {
		return err
	}

	gen.RegisterTestApiServer(controller.GRPC, testapi.NewAPI())

	return controller.Run(ctx)
}
