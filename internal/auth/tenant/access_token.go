package tenant

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"github.com/go-jose/go-jose/v3/jwt"
	"io"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/vanti-dev/bsp-ew/internal/auth"
)

type tokenPayload struct {
	jwt.Claims
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
		Claims: jwt.Claims{
			Issuer:    ts.Issuer,
			Subject:   data.TenantID,
			Audience:  jwt.Audience{ts.Issuer},
			Expiry:    jwt.NewNumericDate(expires),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
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

	payloadBytes, err := jws.Verify(ts.Key.Key)
	if err != nil {
		return nil, err
	}

	var payload tokenPayload
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return nil, err
	}

	err = payload.Claims.Validate(jwt.Expected{
		Audience: jwt.Audience{ts.Issuer},
		Issuer:   ts.Issuer,
	})
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
		Key:       key,
	}, nil
}
