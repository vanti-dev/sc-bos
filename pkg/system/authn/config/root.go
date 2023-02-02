package config

import (
	"encoding/json"

	"github.com/vanti-dev/sc-bos/internal/auth/keycloak"
	"github.com/vanti-dev/sc-bos/pkg/system"
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
)

type Root struct {
	system.Config
	User   *User   `json:"user,omitempty"`
	System *System `json:"system,omitempty"`
}

type User struct {
	Validity *jsontypes.Duration `json:"validity,omitempty"`

	FileAccounts *Identities      `json:"fileAccounts,omitempty"`
	Keycloak     *keycloak.Config `json:"keycloakAccounts,omitempty"`
}

type System struct {
	Validity *jsontypes.Duration `json:"validity,omitempty"`

	// FileAccounts, when non-nil, causes the system to validate system tokens using a local
	// file of identities and secrets.
	// See Identities for how this field can be represented/configured in JSON.
	FileAccounts *Identities `json:"fileAccounts,omitempty"`
	// TenantAccounts causes the system to validate system tokens using the tenants system.
	// All tokens are deemed invalid if the tenants system is not available.
	TenantAccounts bool `json:"tenantAccounts,omitempty"`
	// CohortAccounts causes the system to validate system tokens using the cohort manager, setup via enrollment.
	// All tokens are deemed invalid if the manager is not known (i.e. the controller is not enrolled), or the manager doesn't support TenantApi.
	CohortAccounts bool `json:"cohortAccounts,omitempty"`
}

type Keycloak struct {
	URL      string `json:"url,omitempty"`
	Realm    string `json:"realm,omitempty"`
	ClientID string `json:"clientID,omitempty"`
}

func ReadConfig(data []byte) (Root, error) {
	root := Default()
	err := json.Unmarshal(data, &root)
	return root, err
}

func Default() Root {
	return Root{}
}
