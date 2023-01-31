package token

import (
	"context"

	"go.uber.org/multierr"
)

// Validator allows you to validate that an access token, typically given via an API request, is valid.
// Validator implementations should return a non-nil error if the validation fails in any way.
// For example if the payload fails to parse, or the expiry date is outside a supported range.
type Validator interface {
	// ValidateAccessToken returns a non-nil error if token is valid.
	// Claims are returned containing any information we know to be true about the token.
	ValidateAccessToken(ctx context.Context, token string) (*Claims, error)
}

// ValidatorFunc implements Validator wrapping a func of the correct signature.
type ValidatorFunc func(ctx context.Context, token string) (*Claims, error)

func (t ValidatorFunc) ValidateAccessToken(ctx context.Context, token string) (*Claims, error) {
	return t(ctx, token)
}

// NewMultiValidator creates a Validator that validates tokens if any validator validates the token.
func NewMultiValidator(validators ...Validator) Validator {
	return multiValidator(validators)
}

type multiValidator []Validator

func (m multiValidator) ValidateAccessToken(ctx context.Context, token string) (*Claims, error) {
	var errs error
	for _, verifier := range m {
		claims, err := verifier.ValidateAccessToken(ctx, token)
		if err == nil {
			return claims, nil
		}

		errs = multierr.Append(errs, err)
	}
	return nil, errs
}
