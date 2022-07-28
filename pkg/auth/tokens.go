package auth

import (
	"context"

	"go.uber.org/multierr"
)

type TokenClaims struct {
	Issuer    string   // A unique identifier for the entity that provided this authorization data. Should be a URL.
	Subject   string   // A unique identifier for the entity that has been authorized access (e.g. a user or service)
	Roles     []string // The names of the roles that the subject has been granted
	Scopes    []string // The scopes that this authorization is limited to
	IsService bool     // True if the subject is an application acting on its own behalf, false if it's a user
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

type TokenVerifier interface {
	VerifyAccessToken(ctx context.Context, token string) (*TokenClaims, error)
}

func NewMultiTokenVerifier(verifiers ...TokenVerifier) TokenVerifier {
	return multiTokenVerifier(verifiers)
}

type multiTokenVerifier []TokenVerifier

func (m multiTokenVerifier) VerifyAccessToken(ctx context.Context, token string) (*TokenClaims, error) {
	var errs error
	for _, verifier := range m {
		authz, err := verifier.VerifyAccessToken(ctx, token)
		if err == nil {
			return authz, nil
		}

		errs = multierr.Append(errs, err)
	}
	return nil, errs
}
