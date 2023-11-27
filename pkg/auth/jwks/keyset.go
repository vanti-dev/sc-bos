// Package jwks provides both local and remote JWT Key Set access token verifiers.
package jwks

import (
	"context"
	"errors"
)

var (
	ErrTokenNotSigned          = errors.New("token is not signed")
	ErrTokenMultipleSignatures = errors.New("token has multiple signatures")
	ErrKeyNotFound             = errors.New("signing key not found in key set")
	ErrUpdateTooSoon           = errors.New("trying to update too soon since last update")
)

type KeySet interface {
	VerifySignature(ctx context.Context, jws string) (payload []byte, err error)
}
