package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"os"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait/onoff"
	"github.com/vanti-dev/bsp-ew/internal/app"
	"github.com/vanti-dev/bsp-ew/internal/auth/policy"
	"github.com/vanti-dev/bsp-ew/internal/testapi"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.uber.org/zap"
)

var (
	flagListenGRPC  string
	flagListenHTTPS string
	flagDataDir     string
	flagStaticDir   string
)

func init() {
	flag.StringVar(&flagListenGRPC, "listen-grpc", ":23557", "address (host:port) to host a Smart Core gRPC server on")
	flag.StringVar(&flagListenHTTPS, "listen-https", ":443", "address (host:port) to host a HTTPS server on")
	flag.StringVar(&flagDataDir, "data-dir", ".data/gateway", "path to local data storage directory")
	flag.StringVar(&flagStaticDir, "static-dir", "ui/dist", "path for HTTP static resources")
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

	config := app.SystemConfig{
		Logger:      zap.NewDevelopmentConfig(),
		DataDir:     flagDataDir,
		ListenGRPC:  flagListenGRPC,
		ListenHTTPS: flagListenHTTPS,
		TenantOAuth: true,
		Policy:      pol,
	}

	controller, err := app.Bootstrap(ctx, config)
	if err != nil {
		return err
	}
	gen.RegisterTestApiServer(controller.GRPC, testapi.NewAPI())
	traits.RegisterOnOffApiServer(controller.GRPC, onoff.NewApiRouter(
		onoff.WithOnOffApiClientFactory(func(name string) (traits.OnOffApiClient, error) {
			conn, err := controller.ManagerConn.Connect(ctx)
			if err != nil {
				return nil, err
			}
			if conn == nil {
				return nil, errors.New("no manager conn")
			}
			return traits.NewOnOffApiClient(conn), nil
		})))

	return controller.Run(ctx)
}

//go:embed policy
var policyFS embed.FS
