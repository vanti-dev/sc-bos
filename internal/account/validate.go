package account

import (
	"regexp"
	"strings"
)

const (
	minPasswordLength    = 10
	maxPasswordLength    = 72 // max password length supported by bcrypt
	minUsernameLength    = 3
	maxUsernameLength    = 100
	minDisplayNameLength = 1
	maxDisplayNameLength = 100
	minDescriptionLength = 0
	maxDescriptionLength = 10000
)

func permitPassword(password string) bool {
	password = normalisePassword(password)
	return len(password) >= minPasswordLength && len(password) <= maxPasswordLength
}

func normalisePassword(password string) string {
	return strings.TrimSpace(password)
}

func validateDisplayName(displayName string) bool {
	return len(displayName) >= minDisplayNameLength && len(displayName) <= maxDisplayNameLength
}

func validateUsername(username string) bool {
	if len(username) < minUsernameLength || len(username) > maxUsernameLength {
		return false
	}
	return usernameRegexp.MatchString(username)
}

var usernameRegexp = regexp.MustCompile(`^[a-zA-Z0-9._@\-]+$`)

func validateDescription(description string) bool {
	return len(description) >= minDescriptionLength && len(description) <= maxDescriptionLength
}
