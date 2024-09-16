package auth

import "github.com/go-jose/go-jose/v4"

type SignedToken interface {
	SetPermittedSignatureAlgorithms(permittedSignatureAlgorithms []jose.SignatureAlgorithm)
}

type TokenSink interface {
	GetSignedToken() SignedToken
}
