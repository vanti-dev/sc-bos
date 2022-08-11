package pki

import (
	"crypto"
	"crypto/tls"
	"log"
	"sync"
	"time"
)

type CertSource interface {
	TLSConfigGetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error)
}

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

	log.Println("generating self-signed TLS certificate")
	// we need to (re)generate the certificate
	validity := 30 * 24 * time.Hour
	refreshAt := time.Now().Add(validity / 2)
	certDER, err := CreateSelfSignedCert(cs.PrivateKey, validity)
	if err != nil {
		log.Printf("failed to (re)generate self-signed certificate: %s", err.Error())
		// keep using the old cert if possible
		if cs.cert != nil {
			return cs.cert, nil
		} else {
			return nil, err
		}
	}

	cert := &tls.Certificate{
		Certificate: [][]byte{certDER},
		PrivateKey:  cs.PrivateKey,
	}
	cs.cert = cert
	cs.refreshAt = refreshAt
	return cert, nil
}
