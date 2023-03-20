package enroll

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"net"
	"net/url"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/proto"

	"github.com/vanti-dev/sc-bos/internal/util/pki"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

// Controller sets up the PKI for a remote Smart Core node.
// This connects to the remote node specified by enrollment.TargetAddress,
// constructs a new client certificate signed using the certificate and key from authority,
// and invokes CreateEnrollment on the target with this information.
// The Certificate and RootCAs will be computed from the authority and will be ignored if provided in enrollment.
func Controller(ctx context.Context, enrollment *gen.Enrollment, authority pki.Source) (*gen.Enrollment, error) {
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

func (cw *certInterceptor) ClientHandshake(ctx context.Context, authority string, rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) {
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
