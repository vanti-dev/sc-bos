package authn

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strconv"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/internal/account"
	"github.com/vanti-dev/sc-bos/internal/account/queries"
	"github.com/vanti-dev/sc-bos/internal/auth/accesstoken"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/system/authn/config"
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
			return accesstoken.ErrInvalidCredentials
		} else if err != nil {
			return err
		}

		err = tx.CheckAccountPassword(ctx, userAccount.AccountID, password)
		if errors.Is(err, account.ErrIncorrectPassword) {
			return accesstoken.ErrInvalidCredentials
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
		slices.Sort(data.Roles) // deterministic order for legacy roles
		data.Title = details.DisplayName
		data.TenantID = strconv.FormatInt(userAccount.AccountID, 10)
		if len(data.Roles) == 0 {
			// no point issuing a token because the user has no roles so they cannot access anything
			return accesstoken.ErrNoRolesAssigned
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

func importIdentities(ctx context.Context, accounts *account.Store, ids []config.Identity, logger *zap.Logger) error {
	err := accounts.Write(ctx, func(tx *account.Tx) error {
		legacyRoleIDs := make(map[string]int64)

		skipCount := 0
		importCount := 0
		for _, id := range ids {
			logger := logger.With(zap.String("username", id.ID))
			_, err := tx.GetAccountByUsername(ctx, id.ID)
			if err == nil {
				// skip import if the account already exists
				skipCount++
				continue
			} else if !errors.Is(err, sql.ErrNoRows) {
				// any other error implies a problem with the database
				return fmt.Errorf("failed to check if account %q exists: %w", id.ID, err)
			}

			// create a new user account
			created, err := tx.CreateAccount(ctx, queries.CreateAccountParams{
				DisplayName: id.Title,
				Type:        gen.Account_USER_ACCOUNT.String(),
			})
			if err != nil {
				return fmt.Errorf("failed to import user account %q: %w", id.ID, err)
			}

			var passwordHash []byte
			switch len(id.Secrets) {
			case 0:
				// no hash to import, user will not be able to log in until a password is set
			case 1:
				passwordHash = []byte(id.Secrets[0].Hash)
			default:
				passwordHash = []byte(id.Secrets[0].Hash) // use the first secret as the password
				logger.Warn("importing user account with multiple secrets, only the first will be imported")
			}

			_, err = tx.CreateUserAccount(ctx, queries.CreateUserAccountParams{
				AccountID:    created.ID,
				Username:     id.ID,
				PasswordHash: passwordHash,
			})
			if err != nil {
				return fmt.Errorf("failed to import user account %q: %w", id.ID, err)
			}

			// for each legacy role assigned to the user, we need to find the ID of a corresponding role in the store
			// so that the user gets the correct permissions.
			for _, legacyRole := range id.Roles {
				roleID, ok := legacyRoleIDs[legacyRole]
				if !ok {
					roleID, err = findRoleIDForLegacyRole(ctx, tx, legacyRole, logger)
					if errors.Is(err, errLegacyRoleNotFound) {
						logger.Warn("cannot assign legacy role to user account because a matching role was not found",
							zap.String("legacyRole", legacyRole))
						continue
					} else if err != nil {
						return fmt.Errorf("failed to find role ID for legacy role %q: %w", legacyRole, err)
					}
					legacyRoleIDs[legacyRole] = roleID
				}

				_, err = tx.CreateRoleAssignment(ctx, queries.CreateRoleAssignmentParams{
					AccountID: created.ID,
					RoleID:    roleID,
				})
				if err != nil {
					return fmt.Errorf("failed to add legacy role %q to user account %q: %w", legacyRole, id.ID, err)
				}
			}
			importCount++
		}
		if importCount > 0 {
			logger.Info("imported user accounts from file into database",
				zap.Int("imported", importCount),
				zap.Int("skipped", skipCount),
			)
		}
		return nil
	})
	return err
}

var errLegacyRoleNotFound = errors.New("legacy role not found")

func findRoleIDForLegacyRole(ctx context.Context, tx *account.Tx, legacyRole string, logger *zap.Logger) (int64, error) {
	roles, err := tx.ListRolesWithLegacyRole(ctx, sql.NullString{Valid: true, String: legacyRole})
	if err != nil {
		return 0, err
	}

	switch {
	case len(roles) == 0:
		return 0, errLegacyRoleNotFound
	case len(roles) > 1:
		logger.Warn("multiple roles map to legacy role, choosing one",
			zap.String("legacyRole", legacyRole),
			zap.Int("count", len(roles)),
			zap.Int64("roleID", roles[0].ID),
		)
	}
	return roles[0].ID, nil
}
