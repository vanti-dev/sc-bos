package accesstoken

import (
	"context"
	"crypto/rand"
	"errors"
	"io"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"

	jose_utils "github.com/vanti-dev/sc-bos/internal/util/jose"
	"github.com/vanti-dev/sc-bos/pkg/auth/token"
)

var ErrUnsupportedTokenVersion = errors.New("unsupported token version")

const claimsVersion = 1

type claims struct {
	Version     int                          `json:"v"` // to detect which schema version this token uses
	Name        string                       `json:"name,omitempty"`
	SystemRoles []string                     `json:"roles,omitempty"` // Named roles in JSON for back-compat
	Permissions []token.PermissionAssignment `json:"permissions,omitempty"`
}

type Source struct {
	Key                 jose.SigningKey
	Issuer              string
	Now                 func() time.Time
	SignatureAlgorithms []string
}

func (ts *Source) GenerateAccessToken(data SecretData, validity time.Duration) (token string, err error) {
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
	customClaims := claims{
		Version:     claimsVersion,
		Name:        data.Title,
		Permissions: data.Permissions,
		SystemRoles: data.SystemRoles,
	}
	return jwt.Signed(signer).
		Claims(jwtClaims).
		Claims(customClaims).
		Serialize()
}

func (ts *Source) ValidateAccessToken(_ context.Context, tokenStr string) (*token.Claims, error) {
	tok, err := jwt.ParseSigned(tokenStr, jose_utils.ConvertToNativeJose(ts.SignatureAlgorithms))
	if err != nil {
		return nil, err
	}
	var jwtClaims jwt.Claims
	var customClaims claims
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
	if customClaims.Version != claimsVersion {
		// token issued using a schema we no longer support
		return nil, ErrUnsupportedTokenVersion
	}
	return &token.Claims{
		SystemRoles: customClaims.SystemRoles,
		IsService:   true,
		Permissions: customClaims.Permissions,
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
