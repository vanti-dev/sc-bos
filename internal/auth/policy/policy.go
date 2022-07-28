package policy

import (
	"crypto/x509"
)

// Attributes is a collection of metadata which can be used by policies to decide whether to accept or reject
// a protected operation.
type Attributes struct {
	Service string `json:"service"` // gRPC service name, fully qualified
	Method  string `json:"method"`  // gRPC method name
	Request any    `json:"request"` // gRPC request message

	CertificateValid bool              `json:"certificate_valid"`
	Certificate      *x509.Certificate `json:"certificate"`
	TokenValid       bool              `json:"token_valid"`
	TokenClaims      any               `json:"token_claims"`
}
