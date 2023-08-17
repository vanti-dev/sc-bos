package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
)

type Signer interface {
	Sign(ctx context.Context, data []byte) (signature []byte, err error)
}

type LocalSigner struct {
	Secret SASKey
}

func (s *LocalSigner) Sign(_ context.Context, data []byte) (signature []byte, err error) {
	mac := hmac.New(sha256.New, s.Secret)
	_, err = mac.Write(data)
	if err != nil {
		return signature, err
	}
	return mac.Sum(nil), nil
}
