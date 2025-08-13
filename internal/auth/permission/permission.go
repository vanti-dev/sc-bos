package permission

import (
	"slices"
	"strings"
)

type ID string

type Details struct {
	ID          ID     // Unique, stable identifier for the permission
	DisplayName string // Human-readable short name for the permission
	Description string // Human-readable longer description of the permission
}

const (
	TraitWrite ID = "trait:write"
	TraitRead  ID = "trait:read"
)

func GetDetails(id ID) (Details, bool) {
	idx, ok := slices.BinarySearchFunc(allPermissions, id, func(a Details, b ID) int {
		return strings.Compare(string(a.ID), string(b))
	})
	if !ok {
		return Details{}, false
	}
	return allPermissions[idx], true
}

func All() []Details {
	return slices.Clone(allPermissions)
}

var allPermissions = []Details{
	{
		ID:          TraitWrite,
		DisplayName: "Traits - Write",
		Description: "Access to read and update device trait data",
	},
	{
		ID:          TraitRead,
		DisplayName: "Traits - Read",
		Description: "Access to read device trait data",
	},
}

func init() {
	slices.SortFunc(allPermissions, func(a, b Details) int {
		return strings.Compare(string(a.ID), string(b.ID))
	})
}
