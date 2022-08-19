package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/vanti-dev/bsp-ew/internal/app"
	"github.com/vanti-dev/bsp-ew/internal/auth/policy"
	"github.com/vanti-dev/bsp-ew/internal/manage/enrollment"
	"github.com/vanti-dev/bsp-ew/internal/testapi"
	"github.com/vanti-dev/bsp-ew/internal/util/pki"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

var (
	flagListenGRPC  string
	flagListenHTTPS string
	flagDataDir     string
	flagStaticDir   string
)

func init() {
	flag.StringVar(&flagListenGRPC, "listen-grpc", ":23557", "address (host:port) to host a Smart Core gRPC server on")
	flag.StringVar(&flagListenHTTPS, "listen-https", ":443", "address (host:port) to host a HTTPS server on")
	flag.StringVar(&flagDataDir, "data-dir", ".data/area-controller-01", "path to local data storage directory")
	flag.StringVar(&flagStaticDir, "static-dir", "ui/dist", "path for HTTP static resources")
}

func run(ctx context.Context) error {
	flag.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	// create data dir if it doesn't exist
	err = os.MkdirAll(flagDataDir, 0750)
	if err != nil {
		return err
	}

	// create private key if it doesn't exist
	key, keyPEM, err := pki.LoadOrGeneratePrivateKey(filepath.Join(flagDataDir, "private-key.pem"), logger)
	if err != nil {
		return err
	}

	enrollServer, err := enrollment.LoadOrCreateServer(filepath.Join(flagDataDir, "enrollment"), keyPEM, logger.Named("enrollment"))
	if err != nil {
		return err
	}

	// if the Area Controller is not already enrolled, we need to start it first in enrollment mode,
	// then restart into normal mode.
	if _, enrolled := enrollServer.Enrollment(); !enrolled {
		logger.Info("not enrolled; switching into enrollment mode")
		err = enrollment.Serve(ctx, logger.Named("enrollment"), enrollServer, key, flagListenGRPC)
		if err != nil {
			return err
		}
		logger.Info("switching from enrollment mode to normal mode")
	}

	return runNormal(ctx, logger, enrollServer)
}

func runNormal(ctx context.Context, logger *zap.Logger, enrollServer *enrollment.Server) error {
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
	gen.RegisterTestApiServer(grpcServer, testapi.NewAPI())

	grpcWebServer := grpcweb.WrapServer(grpcServer)
	staticFileHandler := http.FileServer(http.Dir(flagStaticDir))
	httpServer := &http.Server{
		Addr:      flagListenHTTPS,
		TLSConfig: tlsServerConfig,
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if grpcWebServer.IsGrpcWebRequest(request) || grpcWebServer.IsAcceptableGrpcCorsRequest(request) {
				grpcWebServer.ServeHTTP(writer, request)
			} else {
				staticFileHandler.ServeHTTP(writer, request)
			}
		}),
	}

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return app.ServeGRPC(ctx, grpcServer, flagListenGRPC, 15*time.Second, logger.Named("server.grpc"))
	})
	group.Go(func() error {
		return app.ServeHTTPS(ctx, httpServer, 15*time.Second, logger.Named("server.https"))
	})
	return group.Wait()
}

func main() {
	os.Exit(app.RunUntilInterrupt(run))
}
