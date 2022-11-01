package app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/vanti-dev/bsp-ew/internal/auth/policy"
	"github.com/vanti-dev/bsp-ew/internal/auth/tenant"
	"github.com/vanti-dev/bsp-ew/internal/auto"
	"github.com/vanti-dev/bsp-ew/internal/driver"
	"github.com/vanti-dev/bsp-ew/internal/manage/enrollment"
	"github.com/vanti-dev/bsp-ew/internal/node"
	"github.com/vanti-dev/bsp-ew/internal/task"
	"github.com/vanti-dev/bsp-ew/internal/util/pki"
	"github.com/vanti-dev/bsp-ew/internal/util/pki/expire"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const LocalConfigFileName = "area-controller.local.json"

type SystemConfig struct {
	Logger      zap.Config
	DataDir     string
	ListenGRPC  string
	ListenHTTPS string

	// TenantOAuth, if true, means the controller will support issuing tokens using the OAuth client_credentials flow
	// with credentials read from tenants.json or verified using the manager node
	TenantOAuth bool
	// LocalOAuth, if true, means the controller will support issuing tokens using the OAuth password flow with
	// credentials read from users.json
	LocalOAuth bool

	Policy        policy.Policy // Override the policy used for RPC calls. Defaults to policy.Default
	DisablePolicy bool          // Unsafe, disables any policy checking for the server

	DriverFactories map[string]driver.Factory // keyed by driver name
	AutoFactories   map[string]auto.Factory   // keyed by automation type
}

