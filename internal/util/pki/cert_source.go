package pki

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"log"
	"sync"
	"time"
)

// A CertSource provides certificates for TLS.
// Using a CertSource instead of a plain tls.Certificate allows a TLS server to rotate certificates without
// restarting or dropping existing connections.
type CertSource interface {
	// TLSConfigGetCertificate can be used as the GetCertificate option in tls.Config
	// This will cause the TLS server to use this CertSource.
	TLSConfigGetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error)
	TLSConfigGetClientCertificate(info *tls.CertificateRequestInfo) (*tls.Certificate, error)
	// RotateNow attempts an immediate certificate rotation attempt, even if the current certificate is still valid.
	// If no new certificate is available, this is not considered an error.
	RotateNow() error
}

// NewSelfSignedCertSource creates a CertSource that issues self-signed certificates using CreateSelfSignedCert.
// The certificate is rotated lazily once half its lifetime has expired.
// If key is nil, a new RSA private key is generated.
func NewSelfSignedCertSource(key crypto.PrivateKey) (CertSource, error) {
	var err error
	if key == nil {
		key, err = rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return nil, err
		}
	}

	return NewCertSource(func(old *tls.Certificate) (new *tls.Certificate, next time.Time, err error) {
		log.Println("generating self-signed TLS certificate")
		// we need to (re)generate the certificate
		validity := 30 * 24 * time.Hour
		next = time.Now().Add(validity / 2)
		certDER, err := CreateSelfSignedCert(key, validity)
		if err != nil {
			return
		}
		new = &tls.Certificate{
			Certificate: [][]byte{certDER},
			PrivateKey:  key,
		}
		return
	})
}

// NewFileCertSource creates a CertSource that loads certificates from disk.
// Certificates are automatically reloaded once the certificate expires.
// The CertSource operates lazily - the files are only reloaded when the in-memory certificate expires, or
// ReloadNow is called.
func NewFileCertSource(certPath string, keyPath string) (CertSource, error) {
	return NewCertSource(func(old *tls.Certificate) (new *tls.Certificate, next time.Time, err error) {
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			return
		}

		// LoadX509KeyPair doesn't populate the Leaf property, but we will need it to check the expiry time next call
		// (populating it also improves performance of connection establishment)
		cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
		if err != nil {
			// failing ParseCertificate is impossible because LoadX509KeyPair calls it too, which must have succeeded if
			// we reached here
			panic(err)
		}
		// try reloading the certificate once the in-memory copy is one hour from expiry
		return &cert, cert.Leaf.NotAfter.Add(-time.Hour), nil
	})
}

type CertRotation func(old *tls.Certificate) (new *tls.Certificate, next time.Time, err error)

type certSource struct {
	m      sync.Mutex
	cert   *tls.Certificate
	next   time.Time
	rotate CertRotation
}

func NewCertSource(rotate CertRotation) (CertSource, error) {
	cert, next, err := rotate(nil)
	if err != nil {
		return nil, err
	}
	return &certSource{
		cert:   cert,
		next:   next,
		rotate: rotate,
	}, nil
}

func (c *certSource) TLSConfigGetCertificate(_ *tls.ClientHelloInfo) (*tls.Certificate, error) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.next.Before(time.Now()) {
		cert, next, err := c.rotate(c.cert)
		if err != nil {
			return nil, err
		}
		c.next = next
		c.cert = cert
	}

	return c.cert, nil
}

func (c *certSource) TLSConfigGetClientCertificate(info *tls.CertificateRequestInfo) (*tls.Certificate, error) {
	return c.TLSConfigGetCertificate(nil)
}

func (c *certSource) RotateNow() error {
	c.m.Lock()
	defer c.m.Unlock()

	cert, next, err := c.rotate(c.cert)
	if err != nil {
		return err
	}
	c.next = next
	c.cert = cert
	return nil
}
