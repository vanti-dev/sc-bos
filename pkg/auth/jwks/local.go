package jwks

import (
	"context"

	"github.com/go-jose/go-jose/v3"
)

type LocalKeySet struct {
	keys jose.JSONWebKeySet
}

func NewLocalKeySet(keys jose.JSONWebKeySet) *LocalKeySet {
	return &LocalKeySet{keys: keys}
}

func (ks *LocalKeySet) VerifySignature(_ context.Context, jws string) (payload []byte, err error) {
	sig, err := jose.ParseSigned(jws)
	if err != nil {
		return nil, err
	}
	return verifyJWSWithKeySet(sig, ks.keys)
}
