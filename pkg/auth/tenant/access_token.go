package tenant

import (
	"encoding/json"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/vanti-dev/ew-auth-poc/pkg/auth"
)

type AccessTokenPayload struct {
	Issuer     string       `json:"iss"`
	Subject    string       `json:"sub"`
	Audience   string       `json:"aud"`
	Expiration auth.JWTTime `json:"exp"`
	NotBefore  auth.JWTTime `json:"nbf"`
	IssuedAt   auth.JWTTime `json:"iat"`
}

type TokenSource struct {
	Key      jose.SigningKey
	Issuer   string
	Validity time.Duration
	Now      func() time.Time
}

func (ts *TokenSource) GenerateAccessToken(subject string, audience string) (token string, err error) {
	signer, err := jose.NewSigner(ts.Key, nil)
	if err != nil {
		return "", err
	}

	var now time.Time
	if ts.Now != nil {
		now = ts.Now()
	} else {
		now = time.Now()
	}

	expires := now.Add(ts.Validity)

	payload := AccessTokenPayload{
		Issuer:     ts.Issuer,
		Subject:    subject,
		Audience:   audience,
		Expiration: auth.JWTTime(expires),
		NotBefore:  auth.JWTTime(now),
		IssuedAt:   auth.JWTTime(now),
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	jws, err := signer.Sign(encoded)
	if err != nil {
		return "", err
	}

	return jws.CompactSerialize()
}
