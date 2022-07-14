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
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/vanti-dev/bsp-ew/pkg/api"
	"github.com/vanti-dev/bsp-ew/pkg/auth"
	"github.com/vanti-dev/bsp-ew/pkg/auth/keycloak"
	"github.com/vanti-dev/bsp-ew/pkg/auth/tenant"
	"github.com/vanti-dev/bsp-ew/pkg/policy"
	"github.com/vanti-dev/bsp-ew/pkg/testgen"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	// Open gRPC Port
	grpcListener, grpcPort := listen(9000)
	go func() {
		<-ctx.Done()
		err := grpcListener.Close()
		if err != nil {
			log.Printf("grpcListener close error: %s", err.Error())
		}
	}()

	// Configure auth
	verifier, err := initKeycloakVerifier(ctx)
	if err != nil {
		return err
	}

	// Create API handler
	apiHandler := api.NewAPI()

	// Register handlers
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("static")))
	mux.Handle("/oauth2/token", tenant.OAuth2TokenHandler(genTenantSecrets(), genTenantTokenSource()))

	// Setup gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(policy.GRPCUnaryInterceptor(verifier)),
		grpc.StreamInterceptor(policy.GRPCStreamingInterceptor(verifier)),
	)
	testgen.RegisterTestApiServer(grpcServer, apiHandler)
	reflection.Register(grpcServer)
	grpcWebWrapper := grpcweb.WrapServer(grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool { return true }),
	)

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		fmt.Printf("Listening on http://localhost:%d/\n", httpPort)
		handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if grpcWebWrapper.IsGrpcWebRequest(request) || grpcWebWrapper.IsAcceptableGrpcCorsRequest(request) {
				grpcWebWrapper.ServeHTTP(writer, request)
			} else {
				mux.ServeHTTP(writer, request)
			}
		})
		return http.Serve(httpListener, handler)
	})

	group.Go(func() error {
		fmt.Printf("gRPC listening on localhost:%d\n", grpcPort)
		return grpcServer.Serve(grpcListener)
	})

	return group.Wait()
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

func initKeycloakVerifier(ctx context.Context) (auth.TokenVerifier, error) {
	authConfig := keycloak.Config{
		URL:      "http://localhost:8888",
		Realm:    "smart-core",
		ClientID: "sc-api",
	}
	authUrls, err := auth.DiscoverOIDCConfig(ctx, authConfig.Issuer())
	if err != nil {
		panic(err)
	}
	keySet := auth.NewRemoteKeySet(ctx, authUrls.JWKSURI)
	return keycloak.NewTokenVerifier(&authConfig, keySet), nil
}
