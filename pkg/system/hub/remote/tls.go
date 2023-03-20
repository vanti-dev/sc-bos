package remote

import (
	"context"
	"crypto/x509"
	"net"
	"sync"

	"google.golang.org/grpc/credentials"
)

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
