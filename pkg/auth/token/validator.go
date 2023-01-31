package token

import (
	"context"

	"go.uber.org/multierr"
)

type Validator interface {
	ValidateAccessToken(ctx context.Context, token string) (*Claims, error)
}

type ValidatorFunc func(ctx context.Context, token string) (*Claims, error)

func (t ValidatorFunc) ValidateAccessToken(ctx context.Context, token string) (*Claims, error) {
	return t(ctx, token)
}

func NewMultiTokenValidator(verifiers ...Validator) Validator {
	return multiValidator(verifiers)
}

type multiValidator []Validator

func (m multiValidator) ValidateAccessToken(ctx context.Context, token string) (*Claims, error) {
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
