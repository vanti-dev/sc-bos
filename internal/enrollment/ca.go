package enrollment

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"net"
	"net/url"
	"time"

	"github.com/vanti-dev/bsp-ew/internal/util/pki"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
)

const DefaultValidity = 30 * 24 * time.Hour

type CA struct {
	Certificate   *x509.Certificate
	PrivateKey    crypto.PrivateKey
	Intermediates [][]byte // Intermediate certificates between Certificate and the root. DER-encoded X.509

	Now      func() time.Time
	Validity time.Duration
}

// CreateEnrollmentCertificate creates and signs an X.509 certificate for a Smart Core node.
// Returns just the new certificate, in DER encoding.
// For the equivalent function that returns the entire certificate chain, see CreateEnrollmentCertificateChain.
func (ca *CA) CreateEnrollmentCertificate(enrollment *gen.Enrollment, pub crypto.PublicKey) (der []byte, err error) {
	serial, err := pki.GenerateSerialNumber()
	if err != nil {
		return nil, err
	}

	var (
		dnsNames    []string
		ipAddresses []net.IP
	)
	hostOrIP, _, err := net.SplitHostPort(enrollment.TargetAddress)
	if err != nil {
		return nil, err
	}
	if ip := net.ParseIP(hostOrIP); ip != nil {
		ipAddresses = append(ipAddresses, ip)
	} else {
		dnsNames = append(dnsNames, hostOrIP)
	}

	now := ca.now()
	template := &x509.Certificate{
		SerialNumber: serial,
		NotBefore:    now,
		NotAfter:     now.Add(ca.validity()),
		Subject:      pkix.Name{CommonName: enrollment.TargetName},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		DNSNames:     dnsNames,
		IPAddresses:  ipAddresses,
		URIs: []*url.URL{{
			Scheme: "smart-core",
			Opaque: enrollment.TargetName,
		}},
	}

	return x509.CreateCertificate(rand.Reader, template, ca.Certificate, pub, ca.PrivateKey)
}

// EncodeCertificateChain will encode a certificate chain in PEM format, containing the provided leaf certificate
// and all the CA's intermediate certificates.
// The certificate returned by ca.CreateEnrollmentCertificate is a suitable leaf certificate.
func (ca *CA) EncodeCertificateChain(leafDER []byte) (pem []byte) {
	chain := [][]byte{leafDER}
	chain = append(chain, ca.Intermediates...)
	return pki.EncodePEMSequence(chain, "CERTIFICATE")
}

func (ca *CA) now() time.Time {
	if now := ca.Now; now != nil {
		return now()
	}
	return time.Now()
}

func (ca *CA) validity() time.Duration {
	validity := ca.Validity
	if validity == 0 {
		return DefaultValidity
	}
	return validity
}
