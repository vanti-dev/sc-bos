// Package account models and stores accounts and associated resources.
//
// Two kinds of accounts are supported: user accounts and service accounts.
// User accounts have a username and a password which can be used to authenticate.
// Service accounts have server-generated client secrets and are used for machine-to-machine communication.
//
// User accounts can have roles assigned to them. These assignments are represented as RoleAssignments.
// A RoleAssignment can optionally be scoped to specific resource(s).
//
// Roles are a named set of permissions.
//
// The resources can be persistently stored using the Store. The Server uses the data in the Store
// to implement the AccountApiService.
package account

import (
	"strconv"
)

func ParseAccountID(id string) (int64, bool) {
	accountID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, false
	}
	return accountID, true
}