// Bootstrap will obtain a Controller in a ready-to-run state.
func Bootstrap(ctx context.Context, config SystemConfig) (*Controller, error) {
	logger, err := config.Logger.Build()
	if err != nil {
		return nil, err
	}

	// create data dir if it doesn't exist
	err = os.MkdirAll(config.DataDir, 0750)
	if err != nil {
		return nil, err
	}

	// load the local config file if possible
	// TODO: pull config from manager publication
	var localConfig ControllerConfig
	localConfigPath := filepath.Join(config.DataDir, LocalConfigFileName)
	rawLocalConfig, err := os.ReadFile(localConfigPath)
	if err == nil {
		err = json.Unmarshal(rawLocalConfig, &localConfig)
		if err != nil {
			return nil, fmt.Errorf("local config JSON unmarshal: %w", err)
		}
	} else {
		if !errors.Is(err, os.ErrNotExist) {
			logger.Warn("failed to load local config from file", zap.Error(err),
				zap.String("path", localConfigPath))
		} else {
			logger.Debug("failed to load local config from file", zap.Error(err), zap.String("path", localConfigPath))
		}
	}

	services := driver.Services{
		Logger: logger,
		Node:   node.New(localConfig.Name),
		Tasks:  &task.Group{},
	}
	services.Node.Logger = logger.Named("node")

	// create private key if it doesn't exist
	key, keyPEM, err := pki.LoadOrGeneratePrivateKey(filepath.Join(config.DataDir, "private-key.pem"), logger)
	if err != nil {
		return nil, err
	}

	enrollServer, err := enrollment.LoadOrCreateServer(filepath.Join(config.DataDir, "enrollment"), keyPEM, logger.Named("enrollment"))
	if err != nil {
		return nil, err
	}

	// We read certificates from a few sources, choosing the first that succeeds.
	// First we attempt to use cohort enrollment as our source of certs/roots.
	// If that fails we attempt to read from files in the data dir (server-cert.pem, private-key.pem, and roots.pem).
	// If all that fails we mint a new self signed certificate.
	certSource := pki.ChainSource(
		enrollServer,
		pki.CacheSource(pki.FuncSource(func() (*tls.Certificate, []*x509.Certificate, error) {
			return readCertAndRoots(config, key)
		}), expire.BeforeInvalid(time.Hour)),
		pki.CacheSource(pki.SelfSignedSource(key, pki.WithExpireAfter(30*24*time.Hour), pki.WithIfaces()), expire.AfterProgress(0.5)),
	)
	tlsGRPCServerConfig := pki.TLSServerConfig(certSource)
	tlsGRPCServerConfig.ClientAuth = tls.VerifyClientCertIfGiven
	tlsGRPCClientConfig := pki.TLSClientConfig(certSource)

	tlsHTTPServerConfig := pki.TLSServerConfig(certSource)

	// manager represents a delayed connection to the cohort manager.
	// Using the manager connection when we aren't enrolled will result in RPC calls returning 'not resolved' errors or similar.
	// As soon as we get enrolled those connections will be updated with the current manager address and will start to work.
	manager := node.DialChan(ctx, enrollServer.ManagerAddress(ctx),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsGRPCClientConfig)))

	mux := http.NewServeMux()

	var grpcOpts []grpc.ServerOption
	grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(tlsGRPCServerConfig)))

	// Configure how we authenticate requests
	policyInterceptorOpts := []policy.InterceptorOption{policy.WithLogger(logger.Named("policy"))}
	if shouldSetupTokenServer(config) {
		// Setup the OAuth server for issuing and validating tokens
		tokenServerOpts := []tenant.TokenServerOption{tenant.WithLogger(logger.Named("token.server"))}
		if config.TenantOAuth {
			verifier, err := clientVerifier(config, manager)
			if err != nil {
				return nil, err
			}
			tokenServerOpts = append(tokenServerOpts, tenant.WithClientCredentialFlow(verifier, 15*time.Minute))
		}
		if config.LocalOAuth {
			verifier, err := passwordVerifier(config)
			if err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					// if the file exists, but we can't read it, we should let someone know
					return nil, err
				}
			} else {
				tokenServerOpts = append(tokenServerOpts, tenant.WithPasswordFlow(verifier, 24*time.Hour))
			}
		}

		tokenServer, err := tenant.NewTokenServer("gateway", tokenServerOpts...)
		if err != nil {
			return nil, err
		}
		mux.Handle("/oauth2/token", tokenServer)
		policyInterceptorOpts = append(policyInterceptorOpts, policy.WithTokenVerifier(tokenServer.TokenValidator()))
	}
	if !config.DisablePolicy {
		pol := policy.Default(false)
		if config.Policy != nil {
			pol = config.Policy
		}

		interceptor := policy.NewInterceptor(pol, policyInterceptorOpts...)
		grpcOpts = append(grpcOpts,
			grpc.ChainUnaryInterceptor(interceptor.GRPCUnaryInterceptor()),
			grpc.ChainStreamInterceptor(interceptor.GRPCStreamingInterceptor()),
		)
	}

	grpcServer := grpc.NewServer(grpcOpts...)
	reflection.Register(grpcServer)
	gen.RegisterEnrollmentApiServer(grpcServer, enrollServer)

	grpcWebServer := grpcweb.WrapServer(grpcServer, grpcweb.WithOriginFunc(func(origin string) bool {
		return true
	}))

	httpServer := &http.Server{
		Addr:      config.ListenHTTPS,
		TLSConfig: tlsHTTPServerConfig,
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if grpcWebServer.IsGrpcWebRequest(request) || grpcWebServer.IsAcceptableGrpcCorsRequest(request) {
				grpcWebServer.ServeHTTP(writer, request)
			} else {
				mux.ServeHTTP(writer, request)
			}
		}),
	}

	c := &Controller{
		Services:         services,
		SystemConfig:     config,
		ControllerConfig: localConfig,
		Enrollment:       enrollServer,
		Mux:              mux,
		GRPC:             grpcServer,
		HTTP:             httpServer,
		ClientTLSConfig:  tlsGRPCClientConfig,
		ManagerConn:      manager,
	}
	c.Defer(manager.Close)
	return c, nil
}

func shouldSetupTokenServer(config SystemConfig) bool {
	return config.TenantOAuth || config.LocalOAuth
}

