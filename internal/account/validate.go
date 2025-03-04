package account

import (
	"strings"
)

const (
	minPasswordLength = 10
	maxPasswordLength = 72 // max password length supported by bcrypt
)

func permitPassword(password string) bool {
	return len(password) >= minPasswordLength && len(password) <= maxPasswordLength
}

func normalisePassword(password string) string {
	return strings.TrimSpace(password)
}

func validateDisplayName(title string) bool {
	return len(title) > 0
}

func validateUsername(username string) bool {
	return len(username) > 0
}
