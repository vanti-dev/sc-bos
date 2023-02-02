package main

import (
	"context"
	"embed"
	"flag"
	"os"

	"github.com/vanti-dev/sc-bos/pkg/app"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/driver/proxy"
	"github.com/vanti-dev/sc-bos/pkg/node/alltraits"
	"github.com/vanti-dev/sc-bos/pkg/system"
	"github.com/vanti-dev/sc-bos/pkg/system/authn"
	"github.com/vanti-dev/sc-bos/pkg/testapi"

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
	systemConfig.SystemFactories = map[string]system.Factory{
		"authn": authn.Factory(),
	}

	controller, err := app.Bootstrap(ctx, systemConfig)
	if err != nil {
		return err
	}

	alltraits.AddSupport(controller.Node)

	gen.RegisterTestApiServer(controller.GRPC, testapi.NewAPI())

	return controller.Run(ctx)
}

//go:embed policy
var policyFS embed.FS
