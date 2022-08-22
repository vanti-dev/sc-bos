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

type Controller struct {
	Logger           *zap.Logger
	DataDir          string
	ListenGRPC       string
	ListenHTTPS      string
	Routes           map[string]http.Handler
	RegisterServices func(server *grpc.Server)
}

func (c *Controller) Run(ctx context.Context) error {
	logger := c.Logger
	if logger == nil {
		var err error
		logger, err = zap.NewDevelopment()
		if err != nil {
			return err
		}
	}

	// create data dir if it doesn't exist
	err := os.MkdirAll(c.DataDir, 0750)
	if err != nil {
		return err
	}

	// create private key if it doesn't exist
	key, keyPEM, err := pki.LoadOrGeneratePrivateKey(filepath.Join(c.DataDir, "private-key.pem"), logger)
	if err != nil {
		return err
	}

	enrollServer, err := enrollment.LoadOrCreateServer(filepath.Join(c.DataDir, "enrollment"), keyPEM, logger.Named("enrollment"))
	if err != nil {
		return err
	}

	// if the Area Controller is not already enrolled, we need to start it first in enrollment mode,
	// then restart into normal mode.
	if _, enrolled := enrollServer.Enrollment(); !enrolled {
		logger.Info("not enrolled; switching into enrollment mode")
		err = ServeEnrollment(ctx, logger.Named("enrollment"), enrollServer, key, c.ListenGRPC)
		if err != nil {
			return err
		}
		logger.Info("switching from enrollment mode to normal mode")
	}

	return c.runNormal(ctx, logger, enrollServer)
}

func (c *Controller) runNormal(ctx context.Context, logger *zap.Logger, enrollServer *enrollment.Server) error {
	en, ok := enrollServer.Enrollment()
	if !ok {
		return enrollment.ErrNotEnrolled
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
	if c.RegisterServices != nil {
		c.RegisterServices(grpcServer)
	}

	grpcWebServer := grpcweb.WrapServer(grpcServer)

	mux := http.NewServeMux()
	for path, handler := range c.Routes {
		mux.Handle(path, handler)
	}
	httpServer := &http.Server{
		Addr:      c.ListenHTTPS,
		TLSConfig: tlsServerConfig,
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if grpcWebServer.IsGrpcWebRequest(request) || grpcWebServer.IsAcceptableGrpcCorsRequest(request) {
				grpcWebServer.ServeHTTP(writer, request)
			} else {
				mux.ServeHTTP(writer, request)
			}
		}),
	}

	group, ctx := errgroup.WithContext(ctx)
	if c.ListenGRPC != "" {
		group.Go(func() error {
			return ServeGRPC(ctx, grpcServer, c.ListenGRPC, 15*time.Second, logger.Named("server.grpc"))
		})
	}
	if c.ListenHTTPS != "" {
		group.Go(func() error {
			return ServeHTTPS(ctx, httpServer, 15*time.Second, logger.Named("server.https"))
		})
	}
	return group.Wait()
}
