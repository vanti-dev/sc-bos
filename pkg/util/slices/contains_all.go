package slices

// ContainsAll returns true if haystack contains all the elements in needle, in any order.
// If both needle and haystack are empty this returns true.
func ContainsAll[S1, S2 ~[]E, E comparable](needle S1, haystack S2) bool {
	unsatisfied := make(map[E]struct{}, len(needle))
	for _, role := range needle {
		unsatisfied[role] = struct{}{}
	}

	// mark off all the roles we have
	for _, role := range haystack {
		delete(unsatisfied, role)
	}

	// Roles are satisfied if none remain in the map
	return len(unsatisfied) == 0
}
