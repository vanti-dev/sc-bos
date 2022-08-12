package pki

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"

	"go.uber.org/multierr"
)

func GetCertificateSmartCoreNames(cert *x509.Certificate) []string {
	var names []string
	for _, uri := range cert.URIs {
		if uri.Scheme == "smart-core" {
			names = append(names, uri.Opaque)
		}
	}
	return names
}

func ParseCertificateChainPEM(pemBytes []byte) (leaf *x509.Certificate, intermediates *x509.CertPool, err error) {
	block, pemBytes := pem.Decode(pemBytes)
	if block == nil {
		return nil, nil, errors.New("invalid leaf certificate PEM")
	}
	if block.Type != "CERTIFICATE" {
		return nil, nil, fmt.Errorf("expected CERTIFICATE block, found %q", block.Type)
	}
	leaf, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		return
	}

	intermediates = x509.NewCertPool()
	ok := intermediates.AppendCertsFromPEM(pemBytes)
	if !ok {
		return nil, nil, errors.New("failed to parse intermediate certificates")
	}

	return
}

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
