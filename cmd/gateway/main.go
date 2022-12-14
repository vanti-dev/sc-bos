package main

import (
	"context"
	"embed"
	"flag"
	"os"

	"github.com/smart-core-os/sc-golang/pkg/trait/airtemperature"
	"github.com/smart-core-os/sc-golang/pkg/trait/brightnesssensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/light"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"

	"github.com/vanti-dev/sc-bos/pkg/app"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/driver/proxy"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/testapi"

	"github.com/smart-core-os/sc-golang/pkg/trait/onoff"
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/auth/policy"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

var systemConfig app.SystemConfig

func init() {
	flag.StringVar(&systemConfig.ListenGRPC, "listen-grpc", ":23557", "address (host:port) to host a Smart Core gRPC server on")
	flag.StringVar(&systemConfig.ListenHTTPS, "listen-https", ":443", "address (host:port) to host a HTTPS server on")
	flag.StringVar(&systemConfig.DataDir, "data-dir", ".data/gateway", "path to local data storage directory")

	flag.BoolVar(&systemConfig.DisablePolicy, "insecure-disable-policy", false, "Insecure! Disable checking requests against the security policy. This option opens up the server to any request.")
	flag.BoolVar(&systemConfig.LocalOAuth, "local-auth", false, "Enable issuing password tokens based on credentials found in users.json")
	flag.BoolVar(&systemConfig.TenantOAuth, "tenant-auth", true, "enable issuing client tokens based on credentials found in tenants.json or verified via the enrollment manager node")
}

func main() {
	os.Exit(app.RunUntilInterrupt(run))
}

func run(ctx context.Context) error {
	flag.Parse()

	pol, err := policy.FromFS(policyFS)
	if err != nil {
		return err
	}

	systemConfig.LocalConfigFileName = "gateway.local.json"
	systemConfig.Logger = zap.NewDevelopmentConfig()
	systemConfig.Policy = pol
	systemConfig.DriverFactories = map[string]driver.Factory{
		proxy.DriverName: proxy.Factory,
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

//go:embed policy
var policyFS embed.FS
