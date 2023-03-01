package slices

// Contains returns true if haystack contains needle.
func Contains[S1 ~[]E, E comparable](needle E, haystack S1) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
