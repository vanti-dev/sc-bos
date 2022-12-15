package expire

import (
	"crypto/tls"
	"crypto/x509"
)

// Always is a pki.Expiry that always returns true.
func Always(cert *tls.Certificate, pool []*x509.Certificate, err error) bool {
	return true
}

// Never is a pki.Expiry that always returns false.
func Never(cert *tls.Certificate, pool []*x509.Certificate, err error) bool {
	return false
}
