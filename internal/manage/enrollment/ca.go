package enrollment

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
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
// This function returns only the leaf certificate - it can be passed to EncodeCertificateChain to obtain the full chain.
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

// CreateLocalCertificate creates and signs an X.509 certificate for this Smart Core node.
// It can be used to obtain a certificate for TLS on the node that manages the Smart Core network CA.
// If loopback is true, then DNS name 'localhost' and loopback IP addresses are added as SANs in the certificate.
// Otherwise, only non-loopback IP addresses are added as SANs.
func (ca *CA) CreateLocalCertificate(name pkix.Name, loopback bool, pub crypto.PublicKey) (der []byte, err error) {
	serial, err := pki.GenerateSerialNumber()
	if err != nil {
		return nil, err
	}

	ips, err := localInterfaceAddrs(loopback)
	var dnsNames []string
	if loopback {
		dnsNames = append(dnsNames, "localhost")
	}

	now := ca.now()
	template := &x509.Certificate{
		SerialNumber: serial,
		NotBefore:    now,
		NotAfter:     now.Add(ca.validity()),
		Subject:      name,
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		DNSNames:     dnsNames,
		IPAddresses:  ips,
	}
	return x509.CreateCertificate(rand.Reader, template, ca.Certificate, pub, ca.PrivateKey)
}

// LocalCertSource creates a pki.CertSource that uses CreateLocalCertificate to create and automatically rotate
// certificates.
func (ca *CA) LocalCertSource(name pkix.Name, loopback bool) (pki.CertSource, error) {
	return pki.NewCachedCertSource(func(old *tls.Certificate) (new *tls.Certificate, next time.Time, err error) {
		key, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return
		}

		next = ca.now().Add(ca.validity() / 2)
		leaf, err := ca.CreateLocalCertificate(name, loopback, key.Public())
		if err != nil {
			return
		}
		// connection establishment is faster is Certificate.Leaf is populated
		leafParsed, err := x509.ParseCertificate(leaf)
		if err != nil {
			return
		}

		chain := [][]byte{leaf, ca.Certificate.Raw}
		chain = append(chain, ca.Intermediates...)

		new = &tls.Certificate{
			Certificate: chain,
			PrivateKey:  key,
			Leaf:        leafParsed,
		}
		return
	})
}

// EncodeCertificateChain will encode a certificate chain in PEM format, containing the provided leaf certificate
// and all the CA's intermediate certificates.
// The certificate returned by ca.CreateEnrollmentCertificate is a suitable leaf certificate.
func (ca *CA) EncodeCertificateChain(leafDER []byte) (pem []byte) {
	chain := [][]byte{leafDER, ca.Certificate.Raw}
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

func localInterfaceAddrs(allowLoopback bool) ([]net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var ips []net.IP
	for _, iface := range ifaces {
		isLoopback := (iface.Flags & net.FlagLoopback) != 0
		if isLoopback && !allowLoopback {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			if addr, ok := addr.(*net.IPNet); ok {
				ips = append(ips, addr.IP)
			}
		}
	}
	return ips, nil
}
