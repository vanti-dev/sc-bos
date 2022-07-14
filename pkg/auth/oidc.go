package auth

import (
	"context"
	"fmt"

	"github.com/vanti-dev/ew-auth-poc/pkg/fetch"
)

type OIDCConfig struct {
	JWKSURI               string `json:"jwks_uri"`
	TokenEndpoint         string `json:"token_endpoint"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
}

func DiscoverOIDCConfig(ctx context.Context, issuer string) (OIDCConfig, error) {
	url := fmt.Sprintf("%s/.well-known/openid-configuration", issuer)

	var config OIDCConfig
	err := fetch.JSON(ctx, url, &config)
	if err != nil {
		return OIDCConfig{}, err
	}

	return config, nil
}
