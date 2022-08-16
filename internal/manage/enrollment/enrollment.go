package enrollment

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"sync"

	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/proto"
)

func EnrollAreaController(ctx context.Context, enrollment *gen.Enrollment, ca *CA) (*gen.Enrollment, error) {
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

	certDER, err := ca.CreateEnrollmentCertificate(enrollment, peerPublicKey)
	if err != nil {
		return nil, err
	}
	enrollment.Certificate = ca.EncodeCertificateChain(certDER)

	client := gen.NewEnrollmentApiClient(conn)
	_, err = client.CreateEnrollment(ctx, &gen.CreateEnrollmentRequest{
		Enrollment: enrollment,
	})
	if err != nil {
		return nil, err
	}

	return enrollment, nil
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
