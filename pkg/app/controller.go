package app

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/rs/cors"
	"github.com/timshannon/bolthold"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/smart-core-os/sc-golang/pkg/middleware/name"
	"github.com/vanti-dev/sc-bos/internal/manage/devices"
	"github.com/vanti-dev/sc-bos/internal/util/pki"
	"github.com/vanti-dev/sc-bos/internal/util/pki/expire"
	"github.com/vanti-dev/sc-bos/pkg/app/appconf"
	http2 "github.com/vanti-dev/sc-bos/pkg/app/http"
	"github.com/vanti-dev/sc-bos/pkg/app/sysconf"
	"github.com/vanti-dev/sc-bos/pkg/auth/policy"
	"github.com/vanti-dev/sc-bos/pkg/auth/token"
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/manage/enrollment"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/system"
	"github.com/vanti-dev/sc-bos/pkg/task"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"github.com/vanti-dev/sc-bos/pkg/task/serviceapi"
	"github.com/vanti-dev/sc-bos/pkg/zone"
)

// Bootstrap will obtain a Controller in a ready-to-run state.
func Bootstrap(ctx context.Context, config sysconf.Config) (*Controller, error) {
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
	localConfigPath := filepath.Join(config.DataDir, config.AppConfigFile)
	localConfig, err := appconf.LoadLocalConfig(config.DataDir, config.AppConfigFile)
	if localConfig == nil && err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logger.Debug("local config file not found", zap.String("path", localConfigPath))
			// continue with default config
			localConfig = &appconf.Config{}
		} else {
			return nil, err
		}
	} else if err != nil {
		// we loaded some config, but had some errors
		logger.Warn("failed to load some config", zap.String("path", localConfigPath), zap.Error(err))
	} else {
		// successfully loaded the config
		logger.Debug("loaded local config", zap.String("path", localConfigPath), zap.Strings("includes", localConfig.Includes))
	}

	// rootNode grants both local (in process) and networked (via grpc.Server) access to controller apis.
	// If you have a device you want expose via a Smart Core API, rootNode is where you'd do that via Announce.
	// If you need to know the brightness of another controller device, rootNode.Clients is how you'd do that.
	cName := localConfig.Name
	if cName == "" {
		cName = config.Name
	}
	rootNode := node.New(cName)
	rootNode.Logger = logger.Named("node")

	// Setup a local database for storing non-critical data.
	// This is made available to automations and systems as a local cache, for example for lighting reports.
	dbPath := filepath.Join(config.DataDir, "db.bolt")
	db, err := bolthold.Open(dbPath, 0750, nil)
	if err != nil {
		logger.Warn("failed to open local database file - some system components may fail", zap.Error(err),
			zap.String("path", dbPath))
	}

	certConfig := config.CertConfig
	// Create a private key if it doesn't exist.
	// This key is used by the controller for incoming and outgoing connections, and as part of the enrolment process.
	key, keyPEM, err := pki.LoadOrGeneratePrivateKey(filepath.Join(config.DataDir, certConfig.KeyFile), logger)
	if err != nil {
		return nil, err
	}

	// enrollServer manages this controllers participation in a cohort.
	// When registered with a grpc.Server, the enrollServer will accept requests like CreateEnrollment which gives
	// this controller a new certificate for use during outgoing TLS connections to other cohort members.
	// In addition the enrollment will also include details of the trusted root certs so this controller can validate
	// incoming connections that contain a client certificate.
	//
	// enrollServer implements pki.Source providing these features without any extra work to setup.
	enrollServer, err := enrollment.LoadOrCreateServer(filepath.Join(config.DataDir, "enrollment"), keyPEM, logger.Named("enrollment"))
	if err != nil {
		return nil, err
	}

	// fileSource attempts to load a certificate and trust roots from disk.
	// The certificates public key must be paired with private key `key` loaded above.
	fileSource := pki.CacheSource(
		pki.FSKeySource(
			filepath.Join(config.DataDir, certConfig.CertFile), key,
			joinIfPresent(config.DataDir, certConfig.RootsFile)),
		expire.BeforeInvalid(time.Hour),
	)

	// systemSource allows systems to contribute certificates to incoming and outgoing gRPC connections.
	systemSource := &pki.SourceSet{}

	// selfSignedSource creates a self signed certificate.
	// The certificates public key will be paired with and signed by `key`.
	selfSignedSource := pki.CacheSource(
		pki.SelfSignedSource(key, pki.WithExpireAfter(30*24*time.Hour), pki.WithIfaces()),
		expire.AfterProgress(0.5),
		pki.WithFSCache(filepath.Join(config.DataDir, "grpc-self-signed.cert.pem"), "", key),
	)

	// grpcSource is used by both incoming and outgoing grpc connections.
	// The server present the sources certificate and any intermediates between it and the roots to clients during TLS handshake.
	// If an incoming connection presents a client cert then it will be validated against the roots.
	// Outgoing connections will present the sources certificate as a client cert for validation by the remote party.
	// Config can indicate that different certs be used for grpc and https (inc grpc-web)
	grpcSource := &pki.SourceSet{
		enrollServer,
		fileSource,
		systemSource,
		selfSignedSource,
	}
	tlsGRPCServerConfig := pki.TLSServerConfig(grpcSource)
	tlsGRPCClientConfig := pki.TLSClientConfig(grpcSource)

	// Certs used for https (hosting and grpc-web) can be different from the Smart Core native grpc endpoint,
	// mostly to allow support for trusted certs on the https interface and cohort managed certs for grpc requests.
	httpCertSource := grpcSource
	if certConfig.HTTPCert {
		fileSource := pki.CacheSource(
			pki.FSSource(
				filepath.Join(config.DataDir, certConfig.HTTPCertFile),
				filepath.Join(config.DataDir, certConfig.HTTPKeyFile),
				""),
			expire.After(30*time.Minute),
		)
		httpCertSource = &pki.SourceSet{
			fileSource,
			selfSignedSource, // reuse the same self signed cert from grpc requests
		}
	}
	tlsHTTPServerConfig := pki.TLSServerConfig(httpCertSource)

	// manager represents a delayed connection to the cohort manager.
	// Using the manager connection when we aren't enrolled will result in RPC calls returning 'not resolved' errors or similar.
	// As soon as we get enrolled those connections will be updated with the current manager address and will start to work.
	manager := node.DialChan(ctx, enrollServer.ManagerAddress(ctx),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsGRPCClientConfig)))

	mux := http.NewServeMux()
	for _, site := range config.StaticHosting {
		handler := http2.NewStaticHandler(site.FilePath)
		mux.Handle(site.Path, http.StripPrefix(site.Path, handler))
		logger.Info("Serving static site", zap.String("path", site.Path), zap.String("filePath", site.FilePath))
	}

	var grpcOpts []grpc.ServerOption
	grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(tlsGRPCServerConfig)))

	// tokenValidator validates access tokens as part of the authorisation of requests to our APIs.
	// Claims associated with the token are presented along with other information when processing policy files.
	// Systems contribute validators to this set supporting different sources of token.
	tokenValidator := &token.ValidatorSet{}

	// configure request authorisation, here we setup grpc interceptors that decide if a request is denied or not.
	if !config.DisablePolicy {
		pol := policy.Default(false)
		if config.Policy != nil {
			pol = config.Policy
		}

		interceptor := policy.NewInterceptor(pol,
			policy.WithLogger(logger.Named("policy")),
			policy.WithTokenVerifier(tokenValidator),
		)
		grpcOpts = append(grpcOpts,
			grpc.ChainUnaryInterceptor(interceptor.GRPCUnaryInterceptor()),
			grpc.ChainStreamInterceptor(interceptor.GRPCStreamingInterceptor()),
		)
	}

	if rootNode.Name() != "" {
		grpcOpts = append(grpcOpts,
			grpc.ChainUnaryInterceptor(name.IfAbsentUnaryInterceptor(rootNode.Name())),
			grpc.ChainStreamInterceptor(name.IfAbsentStreamInterceptor(rootNode.Name())),
		)
	}

	grpcServer := grpc.NewServer(grpcOpts...)
	reflection.Register(grpcServer)
	gen.RegisterEnrollmentApiServer(grpcServer, enrollServer)
	devices.NewServer(rootNode).Register(grpcServer)
	// support the services api for managing drivers, automations, and systems
	serviceRouter := gen.NewServicesApiRouter()
	rootNode.Support(node.Routing(serviceRouter), node.Clients(gen.WrapServicesApi(serviceRouter)))

	grpcWebServer := grpcweb.WrapServer(grpcServer, grpcweb.WithOriginFunc(func(origin string) bool {
		return true
	}))

	// configure CORS setup
	co := cors.New(cors.Options{
		AllowedOrigins:   config.Cors.CorsOrigins,
		AllowCredentials: true,
		AllowedHeaders:   []string{http2.HeaderAllowOrigin, http2.HeaderAuthorization, http2.HeaderContentType},
		AllowedMethods:   []string{http.MethodHead, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		Debug:            config.Cors.DebugMode,
	})
	corsWrap := co.Handler(mux)

	httpServer := &http.Server{
		Addr:      config.ListenHTTPS,
		TLSConfig: tlsHTTPServerConfig,
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if grpcWebServer.IsGrpcWebRequest(request) || grpcWebServer.IsAcceptableGrpcCorsRequest(request) {
				grpcWebServer.ServeHTTP(writer, request)
			} else {
				corsWrap.ServeHTTP(writer, request)
			}
		}),
	}

	c := &Controller{
		SystemConfig:     config,
		ControllerConfig: *localConfig,
		Enrollment:       enrollServer,
		Logger:           logger,
		Node:             rootNode,
		Tasks:            &task.Group{},
		Database:         db,
		TokenValidators:  tokenValidator,
		GRPCCerts:        systemSource,
		PrivateKey:       key,
		Mux:              mux,
		GRPC:             grpcServer,
		HTTP:             httpServer,
		ClientTLSConfig:  tlsGRPCClientConfig,
		ManagerConn:      manager,
	}
	c.Defer(manager.Close)
	return c, nil
}

