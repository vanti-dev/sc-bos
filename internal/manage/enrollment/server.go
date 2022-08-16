package enrollment

import (
	"context"
	"crypto/tls"
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
	m      sync.Mutex
	done   chan struct{}
}

func NewServer(dir string, keyPEM []byte, logger *zap.Logger) *Server {
	return &Server{
		logger: logger,
		dir:    dir,
		keyPEM: keyPEM,
		done:   make(chan struct{}),
	}
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

	close(es.done)
	return request.GetEnrollment(), nil
}

func (es *Server) Wait(ctx context.Context) (done bool) {
	select {
	case <-es.done:
		return true
	case <-ctx.Done():
		return false
	}
}
