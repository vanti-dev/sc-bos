package keycloak

import (
	"fmt"

	"github.com/go-jose/go-jose/v4"

	"github.com/vanti-dev/sc-bos/internal/auth"
)

type Config struct {
	URL      string `json:"url,omitempty"` // Root URL of Keycloak server
	Realm    string `json:"realm,omitempty"`
	ClientID string `json:"clientId,omitempty"`
}

func (c *Config) Issuer() string {
	return fmt.Sprintf("%s/realms/%s", c.URL, c.Realm)
}

// DefaultPermittedSignatureAlgorithms
// TODO: reduce the number of permitted signature algorithms for all keycloak installations
// KeyCloak will select "a reasonable default" cipher suite if none is specified by the installation
var DefaultPermittedSignatureAlgorithms = []jose.SignatureAlgorithm{
	jose.RS256,
	jose.RS384,
	jose.RS512,
	jose.ES256,
	jose.ES384,
	jose.ES512,
	jose.PS256,
	jose.PS384,
	jose.PS512,
	jose.HS256,
}

func UsePermittedSignatureAlgorithms(permittedSignatureAlgorithms []jose.SignatureAlgorithm) func(ts auth.TokenSink) {
	return func(ts auth.TokenSink) {
		ts.GetSignedToken().SetPermittedSignatureAlgorithms(permittedSignatureAlgorithms)
	}
}
