package keycloak

import (
	"fmt"
)

type Config struct {
	URL      string `json:"url,omitempty"` // Root URL of Keycloak server
	Realm    string `json:"realm,omitempty"`
	ClientID string `json:"clientId,omitempty"`
}

func (c *Config) Issuer() string {
	return fmt.Sprintf("%s/realms/%s", c.URL, c.Realm)
}