func readCertAndRoots(config SystemConfig, key pki.PrivateKey) (*tls.Certificate, []*x509.Certificate, error) {
	certPath := filepath.Join(config.DataDir, "server-cert.pem")
	cert, err := pki.LoadX509Cert(certPath, key)
	if err != nil {
		return nil, nil, err
	}
	rootsPem, err := os.ReadFile(filepath.Join(config.DataDir, "roots.pem"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// we ignore that roots doesn't exist, this just means we don't trust other nodes
			return &cert, nil, nil
		}
		return nil, nil, err
	}
	roots, err := pki.ParseCertificatesPEM(rootsPem)
	if err != nil {
		return nil, nil, err
	}
	return &cert, roots, nil
}

type Controller struct {
	driver.Services
	SystemConfig     SystemConfig
	ControllerConfig ControllerConfig
	Enrollment       *enrollment.Server

	Mux  *http.ServeMux
	GRPC *grpc.Server
	HTTP *http.Server

	ClientTLSConfig *tls.Config
	ManagerConn     node.Remote

	deferred []Deferred
}

type Deferred func() error

// Defer indicates that the given Deferred should be executed when the Controllers Run method returns.
func (c *Controller) Defer(d Deferred) {
	c.deferred = append(c.deferred, d)
}

func (c *Controller) Run(ctx context.Context) (err error) {
	defer func() {
		for _, d := range c.deferred {
			err = multierr.Append(err, d())
		}
	}()

	// if any automation exposes apis on the node, allow them to add that support
	for _, factory := range c.SystemConfig.AutoFactories {
		if api, ok := factory.(node.SelfSupporter); ok {
			api.AddSupport(c.Node)
		}
	}

	// we delay registering the node servers until now, so that the caller can call c.Node.Support in between
	// Bootstrap and Run and have all these added correctly.
	c.Node.Register(c.GRPC)

	group, ctx := errgroup.WithContext(ctx)
	if c.SystemConfig.ListenGRPC != "" {
		group.Go(func() error {
			return ServeGRPC(ctx, c.GRPC, c.SystemConfig.ListenGRPC, 15*time.Second, c.Logger.Named("server.grpc"))
		})
	}
	if c.SystemConfig.ListenHTTPS != "" {
		group.Go(func() error {
			return ServeHTTPS(ctx, c.HTTP, 15*time.Second, c.Logger.Named("server.https"))
		})
	}

	// load and start the drivers
	results := driver.Build(ctx, c.Services, c.SystemConfig.DriverFactories, c.ControllerConfig.Drivers)
	loaded, failed := summariseResults(results)
	c.Logger.Named("driver").Info("driver loading complete", zap.Int("loaded", loaded), zap.Int("failed", failed))

	// load and start the automations
	if err := c.startAutomations(ctx); err != nil {
		return err
	}

	err = multierr.Append(err, group.Wait())
	return
}

func summariseResults(results map[string]driver.BuildResult) (loaded int, failed int) {
	for _, result := range results {
		if result.Err == nil {
			loaded++
		} else {
			failed++
		}
	}
	return
}

func (c *Controller) startAutomations(ctx context.Context) (err error) {
	autoServices := auto.Services{
		Logger: c.Logger.Named("auto"),
		Node:   c.Node,
	}
	for _, autoConfig := range c.ControllerConfig.Automation {
		f, ok := c.SystemConfig.AutoFactories[autoConfig.Type]
		if !ok {
			err = multierr.Append(err, fmt.Errorf("unsupported automation type %v", autoConfig.Type))
			continue
		}
		impl := f.New(autoServices)
		// todo: keep track of running automations so we can update config and/or stop them
		if e := impl.Start(ctx); e != nil {
			err = multierr.Append(err, fmt.Errorf("start %s %w", autoConfig.Name, e))
		}
		if task.Configurable(impl) {
			if e := task.Configure(impl, autoConfig.Raw); e != nil {
				err = multierr.Append(err, fmt.Errorf("configure %s %w", autoConfig.Name, e))
			}
		}
	}
	return err
}
