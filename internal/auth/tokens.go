package auth

import (
	"context"

	"go.uber.org/multierr"
)

const (
	RoleTenant     = "tenant"
	RoleController = "controller"
	RoleUser       = "user"
)

type Authorization struct {
	Roles     []string `json:"roles"`      // The names of the roles that the subject has been granted
	Scopes    []string `json:"scopes"`     // The scopes that this authorization is limited to
	Zones     []string `json:"zones"`      // The zones that this token is authorized for, for tenant tokens
	IsService bool     `json:"is_service"` // True if the subject is an application acting on its own behalf, false if it's a user
}

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

type TokenValidator interface {
	ValidateAccessToken(ctx context.Context, token string) (*Authorization, error)
}

func NewMultiTokenValidator(verifiers ...TokenValidator) TokenValidator {
	return multiTokenValidator(verifiers)
}

type multiTokenValidator []TokenValidator

func (m multiTokenValidator) ValidateAccessToken(ctx context.Context, token string) (*Authorization, error) {
	var errs error
	for _, verifier := range m {
		authz, err := verifier.ValidateAccessToken(ctx, token)
		if err == nil {
			return authz, nil
		}

		errs = multierr.Append(errs, err)
	}
	return nil, errs
}
