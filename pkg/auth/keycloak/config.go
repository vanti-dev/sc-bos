package keycloak

import (
	"fmt"
)

type Config struct {
	URL      string // Root URL of Keycloak server
	Realm    string
	ClientID string
}

func (c *Config) Issuer() string {
	return fmt.Sprintf("%s/realms/%s", c.URL, c.Realm)
}
