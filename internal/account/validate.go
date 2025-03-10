package account

import (
	"regexp"
	"strings"

	"github.com/vanti-dev/sc-bos/pkg/gen"
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

func validateResourceType(rt gen.RoleAssignment_ResourceType) bool {
	if rt == gen.RoleAssignment_RESOURCE_TYPE_UNSPECIFIED {
		return false
	}
	_, ok := gen.RoleAssignment_ResourceType_name[int32(rt)]
	return ok
}

func validateResource(resource string) bool {
	return len(resource) > 0
}
