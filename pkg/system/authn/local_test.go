package authn

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/internal/account"
	"github.com/vanti-dev/sc-bos/internal/auth/accesstoken"
	"github.com/vanti-dev/sc-bos/internal/util/pass"
	"github.com/vanti-dev/sc-bos/pkg/system/authn/config"
)

func TestLocalUserVerifier_Verify(t *testing.T) {
	ctx := context.Background()
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	store := account.NewMemoryStore(logger)
	verifier := newLocalUserVerifier(store)

	identities := []config.Identity{
		{
			ID:    "user1",
			Title: "User 1 - Basic",
			Secrets: []config.Secret{
				{Hash: testHash(t, "password123")},
			},
			Roles: []string{"admin"},
		},
		{
			ID:    "user1",
			Title: "User 1 - Duplicate",
			Secrets: []config.Secret{
				{Hash: testHash(t, "password456")},
			},
		},
		{
			ID:    "multi-role-user",
			Title: "User with Multiple Roles",
			Secrets: []config.Secret{
				{Hash: testHash(t, "secure456")},
			},
			Roles: []string{"admin", "commissioner", "operator"},
		},
		{
			ID:    "invalid-role-user",
			Title: "User with Invalid Role",
			Secrets: []config.Secret{
				{Hash: testHash(t, "pass789")},
			},
			Roles: []string{"unknown-role"},
		},
		{
			ID:    "no-role-user",
			Title: "User with No Roles",
			Secrets: []config.Secret{
				{Hash: testHash(t, "password")},
			},
			Roles: []string{},
		},
		{
			ID:      "no-password-user",
			Title:   "User with No Password",
			Secrets: []config.Secret{},
			Roles:   []string{"viewer"},
		},
		{
			ID:    "multi-password-user",
			Title: "User with Multiple Passwords",
			Secrets: []config.Secret{
				{Hash: testHash(t, "firstpass")},
				{Hash: testHash(t, "secondpass")},
			},
			Roles: []string{"viewer", "operator"},
		},
		{
			ID:    "super-admin",
			Title: "Super Admin User",
			Secrets: []config.Secret{
				{Hash: testHash(t, "superadmin!")},
			},
			Roles: []string{"super-admin"},
		},
	}
	err = importIdentities(ctx, store, identities, logger)
	if err != nil {
		t.Fatalf("failed to import identities: %v", err)
	}

	type testCase struct {
		username      string
		password      string
		expectedError error
		expect        accesstoken.SecretData
	}
	testCases := map[string]testCase{
		"correct_password_and_roles": {
			username: "user1",
			password: "password123",
			expect: accesstoken.SecretData{
				Title:       "User 1 - Basic",
				SystemRoles: []string{"admin"},
			},
		},
		"wrong_password": {
			username:      "user1",
			password:      "wrong-password",
			expectedError: accesstoken.ErrInvalidCredentials,
		},
		"non_existent_user": {
			username:      "nonexistent",
			password:      "anything",
			expectedError: accesstoken.ErrInvalidCredentials,
		},
		"multiple_valid_roles": {
			username: "multi-role-user",
			password: "secure456",
			expect: accesstoken.SecretData{
				Title:       "User with Multiple Roles",
				SystemRoles: []string{"admin", "commissioner", "operator"},
			},
		},
		"no_roles_assigned": {
			username:      "no-role-user",
			password:      "password",
			expectedError: accesstoken.ErrNoRolesAssigned,
		},
		"invalid_role_missing": {
			username:      "invalid-role-user",
			password:      "pass789",
			expectedError: accesstoken.ErrNoRolesAssigned, // Since the role wasn't imported, user will have no valid roles
		},
		"no_password": {
			username:      "no-password-user",
			password:      "", // Try with empty password
			expectedError: accesstoken.ErrInvalidCredentials,
		},
		"multiple_passwords_first_works": {
			username: "multi-password-user",
			password: "firstpass",
			expect: accesstoken.SecretData{
				Title:       "User with Multiple Passwords",
				SystemRoles: []string{"operator", "viewer"},
			},
		},
		"multiple_passwords_second_fails": {
			username:      "multi-password-user",
			password:      "secondpass", // System only imports first password
			expectedError: accesstoken.ErrInvalidCredentials,
		},
		"super_admin_role": {
			username: "super-admin",
			password: "superadmin!",
			expect: accesstoken.SecretData{
				Title:       "Super Admin User",
				SystemRoles: []string{"super-admin"},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result, err := verifier.Verify(ctx, tc.username, tc.password)
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error:\n\t%v\ngot:\n\t%v", tc.expectedError, err)
			}
			if tc.expectedError == nil {
				diff := cmp.Diff(tc.expect, result,
					cmpopts.IgnoreFields(accesstoken.SecretData{}, "TenantID"),
				)
				if diff != "" {
					t.Errorf("unexpected result (-want +got):\n%s", diff)
				}
				if result.TenantID == "" {
					t.Errorf("expected non-empty TenantID, got empty")
				}
			}
		})
	}
}

func testHash(t *testing.T, password string) string {
	t.Helper()
	hashed, err := pass.Hash([]byte(password))
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	return string(hashed)
}
