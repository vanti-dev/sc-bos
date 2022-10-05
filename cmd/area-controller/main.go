package main

import (
	"context"
	"flag"
	"os"

	"github.com/vanti-dev/bsp-ew/internal/app"
	"github.com/vanti-dev/bsp-ew/internal/testapi"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.uber.org/zap"
)

var (
	flagListenGRPC  string
	flagListenHTTPS string
	flagDataDir     string
	flagStaticDir   string

	flagDisablePolicy bool
)

func init() {
	flag.StringVar(&flagListenGRPC, "listen-grpc", ":23557", "address (host:port) to host a Smart Core gRPC server on")
	flag.StringVar(&flagListenHTTPS, "listen-https", ":443", "address (host:port) to host a HTTPS server on")
	flag.StringVar(&flagDataDir, "data-dir", ".data/area-controller-01", "path to local data storage directory")
	flag.StringVar(&flagStaticDir, "static-dir", "ui/dist", "path for HTTP static resources")

	flag.BoolVar(&flagDisablePolicy, "insecure-disable-policy", false, "Insecure! Disable checking requests against the security policy. This option opens up the server to any request.")
}

func main() {
	os.Exit(app.RunUntilInterrupt(run))
}

func run(ctx context.Context) error {
	flag.Parse()
	config := app.SystemConfig{
		Logger:        zap.NewDevelopmentConfig(),
		DataDir:       flagDataDir,
		ListenGRPC:    flagListenGRPC,
		ListenHTTPS:   flagListenHTTPS,
		DisablePolicy: flagDisablePolicy,
	}

	controller, err := app.Bootstrap(config)
	if err != nil {
		return err
	}

	gen.RegisterTestApiServer(controller.GRPC, testapi.NewAPI())

	return controller.Run(ctx)
}
