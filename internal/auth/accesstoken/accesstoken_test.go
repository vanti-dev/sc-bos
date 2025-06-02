package accesstoken

import (
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
)

func TestTokenSource_createAndVerify(t *testing.T) {
	key, err := generateKey()
	if err != nil {
		t.Fatal(err)
	}
	ts := &Source{
		Key:    key,
		Issuer: "test",
		Now:    time.Now,
	}

	ts.SignatureAlgorithms = []string{string(jose.HS256)}

	token, err := ts.GenerateAccessToken(SecretData{TenantID: "Foo"}, 10*time.Minute)
	if err != nil {
		t.Fatalf("GenerateAccessToken %v", err)
	}

	_, err = ts.ValidateAccessToken(nil, token)
	if err != nil {
		t.Fatalf("ValidateAccessToken %v", err)
	}
}
