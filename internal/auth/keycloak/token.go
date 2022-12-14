package keycloak

import (
	"context"
	"encoding/json"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/vanti-dev/sc-bos/internal/auth"
)

type accessTokenPayload struct {
	jwt.Claims
	Roles          []string                  `json:"roles"`
	Scopes         auth.JWTScopes            `json:"scope"`
	ResourceAccess map[string]resourceAccess `json:"resource_access"`
}

func (a *accessTokenPayload) AllRoles() []string {
	var roles []string
	roles = append(roles, a.Roles...)
	for _, res := range a.ResourceAccess {
		roles = append(roles, res.Roles...)
	}
	return roles
}

func (a *accessTokenPayload) IsAppOnly() bool {
	return false
}

type resourceAccess struct {
	Roles []string `json:"roles"`
}

func NewTokenVerifier(config *Config, keySet auth.KeySet) *TokenValidator {
	return &TokenValidator{
		keySet: keySet,
		expected: jwt.Expected{
			Audience: jwt.Audience{config.ClientID},
			Issuer:   config.Issuer(),
		},
	}
}

type TokenValidator struct {
	keySet   auth.KeySet
	expected jwt.Expected
}

func (v *TokenValidator) ValidateAccessToken(ctx context.Context, token string) (*auth.Authorization, error) {
	payloadBytes, err := v.keySet.VerifySignature(ctx, token)
	if err != nil {
		return nil, err
	}

	var payload accessTokenPayload
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return nil, err
	}

	err = payload.Claims.Validate(v.expected)
	if err != nil {
		return nil, err
	}

	return &auth.Authorization{
		Roles:     payload.AllRoles(),
		Scopes:    payload.Scopes,
		IsService: payload.IsAppOnly(),
	}, nil
}
