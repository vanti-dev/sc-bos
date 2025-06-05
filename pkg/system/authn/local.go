package authn

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vanti-dev/sc-bos/internal/account"
	"github.com/vanti-dev/sc-bos/internal/auth/accesstoken"
	"github.com/vanti-dev/sc-bos/internal/util/pass"
)

var (
	errUsernamePassword = status.Error(codes.Unauthenticated, "invalid username or password")
	errNoRoles          = status.Error(codes.Unauthenticated, "cannot login because user has no roles assigned")
)

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

		legacyRoles, err := tx.ListLegacyRolesForAccount(ctx, userAccount.AccountID)
		if err != nil {
			return err
		}
		for _, role := range legacyRoles {
			if role.Valid {
				data.Roles = append(data.Roles, role.String)
			}
		}
		if len(data.Roles) == 0 {
			// no point issuing a token because the user has no roles so they cannot access anything
			return errNoRoles
		}
		slices.Sort(data.Roles) // deterministic order for legacy roles
		data.Title = details.DisplayName
		data.TenantID = strconv.FormatInt(userAccount.AccountID, 10)
		return nil
	})
	if err != nil {
		return accesstoken.SecretData{}, err
	}

	return data, nil
}
