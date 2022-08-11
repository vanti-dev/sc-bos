package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/jackc/pgx/v4"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/bsp-ew/internal/auth"
	"github.com/vanti-dev/bsp-ew/internal/auth/keycloak"
	"github.com/vanti-dev/bsp-ew/internal/auth/policy"
	"github.com/vanti-dev/bsp-ew/internal/db"
	"github.com/vanti-dev/bsp-ew/internal/testapi"
	"github.com/vanti-dev/bsp-ew/internal/util/pki"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

var (
	flagListenGRPC      string
	flagListenHTTPS     string
	flagStaticDir       string
	flagConfigDir       string
	flagDisableAuth     bool
	flagDBURL           string
	flagDBPasswordFile  string
	flagPopulateDB      bool
	flagKeycloakAddress string
	flagKeycloakRealm   string
)

func init() {
	flag.StringVar(&flagListenGRPC, "listen-grpc", ":23557", "address (host:port) to host gRPC server on")
	flag.StringVar(&flagListenHTTPS, "listen-https", ":80", "address (host:port) to host HTTPS server on")
	flag.StringVar(&flagStaticDir, "static-dir", "ui/dist", "path to static files directory")
	flag.StringVar(&flagConfigDir, "config-dir", ".data/app-server", "path to the configuration directory")
	flag.BoolVar(&flagDisableAuth, "disable-auth", false, "[INSECURE!] disable API call authorization checks")
	flag.StringVar(&flagDBURL, "db-url", "postgres://postgres:postgres@localhost:5432/smart_core", "PostgreSQL connection URL in libpq style")
	flag.StringVar(&flagDBPasswordFile, "db-password-file", "", "path to a file containing the PostgreSQL password")
	flag.BoolVar(&flagPopulateDB, "populate-db", false, "inserts some test data into the database and exits")
	flag.StringVar(&flagKeycloakAddress, "keycloak-url", "http://localhost:8888", "root URL of Keycloak server")
	flag.StringVar(&flagKeycloakRealm, "keycloak-realm", "smart-core", "realm ID to use for Keycloak authentication")
}

func run(ctx context.Context) error {
	flag.Parse()

	privateKey, _, err := pki.LoadOrGeneratePrivateKey(filepath.Join(flagConfigDir, "private-key.pem"))
	if err != nil {
		return err
	}

	certSource := &pki.SelfSignedCertSource{PrivateKey: privateKey}

	tlsConfig := &tls.Config{
		GetCertificate: certSource.TLSConfigGetCertificate,
	}

	dbConn, err := connectDB(ctx)
	if err != nil {
		return err
	}
	if flagPopulateDB {
		return populateDB(ctx, dbConn)
	}

	grpcServerOptions := []grpc.ServerOption{
		grpc.Creds(credentials.NewTLS(tlsConfig)),
	}
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

	servers := &Servers{
		ShutdownTimeout: 15 * time.Second,
		GRPC:            grpc.NewServer(grpcServerOptions...),
		GRPCAddress:     flagListenGRPC,
		HTTP: &http.Server{
			Addr:      flagListenHTTPS,
			TLSConfig: tlsConfig,
		},
	}

	grpcServer := grpc.NewServer(grpcServerOptions...)
	traits.RegisterPublicationApiServer(grpcServer, &PublicationServer{conn: dbConn})
	gen.RegisterTestApiServer(grpcServer, testapi.NewAPI())
	reflection.Register(grpcServer)

	grpcWebWrapper := grpcweb.WrapServer(grpcServer)
	staticFiles := http.FileServer(http.Dir(flagStaticDir))
	servers.HTTP.Handler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if grpcWebWrapper.IsGrpcWebRequest(request) || grpcWebWrapper.IsAcceptableGrpcCorsRequest(request) {
			grpcWebWrapper.ServeHTTP(writer, request)
		} else {
			staticFiles.ServeHTTP(writer, request)
		}
	})

	return servers.Serve(ctx)
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
	connConfig, err := pgx.ParseConfig(flagDBURL)
	if err != nil {
		return nil, err
	}

	if flagDBPasswordFile != "" {
		passwordFile, err := ioutil.ReadFile(flagDBPasswordFile)
		if err != nil {
			return nil, err
		}

		connConfig.Password = strings.TrimSpace(string(passwordFile))
	}

	return pgx.ConnectConfig(ctx, connConfig)
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
