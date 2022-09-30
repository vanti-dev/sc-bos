package pki

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"io"
	"math/big"

	"go.uber.org/multierr"
)

// SaveCertificateChain writes derCerts as pem encoded CERTIFICATE blocks to certFile.
// Returns the written bytes.
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

// ParseCertificatesPEM parses and returns any CERTIFICATE blocks found in the pem encoded pemBytes.
func ParseCertificatesPEM(pemBytes []byte) (certs []*x509.Certificate, errs error) {
	for len(pemBytes) > 0 {
		var block *pem.Block
		block, pemBytes = pem.Decode(pemBytes)
		if block == nil {
			break
		}

		if block.Type != "CERTIFICATE" {
			return
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}

		certs = append(certs, cert)
	}
	return
}

func EncodePEMSequence(contents [][]byte, blockType string) (encoded []byte) {
	var buf bytes.Buffer

	for _, blockData := range contents {
		block := &pem.Block{
			Type:  blockType,
			Bytes: blockData,
		}

		err := pem.Encode(&buf, block)
		if err != nil {
			// buf.Write never returns an error, so pem.Encode won't either
			panic(err)
		}
	}

	return buf.Bytes()
}

func EncodeCertificates(certs []*x509.Certificate) []byte {
	out := &bytes.Buffer{}
	block := &pem.Block{Type: "CERTIFICATE"}
	for _, cert := range certs {
		if len(cert.Raw) == 0 {
			continue
		}
		block.Bytes = cert.Raw
		// writing to a memory buffer shouldn't return any errors
		err := pem.Encode(out, block)
		if err != nil {
			panic(err)
		}
	}
	return out.Bytes()
}

// GenerateSerialNumber generates a random unsigned 128-bit integer using a cryptographically secure source of random
// numbers.
// The returned integer is suitable for use as an X.509 certificate serial number.
func GenerateSerialNumber() (*big.Int, error) {
	blob := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, blob)
	if err != nil {
		return nil, err
	}
	return big.NewInt(0).SetBytes(blob), nil
}
