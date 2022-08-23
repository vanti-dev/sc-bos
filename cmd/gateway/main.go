package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/vanti-dev/bsp-ew/internal/app"
	"github.com/vanti-dev/bsp-ew/internal/auth/tenant"
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

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	tokenServer, err := tenant.NewTokenSever(genTenantSecrets(logger), "gateway", 5*time.Minute, logger.Named("tenant.oauth"))
	if err != nil {
		return err
	}

	c := &app.Controller{
		Logger:      logger,
		DataDir:     flagDataDir,
		ListenGRPC:  flagListenGRPC,
		ListenHTTPS: flagListenHTTPS,
		Routes: map[string]http.Handler{
			"/oauth2/token": tokenServer,
		},
	}

	return c.Run(ctx)
}

func genTenantSecrets(logger *zap.Logger) tenant.SecretSource {
	store := tenant.NewMemorySecretStore(nil)
	for i := 1; i <= 3; i++ {
		clientId := fmt.Sprintf("tenant-%d", i)
		data := tenant.SecretData{TenantID: clientId}
		secret, err := store.Enroll(context.TODO(), data)
		if err != nil {
			panic(err)
		}

		logger.Info("created new tenant",
			zap.String("clientId", clientId),
			zap.String("secret", secret),
		)
	}
	return store
}
