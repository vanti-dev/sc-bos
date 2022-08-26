package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait/onoff"
	"github.com/vanti-dev/bsp-ew/internal/app"
	"github.com/vanti-dev/bsp-ew/internal/auth"
	"github.com/vanti-dev/bsp-ew/internal/auth/keycloak"
	"github.com/vanti-dev/bsp-ew/internal/auth/policy"
	"github.com/vanti-dev/bsp-ew/internal/db"
	"github.com/vanti-dev/bsp-ew/internal/manage/enrollment"
	"github.com/vanti-dev/bsp-ew/internal/manage/tenantapi"
	"github.com/vanti-dev/bsp-ew/internal/testapi"
	"github.com/vanti-dev/bsp-ew/internal/util/pki"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

var (
	flagConfigDir  string
	flagPopulateDB bool
)

func init() {
	flag.StringVar(&flagConfigDir, "config-dir", ".data/building-controller", "path to the configuration directory")
	flag.BoolVar(&flagPopulateDB, "populate-db", false, "inserts some test data into the database and exits")
}

func run(ctx context.Context) error {
	flag.Parse()
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	// load system config file
	sysConf, err := readSystemConfig()
	if err != nil {
		return err
	}

	// connect (& optionally initialise) DB
	dbConn, err := connectDB(ctx, sysConf)
	if err != nil {
		return err
	}
	if flagPopulateDB {
		return populateDB(ctx, logger, dbConn)
	}

	// load certificates
	rootsPEM, err := os.ReadFile(filepath.Join(flagConfigDir, "pki", "roots.cert.pem"))
	if err != nil {
		return err
	}
	rootsPool := x509.NewCertPool()
	if !rootsPool.AppendCertsFromPEM(rootsPEM) {
		return errors.New("unable to parse any Root CA certificates")
	}
	ca, err := loadEnrollmentCA(sysConf)
	if err != nil {
		return err
	}
	grpcCertSource, err := ca.LocalCertSource(pkix.Name{CommonName: "building-controller"}, true)
	if err != nil {
		return err
	}
	grpcTlsConfig := &tls.Config{
		GetCertificate:       grpcCertSource.TLSConfigGetCertificate,
		GetClientCertificate: grpcCertSource.TLSConfigGetClientCertificate,
		RootCAs:              rootsPool,
		ClientCAs:            rootsPool,
		ClientAuth:           tls.VerifyClientCertIfGiven,
	}
	httpsCertSource, err := loadHTTPSCertSource(sysConf, logger)
	if err != nil {
		return err
	}
	httpsTlsConfig := &tls.Config{GetCertificate: httpsCertSource.TLSConfigGetCertificate}

	grpcServerOptions := []grpc.ServerOption{
		grpc.Creds(credentials.NewTLS(grpcTlsConfig)),
	}
	if !sysConf.DisableAuth {
		verifier, err := initKeycloakValidator(ctx, sysConf)
		if err != nil {
			return fmt.Errorf("init keycloak token verifier: %w", err)
		}
		interceptor := policy.NewInterceptor(policy.Default(),
			policy.WithTokenVerifier(verifier),
			policy.WithLogger(logger.Named("policy")),
		)
		grpcServerOptions = append(grpcServerOptions,
			grpc.UnaryInterceptor(interceptor.GRPCUnaryInterceptor()),
			grpc.StreamInterceptor(interceptor.GRPCStreamingInterceptor()),
		)
	}

	grpcServer := grpc.NewServer(grpcServerOptions...)
	reflection.Register(grpcServer)
	traits.RegisterPublicationApiServer(grpcServer, &PublicationServer{conn: dbConn})
	gen.RegisterTestApiServer(grpcServer, testapi.NewAPI())
	gen.RegisterNodeApiServer(grpcServer, &NodeServer{
		logger:        logger.Named("NodeServer"),
		db:            dbConn,
		ca:            ca,
		managerName:   "building-controller",
		managerAddr:   sysConf.CanonicalAddress,
		rootsPEM:      rootsPEM,
		testTLSConfig: grpcTlsConfig,
	})
	gen.RegisterTenantApiServer(grpcServer, tenantapi.NewServer(dbConn,
		tenantapi.WithLogger(logger.Named("tenantapi"))))
	traits.RegisterOnOffApiServer(grpcServer, onoff.NewApiRouter(
		onoff.WithOnOffApiClientFactory(func(name string) (traits.OnOffApiClient, error) {
			model := onoff.NewModel(traits.OnOff_OFF)
			return onoff.WrapApi(onoff.NewModelServer(model)), nil
		}),
	))

	grpcWebWrapper := grpcweb.WrapServer(grpcServer, grpcweb.WithOriginFunc(func(origin string) bool {
		return true
	}))
	staticFiles := http.FileServer(http.Dir(sysConf.StaticDir))
	httpServer := &http.Server{
		Addr:      sysConf.ListenHTTPS,
		TLSConfig: httpsTlsConfig,
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if grpcWebWrapper.IsGrpcWebRequest(request) || grpcWebWrapper.IsAcceptableGrpcCorsRequest(request) {
				grpcWebWrapper.ServeHTTP(writer, request)
			} else {
				staticFiles.ServeHTTP(writer, request)
			}
		}),
	}

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return app.ServeGRPC(ctx, grpcServer, sysConf.ListenGRPC, 15*time.Second, logger.Named("server.grpc"))
	})
	group.Go(func() error {
		return app.ServeHTTPS(ctx, httpServer, 15*time.Second, logger.Named("server.https"))
	})
	return group.Wait()
}

