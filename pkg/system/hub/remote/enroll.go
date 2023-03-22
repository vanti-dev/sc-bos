package remote

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"net/url"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/proto"

	"github.com/vanti-dev/sc-bos/internal/util/pki"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

// Enroll sets up the PKI for a remote Smart Core node.
// This connects to the remote node specified by enrollment.TargetAddress,
// constructs a new client certificate signed using the certificate and key from authority,
// and invokes CreateEnrollment on the target with this information.
// The Certificate and RootCAs will be computed from the authority and will be ignored if provided in enrollment.
// If any remoteRoots are provided, the remote server will be checked using these as trust roots, otherwise any remote certificate will be allowed.
func Enroll(ctx context.Context, enrollment *gen.Enrollment, authority pki.Source, remoteRoots ...string) (*gen.Enrollment, error) {
	enrollment = proto.Clone(enrollment).(*gen.Enrollment)

	// when in enrollment mode, the target node will be using a self-signed cert we won't be able to
	// automatically verify.
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	if len(remoteRoots) > 0 {
		// when we have explicit remoteRoots we verify the remote certs!
		tlsConfig.InsecureSkipVerify = false
		tlsConfig.RootCAs = x509.NewCertPool()
		for _, root := range remoteRoots {
			tlsConfig.RootCAs.AppendCertsFromPEM([]byte(root))
		}

		// Make sure there's a timeout to avoid infinite (or 120s) connection loops in the case when tls fails.
		// There's no way to actually get the grpc.Dial func to return on first tls error, it will continue to retry
		// until ctx expires.
		var cleanUp context.CancelFunc
		ctx, cleanUp = context.WithTimeout(ctx, 30*time.Second)
		defer cleanUp()
	}

	// the certInterceptor captures and saves the certificate presented by the server when the connection is opened
	creds := &certInterceptor{TransportCredentials: credentials.NewTLS(tlsConfig)}
	conn, err := grpc.DialContext(ctx, enrollment.TargetAddress,
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
		grpc.WithReturnConnectionError(),
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

// Renew updates the PKI for a remote Smart Core node.
// This connects to the remote node specified by enrollment.TargetAddress using tlsConfig,
// signs the servers public certificate using authority,
// and calls EnrollmentApi.UpdateEnrollment on the remote node.
func Renew(ctx context.Context, enrollment *gen.Enrollment, authority pki.Source, tlsConfig *tls.Config) (*gen.Enrollment, error) {
	enrollment = proto.Clone(enrollment).(*gen.Enrollment)

	// the certInterceptor captures and saves the certificate presented by the server when the connection is opened
	creds := &certInterceptor{TransportCredentials: credentials.NewTLS(tlsConfig)}
	conn, err := grpc.DialContext(ctx, enrollment.TargetAddress,
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
		grpc.WithReturnConnectionError(),
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
	_, err = client.UpdateEnrollment(ctx, &gen.UpdateEnrollmentRequest{
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
