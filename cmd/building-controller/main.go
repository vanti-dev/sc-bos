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
	"github.com/vanti-dev/bsp-ew/internal/app"
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
	flagConfigDir  string
	flagPopulateDB bool
)

func init() {
	flag.StringVar(&flagConfigDir, "config-dir", ".data/app-server", "path to the configuration directory")
	flag.BoolVar(&flagPopulateDB, "populate-db", false, "inserts some test data into the database and exits")
}

func run(ctx context.Context) error {
	flag.Parse()
	// load system config file
	sysConfJSON, err := os.ReadFile(filepath.Join(flagConfigDir, "system.json"))
	if err != nil {
		return err
	}
	sysConf := DefaultSystemConfig()
	err = json.Unmarshal(sysConfJSON, &sysConf)
	if err != nil {
		return err
	}

	certSource, err := pki.NewSelfSignedCertSource(nil)
	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{
		GetCertificate: certSource.TLSConfigGetCertificate,
	}

	dbConn, err := connectDB(ctx, sysConf)
	if err != nil {
		return err
	}
	if flagPopulateDB {
		return populateDB(ctx, dbConn)
	}

	grpcServerOptions := []grpc.ServerOption{
		grpc.Creds(credentials.NewTLS(tlsConfig)),
	}
	if !sysConf.DisableAuth {
		verifier, err := initKeycloakVerifier(ctx, sysConf)
		if err != nil {
			return fmt.Errorf("init keycloak token verifier: %w", err)
		}
		grpcServerOptions = append(grpcServerOptions,
			grpc.UnaryInterceptor(policy.GRPCUnaryInterceptor(verifier)),
			grpc.StreamInterceptor(policy.GRPCStreamingInterceptor(verifier)),
		)
	}

	servers := &app.Servers{
		ShutdownTimeout: 15 * time.Second,
		GRPC:            grpc.NewServer(grpcServerOptions...),
		GRPCAddress:     sysConf.ListenGRPC,
		HTTP: &http.Server{
			Addr:      sysConf.ListenHTTPS,
			TLSConfig: tlsConfig,
		},
	}

	grpcServer := grpc.NewServer(grpcServerOptions...)
	traits.RegisterPublicationApiServer(grpcServer, &PublicationServer{conn: dbConn})
	gen.RegisterTestApiServer(grpcServer, testapi.NewAPI())
	reflection.Register(grpcServer)

	grpcWebWrapper := grpcweb.WrapServer(grpcServer)
	staticFiles := http.FileServer(http.Dir(sysConf.StaticDir))
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

func connectDB(ctx context.Context, sysConf SystemConfig) (*pgx.Conn, error) {
	connConfig, err := pgx.ParseConfig(sysConf.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if sysConf.DatabasePasswordFile != "" {
		passwordFile, err := ioutil.ReadFile(sysConf.DatabasePasswordFile)
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

func initKeycloakVerifier(ctx context.Context, sysConf SystemConfig) (auth.TokenVerifier, error) {
	authConfig := keycloak.Config{
		URL:      sysConf.KeycloakAddress,
		Realm:    sysConf.KeycloakRealm,
		ClientID: "sc-api",
	}
	authUrls, err := auth.DiscoverOIDCConfig(ctx, authConfig.Issuer())
	if err != nil {
		panic(err)
	}
	keySet := auth.NewRemoteKeySet(ctx, authUrls.JWKSURI)
	return keycloak.NewTokenVerifier(&authConfig, keySet), nil
}
