package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/jackc/pgx/v4"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/bsp-ew/internal/auth"
	"github.com/vanti-dev/bsp-ew/internal/auth/keycloak"
	"github.com/vanti-dev/bsp-ew/internal/auth/policy"
	"github.com/vanti-dev/bsp-ew/internal/db"
	"github.com/vanti-dev/bsp-ew/internal/testapi"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	flagListenGRPC       string
	flagListenHTTP       string
	flagStaticDir        string
	flagDisableAuth      bool
	flagPostgresAddress  string
	flagPostgresUsername string
	flagPostgresPassword string
	flagPostgresDatabase string
	flagPopulateDatabase bool
	flagKeycloakAddress  string
	flagKeycloakRealm    string
)

func init() {
	flag.StringVar(&flagListenGRPC, "listen-grpc", "localhost:23557", "port to host gRPC server on")
	flag.StringVar(&flagListenHTTP, "listen-http", "localhost:8080", "port to host HTTP server on")
	flag.StringVar(&flagStaticDir, "static-dir", "ui/dist", "path to static files directory")
	flag.BoolVar(&flagDisableAuth, "disable-auth", false, "[INSECURE!] disable API call authorization checks")
	flag.StringVar(&flagPostgresAddress, "postgres-addr", "localhost:5432", "host:port to connect to postgres on")
	flag.StringVar(&flagPostgresUsername, "postgres-user", "postgres", "username for authenticating with postgres")
	flag.StringVar(&flagPostgresPassword, "postgres-password", "postgres", "password for authenticating with postgres")
	flag.StringVar(&flagPostgresDatabase, "postgres-db", "smart_core", "database name for connecting to postgres")
	flag.BoolVar(&flagPopulateDatabase, "populate-db", false, "inserts some test data into the database and exits")
	flag.StringVar(&flagKeycloakAddress, "keycloak-url", "http://localhost:8888", "root URL of Keycloak server")
	flag.StringVar(&flagKeycloakRealm, "keycloak-realm", "smart-core", "realm ID to use for Keycloak authentication")
}

func run(ctx context.Context) error {
	flag.Parse()

	group, ctx := errgroup.WithContext(ctx)

	dbConn, err := connectDB(ctx)
	if err != nil {
		return err
	}
	if flagPopulateDatabase {
		return populateDB(ctx, dbConn)
	}

	pubServer := &PublicationServer{conn: dbConn}

	grpcListener, err := net.Listen("tcp", flagListenGRPC)
	if err != nil {
		return fmt.Errorf("can't listen on %q: %w", flagListenGRPC, err)
	}

	httpListener, err := net.Listen("tcp", flagListenHTTP)
	if err != nil {
		return fmt.Errorf("can't listen on %q: %w", flagListenHTTP, err)
	}

	// ===== Serve gGRPC ======
	var grpcServerOptions []grpc.ServerOption
	if !flagDisableAuth {
		verifier, err := initKeycloakVerifier(ctx)
		if err != nil {
			return fmt.Errorf("init keycloak token verifier: %w", err)
		}
		grpcServerOptions = append(grpcServerOptions,
			grpc.UnaryInterceptor(policy.GRPCUnaryInterceptor(verifier)),
			grpc.StreamInterceptor(policy.GRPCStreamingInterceptor(verifier)),
		)
	}

	grpcServer := grpc.NewServer(grpcServerOptions...)
	traits.RegisterPublicationApiServer(grpcServer, pubServer)
	gen.RegisterTestApiServer(grpcServer, testapi.NewAPI())
	reflection.Register(grpcServer)
	grpcWebWrapper := grpcweb.WrapServer(grpcServer)

	// serve gRPC on the grpcListener
	group.Go(func() error {
		return grpcServer.Serve(grpcListener)
	})
	log.Printf("insecure gRPC server listening on %s", grpcListener.Addr().String())

	// ===== Serve HTTP =====
	staticFiles := http.FileServer(http.Dir(flagStaticDir))
	httpHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if grpcWebWrapper.IsGrpcWebRequest(request) || grpcWebWrapper.IsAcceptableGrpcCorsRequest(request) {
			grpcWebWrapper.ServeHTTP(writer, request)
		} else {
			staticFiles.ServeHTTP(writer, request)
		}
	})
	group.Go(func() error {
		return http.Serve(httpListener, httpHandler)
	})

	// ===== Handle Shutdowns =====
	// immediately attempt graceful shutdown when context cancelled
	stopped := make(chan struct{})
	group.Go(func() error {
		<-ctx.Done()
		grpcServer.GracefulStop()
		close(stopped)
		return nil
	})

	// force shutdown 5s after context cancelled
	group.Go(func() error {
		<-ctx.Done()
		log.Println("gRPC server will be force-closed in 5 seconds")
		select {
		case <-time.After(5 * time.Second):
			grpcServer.Stop()
		case <-stopped:
		}
		log.Println("gRPC server shut down")
		err := httpListener.Close()
		log.Println("HTTP server shut down")
		return err
	})

	return group.Wait()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	errs := multierr.Errors(run(ctx))

	var code int
	switch len(errs) {
	case 0:
	case 1:
		_, _ = fmt.Fprintf(os.Stderr, "fatal error: %s\n", errs[0].Error())
		code = 1
	default:
		_, _ = fmt.Fprintln(os.Stderr, "fatal errors:")
		for _, err := range errs {
			_, _ = fmt.Fprintf(os.Stderr, "\t%s\n", err.Error())
		}
		code = 1
	}

	os.Exit(code)
}

