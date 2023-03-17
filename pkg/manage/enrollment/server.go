package enrollment

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vanti-dev/sc-bos/internal/util/pki"
	"github.com/vanti-dev/sc-bos/internal/util/rpcutil"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/minibus"
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
	go es.managerAddressChanged.Send(context.Background(), enrollment.ManagerAddress)
	close(es.done)
	return request.GetEnrollment(), nil
}

func (es *Server) DeleteEnrollment(ctx context.Context, request *gen.DeleteEnrollmentRequest) (*gen.Enrollment, error) {
	// delete only if we are enrolled already
	es.m.Lock()
	defer es.m.Unlock()

	select {
	case <-es.done:
	default:
		return nil, status.Error(codes.NotFound, "not enrolled")
	}

	// 0. remember state needed to undo this delete if needed
	// 1. clear the enrollment pki.Source
	// 2. delete the enrollment files
	// 3. reset tracking state

	en := es.enrollment
	es.enrollment = Enrollment{}
	if err := DeleteEnrollment(es.dir); err != nil {
		es.enrollment = en
		// try our best to save the enrollment again
		var saveErr error
		for i := 0; i < 5; i++ {
			saveErr = SaveEnrollment(es.dir, en)
			if saveErr == nil {
				break
			}
		}
		if saveErr != nil {
			es.logger.Error("delete failed, rollback failed - manual intervention required", zap.NamedError("delErr", err), zap.NamedError("rollbackErr", saveErr))
			return nil, status.Errorf(codes.DataLoss, "failed to delete, failed to rollback")
		}
		return nil, status.Errorf(codes.Aborted, err.Error())
	}
	es.done = make(chan struct{})
	go es.managerAddressChanged.Send(context.Background(), "")

	return &gen.Enrollment{
		TargetName:     en.RootDeviceName,
		ManagerName:    en.ManagerName,
		ManagerAddress: en.ManagerAddress,
	}, nil
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
