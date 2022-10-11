package enrollment

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/vanti-dev/bsp-ew/internal/util/minibus"
	"sync"

	"github.com/vanti-dev/bsp-ew/internal/util/pki"
	"github.com/vanti-dev/bsp-ew/internal/util/rpcutil"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	gen.UnimplementedEnrollmentApiServer
	logger *zap.Logger
	dir    string
	keyPEM []byte

	m          sync.Mutex
	enrollment Enrollment
	done       chan struct{}

	managerAddressChanged minibus.Bus[string]
}

// NewServer creates an enrollment server, without attempting to load an existing enrollment.
// The new server will be in an un-enrolled state.
// New enrollments will be saved in the provided directory.
func NewServer(dir string, keyPEM []byte, logger *zap.Logger) *Server {
	es := &Server{
		logger: logger,
		dir:    dir,
		keyPEM: keyPEM,
		done:   make(chan struct{}),
	}
	return es
}

// LoadOrCreateServer will try to load an enrollment from disk. If successful, a server in the enrolled state is
// returned. Otherwise, a server in the unenrolled state is returned and new enrollments will be saved in the
// provided directory.
func LoadOrCreateServer(dir string, keyPEM []byte, logger *zap.Logger) (*Server, error) {
	es := NewServer(dir, keyPEM, logger)
	enrollment, err := LoadEnrollment(dir, keyPEM)
	if err == nil {
		es.enrollment = enrollment
		close(es.done)
	} else if !errors.Is(err, ErrNotEnrolled) {
		return nil, err
	}

	return es, nil
}

func (es *Server) CreateEnrollment(ctx context.Context, request *gen.CreateEnrollmentRequest) (*gen.Enrollment, error) {
	logger := rpcutil.ServerLogger(ctx, es.logger)

	// only allow one enrollment at a time
	es.m.Lock()
	defer es.m.Unlock()

	select {
	case <-es.done:
		return nil, status.Error(codes.AlreadyExists, "already enrolled")
	default:
	}

	cert, err := tls.X509KeyPair(request.GetEnrollment().GetCertificate(), es.keyPEM)
	if err != nil {
		logger.Error("invalid enrollment certificate", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid certificate")
	}

	roots, err := pki.ParseCertificatesPEM(request.GetEnrollment().GetRootCas())
	if err != nil {
		logger.Error("invalid enrollment root", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid root certificate(s)")
	}
	if len(roots) != 1 {
		return nil, status.Error(codes.InvalidArgument, "only 1 root CA is supported")
	}

	enrollment := Enrollment{
		RootDeviceName: request.GetEnrollment().GetTargetName(),
		ManagerName:    request.GetEnrollment().GetManagerName(),
		ManagerAddress: request.GetEnrollment().GetManagerAddress(),
		RootCA:         roots[0],
		Cert:           cert,
	}
	err = SaveEnrollment(es.dir, enrollment)
	if err != nil {
		logger.Error("failed to save enrollment", zap.Error(err), zap.String("dir", es.dir))
		return nil, status.Error(codes.Internal, "failed to save enrollment")
	}

	es.enrollment = enrollment
	es.managerAddressChanged.Send(context.Background(), enrollment.ManagerAddress)
	close(es.done)
	return request.GetEnrollment(), nil
}

func (es *Server) Wait(ctx context.Context) (enrollment Enrollment, done bool) {
	select {
	case <-es.done:
		es.m.Lock()
		defer es.m.Unlock()
		return es.enrollment, true
	case <-ctx.Done():
		return es.enrollment, false
	}
}

func (es *Server) Enrollment() (enrollment Enrollment, ok bool) {
	select {
	case <-es.done:
		es.m.Lock()
		defer es.m.Unlock()
		return es.enrollment, true
	default:
		return Enrollment{}, false
	}
}

// Certs implements pki.Source and provides a certificate source that provides the latest known enrollment certificate.
// If the certificate source is used while this Server has no enrollment, an error will be returned.
// This is therefore not suitable for use in enrollment mode - use a self-signed certificate source (with the enrollment
// private key) instead.
func (es *Server) Certs() (*tls.Certificate, []*x509.Certificate, error) {
	es.m.Lock()
	defer es.m.Unlock()

	// check that we are enrolled
	select {
	case <-es.done:
	default:
		return nil, nil, ErrNotEnrolled
	}

	cert := es.enrollment.Cert
	roots := []*x509.Certificate{es.enrollment.RootCA}
	return &cert, roots, nil
}

// ManagerAddress returns a chan that emits the manager address whenever it changes.
// Cancel the given context to stop listening for changes.
func (es *Server) ManagerAddress(ctx context.Context) <-chan string {
	changes := es.managerAddressChanged.Listen(ctx)
	en, ok := es.Enrollment()
	if !ok {
		return changes
	}

	// send the initial data right away
	out := make(chan string, 1)
	out <- en.ManagerAddress
	go func() {
		defer close(out)
		for addr := range changes {
			out <- addr
		}
	}()
	return out
}
