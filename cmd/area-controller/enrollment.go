package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/vanti-dev/bsp-ew/internal/util/pki"
	"github.com/vanti-dev/bsp-ew/internal/util/rpcutil"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Enrollment struct {
	RootDeviceName string `json:"root_device_name"`
	ManagerName    string `json:"manager_name"`
	ManagerAddress string `json:"manager_address"`

	RootCA *x509.Certificate `json:"-"`
	Cert   tls.Certificate   `json:"-"`
}

var ErrNotEnrolled = errors.New("node is not enrolled")

const (
	enrollmentFile = "enrollment.json"
	rootCaCertFile = "root-ca.cert.pem"
	certFile       = "enrollment.cert.pem"
)

// LoadEnrollment will load a previously saved Enrollment from a directory on disk. The directory should have
// the following structure:
//
//     <root>
//       - enrollment.json - JSON-encoded Enrollment structure
//       - root-ca.crt - Root CA for the enrollment, PEM-encoded X.509 certificate
//       - cert.crt - Certificate chain for keyPEM, with the Root CA at the top of the chain
//
// This node's private key must be passed in, in PEM-wrapped PKCS#1 or PKCS#8 format.
func LoadEnrollment(dir string, keyPEM []byte) (Enrollment, error) {
	// read and unmarshal enrollment metadata
	var enrollment Enrollment
	raw, err := os.ReadFile(filepath.Join(dir, enrollmentFile))
	if errors.Is(err, os.ErrNotExist) {
		return Enrollment{}, ErrNotEnrolled
	} else if err != nil {
		return Enrollment{}, err
	}
	err = json.Unmarshal(raw, &enrollment)
	if err != nil {
		return Enrollment{}, err
	}

	// load root CA certificate
	rootCaPem, err := os.ReadFile(filepath.Join(dir, rootCaCertFile))
	if err != nil {
		return enrollment, err
	}
	rootCaBlock, _ := pem.Decode(rootCaPem)
	if rootCaBlock == nil {
		return enrollment, errors.New("ca.pem is not valid PEM data")
	}
	if rootCaBlock.Type != "CERTIFICATE" {
		return enrollment, fmt.Errorf("expected PEM type 'CERTIFICATE' but got %q", rootCaBlock.Type)
	}
	rootCa, err := x509.ParseCertificate(rootCaBlock.Bytes)
	if err != nil {
		return enrollment, err
	}
	enrollment.RootCA = rootCa

	// load this node's certificate
	certPem, err := os.ReadFile(filepath.Join(dir, certFile))
	if err != nil {
		return enrollment, err
	}
	cert, err := tls.X509KeyPair(certPem, keyPEM)
	if err != nil {
		return enrollment, err
	}
	enrollment.Cert = cert

	return enrollment, nil
}

func SaveEnrollment(dir string, enrollment Enrollment) error {
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}

	jsonBytes, err := json.Marshal(enrollment)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(dir, enrollmentFile), jsonBytes, 0640)
	if err != nil {
		return err
	}

	_, err = pki.SaveCertificateChain(filepath.Join(dir, rootCaCertFile), [][]byte{enrollment.RootCA.Raw})
	if err != nil {
		return err
	}

	_, err = pki.SaveCertificateChain(filepath.Join(dir, certFile), enrollment.Cert.Certificate)
	return err
}

type EnrollmentServer struct {
	gen.UnimplementedEnrollmentApiServer
	logger *zap.Logger
	dir    string
	keyPEM []byte
	m      sync.Mutex
	done   chan struct{}
}

func NewEnrollmentServer(dir string, keyPEM []byte) *EnrollmentServer {
	return &EnrollmentServer{
		dir:    dir,
		keyPEM: keyPEM,
		done:   make(chan struct{}),
	}
}

func (es *EnrollmentServer) CreateEnrollment(ctx context.Context, request *gen.CreateEnrollmentRequest) (*gen.Enrollment, error) {
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

func (es *EnrollmentServer) Wait(ctx context.Context) (done bool) {
	select {
	case <-es.done:
		return true
	case <-ctx.Done():
		return false
	}
}
