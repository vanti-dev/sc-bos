package main

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/bsp-ew/internal/app"
	"github.com/vanti-dev/bsp-ew/internal/auth/policy"
	"github.com/vanti-dev/bsp-ew/internal/testapi"
	"github.com/vanti-dev/bsp-ew/internal/util/pki"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.uber.org/multierr"
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

func run(ctx context.Context) (errs error) {
	flag.Parse()

	// create data dir if it doesn't exist
	err := os.MkdirAll(flagDataDir, 0750)
	if err != nil {
		errs = multierr.Append(errs, err)
	}

	// create private key if it doesn't exist
	key, keyPEM, err := pki.LoadOrGeneratePrivateKey(filepath.Join(flagDataDir, "private-key.pem"))
	if err != nil {
		errs = multierr.Append(errs, err)
		return
	}

	// try to load an enrollment from disk
	enrollment, err := LoadEnrollment(filepath.Join(flagDataDir, "enrollment"), keyPEM)
	if errors.Is(err, ErrNotEnrolled) {
		// switch to enrollment mode, so this node can be enrolled with a Smart Core app server
		return runEnrollment(ctx, key, keyPEM)
	} else if err != nil {
		return err
	}

	return runNormal(ctx, enrollment)
}

func runEnrollment(ctx context.Context, key crypto.PrivateKey, keyPEM []byte) error {
	enrollmentServer := NewEnrollmentServer(filepath.Join(flagDataDir, "enrollment"), keyPEM)
	certSource, err := pki.NewSelfSignedCertSource(key)
	if err != nil {
		return err
	}
	tlsConfig := &tls.Config{
		GetCertificate: certSource.TLSConfigGetCertificate,
	}

	grpcServer := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
	reflection.Register(grpcServer)
	gen.RegisterEnrollmentApiServer(grpcServer, enrollmentServer)

	srv := &app.Servers{
		GRPC:        grpcServer,
		GRPCAddress: flagListenGRPC,
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		err := srv.Serve(ctx)
		if err != nil {
			log.Printf("server stopped: %s", err.Error())
		}
	}()
	log.Println("gRPC serving; waiting for enrollment")

	ok := enrollmentServer.Wait(ctx)
	if ok {
		log.Printf("area controller is now enrolled")
	} else {
		return errors.New("server stopped without an enrollment")
	}
	return nil
}

func runNormal(ctx context.Context, enrollment Enrollment) error {
	clientRoot := x509.NewCertPool()
	clientRoot.AddCert(enrollment.RootCA)

	tlsServerConfig := &tls.Config{
		Certificates: []tls.Certificate{enrollment.Cert},
		ClientAuth:   tls.VerifyClientCertIfGiven,
		ClientCAs:    clientRoot,
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsServerConfig)),
		grpc.UnaryInterceptor(policy.GRPCUnaryInterceptor(nil)),
		grpc.StreamInterceptor(policy.GRPCStreamingInterceptor(nil)),
	)
	reflection.Register(grpcServer)
	gen.RegisterTestApiServer(grpcServer, testapi.NewAPI())

	grpcWebServer := grpcweb.WrapServer(grpcServer)
	staticFileHandler := http.FileServer(http.Dir(flagStaticDir))

	httpsHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if grpcWebServer.IsGrpcWebRequest(request) || grpcWebServer.IsAcceptableGrpcCorsRequest(request) {
			grpcWebServer.ServeHTTP(writer, request)
		} else {
			staticFileHandler.ServeHTTP(writer, request)
		}
	})
	httpsServer := &http.Server{
		Addr:      flagListenHTTPS,
		Handler:   httpsHandler,
		TLSConfig: tlsServerConfig,
	}

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		grpcListener, err := net.Listen("tcp", flagListenGRPC)
		if err != nil {
			return err
		}

		return grpcServer.Serve(grpcListener)
	})
	group.Go(func() error {
		// don't need to supply certs here because httpServer.TLSConfig is populated
		return httpsServer.ListenAndServeTLS("", "")
	})

	return group.Wait()
}

func logPublication(pub *traits.Publication) {
	fmt.Printf("\tAudience: %q\n", pub.GetAudience())
	fmt.Printf("\tMedia Type: %q\n", pub.GetMediaType())
	fmt.Printf("\tVersion: %q\n", pub.GetVersion())
	body := pub.GetBody()
	fmt.Printf("\tBody (%d bytes):\n", len(body))

	bodyRunes := []rune(strings.ToValidUTF8(string(body), "."))
	for len(bodyRunes) > 0 {
		var lineRunes []rune
		if len(bodyRunes) >= 64 {
			lineRunes = bodyRunes[:64]
			bodyRunes = bodyRunes[64:]
		} else {
			lineRunes = bodyRunes
			bodyRunes = nil
		}

		fmt.Printf("\t\t%s\n", string(lineRunes))
	}
	fmt.Println()
}

func main() {
	app.RunUntilInterrupt(run)
}
