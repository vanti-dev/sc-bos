package auth

const (
	RoleTenant     = "tenant"
	RoleController = "controller"
	RoleUser       = "user"
)

func RequireAll(want []string, have []string) bool {
	unsatisfied := make(map[string]struct{}, len(want))
	for _, role := range want {
		unsatisfied[role] = struct{}{}
	}

	// mark off all the roles we have
	for _, role := range have {
		delete(unsatisfied, role)
	}

	// Roles are satisfied if none remain in the map
	return len(unsatisfied) == 0
}
