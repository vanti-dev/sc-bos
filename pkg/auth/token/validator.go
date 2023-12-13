// Package token provides mechanisms for validating access tokens and extracting claims.
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

// ValidatorSet is a collection of Validators where a token is deemed valid if any member Validator deems it valid.
type ValidatorSet []Validator

func (m *ValidatorSet) ValidateAccessToken(ctx context.Context, token string) (*Claims, error) {
	var errs error
	for _, verifier := range *m {
		claims, err := verifier.ValidateAccessToken(ctx, token)
		if err == nil {
			return claims, nil
		}

		errs = multierr.Append(errs, err)
	}
	return nil, errs
}

func (m *ValidatorSet) Append(v Validator) {
	*m = append(*m, v)
}

func (m *ValidatorSet) Delete(v Validator) {
	for i, validator := range *m {
		if validator == v {
			*m = append((*m)[:i], (*m)[i+1:]...)
			return
		}
	}
}

// NeverValid returns a Validator that always returns err.
func NeverValid(err error) Validator {
	return neverValid{err: err}
}

type neverValid struct {
	err error
}

func (nv neverValid) ValidateAccessToken(ctx context.Context, token string) (*Claims, error) {
	return nil, nv.err
}

// AlwaysValid returns a Validator that always returns claims.
func AlwaysValid(claims *Claims) Validator {
	return alwaysValid{claims: claims}
}

type alwaysValid struct {
	claims *Claims
}

func (av alwaysValid) ValidateAccessToken(ctx context.Context, token string) (*Claims, error) {
	return av.claims, nil
}
