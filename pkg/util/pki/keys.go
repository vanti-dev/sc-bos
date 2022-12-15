package pki

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/fs"
	"os"

	"go.uber.org/multierr"
	"go.uber.org/zap"
)

// PublicKey contains the method set that the standard library public key types (from crypto/*) all implement.
// See the docs for crypto.PublicKey
type PublicKey interface {
	crypto.PublicKey
	Equal(key crypto.PublicKey) bool
}

// PrivateKey contains the method set that the standard library public key types (from crypto/*) all implement.
// See the docs for crypto.PrivateKey
type PrivateKey interface {
	crypto.PrivateKey
	Public() crypto.PublicKey
	Equal(key crypto.PrivateKey) bool
}

var (
	ErrInvalidPEM       = errors.New("invalid PEM data")
	ErrUnknownBlockType = errors.New("unknown PEM block type")
)

// LoadPrivateKey will load a private key from a PEM file, in PKCS#8, PKCS#1 or EC format.
// Both the loaded key and the raw PEM read from the file are returned.
func LoadPrivateKey(keyFile string) (key PrivateKey, pemBytes []byte, err error) {
	pemBytes, err = os.ReadFile(keyFile)
	if err != nil {
		return
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, nil, ErrInvalidPEM
	}

	switch block.Type {
	case "PRIVATE KEY":
		var k any
		k, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		key = k.(PrivateKey)
	case "RSA PRIVATE KEY":
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	case "EC PRIVATE KEY":
		key, err = x509.ParseECPrivateKey(block.Bytes)
	default:
		err = ErrUnknownBlockType
	}
	return
}

func EncodePrivateKey(key crypto.PrivateKey) (pemBytes []byte, err error) {
	derBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return
	}

	pemBytes = pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derBytes,
	})
	if pemBytes == nil {
		// PEM encoding only fails if there are invalid headers, which cannot be the case here
		panic("PEM encode failed")
	}
	return
}

func SavePrivateKey(keyFile string, key crypto.PrivateKey) (pemBytes []byte, err error) {
	pemBytes, err = EncodePrivateKey(key)
	if err != nil {
		return
	}

	err = writeFileNoTruncate(keyFile, pemBytes, 0600)
	return
}

// LoadOrGeneratePrivateKey will load a private key in PEM-encoded PKCS#8, PKCS#1 or EC format from a file.
// If the file does not exist, a new 4096-bit RSA key in PKCS#8 format is generated and saved to the file.
func LoadOrGeneratePrivateKey(keyFile string, logger *zap.Logger) (key PrivateKey, pemBytes []byte, err error) {
	key, pemBytes, err = LoadPrivateKey(keyFile)
	if err == nil || !errors.Is(err, fs.ErrNotExist) {
		return
	}

	logger.Info("generating new RSA private key", zap.String("path", keyFile))

	key, err = rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return
	}

	pemBytes, err = SavePrivateKey(keyFile, key)
	return
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
