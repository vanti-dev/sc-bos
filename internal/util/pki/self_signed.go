package pki

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"log"
	"net"
	"time"
)

var ErrInvalidPrivateKey = errors.New("invalid private key")

// CreateSelfSignedCert generates a self-signed DER-encoded X.509 certificate suitable for TLS server authentication.
// It is issued with a CN of localhost and SANs of localhost, plus all local interface IP addresses.
func CreateSelfSignedCert(key crypto.PrivateKey, validity time.Duration) (derBytes []byte, err error) {
	private, ok := key.(PrivateKey)
	if !ok {
		return nil, ErrInvalidPrivateKey
	}
	public := private.Public()

	interfaceIPs, err := localInterfaceAddresses()
	if err != nil {
		log.Printf("can't get network interface IP addresses; certificate will be signed without any: %s", err.Error())
	}

	serial, err := GenerateSerialNumber()
	if err != nil {
		return nil, err
	}
	now := time.Now()

	template := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: "localhost"},
		DNSNames:              []string{"localhost"},
		IPAddresses:           interfaceIPs,
		NotBefore:             now.Add(-time.Hour), // a bit of leeway for clock skew
		NotAfter:              now.Add(validity),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err = x509.CreateCertificate(rand.Reader, template, template, public, private)
	return
}

func SaveCertificateChain(certFile string, derCerts [][]byte) (pemBytes []byte, err error) {
	var buf bytes.Buffer

	for _, derCert := range derCerts {
		err = pem.Encode(&buf, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: derCert,
		})
		if err != nil {
			return
		}
	}

	pemBytes = buf.Bytes()
	err = writeFileNoTruncate(certFile, pemBytes, 0640)
	return
}

func localInterfaceAddresses() ([]net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	var ips []net.IP
	for _, addr := range addrs {
		if addr, ok := addr.(*net.IPNet); ok {
			ips = append(ips, addr.IP)
		}
	}
	return ips, nil
}
