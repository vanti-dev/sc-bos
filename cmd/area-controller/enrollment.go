package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Enrollment struct {
	RootDeviceName string `json:"root_device_name"`
	ManagerName    string `json:"manager_name"`
	ManagerAddress string `json:"manager_address"`

	RootCA *x509.Certificate `json:"-"`
	Cert   tls.Certificate   `json:"-"`
}

var ErrNotEnrolled = errors.New("node is not enrolled")

// LoadEnrollment will load a previously saved Enrollment from a directory on disk. The directory should have
// the following structure:
//
//     <root>
//       - enrollment.json - JSON-encoded Enrollment structure
//       - root-ca.crt - Root CA for the enrollment, PEM-encoded X.509 certificate
//       - cert.crt - Certificate chain for keyPEM, with the Root CA at the top of the chain
//
// This node's private key must be passed in, in PEM-wrapped PKCS#1 or PKCS#8 format.
func LoadEnrollment(dir string, keyPEM []byte) (Enrollment, error) {
	// check that enrollment dir exists
	_, err := os.Stat(dir)
	if errors.Is(err, os.ErrNotExist) {
		return Enrollment{}, ErrNotEnrolled
	} else if err != nil {
		return Enrollment{}, err
	}

	// read and unmarshal enrollment metadata
	var enrollment Enrollment
	raw, err := os.ReadFile(filepath.Join(dir, "enrollment.json"))
	if err != nil {
		return Enrollment{}, err
	}
	err = json.Unmarshal(raw, &enrollment)
	if err != nil {
		return Enrollment{}, err
	}

	// load root CA certificate
	rootCaPem, err := os.ReadFile(filepath.Join(dir, "root-ca.crt"))
	if err != nil {
		return enrollment, err
	}
	rootCaBlock, _ := pem.Decode(rootCaPem)
	if rootCaBlock == nil {
		return enrollment, errors.New("ca.pem is not valid PEM data")
	}
	if rootCaBlock.Type != "CERTIFICATE" {
		return enrollment, fmt.Errorf("expected PEM type 'CERTIFICATE' but got %q", rootCaBlock.Type)
	}
	rootCa, err := x509.ParseCertificate(rootCaBlock.Bytes)
	if err != nil {
		return enrollment, err
	}
	enrollment.RootCA = rootCa

	// load this node's certificate
	certPem, err := os.ReadFile(filepath.Join(dir, "cert.crt"))
	if err != nil {
		return enrollment, err
	}
	cert, err := tls.X509KeyPair(certPem, keyPEM)
	if err != nil {
		return enrollment, err
	}
	enrollment.Cert = cert

	return enrollment, nil
}
