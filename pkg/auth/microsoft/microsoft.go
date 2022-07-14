package microsoft

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/vanti-dev/bsp-ew/pkg/auth"
)

type ClientConfig struct {
	Tenant   string
	ClientId string
}

func (c ClientConfig) Issuer() string {
	return fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", c.Tenant)
}

// TokenVerifier verifies Access Tokens issued by the Microsoft Identity Platform in its JWT version 2 format.
type TokenVerifier struct {
	keySet        auth.KeySet
	claimVerifier auth.JWTClaimVerifier
}

func NewTokenVerifier(keySet auth.KeySet, config *ClientConfig, options ...Option) *TokenVerifier {
	o := resolveOpts(options...)

	return &TokenVerifier{
		keySet: keySet,
		claimVerifier: auth.JWTClaimVerifier{
			Now:      o.now,
			Audience: config.ClientId,
		},
	}
}

func (v *TokenVerifier) VerifyAccessToken(ctx context.Context, token string) (*auth.Authorization, error) {
	payload, err := v.keySet.VerifySignature(ctx, token)
	if err != nil {
		return nil, err
	}

	var decoded AccessTokenPayload
	err = json.Unmarshal(payload, &decoded)
	if err != nil {
		return nil, err
	}

	err = v.claimVerifier.Verify(decoded.JWTCommonClaims)
	if err != nil {
		return nil, err
	}

	return &auth.Authorization{
		Issuer:    decoded.Issuer,
		Subject:   decoded.Subject,
		Roles:     decoded.Roles,
		Scopes:    decoded.Scp,
		IsService: decoded.IsAppOnly(),
	}, nil
}

type AccessTokenPayload struct {
	auth.JWTCommonClaims

	// Basic Claims
	Idp               string         `json:"idp"`
	Azp               string         `json:"azp"`
	AzpAcr            string         `json:"azpacr"`
	PreferredUsername string         `json:"preferred_username"`
	Name              string         `json:"name"`
	Scp               auth.JWTScopes `json:"scp"` // space-separated list of scope names
	Roles             []string       `json:"roles"`
	Wids              []string       `json:"wids"`
	Groups            []string       `json:"groups"`
	HasGroups         bool           `json:"hasgroups"`
	OId               string         `json:"oid"`
	TId               string         `json:"tid"`
	Uti               string         `json:"uti"`
	Version           string         `json:"version"`

	// Optional claims
	IdTyp string `json:"idtyp"`
}

// IsAppOnly returns true if the token was issued to an app operating under its own identity, rather than on behalf
// of a user account.
func (p *AccessTokenPayload) IsAppOnly() bool {
	return p.IdTyp == "app"
}