func connectDB(ctx context.Context) (*pgx.Conn, error) {
	// The only valid way to construct a pgconn.Config is using a URL string, so we need to manually construct it even
	// though we have all the parts separately.
	// Constructing the config manually will cause a panic when we attempt to connect.
	// See: documentation for pgconn.Config
	connectURL := url.URL{
		Scheme: "postgres",
		Host:   flagPostgresAddress,
		User:   url.UserPassword(flagPostgresUsername, flagPostgresPassword),
		Path:   "/" + flagPostgresDatabase,
	}
	conn, err := pgx.Connect(ctx, connectURL.String())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func populateDB(ctx context.Context, conn *pgx.Conn) error {
	deviceNames := []string{
		"test/area-controller-1",
		"test/area-controller-2",
		"test/area-controller-3",
	}

	baseTime := time.Date(2022, 7, 6, 11, 18, 0, 0, time.UTC)

	err := conn.BeginFunc(ctx, func(tx pgx.Tx) error {
		var errs error
		for _, name := range deviceNames {
			// register a publication
			id := name + ":config"
			errs = multierr.Append(errs, db.RegisterPublication(ctx, tx, id, name))

			// add some versions to it
			for i := 1; i <= 3; i++ {
				payload := struct {
					Device      string `json:"device"`
					Publication string `json:"publication"`
					Sequence    int    `json:"sequence"`
				}{
					Device:      name,
					Publication: id,
					Sequence:    i,
				}

				encoded, err := json.Marshal(payload)
				if err != nil {
					errs = multierr.Append(errs, err)
					continue
				}

				_, err = db.AddPublicationVersion(ctx, tx, db.PublicationVersion{
					PublicationID: id,
					PublishTime:   baseTime.Add(time.Duration(i) * time.Hour),
					Body:          encoded,
					MediaType:     "application/json",
					Changelog:     fmt.Sprintf("auto-populated revision %d", i),
				})
				errs = multierr.Append(errs, err)
			}
		}

		return errs
	})

	if err != nil {
		log.Printf("failed to populate database: %s", err.Error())
	} else {
		log.Println("database populated")
	}
	return err
}

func initKeycloakVerifier(ctx context.Context) (auth.TokenVerifier, error) {
	authConfig := keycloak.Config{
		URL:      flagKeycloakAddress,
		Realm:    flagKeycloakRealm,
		ClientID: "sc-api",
	}
	authUrls, err := auth.DiscoverOIDCConfig(ctx, authConfig.Issuer())
	if err != nil {
		panic(err)
	}
	keySet := auth.NewRemoteKeySet(ctx, authUrls.JWKSURI)
	return keycloak.NewTokenVerifier(&authConfig, keySet), nil
}
