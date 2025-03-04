package account

import (
	"regexp"
	"strings"
)

const (
	minPasswordLength = 10
	maxPasswordLength = 72 // max password length supported by bcrypt
	minUsernameLength = 3
	maxUsernameLength = 100
)

func permitPassword(password string) bool {
	password = normalisePassword(password)
	return len(password) >= minPasswordLength && len(password) <= maxPasswordLength
}

func normalisePassword(password string) string {
	return strings.TrimSpace(password)
}

func validateDisplayName(title string) bool {
	return len(title) > 0
}

func validateUsername(username string) bool {
	if len(username) < minUsernameLength || len(username) > maxUsernameLength {
		return false
	}
	return usernameRegexp.MatchString(username)
}

var usernameRegexp = regexp.MustCompile(`^[a-zA-Z0-9._@\-]+$`)
