package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-jose/go-jose/v3"
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
	flag.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	tenantLogger := logger.Named("tenant.oauth")

	c := &app.Controller{
		Logger:      logger,
		DataDir:     flagDataDir,
		ListenGRPC:  flagListenGRPC,
		ListenHTTPS: flagListenHTTPS,
		Routes: map[string]http.Handler{
			"/oauth2/token": tenant.OAuth2TokenHandler(
				genTenantSecrets(tenantLogger),
				genTenantTokenSource(tenantLogger),
			),
		},
	}

	os.Exit(app.RunUntilInterrupt(c.Run))
}

func genTenantSecrets(logger *zap.Logger) tenant.SecretStore {
	store := tenant.NewMemorySecretStore(nil)
	for i := 1; i <= 3; i++ {
		clientId := fmt.Sprintf("tenant-%d", i)
		data := tenant.SecretData{ClientID: clientId}
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

func genTenantTokenSource(logger *zap.Logger) *tenant.TokenSource {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	signingKey := jose.SigningKey{
		Algorithm: jose.RS256,
		Key:       key,
	}

	jwk := jose.JSONWebKey{
		Key:       key.Public(),
		Algorithm: string(signingKey.Algorithm),
		Use:       "sig",
	}
	jwkBytes, err := json.Marshal(jwk)
	if err != nil {
		panic(err)
	}
	logger.Debug("generated signing key", zap.Any("key", json.RawMessage(jwkBytes)))

	return &tenant.TokenSource{
		Key:      signingKey,
		Issuer:   "http://localhost:8080",
		Validity: 5 * time.Minute,
	}
}
