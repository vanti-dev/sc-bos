package token

import (
	"encoding/json"

	"github.com/vanti-dev/sc-bos/internal/auth/permission"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type Claims struct {
	SystemRoles []string               `json:"system_roles"` // The built-in system roles that this token is authorized for
	Zones       []string               `json:"zones"`        // The zones that this token is authorized for, for tenant tokens
	IsService   bool                   `json:"is_service"`   // True if the subject is an application acting on its own behalf, false if it's a user
	Permissions []PermissionAssignment `json:"permissions"`
}

type PermissionAssignment struct {
	Permission   permission.ID `json:"permission"`    // The name of the permission, e.g. trait:read:*
	Scoped       bool          `json:"scoped"`        // True if the permission is scoped to a specific resource
	ResourceType ResourceType  `json:"resource_type"` // The type of resource this permission is scoped to
	Resource     string        `json:"resource"`      // The resource this permission is scoped to - its meaning depends on the resource type
}

type ResourceType gen.RoleAssignment_ResourceType

func ParseResourceType(s string) (ResourceType, bool) {
	rt, ok := gen.RoleAssignment_ResourceType_value[s]
	if !ok {
		return 0, false
	}
	return ResourceType(rt), true
}

func (rt ResourceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(gen.RoleAssignment_ResourceType(rt))
}

func (rt ResourceType) String() string {
	return gen.RoleAssignment_ResourceType(rt).String()
}
