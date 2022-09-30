package pki

import (
	"crypto/tls"
	"crypto/x509"
)

// TLSConfig returns a *tls.Config whose client and server certs and roots are backed by source.
func TLSConfig(source Source) *tls.Config {
	r := &resolver{
		source: source,
		config: new(tls.Config),
	}

	r.config.GetCertificate = r.getCertificate
	r.config.GetClientCertificate = r.getClientCertificate
	r.config.InsecureSkipVerify = true // we do the verify ourselves via verifyConnection
	r.config.VerifyConnection = r.verifyConnection

	return r.config
}

type resolver struct {
	config *tls.Config
	source Source
}

func (r *resolver) getCertificate(_ *tls.ClientHelloInfo) (*tls.Certificate, error) {
	cert, _, err := r.source.Certs()
	return cert, err
}

func (r *resolver) getClientCertificate(_ *tls.CertificateRequestInfo) (*tls.Certificate, error) {
	cert, _, err := r.source.Certs()
	return cert, err
}

func (r *resolver) verifyConnection(cs tls.ConnectionState) error {
	_, roots, err := r.source.Certs()
	if err != nil {
		return err
	}

	if r.config.ClientAuth < tls.VerifyClientCertIfGiven &&
		len(cs.PeerCertificates) == 0 {
		return nil
	}

	pool := x509.NewCertPool()
	for _, root := range roots {
		pool.AddCert(root)
	}
	opts := x509.VerifyOptions{
		Roots:         pool,
		DNSName:       r.config.ServerName,
		Intermediates: x509.NewCertPool(),
	}

	if r.config.Time != nil {
		opts.CurrentTime = r.config.Time()
	}

	if r.config.ClientAuth >= tls.VerifyClientCertIfGiven {
		opts.KeyUsages = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	}

	// Copy intermediates certificates to verify options from cs if needed.
	// ignore cs.PeerCertificates[0] it refer to client certificates.
	for _, inter := range cs.PeerCertificates[1:] {
		opts.Intermediates.AddCert(inter)
	}

	_, err = cs.PeerCertificates[0].Verify(opts)
	return err
}
