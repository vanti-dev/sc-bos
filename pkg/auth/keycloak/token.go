package keycloak

import (
	"context"
	"encoding/json"

	"github.com/vanti-dev/ew-auth-poc/pkg/auth"
)

type AccessTokenPayload struct {
	auth.JWTCommonClaims
	Roles          []string                  `json:"roles"`
	Scopes         auth.JWTScopes            `json:"scope"`
	ResourceAccess map[string]ResourceAccess `json:"resource_access"`
}

func (a *AccessTokenPayload) AllRoles() []string {
	var roles []string
	roles = append(roles, a.Roles...)
	for _, resourceAccess := range a.ResourceAccess {
		roles = append(roles, resourceAccess.Roles...)
	}
	return roles
}

func (a *AccessTokenPayload) IsAppOnly() bool {
	return false
}

type ResourceAccess struct {
	Roles []string `json:"roles"`
}

func NewTokenVerifier(config *Config, keySet auth.KeySet) *TokenVerifier {
	return &TokenVerifier{
		keySet: keySet,
		claimVerifier: &auth.JWTClaimVerifier{
			Audience: config.ClientID,
			Issuer:   config.Issuer(),
		},
	}
}

type TokenVerifier struct {
	keySet        auth.KeySet
	claimVerifier *auth.JWTClaimVerifier
}

func (v *TokenVerifier) VerifyAccessToken(ctx context.Context, token string) (*auth.Authorization, error) {
	payloadBytes, err := v.keySet.VerifySignature(ctx, token)
	if err != nil {
		return nil, err
	}

	var payload AccessTokenPayload
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return nil, err
	}

	err = v.claimVerifier.Verify(payload.JWTCommonClaims)
	if err != nil {
		return nil, err
	}

	return &auth.Authorization{
		Issuer:    payload.Issuer,
		Subject:   payload.Subject,
		Roles:     payload.AllRoles(),
		Scopes:    payload.Scopes,
		IsService: payload.IsAppOnly(),
	}, nil
}
