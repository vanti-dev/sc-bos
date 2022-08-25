package app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/vanti-dev/bsp-ew/internal/auth/policy"
	"github.com/vanti-dev/bsp-ew/internal/manage/enrollment"
	"github.com/vanti-dev/bsp-ew/internal/util/pki"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type SystemConfig struct {
	Logger      zap.Config
	DataDir     string
	ListenGRPC  string
	ListenHTTPS string
}

// Bootstrap will obtain a Controller in a ready-to-run state.
// If there is no saved enrollment, then Bootstrap will start an enrollment server and wait for the enrollment to
// complete.
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

	// create private key if it doesn't exist
	key, keyPEM, err := pki.LoadOrGeneratePrivateKey(filepath.Join(config.DataDir, "private-key.pem"), logger)
	if err != nil {
		return nil, err
	}

	enrollServer, err := enrollment.LoadOrCreateServer(filepath.Join(config.DataDir, "enrollment"), keyPEM, logger.Named("enrollment"))
	if err != nil {
		return nil, err
	}

	// if the Area Controller is not already enrolled, we need to start it first in enrollment mode,
	// then restart into normal mode.
	en, enrolled := enrollServer.Enrollment()
	if !enrolled {
		logger.Info("not enrolled; switching into enrollment mode")
		err = ServeEnrollment(ctx, logger.Named("enrollment"), enrollServer, key, config.ListenGRPC)
		if err != nil {
			return nil, err
		}
		logger.Info("switching from enrollment mode to normal mode")
		en, enrolled = enrollServer.Enrollment()
		if !enrolled {
			panic("we just enrolled successfully, but it's somehow not saved")
		}
	}

	clientRoot := x509.NewCertPool()
	clientRoot.AddCert(en.RootCA)
	tlsServerConfig := &tls.Config{
		GetCertificate: enrollServer.CertSource().TLSConfigGetCertificate,
		ClientAuth:     tls.VerifyClientCertIfGiven,
		ClientCAs:      clientRoot,
	}

	interceptor := policy.NewInterceptor(policy.Default(), policy.WithLogger(logger.Named("policy")))
	grpcServer := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsServerConfig)),
		grpc.UnaryInterceptor(interceptor.GRPCUnaryInterceptor()),
		grpc.StreamInterceptor(interceptor.GRPCStreamingInterceptor()),
	)
	reflection.Register(grpcServer)
	gen.RegisterEnrollmentApiServer(grpcServer, enrollServer)

	grpcWebServer := grpcweb.WrapServer(grpcServer, grpcweb.WithOriginFunc(func(origin string) bool {
		return true
	}))
	mux := http.NewServeMux()

	httpServer := &http.Server{
		Addr:      config.ListenHTTPS,
		TLSConfig: tlsServerConfig,
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if grpcWebServer.IsGrpcWebRequest(request) || grpcWebServer.IsAcceptableGrpcCorsRequest(request) {
				grpcWebServer.ServeHTTP(writer, request)
			} else {
				mux.ServeHTTP(writer, request)
			}
		}),
	}

	return &Controller{
		Logger:     logger,
		Config:     config,
		Enrollment: enrollServer,
		Mux:        mux,
		GRPC:       grpcServer,
		HTTP:       httpServer,
	}, nil
}

type Controller struct {
	Logger     *zap.Logger
	Config     SystemConfig
	Enrollment *enrollment.Server

	Mux  *http.ServeMux
	GRPC *grpc.Server
	HTTP *http.Server
}

func (c *Controller) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)
	if c.Config.ListenGRPC != "" {
		group.Go(func() error {
			return ServeGRPC(ctx, c.GRPC, c.Config.ListenGRPC, 15*time.Second, c.Logger.Named("server.grpc"))
		})
	}
	if c.Config.ListenHTTPS != "" {
		group.Go(func() error {
			return ServeHTTPS(ctx, c.HTTP, 15*time.Second, c.Logger.Named("server.https"))
		})
	}
	return group.Wait()
}
