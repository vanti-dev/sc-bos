package tenant

import (
	"testing"
	"time"
)

func TestTokenSource_createAndVerify(t *testing.T) {
	key, err := generateKey()
	if err != nil {
		t.Fatal(err)
	}
	ts := &TokenSource{
		Key:    key,
		Issuer: "test",
		Now:    time.Now,
	}

	token, err := ts.GenerateAccessToken(SecretData{TenantID: "Foo"}, 10*time.Minute)
	if err != nil {
		t.Fatalf("GenerateAccessToken %v", err)
	}

	_, err = ts.ValidateAccessToken(nil, token)
	if err != nil {
		t.Fatalf("ValidateAccessToken %v", err)
	}
}
