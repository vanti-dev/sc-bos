package enrollment

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/proto"

	"github.com/vanti-dev/sc-bos/internal/util/pki"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

// EnrollAreaController sets up the PKI for a remote Smart Core node.
// This connects to the remote node specified by enrollment.TargetAddress,
// constructs a new client certificate signed using the certificate and key from authority,
// and invokes CreateEnrollment on the target with this information.
// The Certificate and RootCAs will be computed from the authority and will be ignored if provided in enrollment.
func EnrollAreaController(
	ctx context.Context, enrollment *gen.Enrollment, authority pki.Source,
) (*gen.Enrollment, error) {
	enrollment = proto.Clone(enrollment).(*gen.Enrollment)

	// when in enrollment mode, the target node will be using a self-signed cert we won't be able to
	// automatically verify.
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	// the certInterceptor captures and saves the certificate presented by the server when the connection is opened
	creds := &certInterceptor{TransportCredentials: credentials.NewTLS(tlsConfig)}
	conn, err := grpc.DialContext(ctx, enrollment.TargetAddress,
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	peerCerts, ok := creds.PeerCertificates()
	if !ok || len(peerCerts) == 0 {
		return nil, errors.New("peer did not present any certificates")
	}
	peerPublicKey := peerCerts[0].PublicKey

	authorityCert, roots, err := authority.Certs()
	if err != nil {
		return nil, err
	}

	enrollment.RootCas = pki.EncodeCertificates(roots)

	certTemplate := newTargetCertificate(enrollment)
	enrollment.Certificate, err = pki.CreateCertificateChain(authorityCert, certTemplate, peerPublicKey,
		pki.WithAuthority(enrollment.TargetAddress),
	)
	if err != nil {
		return nil, err
	}

	client := gen.NewEnrollmentApiClient(conn)
	_, err = client.CreateEnrollment(ctx, &gen.CreateEnrollmentRequest{
		Enrollment: enrollment,
	})
	if err != nil {
		return nil, err
	}

	return enrollment, nil
}

func newTargetCertificate(enrollment *gen.Enrollment) *x509.Certificate {
	return &x509.Certificate{
		Subject:     pkix.Name{CommonName: enrollment.TargetName},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		URIs: []*url.URL{{
			Scheme: "smart-core",
			Opaque: enrollment.TargetName,
		}},
	}
}

type certInterceptor struct {
	credentials.TransportCredentials

	m                sync.Mutex
	populated        bool
	peerCertificates []*x509.Certificate
}

func (cw *certInterceptor) ClientHandshake(
	ctx context.Context, authority string, rawConn net.Conn,
) (net.Conn, credentials.AuthInfo, error) {
	conn, info, err := cw.TransportCredentials.ClientHandshake(ctx, authority, rawConn)
	if info, ok := info.(credentials.TLSInfo); ok {
		cw.m.Lock()
		defer cw.m.Unlock()

		if !cw.populated {
			cw.peerCertificates = info.State.PeerCertificates
			cw.populated = true
		}
	}

	return conn, info, err
}

func (cw *certInterceptor) PeerCertificates() (certs []*x509.Certificate, ok bool) {
	cw.m.Lock()
	defer cw.m.Unlock()

	return cw.peerCertificates, cw.populated
}

type Enrollment struct {
	RootDeviceName string `json:"root_device_name"`
	ManagerName    string `json:"manager_name"`
	ManagerAddress string `json:"manager_address"`

	RootCA *x509.Certificate `json:"-"`
	Cert   tls.Certificate   `json:"-"`
}

func (e Enrollment) Equal(other Enrollment) bool {
	if e.RootDeviceName != other.RootDeviceName ||
		e.ManagerName != other.ManagerName ||
		e.ManagerAddress != other.ManagerAddress {
		return false
	}
	return string(e.Cert.Certificate[0]) == string(other.Cert.Certificate[0])
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
//	<root>
//	  - enrollment.json - JSON-encoded Enrollment structure
//	  - root-ca.crt - Root CA for the enrollment, PEM-encoded X.509 certificate
//	  - cert.crt - Certificate chain for keyPEM, with the Root CA at the top of the chain
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
		return enrollment, errors.New("CA.pem is not valid PEM data")
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
