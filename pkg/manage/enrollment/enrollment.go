package enrollment

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vanti-dev/sc-bos/internal/util/pki"
)

type Enrollment struct {
	RootDeviceName string `json:"root_device_name"`
	ManagerName    string `json:"manager_name"`
	ManagerAddress string `json:"manager_address"`

	RootCA *x509.Certificate `json:"-"`
	Cert   tls.Certificate   `json:"-"`
}

func (e Enrollment) Equal(other Enrollment) bool {
	if e.RootDeviceName != other.RootDeviceName ||
		e.ManagerName != other.ManagerName ||
		e.ManagerAddress != other.ManagerAddress {
		return false
	}
	return string(e.Cert.Certificate[0]) == string(other.Cert.Certificate[0])
}

var ErrNotEnrolled = errors.New("node is not enrolled")

const (
	enrollmentFile = "enrollment.json"
	rootCaCertFile = "root-ca.cert.pem"
	certFile       = "enrollment.cert.pem"
)

// LoadEnrollment will load a previously saved Enrollment from a directory on disk. The directory should have
// the following structure:
//
//	<root>
//	  - enrollment.json - JSON-encoded Enrollment structure
//	  - root-ca.crt - Root CA for the enrollment, PEM-encoded X.509 certificate
//	  - cert.crt - Certificate chain for keyPEM, with the Root CA at the top of the chain
//
// This node's private key must be passed in, in PEM-wrapped PKCS#1 or PKCS#8 format.
func LoadEnrollment(dir string, keyPEM []byte) (Enrollment, error) {
	// read and unmarshal enrollment metadata
	var enrollment Enrollment
	raw, err := os.ReadFile(filepath.Join(dir, enrollmentFile))
	if errors.Is(err, os.ErrNotExist) {
		return Enrollment{}, ErrNotEnrolled
	} else if err != nil {
		return Enrollment{}, err
	}
	err = json.Unmarshal(raw, &enrollment)
	if err != nil {
		return Enrollment{}, err
	}

	// load root CA certificate
	rootCaPem, err := os.ReadFile(filepath.Join(dir, rootCaCertFile))
	if err != nil {
		return enrollment, err
	}
	rootCaBlock, _ := pem.Decode(rootCaPem)
	if rootCaBlock == nil {
		return enrollment, errors.New("CA.pem is not valid PEM data")
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
	certPem, err := os.ReadFile(filepath.Join(dir, certFile))
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

func SaveEnrollment(dir string, enrollment Enrollment) error {
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}

	jsonBytes, err := json.Marshal(enrollment)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(dir, enrollmentFile), jsonBytes, 0640)
	if err != nil {
		return err
	}

	_, err = pki.SaveCertificateChain(filepath.Join(dir, rootCaCertFile), [][]byte{enrollment.RootCA.Raw})
	if err != nil {
		return err
	}

	_, err = pki.SaveCertificateChain(filepath.Join(dir, certFile), enrollment.Cert.Certificate)
	return err
}

func DeleteEnrollment(dir string) error {
	return os.RemoveAll(dir)
}
