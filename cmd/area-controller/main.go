package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/bsp-ew/internal/auth/policy"
	"github.com/vanti-dev/bsp-ew/internal/pki"
	"github.com/vanti-dev/bsp-ew/internal/testapi"
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
	keyPEM, err := pki.LoadOrGenerateKeyPair(filepath.Join(flagDataDir, "private-key.pem"))
	if err != nil {
		errs = multierr.Append(errs, err)
		return
	}

	// try to load an enrollment from disk
	enrollment, err := LoadEnrollment(filepath.Join(flagDataDir, "enrollment"), keyPEM)
	if errors.Is(err, ErrNotEnrolled) {
		// switch to enrollment mode, so this node can be enrolled with a Smart Core app server
		return runEnrollment(ctx, keyPEM)
	} else if err != nil {
		return err
	}

	return runNormal(ctx, enrollment)
}

func runEnrollment(ctx context.Context, keyPEM []byte) error {
	fmt.Println("TODO: enrollment mode")
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
