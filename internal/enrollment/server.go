package enrollment

import (
	"context"
	"crypto/x509"
	"log"

	"github.com/vanti-dev/bsp-ew/internal/pki"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	gen.UnimplementedEnrollmentApiServer

	pubKey pki.PublicKey
}

func (e *Server) GetEnrollment(ctx context.Context, request *gen.GetEnrollmentRequest) (*gen.Enrollment, error) {
	return nil, status.Error(codes.NotFound, "node is not enrolled")
}

func (e *Server) CreateEnrollment(ctx context.Context, request *gen.CreateEnrollmentRequest) (*gen.Enrollment, error) {
	enrollment := request.GetEnrollment()
	if enrollment == nil {
		return nil, status.Error(codes.InvalidArgument, "enrollment must be provided")
	}

	leaf, intermediates, err := pki.ParseCertificateChainPEM(enrollment.Certificate)
	if err != nil {
		log.Printf("failed to parse certificate as a PEM cert chain: %s", err.Error())
		return nil, status.Error(codes.InvalidArgument, "certificate chain failed to parse")
	}
	if !slices.Contains(pki.GetCertificateSmartCoreNames(leaf), enrollment.TargetName) {
		return nil, status.Error(codes.InvalidArgument, "leaf certificate does not contain 'target_name' as a Subject Alternative Name")
	}
	// check the new cert is actually signed for our public key
	if !e.pubKey.Equal(leaf.PublicKey) {
		return nil, status.Error(codes.InvalidArgument, "leaf certificate is not signed for the correct public key")
	}

	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(enrollment.RootCas) {
		log.Printf("failed to parse root CA certificates")
		return nil, status.Error(codes.InvalidArgument, "root_cas failed to parse")
	}

	_, err = leaf.Verify(x509.VerifyOptions{
		Intermediates: intermediates,
		Roots:         roots,
	})
	if err != nil {
		log.Printf("peer leaf certificate failed verification: %s", err.Error())
		return nil, status.Error(codes.InvalidArgument, "leaf certificate failed verification")
	}

	return enrollment, nil
}

func (e *Server) DeleteEnrollment(ctx context.Context, request *gen.DeleteEnrollmentRequest) (*gen.Enrollment, error) {
	return nil, status.Error(codes.NotFound, "node is not enrolled")
}

var _ gen.EnrollmentApiServer = (*Server)(nil)
