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

	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/manage/enrollment"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/system"
	"github.com/vanti-dev/sc-bos/pkg/task"
	pki2 "github.com/vanti-dev/sc-bos/pkg/util/pki"
	"github.com/vanti-dev/sc-bos/pkg/util/pki/expire"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/rs/cors"
	"github.com/timshannon/bolthold"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/vanti-dev/sc-bos/pkg/auth/policy"
	"github.com/vanti-dev/sc-bos/pkg/auth/tenant"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

const LocalConfigFileName = "area-controller.local.json"

type SystemConfig struct {
	Logger      zap.Config
	ListenGRPC  string
	ListenHTTPS string

	DataDir             string
	StaticDir           string // hosts static files from this directory over HTTP if StaticDir is non-empty
	LocalConfigFileName string // defaults to LocalConfigFileName

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
	SystemFactories map[string]system.Factory // keyed by system type
}

func (sc SystemConfig) LocalConfigPath() string {
	s := sc.LocalConfigFileName
	if s == "" {
		s = LocalConfigFileName
	}
	return filepath.Join(sc.DataDir, s)
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
	localConfigPath := config.LocalConfigPath()
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

	// initialise services
	rootNode := node.New(localConfig.Name)
	rootNode.Logger = logger.Named("node")
	dbPath := filepath.Join(config.DataDir, "db.bolt")
	db, err := bolthold.Open(dbPath, 0750, nil)
	if err != nil {
		logger.Warn("failed to open local database file - some system components may fail", zap.Error(err),
			zap.String("path", dbPath))
	}

	// create private key if it doesn't exist
	key, keyPEM, err := pki2.LoadOrGeneratePrivateKey(filepath.Join(config.DataDir, "private-key.pem"), logger)
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
	certSource := pki2.ChainSource(
		enrollServer,
		pki2.CacheSource(pki2.FuncSource(func() (*tls.Certificate, []*x509.Certificate, error) {
			return readCertAndRoots(config, key)
		}), expire.BeforeInvalid(time.Hour)),
		pki2.CacheSource(pki2.SelfSignedSource(key, pki2.WithExpireAfter(30*24*time.Hour), pki2.WithIfaces()), expire.AfterProgress(0.5)),
	)
	tlsGRPCServerConfig := pki2.TLSServerConfig(certSource)
	tlsGRPCServerConfig.ClientAuth = tls.VerifyClientCertIfGiven
	tlsGRPCClientConfig := pki2.TLSClientConfig(certSource)

	tlsHTTPServerConfig := pki2.TLSServerConfig(certSource)

	// manager represents a delayed connection to the cohort manager.
	// Using the manager connection when we aren't enrolled will result in RPC calls returning 'not resolved' errors or similar.
	// As soon as we get enrolled those connections will be updated with the current manager address and will start to work.
	manager := node.DialChan(ctx, enrollServer.ManagerAddress(ctx),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsGRPCClientConfig)))

	mux := http.NewServeMux()
	if config.StaticDir != "" {
		mux.Handle("/", http.FileServer(http.Dir(config.StaticDir)))
	}

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
		mux.Handle("/oauth2/token", cors.Default().Handler(tokenServer))
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
		SystemConfig:     config,
		ControllerConfig: localConfig,
		Enrollment:       enrollServer,
		Logger:           logger,
		Node:             rootNode,
		Tasks:            &task.Group{},
		Database:         db,
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

func readCertAndRoots(config SystemConfig, key pki2.PrivateKey) (*tls.Certificate, []*x509.Certificate, error) {
	certPath := filepath.Join(config.DataDir, "server-cert.pem")
	cert, err := pki2.LoadX509Cert(certPath, key)
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
	roots, err := pki2.ParseCertificatesPEM(rootsPem)
	if err != nil {
		return nil, nil, err
	}
	return &cert, roots, nil
}

type Controller struct {
	SystemConfig     SystemConfig
	ControllerConfig ControllerConfig
	Enrollment       *enrollment.Server

	// services for drivers/automations
	Logger   *zap.Logger
	Node     *node.Node
	Tasks    *task.Group
	Database *bolthold.Store

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

	addFactorySupport(c.Node, c.SystemConfig.DriverFactories)
	addFactorySupport(c.Node, c.SystemConfig.AutoFactories)
	addFactorySupport(c.Node, c.SystemConfig.SystemFactories)

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

	// load and start the systems
	if err := c.startSystems(ctx); err != nil {
		return err
	}
	// load and start the drivers
	if err := c.startDrivers(ctx); err != nil {
		return err
	}
	// load and start the automations
	if err := c.startAutomations(ctx); err != nil {
		return err
	}

	err = multierr.Append(err, group.Wait())
	return
}

