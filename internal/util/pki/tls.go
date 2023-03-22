package pki

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
)

// TLSServerConfig returns a *tls.Config for use by a server using source to provide the server cert.
// If the returned tls.Config requires validation of client certificates then sources roots will be used to validate
// the client certificates.
//
// Warning: if tls.Config.ClientAuth is configured to verify client certificates then this happens _before_ source is
// consulted for ClientCAs which will likely cause verification to fail.
func TLSServerConfig(source Source) *tls.Config {
	r := &resolver{
		server: true,
		source: source,
		config: new(tls.Config),
	}

	// We want to set config.ClientAuth to tls.VerifyClientCertIfGiven but this triggers cert validation that doesn't use source.
	// So we have to use non-verifying ClientAuth settings and take over the verification ourselves using r.clientAuth.
	r.config.ClientAuth = tls.RequestClientCert
	r.clientAuth = tls.VerifyClientCertIfGiven

	r.config.GetCertificate = r.getCertificate
	r.config.VerifyConnection = r.verifyConnection

	return r.config
}

// TLSClientConfig returns a *tls.Config for use by a client using sources roots to validate the server certificate.
// If the server requests a client certificate then sources cert will be used.
func TLSClientConfig(source Source) *tls.Config {
	r := &resolver{
		source: source,
		config: new(tls.Config),
	}

	r.config.GetClientCertificate = r.getClientCertificate
	r.config.InsecureSkipVerify = true
	r.config.VerifyConnection = r.verifyConnection

	return r.config
}

type resolver struct {
	server bool
	config *tls.Config
	source Source

	clientAuth tls.ClientAuthType
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

	if skip, err := r.skipVerify(cs); skip || err != nil {
		return err
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

	if r.server {
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

func (r *resolver) skipVerify(cs tls.ConnectionState) (bool, error) {
	// If we represent a client and the server has no certs then this is an error
	if !r.server && len(cs.PeerCertificates) == 0 {
		return false, errors.New("tls: no server certs")
	}

	if r.server {
		switch r.clientAuth {
		case tls.NoClientCert, tls.RequestClientCert:
			return true, nil
		case tls.RequireAnyClientCert:
			if len(cs.PeerCertificates) == 0 {
				return true, errors.New("tls: no client cert")
			}
			return true, nil
		case tls.VerifyClientCertIfGiven:
			return len(cs.PeerCertificates) == 0, nil
		case tls.RequireAndVerifyClientCert:
			if len(cs.PeerCertificates) == 0 {
				return false, errors.New("tls: no client cert")
			}
			return false, nil
		}
	}

	return false, nil
}

// TLSLeaf returns either cert.Leaf or parses the first of cert.Certificate.
// This does not set cert.Leaf.
func TLSLeaf(cert *tls.Certificate) (*x509.Certificate, error) {
	if cert.Leaf != nil {
		return cert.Leaf, nil
	}
	return x509.ParseCertificate(cert.Certificate[0])
}
