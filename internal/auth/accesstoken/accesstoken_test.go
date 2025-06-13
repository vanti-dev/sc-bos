package accesstoken

import (
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
)

func TestTokenSource_createAndVerify(t *testing.T) {
	ts := newTestSource(t)
	token, err := ts.GenerateAccessToken(SecretData{TenantID: "Foo"}, 10*time.Minute)
	if err != nil {
		t.Fatalf("GenerateAccessToken %v", err)
	}

	_, err = ts.ValidateAccessToken(nil, token)
	if err != nil {
		t.Fatalf("ValidateAccessToken %v", err)
	}
}

func newTestSource(t *testing.T) *Source {
	t.Helper()
	key, err := generateKey()
	if err != nil {
		t.Fatal(err)
	}
	return &Source{
		Key:                 key,
		Issuer:              "test",
		Now:                 time.Now,
		SignatureAlgorithms: []string{string(jose.HS256)},
	}
}
