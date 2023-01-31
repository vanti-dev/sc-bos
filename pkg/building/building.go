package building

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait/onoff"
	"github.com/vanti-dev/sc-bos/pkg/auth/jwks"
	"github.com/vanti-dev/sc-bos/pkg/auth/oidc"
	"github.com/vanti-dev/sc-bos/pkg/auth/token"
	"github.com/vanti-dev/sc-bos/pkg/system/publications/pgxpublications"
	"github.com/vanti-dev/sc-bos/pkg/system/tenants/pgxtenants"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/vanti-dev/sc-bos/internal/auth/keycloak"
	"github.com/vanti-dev/sc-bos/internal/util/pgxutil"
	"github.com/vanti-dev/sc-bos/internal/util/pki"
	"github.com/vanti-dev/sc-bos/internal/util/pki/expire"
	"github.com/vanti-dev/sc-bos/pkg/app"
	"github.com/vanti-dev/sc-bos/pkg/auth/policy"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/testapi"
)

func RunController(ctx context.Context, logger *zap.Logger, configDir string) error {
	// load system config file
	sysConf, err := readSystemConfig(configDir)
	if err != nil {
		return err
	}

	// connect (& optionally initialise) DB
	dbConn, err := connectDB(ctx, sysConf)
	if err != nil {
		return err
	}

	// load certificates

	// serverAuthority is used as a source for trust and certificates between smart core nodes.
	// The certs this source produces should be CA certs (have the CA flag set) to allow them to sign child certificates.
	// These child certs will actually be used for our server and for issuance to other nodes as part of enrollment
	serverAuthority := loadServerAuthority(configDir)

	// Setup tls.Config for our server apis
	grpcTlsServerConfig := pki.TLSServerConfig(serverCerts(sysConf, serverAuthority))

	// Setup the tls.Config for serving https apis - including grpc-web and hosting
	httpsCertSource, err := loadHTTPSCertSource(configDir, sysConf, serverAuthority)
	if err != nil {
		return err
	}
	httpsTlsConfig := pki.TLSServerConfig(httpsCertSource)

	grpcServerOptions := []grpc.ServerOption{
		grpc.Creds(credentials.NewTLS(grpcTlsServerConfig)),
	}
	if !sysConf.DisableAuth {
		verifier, err := initKeycloakValidator(ctx, sysConf)
		if err != nil {
			return fmt.Errorf("init keycloak token verifier: %w", err)
		}
		interceptor := policy.NewInterceptor(policy.Default(false),
			policy.WithTokenVerifier(verifier),
			policy.WithLogger(logger.Named("policy")),
		)
		grpcServerOptions = append(grpcServerOptions,
			grpc.ChainUnaryInterceptor(interceptor.GRPCUnaryInterceptor()),
			grpc.ChainStreamInterceptor(interceptor.GRPCStreamingInterceptor()),
		)
	}

	grpcServer := grpc.NewServer(grpcServerOptions...)
	reflection.Register(grpcServer)
	// todo: replace this with the systems package
	publicationApi, err := pgxpublications.NewServerFromPool(ctx, dbConn, pgxpublications.WithLogger(logger.Named("publications")))
	if err != nil {
		return fmt.Errorf("publication api %w", err)
	}
	traits.RegisterPublicationApiServer(grpcServer, publicationApi)
	gen.RegisterTestApiServer(grpcServer, testapi.NewAPI())
	gen.RegisterNodeApiServer(grpcServer, &NodeServer{
		logger:        logger.Named("NodeServer"),
		db:            dbConn,
		authority:     serverAuthority,
		managerName:   "building-controller",
		managerAddr:   sysConf.CanonicalAddress,
		testTLSConfig: grpcTlsServerConfig,
	})
	// todo: replace this with the systems package
	tenantsApi, err := pgxtenants.NewServerFromPool(ctx, dbConn,
		pgxtenants.WithLogger(logger.Named("tenantapi")))
	if err != nil {
		return fmt.Errorf("tenant api %w", err)
	}
	gen.RegisterTenantApiServer(grpcServer, tenantsApi)
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

func readSystemConfig(configDir string) (SystemConfig, error) {
	sysConfJSON, err := os.ReadFile(filepath.Join(configDir, "system.json"))
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
	return pgxutil.Connect(ctx, pgxutil.ConnectConfig{URI: sysConf.DatabaseURL, PasswordFile: sysConf.DatabasePasswordFile})
}

func initKeycloakValidator(ctx context.Context, sysConf SystemConfig) (token.Validator, error) {
	authConfig := keycloak.Config{
		URL:      sysConf.KeycloakAddress,
		Realm:    sysConf.KeycloakRealm,
		ClientID: "sc-api",
	}
	authUrls, err := oidc.FetchConfig(ctx, authConfig.Issuer())
	if err != nil {
		panic(err)
	}
	keySet := jwks.NewRemoteKeySet(ctx, authUrls.JWKSURI)
	return keycloak.NewTokenValidator(&authConfig, keySet), nil
}

func loadServerAuthority(configDir string) pki.Source {
	return pki.FSSource(
		filepath.Join(configDir, "pki", "enrollment-ca.cert.pem"),
		filepath.Join(configDir, "pki", "enrollment-ca.key.pem"),
		filepath.Join(configDir, "pki", "roots.cert.pem"),
	)
}

// serverCerts creates a new pki.Source that mints new certificates (with server auth usage) using the given authority.
func serverCerts(sysConf SystemConfig, authority pki.Source) pki.Source {
	validity := time.Duration(sysConf.EnrollmentValidityDays) * 24 * time.Hour
	keyPair := func() (*x509.Certificate, pki.PrivateKey, error) {
		key, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return nil, nil, err
		}
		cert := &x509.Certificate{
			Subject:     pkix.Name{CommonName: "building-controller"},
			KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		}
		return cert, key, nil
	}

	return pki.CacheSource(
		pki.AuthoritySourceFn(authority, keyPair, pki.WithIfaces(), pki.WithExpireAfter(validity)),
		expire.AfterProgress(0.5),
	)
}

func loadHTTPSCertSource(configDir string, sysConf SystemConfig, authority pki.Source) (pki.Source, error) {
	if sysConf.SelfSignedHTTPS {
		key, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return nil, err
		}
		template := &x509.Certificate{
			Subject:     pkix.Name{CommonName: "localhost"},
			KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		}
		source := pki.AuthoritySource(authority, template, key, pki.WithIfaces())
		source = pki.CacheSource(source, expire.BeforeInvalid(30*time.Minute))
		return source, nil
	} else {
		certPath := filepath.Join(configDir, "pki", "https.cert.pem")
		keyPath := filepath.Join(configDir, "pki", "https.key.pem")
		source := pki.FSSource(certPath, keyPath, "")
		source = pki.CacheSource(source, expire.BeforeInvalid(30*time.Minute))
		return source, nil
	}
}