// addFactorySupport is used to register factories with a node to expose custom factory APIs.
// This checks each value in m and if that value has an API, via node.SelfSupporter, then it is registered with s.
func addFactorySupport[M ~map[K]F, K comparable, F any](s node.Supporter, m M) {
	for _, factory := range m {
		if api, ok := any(factory).(node.SelfSupporter); ok {
			api.AddSupport(s)
		}
	}
}

func (c *Controller) startDrivers(ctx context.Context) (err error) {
	driverServices := driver.Services{
		Logger:          c.Logger.Named("driver"),
		Node:            c.Node,
		Tasks:           &task.Group{},
		ClientTLSConfig: c.ClientTLSConfig,
	}

	var started []task.Starter
	go func() {
		<-ctx.Done()
		var err error
		for _, impl := range started {
			if task.Stoppable(impl) {
				err = multierr.Append(err, task.Stop(impl))
			}
		}
		if err != nil {
			c.Logger.Warn("Failed to cleanly stop drivers after ctx done", zap.Error(err))
		}
	}()

	for _, driverConfig := range c.ControllerConfig.Drivers {
		f, ok := c.SystemConfig.DriverFactories[driverConfig.Type]
		if !ok {
			err = multierr.Append(err, fmt.Errorf("unsupported driver type %v", driverConfig.Type))
			continue
		}
		impl := f.New(driverServices)
		if e := impl.Start(ctx); e != nil {
			err = multierr.Append(err, fmt.Errorf("start %s %w", driverConfig.Name, e))
		}
		// keep track so we can stop them if ctx ends
		started = append(started, impl)

		if task.Configurable(impl) {
			if e := task.Configure(impl, driverConfig.Raw); e != nil {
				err = multierr.Append(err, fmt.Errorf("configure %s %w", driverConfig.Name, e))
			}
		}
	}
	return err
}

func (c *Controller) startAutomations(ctx context.Context) (err error) {
	autoServices := auto.Services{
		Logger:       c.Logger.Named("auto"),
		Node:         c.Node,
		Database:     c.Database,
		GRPCServices: c.GRPC,
	}

	var started []task.Starter
	go func() {
		<-ctx.Done()
		var err error
		for _, impl := range started {
			if task.Stoppable(impl) {
				err = multierr.Append(err, task.Stop(impl))
			}
		}
		if err != nil {
			c.Logger.Warn("Failed to cleanly stop automations after ctx done", zap.Error(err))
		}
	}()

	for _, autoConfig := range c.ControllerConfig.Automation {
		f, ok := c.SystemConfig.AutoFactories[autoConfig.Type]
		if !ok {
			err = multierr.Append(err, fmt.Errorf("unsupported automation type %v", autoConfig.Type))
			continue
		}
		impl := f.New(autoServices)
		if e := impl.Start(ctx); e != nil {
			err = multierr.Append(err, fmt.Errorf("start %s %w", autoConfig.Name, e))
		}
		// keep track so we can stop them if ctx ends
		started = append(started, impl)

		if task.Configurable(impl) {
			if e := task.Configure(impl, autoConfig.Raw); e != nil {
				err = multierr.Append(err, fmt.Errorf("configure %s %w", autoConfig.Name, e))
			}
		}
	}
	return err
}

func (c *Controller) startSystems(ctx context.Context) (err error) {
	services := system.Services{
		Logger: c.Logger.Named("system"),
		Node:   c.Node,
	}

	var started []task.Starter
	go func() {
		<-ctx.Done()
		var err error
		for _, impl := range started {
			if task.Stoppable(impl) {
				err = multierr.Append(err, task.Stop(impl))
			}
		}
		if err != nil {
			c.Logger.Warn("Failed to cleanly stop systems after ctx done", zap.Error(err))
		}
	}()

	for k, cfg := range c.ControllerConfig.Systems {
		f, ok := c.SystemConfig.SystemFactories[k]
		if !ok {
			err = multierr.Append(err, fmt.Errorf("unsupported system type %v", k))
			continue
		}
		impl := f.New(services)
		if e := impl.Start(ctx); e != nil {
			err = multierr.Append(err, fmt.Errorf("start %s %w", k, e))
		}
		// keep track so we can stop them if ctx ends
		started = append(started, impl)

		if task.Configurable(impl) {
			if e := task.Configure(impl, cfg); e != nil {
				err = multierr.Append(err, fmt.Errorf("configure %s %w", k, e))
			}
		}
	}
	return err
}
