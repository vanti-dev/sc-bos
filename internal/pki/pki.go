package pki

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"

	"go.uber.org/multierr"
)

// PublicKey contains the method set that the standard library public key types (from crypto/*) all implement.
// See the docs for crypto.PublicKey
type PublicKey interface {
	crypto.PublicKey
	Equal(key crypto.PublicKey) bool
}

// LoadOrGenerateKeyPair loads a keypair from a private key file. If the file is not present, a new keypair is generated
// and saved into the file.
// The file is PEM-encoded PKCS#8 file.
// Generated keys use 4096-bit RSA.
func LoadOrGenerateKeyPair(path string) (pkcs8PEM []byte, err error) {
	pkcs8PEM, err = os.ReadFile(path)
	// on success or unexpected error, we don't need to generate a new key
	if err == nil || !errors.Is(err, os.ErrNotExist) {
		return
	}

	log.Printf("generating new RSA private key in %q", path)

	// generate RSA key using the global CSPRNG
	private, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("generate private key: %w", err)
	}

	// encode PKCS#8 DER
	der, err := x509.MarshalPKCS8PrivateKey(private)
	if err != nil {
		return nil, fmt.Errorf("save private key: %w", err)
	}

	// wrap PKCS#8 in PEM
	pkcs8PEM = pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: der,
	})
	if pkcs8PEM == nil {
		return nil, errors.New("save private key: PEM encoding failed")
	}

	// write PEM to file
	err = writeFileNoTruncate(path, pkcs8PEM, 0600)
	if err != nil {
		return nil, fmt.Errorf("save private key: %w", err)
	}

	return
}

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

// works like os.WriteFile, but uses os.O_EXCL to produce an error if the file already exists
func writeFileNoTruncate(path string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	_, writeErr := f.Write(data)
	return multierr.Combine(writeErr, f.Close())
}
