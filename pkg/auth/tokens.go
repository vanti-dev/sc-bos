package auth

import (
	"context"

	"go.uber.org/multierr"
)

type TokenVerifier interface {
	VerifyAccessToken(ctx context.Context, token string) (*Authorization, error)
}

func NewMultiTokenVerifier(verifiers ...TokenVerifier) TokenVerifier {
	return multiTokenVerifier(verifiers)
}

type multiTokenVerifier []TokenVerifier

func (m multiTokenVerifier) VerifyAccessToken(ctx context.Context, token string) (*Authorization, error) {
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
