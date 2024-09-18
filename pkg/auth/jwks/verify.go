package jwks

import (
	"github.com/go-jose/go-jose/v4"
	"go.uber.org/multierr"
)

// verifies that a JWS is signed by a key from a key set.
// The key ID recorded in the JWS must match one of the keys in the key set.
func verifyJWSWithKeySet(jws *jose.JSONWebSignature, keySet jose.JSONWebKeySet) (payload []byte, err error) {
	// check there is exactly one signature
	switch len(jws.Signatures) {
	case 0:
		return nil, ErrTokenNotSigned
	case 1: // this is fine
	default:
		return nil, ErrTokenMultipleSignatures
	}

	var errs error

	// search through all keys matching the key ID
	keyId := jws.Signatures[0].Header.KeyID
	for _, key := range keySet.Key(keyId) {
		payload, err = jws.Verify(key)
		if err == nil {
			return payload, nil
		}
		// record the reason this key failed to verify
		errs = multierr.Append(errs, err)
	}

	if errs != nil {
		// some key ID(s) matched, but failed to verify
		return nil, errs
	} else {
		// no keys matched
		return nil, ErrKeyNotFound
	}
}
