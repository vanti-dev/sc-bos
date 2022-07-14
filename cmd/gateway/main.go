package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/vanti-dev/bsp-ew/pkg/auth/tenant"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt)
	defer done()
	err := run(ctx)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// Open HTTP Port
	httpListener, httpPort := listen(8000)
	go func() {
		<-ctx.Done()
		err := httpListener.Close()
		if err != nil {
			log.Printf("httpListener close error: %s", err.Error())
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/oauth2/token", tenant.OAuth2TokenHandler(genTenantSecrets(), genTenantTokenSource()))

	fmt.Printf("HTTP server listening on :%d\n", httpPort)
	return http.Serve(httpListener, mux)
}

func listen(desiredPort int) (listener net.Listener, port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", desiredPort))
	if err != nil {
		panic(err)
	}

	port = listener.Addr().(*net.TCPAddr).Port

	return
}

func genTenantSecrets() tenant.SecretStore {
	store := tenant.NewMemorySecretStore(nil)
	for i := 1; i <= 3; i++ {
		clientId := fmt.Sprintf("tenant-%d", i)
		data := tenant.SecretData{ClientID: clientId}
		secret, err := store.Enroll(context.TODO(), data)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Created new tenant %s with secret %s\n", clientId, secret)
	}
	return store
}

func genTenantTokenSource() *tenant.TokenSource {
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
	fmt.Printf("generated signing key %s\n", string(jwkBytes))

	return &tenant.TokenSource{
		Key:      signingKey,
		Issuer:   "http://localhost:8080",
		Validity: 5 * time.Minute,
	}
}