func main() {
	os.Exit(app.RunUntilInterrupt(run))
}

func readSystemConfig() (SystemConfig, error) {
	sysConfJSON, err := os.ReadFile(filepath.Join(flagConfigDir, "system.json"))
	if err != nil {
		return SystemConfig{}, err
	}
	sysConf := DefaultSystemConfig()
	err = json.Unmarshal(sysConfJSON, &sysConf)
	if err != nil {
		return SystemConfig{}, err
	}
	return sysConf, nil
}

func connectDB(ctx context.Context, sysConf SystemConfig) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(sysConf.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if sysConf.DatabasePasswordFile != "" {
		passwordFile, err := ioutil.ReadFile(sysConf.DatabasePasswordFile)
		if err != nil {
			return nil, err
		}

		poolConfig.ConnConfig.Password = strings.TrimSpace(string(passwordFile))
	}

	return pgxpool.ConnectConfig(ctx, poolConfig)
}

func populateDB(ctx context.Context, logger *zap.Logger, conn *pgxpool.Pool) error {
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
			errs = multierr.Append(errs, db.CreatePublication(ctx, tx, id, name))

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

				_, err = db.CreatePublicationVersion(ctx, tx, db.PublicationVersion{
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
		logger.Error("failed to populate database", zap.Error(err))
	} else {
		logger.Info("database populated")
	}
	return err
}

func initKeycloakValidator(ctx context.Context, sysConf SystemConfig) (auth.TokenValidator, error) {
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

func loadEnrollmentCA(sysConf SystemConfig) (*enrollment.CA, error) {
	certPath := filepath.Join(flagConfigDir, "pki", "enrollment-ca.cert.pem")
	keyPath := filepath.Join(flagConfigDir, "pki", "enrollment-ca.key.pem")

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, err
	}
	return &enrollment.CA{
		Certificate:   leaf,
		PrivateKey:    cert.PrivateKey,
		Intermediates: cert.Certificate[1:],
		Now:           time.Now,
		Validity:      time.Duration(sysConf.EnrollmentValidityDays) * 24 * time.Hour,
	}, nil
}

func loadHTTPSCertSource(sysConf SystemConfig, logger *zap.Logger) (pki.CertSource, error) {
	if sysConf.SelfSignedHTTPS {
		return pki.NewSelfSignedCertSource(nil, logger)
	} else {
		certPath := filepath.Join(flagConfigDir, "pki", "https.cert.pem")
		keyPath := filepath.Join(flagConfigDir, "pki", "https.key.pem")
		return pki.NewFileCertSource(certPath, keyPath)
	}
}
