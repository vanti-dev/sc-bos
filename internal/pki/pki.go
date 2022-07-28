package pki

import (
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

// works like os.WriteFile, but uses os.O_EXCL to produce an error if the file already exists
func writeFileNoTruncate(path string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	_, writeErr := f.Write(data)
	return multierr.Combine(writeErr, f.Close())
}
