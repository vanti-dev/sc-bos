package tenant

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"io"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/vanti-dev/bsp-ew/internal/auth"
)

type tokenPayload struct {
	auth.JWTCommonClaims
	Zones []string `json:"zones"`
}

type TokenSource struct {
	Key      jose.SigningKey
	Issuer   string
	Validity time.Duration
	Now      func() time.Time
}

func (ts *TokenSource) GenerateAccessToken(data SecretData) (token string, err error) {
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

	payload := tokenPayload{
		JWTCommonClaims: auth.JWTCommonClaims{
			Issuer:     ts.Issuer,
			Subject:    data.TenantID,
			Audience:   ts.Issuer,
			Expiration: auth.JWTTime(expires),
			NotBefore:  auth.JWTTime(now),
			IssuedAt:   auth.JWTTime(now),
		},
		Zones: data.Zones,
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

func (ts *TokenSource) ValidateAccessToken(_ context.Context, token string) (*auth.Authorization, error) {
	jws, err := jose.ParseSigned(token)
	if err != nil {
		return nil, err
	}

	payloadBytes, err := jws.Verify(ts.Key)
	if err != nil {
		return nil, err
	}

	var payload tokenPayload
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return nil, err
	}

	err = (&auth.JWTClaimValidator{
		Audience: ts.Issuer,
		Issuer:   ts.Issuer,
	}).ValidateClaims(payload.JWTCommonClaims)
	if err != nil {
		return nil, err
	}

	return &auth.Authorization{
		Roles:     []string{auth.RoleTenant},
		Zones:     payload.Zones,
		IsService: true,
	}, nil
}

func generateKey() (jose.SigningKey, error) {
	key := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return jose.SigningKey{}, err
	}

	return jose.SigningKey{
		Algorithm: jose.HS256,
		Key:       nil,
	}, nil
}
