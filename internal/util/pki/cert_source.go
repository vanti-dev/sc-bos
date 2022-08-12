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
type CertSource interface {
	// TLSConfigGetCertificate can be used as the GetCertificate option in tls.Config
	// This will cause the TLS server to use this CertSource.
	TLSConfigGetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error)
	// RotateNow attempts an immediate certificate rotation attempt, even if the current certificate is still valid.
	// If no new certificate is available, this is not considered an error.
	RotateNow() error
}

// NewSelfSignedCertSource creates
func NewSelfSignedCertSource(key crypto.PrivateKey) (*SelfSignedCertSource, error) {
	var err error
	if key == nil {
		key, err = rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return nil, err
		}
	}

	certSource := &SelfSignedCertSource{
		PrivateKey: key,
	}
	// pre-populate the certificate
	err = certSource.RotateNow()
	if err != nil {
		return nil, err
	}
	return certSource, nil
}

// SelfSignedCertSource is a CertSource that issues self-signed certificates.
// The certificate is issued to localhost and all local network interface addresses.
// If key is nil, a new RSA private key is generated.
type SelfSignedCertSource struct {
	PrivateKey crypto.PrivateKey

	m         sync.Mutex
	cert      *tls.Certificate
	refreshAt time.Time
}

func (cs *SelfSignedCertSource) TLSConfigGetCertificate(_ *tls.ClientHelloInfo) (*tls.Certificate, error) {
	cs.m.Lock()
	defer cs.m.Unlock()

	if cs.cert != nil && time.Now().Before(cs.refreshAt) {
		// cert is present and still valid
		return cs.cert, nil
	}

	err := cs.rotateNow()
	if err != nil {
		log.Printf("failed to (re)generate self-signed certificate: %s", err.Error())
		// keep using the old cert if possible
		if cs.cert != nil {
			return cs.cert, nil
		} else {
			return nil, err
		}
	}

	return cs.cert, nil
}

func (cs *SelfSignedCertSource) RotateNow() error {
	cs.m.Lock()
	defer cs.m.Unlock()
	return cs.rotateNow()
}

func (cs *SelfSignedCertSource) rotateNow() error {
	log.Println("generating self-signed TLS certificate")
	// we need to (re)generate the certificate
	validity := 30 * 24 * time.Hour
	refreshAt := time.Now().Add(validity / 2)
	certDER, err := CreateSelfSignedCert(cs.PrivateKey, validity)
	if err != nil {
		return err
	}
	cert := &tls.Certificate{
		Certificate: [][]byte{certDER},
		PrivateKey:  cs.PrivateKey,
	}
	cs.cert = cert
	cs.refreshAt = refreshAt
	return nil
}

// NewFileCertSource creates a FileCertSource and loads the certificate.
func NewFileCertSource(certPath string, keyPath string) (*FileCertSource, error) {
	cs := &FileCertSource{
		CertificatePath: certPath,
		PrivateKeyPath:  keyPath,
	}
	err := cs.RotateNow()
	if err != nil {
		return nil, err
	}
	return cs, nil
}

// FileCertSource is a CertSource that loads certificates from disk.
// Certificates are automatically reloaded once the certificate expires.
type FileCertSource struct {
	PrivateKeyPath  string
	CertificatePath string

	m    sync.Mutex
	cert *tls.Certificate
}

func (cs *FileCertSource) TLSConfigGetCertificate(_ *tls.ClientHelloInfo) (*tls.Certificate, error) {
	cs.m.Lock()
	defer cs.m.Unlock()

	// the cached certificate is not expired yet
	if cs.cert != nil && time.Now().Before(cs.cert.Leaf.NotBefore) {
		return cs.cert, nil
	}

	err := cs.rotateNow()
	if err != nil {
		log.Printf("failed to load TLS cert: %s", err.Error())
		// on failure, return the old certificate if there is one
		if cs.cert != nil {
			return cs.cert, nil
		} else {
			return nil, err
		}
	}

	return cs.cert, nil
}

func (cs *FileCertSource) RotateNow() error {
	cs.m.Lock()
	defer cs.m.Unlock()
	return cs.rotateNow()
}

func (cs *FileCertSource) rotateNow() error {
	cert, err := tls.LoadX509KeyPair(cs.CertificatePath, cs.PrivateKeyPath)
	if err != nil {
		return err
	}
	// LoadX509KeyPair doesn't populate the Leaf property, but we will need it to check the expiry time next call
	// (populating it also improves performance of connection establishment)
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		// failing ParseCertificate is impossible because LoadX509KeyPair calls it too, which must have succeeded if
		// we reached here
		panic(err)
	}
	cs.cert = &cert
	return nil
}
