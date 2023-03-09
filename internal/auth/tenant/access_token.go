package tenant

import (
	"context"
	"crypto/rand"
	"io"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"

	"github.com/vanti-dev/sc-bos/pkg/auth/role"
	"github.com/vanti-dev/sc-bos/pkg/auth/token"
)

type tokenClaims struct {
	Zones []string `json:"zones,omitempty"`
}

type TokenSource struct {
	Key    jose.SigningKey
	Issuer string
	Now    func() time.Time
}

func (ts *TokenSource) GenerateAccessToken(data SecretData, validity time.Duration) (token string, err error) {
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

	expires := now.Add(validity)

	jwtClaims := jwt.Claims{
		Issuer:    ts.Issuer,
		Subject:   data.TenantID,
		Audience:  jwt.Audience{ts.Issuer},
		Expiry:    jwt.NewNumericDate(expires),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	customClaims := tokenClaims{Zones: data.Zones}
	return jwt.Signed(signer).
		Claims(jwtClaims).
		Claims(customClaims).
		CompactSerialize()
}

func (ts *TokenSource) ValidateAccessToken(_ context.Context, tokenStr string) (*token.Claims, error) {
	tok, err := jwt.ParseSigned(tokenStr)
	if err != nil {
		return nil, err
	}
	var jwtClaims jwt.Claims
	var customClaims tokenClaims
	err = tok.Claims(ts.Key.Key, &jwtClaims, &customClaims)
	if err != nil {
		return nil, err
	}
	err = jwtClaims.Validate(jwt.Expected{
		Audience: jwt.Audience{ts.Issuer},
		Issuer:   ts.Issuer,
	})
	if err != nil {
		return nil, err
	}
	return &token.Claims{
		Roles:     []string{role.Tenant},
		Zones:     customClaims.Zones,
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
