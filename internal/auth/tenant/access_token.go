package tenant

import (
	"context"
	"crypto/rand"
	"io"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"

	jose_utils "github.com/vanti-dev/sc-bos/internal/util/jose"
	"github.com/vanti-dev/sc-bos/pkg/auth/token"
)

type tokenClaims struct {
	Name  string   `json:"name,omitempty"`
	Zones []string `json:"zones,omitempty"`
	Roles []string `json:"roles,omitempty"`
}

type TokenSource struct {
	Key                 jose.SigningKey
	Issuer              string
	Now                 func() time.Time
	SignatureAlgorithms []string
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
	customClaims := tokenClaims{Name: data.Title, Zones: data.Zones, Roles: data.Roles}
	return jwt.Signed(signer).
		Claims(jwtClaims).
		Claims(customClaims).
		Serialize()
}

func (ts *TokenSource) ValidateAccessToken(_ context.Context, tokenStr string) (*token.Claims, error) {
	tok, err := jwt.ParseSigned(tokenStr, jose_utils.ConvertToNativeJose(ts.SignatureAlgorithms))
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
		AnyAudience: jwt.Audience{ts.Issuer},
		Issuer:      ts.Issuer,
	})
	if err != nil {
		return nil, err
	}
	return &token.Claims{
		Roles:     customClaims.Roles,
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
