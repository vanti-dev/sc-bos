package jwks

import (
	"context"
	"errors"
	"testing"

	"github.com/go-jose/go-jose/v4"
	"github.com/google/go-cmp/cmp"
)

func TestLocalKeySet_VerifySignature(t *testing.T) {
	inputPayload := []byte("TestLocalKeySet_VerifySignature")

	// sign a test message using the key we will use
	sig1 := signJWS(t, testJWK1, inputPayload)
	// sign again using the other key that's not in our key set
	sig2 := signJWS(t, testJWK2, inputPayload)

	// verify the first signature using the JWKS, which should succeed
	jwks := jose.JSONWebKeySet{Keys: []jose.JSONWebKey{testJWK1.Public()}}
	localKeySet := NewLocalKeySet(jwks, []jose.SignatureAlgorithm{jose.RS256})
	outputPayload, err := localKeySet.VerifySignature(context.Background(), sig1)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(inputPayload, outputPayload) {
		t.Error("payloads different")
	}

	// attempt to verify the second signature using the JWKS, which should fail as that key's not in the set
	_, err = localKeySet.VerifySignature(context.Background(), sig2)
	if !errors.Is(err, ErrKeyNotFound) {
		t.Errorf("verification didn't fail as expected: %s", err.Error())
	}
}
