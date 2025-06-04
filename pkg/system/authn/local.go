package authn

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vanti-dev/sc-bos/internal/account"
	"github.com/vanti-dev/sc-bos/internal/auth/accesstoken"
	"github.com/vanti-dev/sc-bos/internal/util/pass"
)

var errUsernamePassword = status.Error(codes.Unauthenticated, "invalid username or password")

type localUserVerifier struct {
	accounts *account.Store
}

func newLocalUserVerifier(accounts *account.Store) *localUserVerifier {
	return &localUserVerifier{
		accounts: accounts,
	}
}

func (l *localUserVerifier) Verify(ctx context.Context, username, password string) (accesstoken.SecretData, error) {
	var data accesstoken.SecretData
	err := l.accounts.Read(ctx, func(tx *account.Tx) error {
		userAccount, err := tx.GetAccountByUsername(ctx, username)
		if errors.Is(err, sql.ErrNoRows) {
			return errUsernamePassword
		} else if err != nil {
			return err
		}

		err = tx.CheckAccountPassword(ctx, userAccount.AccountID, password)
		if errors.Is(err, pass.ErrMismatchedHashAndPassword) {
			return errUsernamePassword
		} else if err != nil {
			return err
		}

		details, err := tx.GetAccountDetails(ctx, userAccount.AccountID)
		if err != nil {
			return err
		}

		data.Title = details.DisplayName
		data.TenantID = strconv.FormatInt(userAccount.AccountID, 10)
		data.Roles = []string{"admin"}
		return nil
	})
	if err != nil {
		return accesstoken.SecretData{}, err
	}

	return data, nil
}

type localServiceVerifier struct {
	accounts *account.Store
}

func newLocalServiceVerifier(accounts *account.Store) *localServiceVerifier {
	return &localServiceVerifier{
		accounts: accounts,
	}
}

func (l *localServiceVerifier) Verify(ctx context.Context, clientID, secret string) (accesstoken.SecretData, error) {
	// TODO implement me
	panic("implement me")
}
