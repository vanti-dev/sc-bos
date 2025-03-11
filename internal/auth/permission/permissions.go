package permission

import (
	"strings"

	"golang.org/x/exp/slices"
)

type ID string

const (
	AccountRead       ID = "account:read"
	AccountWrite      ID = "account:write"
	AccountSelfRead   ID = "account:self:read"
	AccountSelfWrite  ID = "account:self:write"
	TraitsRead        ID = "traits:read"
	TraitsWrite       ID = "traits:write"
	TraitsHistoryRead ID = "traits:history:read"
)

type Permission struct {
	ID           ID
	DisplayName  string
	Description  string
	InheritsFrom []ID
}

var allPermissions = []Permission{
	{
		ID:           AccountRead,
		DisplayName:  "Account - Read All",
		Description:  "Read-only access to all account details",
		InheritsFrom: []ID{AccountSelfRead},
	},
	{
		ID:           AccountWrite,
		DisplayName:  "Account - Write All",
		Description:  "Manage all accounts; Create and assign roles",
		InheritsFrom: []ID{AccountRead},
	},
	{
		ID:          AccountSelfRead,
		DisplayName: "Account - Read Own Account",
		Description: "Read details of your own account",
	},
	{
		ID:           AccountSelfWrite,
		DisplayName:  "Account - Write Own Account",
		Description:  "Update your own account details, including credentials",
		InheritsFrom: []ID{AccountSelfRead},
	},
	{
		ID:          TraitsRead,
		DisplayName: "Traits - Read",
		Description: "Read all live data from devices",
	},
	{
		ID:           TraitsWrite,
		DisplayName:  "Traits - Write",
		Description:  "Send commands to devices and update device state",
		InheritsFrom: []ID{TraitsRead},
	},
	{
		ID:          TraitsHistoryRead,
		DisplayName: "Traits - Read History",
		Description: "Read all historical data from devices",
	},
}

func init() {
	sortPermissions(allPermissions)
}

// ByID retrieves a single permission from the registry.
func ByID(id ID) (Permission, bool) {
	idx, ok := findID(id)
	if !ok {
		return Permission{}, false
	}
	return allPermissions[idx], true
}

// Range returns a subset of the available permissions.
// Returns permissions sorted by ID. Only Permissions whose ID is lexically after the provided ID are included.
// The number of returned Permissions is limited by limit.
// If after is the empty string, returns Permissions starting from the (lexically) first ID.
// If limit is negative, returns all following Permissions.
//
// Return value more indicates if the returned slice was truncated to the limit.
func Range(after ID, limit int) (permissions []Permission, more bool) {
	idx := 0
	if after != "" {
		idx, _ = findID(after)
		idx += 1
	}
	if idx > len(allPermissions) {
		return nil, false
	}
	selected := allPermissions[idx:]
	if limit >= 0 && len(selected) > limit {
		selected = selected[:limit]
		more = true
	}
	return slices.Clone(selected), more
}

// All returns a list of all valid permissions, sorted by ID.
func All() []Permission {
	all, _ := Range("", -1)
	return all
}

func Count() int {
	return len(allPermissions)
}

func comparePermissions(a, b Permission) int {
	return strings.Compare(string(a.ID), string(b.ID))
}

func sortPermissions(s []Permission) {
	slices.SortFunc(s, comparePermissions)
}

// finds the index in allPermissions where the permission with the given ID would be.
// Returns true if that permission actually exists.
func findID(id ID) (int, bool) {
	return slices.BinarySearchFunc(allPermissions, id, func(permission Permission, id ID) int {
		return strings.Compare(string(permission.ID), string(id))
	})
}
