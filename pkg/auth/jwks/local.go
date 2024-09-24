package jwks

import (
	"context"

	"github.com/go-jose/go-jose/v4"
)

type LocalKeySet struct {
	keys                         jose.JSONWebKeySet
	permittedSignatureAlgorithms []jose.SignatureAlgorithm
}

func NewLocalKeySet(keys jose.JSONWebKeySet, permittedSignatureAlgorithms []jose.SignatureAlgorithm) *LocalKeySet {
	return &LocalKeySet{
		keys:                         keys,
		permittedSignatureAlgorithms: permittedSignatureAlgorithms,
	}
}

func (ks *LocalKeySet) VerifySignature(_ context.Context, jws string) (payload []byte, err error) {
	sig, err := jose.ParseSigned(jws, ks.permittedSignatureAlgorithms)
	if err != nil {
		return nil, err
	}
	return verifyJWSWithKeySet(sig, ks.keys)
}