type Controller struct {
	SystemConfig     sysconf.Config
	ControllerConfig appconf.Config
	Enrollment       *enrollment.Server

	// services for drivers/automations
	Logger          *zap.Logger
	Node            *node.Node
	Tasks           *task.Group
	Database        *bolthold.Store
	TokenValidators *token.ValidatorSet
	GRPCCerts       *pki.SourceSet

	Mux  *http.ServeMux
	GRPC *grpc.Server
	HTTP *http.Server

	PrivateKey      pki.PrivateKey
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
	systemServices, err := c.startSystems()
	if err != nil {
		return err
	}
	c.Node.Announce("systems", node.HasClient(gen.WrapServicesApi(serviceapi.NewApi(systemServices, serviceapi.WithKnownTypesFromMapKeys(c.SystemConfig.SystemFactories)))))
	go logServiceMapChanges(ctx, c.Logger.Named("system"), systemServices)
	// load and start the drivers
	driverServices, err := c.startDrivers()
	if err != nil {
		return err
	}
	c.Node.Announce("drivers", node.HasClient(gen.WrapServicesApi(serviceapi.NewApi(driverServices, serviceapi.WithKnownTypesFromMapKeys(c.SystemConfig.DriverFactories)))))
	go logServiceMapChanges(ctx, c.Logger.Named("driver"), driverServices)
	// load and start the automations
	autoServices, err := c.startAutomations()
	if err != nil {
		return err
	}
	c.Node.Announce("automations", node.HasClient(gen.WrapServicesApi(serviceapi.NewApi(autoServices, serviceapi.WithKnownTypesFromMapKeys(c.SystemConfig.AutoFactories)))))
	go logServiceMapChanges(ctx, c.Logger.Named("auto"), autoServices)
	// load and start the zones
	zoneServices, err := c.startZones()
	if err != nil {
		return err
	}
	c.Node.Announce("zones", node.HasClient(gen.WrapServicesApi(serviceapi.NewApi(zoneServices, serviceapi.WithKnownTypesFromMapKeys(c.SystemConfig.ZoneFactories)))))
	go logServiceMapChanges(ctx, c.Logger.Named("zone"), zoneServices)

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

func (c *Controller) startDrivers() (*service.Map, error) {
	ctxServices := driver.Services{
		Logger:          c.Logger.Named("driver"),
		Node:            c.Node,
		ClientTLSConfig: c.ClientTLSConfig,
		HTTPMux:         c.Mux,
	}

	m := service.NewMap(func(kind string) (service.Lifecycle, error) {
		f, ok := c.SystemConfig.DriverFactories[kind]
		if !ok {
			return nil, fmt.Errorf("unsupported driver type %v", kind)
		}
		return f.New(ctxServices), nil
	}, service.IdIsRequired)

	var allErrs error
	for _, cfg := range c.ControllerConfig.Drivers {
		_, _, err := m.Create(cfg.Name, cfg.Type, service.State{Active: !cfg.Disabled, Config: cfg.Raw})
		allErrs = multierr.Append(allErrs, err)
	}
	return m, allErrs
}

func (c *Controller) startAutomations() (*service.Map, error) {
	ctxServices := auto.Services{
		Logger:       c.Logger.Named("auto"),
		Node:         c.Node,
		Database:     c.Database,
		GRPCServices: c.GRPC,
	}

	m := service.NewMap(func(kind string) (service.Lifecycle, error) {
		f, ok := c.SystemConfig.AutoFactories[kind]
		if !ok {
			return nil, fmt.Errorf("unsupported automation type %v", kind)
		}
		return f.New(ctxServices), nil
	}, service.IdIsRequired)

	var allErrs error
	for _, cfg := range c.ControllerConfig.Automation {
		_, _, err := m.Create(cfg.Name, cfg.Type, service.State{Active: !cfg.Disabled, Config: cfg.Raw})
		allErrs = multierr.Append(allErrs, err)
	}
	return m, allErrs
}

func (c *Controller) startSystems() (*service.Map, error) {
	grpcEndpoint, err := c.grpcEndpoint()
	if err != nil {
		return nil, err
	}
	ctxServices := system.Services{
		DataDir:         c.SystemConfig.DataDir,
		Logger:          c.Logger.Named("system"),
		Node:            c.Node,
		GRPCEndpoint:    grpcEndpoint,
		Database:        c.Database,
		HTTPMux:         c.Mux,
		TokenValidators: c.TokenValidators,
		GRPCCerts:       c.GRPCCerts,
		PrivateKey:      c.PrivateKey,
		CohortManager:   c.ManagerConn,
		ClientTLSConfig: c.ClientTLSConfig,
	}
	m := service.NewMap(func(kind string) (service.Lifecycle, error) {
		f, ok := c.SystemConfig.SystemFactories[kind]
		if !ok {
			return nil, fmt.Errorf("unsupported system type %v", kind)
		}
		return f.New(ctxServices), nil
	}, service.IdIsKind)

	var allErrs error
	for kind, cfg := range c.SystemConfig.Systems {
		_, _, err := m.Create("", kind, service.State{Active: !cfg.Disabled, Config: cfg.Raw})
		allErrs = multierr.Append(allErrs, err)
	}
	return m, allErrs
}
func (c *Controller) startZones() (*service.Map, error) {
	ctxServices := zone.Services{
		Logger: c.Logger.Named("auto"),
		Node:   c.Node,
	}

	m := service.NewMap(func(kind string) (service.Lifecycle, error) {
		f, ok := c.SystemConfig.ZoneFactories[kind]
		if !ok {
			return nil, fmt.Errorf("unsupported zone type %v", kind)
		}
		return f.New(ctxServices), nil
	}, service.IdIsRequired)

	var allErrs error
	for _, cfg := range c.ControllerConfig.Zones {
		_, _, err := m.Create(cfg.Name, cfg.Type, service.State{Active: !cfg.Disabled, Config: cfg.Raw})
		allErrs = multierr.Append(allErrs, err)
	}
	return m, allErrs
}

func (c *Controller) grpcEndpoint() (string, error) {
	lisAddr := c.SystemConfig.ListenGRPC
	addr := c.SystemConfig.GRPCAddr
	_, p, err := net.SplitHostPort(lisAddr)
	if err != nil {
		return "", err
	}
	return net.JoinHostPort(addr, p), nil
}

func (c *Controller) httpEndpoint() (string, error) {
	lisAddr := c.SystemConfig.ListenHTTPS
	addr := c.SystemConfig.HTTPAddr
	_, p, err := net.SplitHostPort(lisAddr)
	if err != nil {
		return "", err
	}
	return net.JoinHostPort(addr, p), nil
}

func logServiceMapChanges(ctx context.Context, logger *zap.Logger, m *service.Map) {
	known := map[string]func(){}
	changes := m.Listen(ctx)
	for _, record := range m.Values() {
		ctx, stop := context.WithCancel(ctx)
		known[record.Id] = stop
		record := record
		go logServiceRecordChanges(ctx, logger, record)
	}
	for change := range changes {
		if change.OldValue == nil && change.NewValue != nil {
			// add
			if _, ok := known[change.NewValue.Id]; ok {
				continue // deal with potential race between Listen and Values
			}
			ctx, stop := context.WithCancel(ctx)
			known[change.NewValue.Id] = stop
			go logServiceRecordChanges(ctx, logger, change.NewValue)
		} else if change.OldValue != nil && change.NewValue == nil {
			// remove
			stop, ok := known[change.OldValue.Id]
			if !ok {
				continue
			}
			delete(known, change.OldValue.Id)
			stop()
		}
	}
}

func logServiceRecordChanges(ctx context.Context, logger *zap.Logger, r *service.Record) {
	logger = logger.With(zap.String("id", r.Id), zap.String("kind", r.Kind))
	state, changes := r.Service.StateAndChanges(ctx)
	lastMode := ""
	logMode := func(change service.State) {
		mode := ""
		switch {
		case !change.Active && change.Err != nil:
			mode = "error"
		case !change.Active:
			mode = "Stopped"
		case change.Loading:
			mode = "Loading"
		case change.Active:
			mode = "Running"
		}
		if mode == lastMode {
			return
		}
		switch mode {
		case "error":
			logger.Warn("Failed to load", zap.Error(change.Err))
		case "":
			return
		case "Stopped":
			if lastMode == "" {
				logger.Debug("Created")
			} else {
				logger.Debug(mode)
			}
		default:
			logger.Debug(mode)
		}
		lastMode = mode
	}
	logMode(state)
	for change := range changes {
		logMode(change)
	}
}

func joinIfPresent(dir, path string) string {
	if path == "" {
		return ""
	}
	return filepath.Join(dir, path)
}
