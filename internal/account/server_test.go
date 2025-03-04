package account

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

func TestServer_CreateAccount(t *testing.T) {
	type testCase struct {
		others []*gen.Account // other accounts that should be created before the test account
		req    *gen.CreateAccountRequest
		expect *gen.Account
		code   codes.Code
	}

	cases := map[string]testCase{
		"user_account_no_password": {
			req: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 1",
					Username:    "user1",
				},
			},
			expect: &gen.Account{
				Type:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
		},
		"user_account_with_password": {
			req: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 2",
					Username:    "user2",
				},
				Password: "user2Password",
			},
			expect: &gen.Account{
				Type:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 2",
				Username:    "user2",
			},
		},
		"user_account_short_password": {
			req: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 3",
					Username:    "user3",
				},
				Password: "short",
			},
			code: codes.InvalidArgument,
		},
		"user_account_long_password": {
			req: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 4",
					Username:    "user4",
				},
				Password: strings.Repeat("a", 101),
			},
			code: codes.InvalidArgument,
		},
		"user_account_short_password_whitespace": {
			req: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 5",
					Username:    "user5",
				},
				Password: " short        ",
			},
			code: codes.InvalidArgument,
		},
		"service_account": {
			req: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_SERVICE_ACCOUNT,
					DisplayName: "Service",
				},
			},
			expect: &gen.Account{
				Type:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service",
			},
		},
		"service_account_password": {
			req: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_SERVICE_ACCOUNT,
					DisplayName: "Service",
				},
				Password: "servicePassword",
			},
			code: codes.InvalidArgument,
		},
		"missing_account_kind": {
			req: &gen.CreateAccountRequest{
				Account: &gen.Account{
					DisplayName: "Missing Kind",
				},
			},
			code: codes.InvalidArgument,
		},
		"missing_display_name": {
			req: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:     gen.Account_USER_ACCOUNT,
					Username: "foo",
				},
			},
			code: codes.InvalidArgument,
		},
		"missing_username_for_user_account": {
			req: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "Missing Username",
				},
			},
			code: codes.InvalidArgument,
		},
		"username_supplied_for_service_account": {
			req: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_SERVICE_ACCOUNT,
					DisplayName: "Service Account",
					Username:    "service",
				},
			},
			code: codes.InvalidArgument,
		},
		"username_conflict": {
			others: []*gen.Account{
				{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 1",
					Username:    "user1",
				},
			},
			req: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 1A",
					Username:    "user1",
				},
			},
			code: codes.AlreadyExists,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			store := NewMemoryStore(zap.NewNop())
			server := NewServer(store, zap.NewNop())

			for _, other := range tc.others {
				_, err := server.CreateAccount(context.Background(), &gen.CreateAccountRequest{Account: other})
				if err != nil {
					t.Fatalf("failed to create other account: %v", err)
				}
			}

			res, err := server.CreateAccount(context.Background(), tc.req)
			if tc.code != codes.OK {
				if status.Code(err) != tc.code {
					t.Fatalf("expected error code %v, got %v", tc.code, status.Code(err))
				}
				return
			}
			checkNilIfErrored(t, res, err)
			t.Helper()
			diff := cmp.Diff(tc.expect, res,
				protocmp.Transform(),
				protocmp.IgnoreFields(&gen.Account{}, "id", "create_time"),
			)
			if diff != "" {
				t.Errorf("unexpected provided account value (-want +got):\n%s", diff)
			}

			// also retrieve using GetAccount and check it matches
			id := res.Id
			expect := proto.Clone(tc.expect).(*gen.Account)
			expect.Id = id
			expect.CreateTime = res.CreateTime
			account, err := server.GetAccount(context.Background(), &gen.GetAccountRequest{Id: id})
			checkNilIfErrored(t, account, err)
			if err != nil {
				t.Fatalf("failed to get account %q: %v", id, err)
			}
			diff = cmp.Diff(expect, account, protocmp.Transform())
			if diff != "" {
				t.Errorf("unexpected retrieved account value (-want +got):\n%s", diff)
			}
		})
	}
}

// tests ordering and pagination of ListAccounts
func TestServer_ListAccounts(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore(zap.NewNop())
	server := NewServer(store, zap.NewNop())

	createAccount := func(ty gen.Account_Type, username, displayName string) string {
		t.Helper()
		res, err := server.CreateAccount(ctx, &gen.CreateAccountRequest{
			Account: &gen.Account{
				Type:        ty,
				Username:    username,
				DisplayName: displayName,
			},
		})
		checkNilIfErrored(t, res, err)
		if err != nil {
			t.Fatalf("failed to create account: %v", err)
		}
		return res.Id
	}

	// we assume that accounts are returned in creation order
	var expected []*gen.Account
	const numAccounts = 200
	for i := range numAccounts {
		username := fmt.Sprintf("account-%03d", i)
		displayName := fmt.Sprintf("Account %d", i)

		if i%2 == 0 {
			// make it a user account
			id := createAccount(gen.Account_USER_ACCOUNT, username, displayName)
			expected = append(expected, &gen.Account{
				Id:          id,
				Type:        gen.Account_USER_ACCOUNT,
				Username:    username,
				DisplayName: displayName,
			})
		} else {
			// make it a service account
			id := createAccount(gen.Account_SERVICE_ACCOUNT, "", displayName)
			expected = append(expected, &gen.Account{
				Id:          id,
				Type:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: displayName,
			})
		}
	}

	const pageSize = 42
	var nextPageToken string
	var got [][]*gen.Account
	for {
		t.Logf("fetching page with token %q", nextPageToken)
		res, err := server.ListAccounts(ctx, &gen.ListAccountsRequest{
			PageToken: nextPageToken,
			PageSize:  pageSize,
		})
		checkNilIfErrored(t, res, err)
		if err != nil {
			t.Fatalf("failed to list accounts: %v", err)
		}

		size := len(res.Accounts)
		if size < pageSize && res.NextPageToken != "" {
			t.Errorf("fewer results (%d) returned than expected (%d), but got a page token", size, pageSize)
		}

		got = append(got, res.Accounts)

		nextPageToken = res.NextPageToken
		if nextPageToken == "" {
			break
		}
	}
	comparePages(t, pageSize, expected, got, "create_time")

}

func TestServer_UpdateAccount(t *testing.T) {
	ctx := context.Background()
	type testCase struct {
		initial  *gen.Account
		update   *gen.UpdateAccountRequest
		expected *gen.Account
		code     codes.Code
	}
	cases := map[string]testCase{
		"empty update": {
			initial: &gen.Account{
				Type:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{},
			},
			expected: &gen.Account{
				Type:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
		},
		"kind_change_prohibited": {
			initial: &gen.Account{
				Type:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					Type: gen.Account_SERVICE_ACCOUNT,
				},
			},
			expected: &gen.Account{
				Type:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
			code: codes.InvalidArgument,
		},
		"same_kind_allowed": {
			initial: &gen.Account{
				Type:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					Type: gen.Account_USER_ACCOUNT,
				},
			},
			expected: &gen.Account{
				Type:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
		},
		"update_display_name": {
			initial: &gen.Account{
				Type:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					DisplayName: "Service MODIFIED",
				},
			},
			expected: &gen.Account{
				Type:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service MODIFIED",
			},
		},
		"update_username": {
			initial: &gen.Account{
				Type:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					Username: "user1-modified",
				},
			},
			expected: &gen.Account{
				Type:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1-modified",
			},
		},
		"update_username_service_account": {
			initial: &gen.Account{
				Type:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					Username: "username",
				},
			},
			expected: &gen.Account{
				Type:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service",
			},
			code: codes.FailedPrecondition,
		},
		"update_display_name_empty": {
			initial: &gen.Account{
				Type:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					DisplayName: "",
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"display_name"}},
			},
			expected: &gen.Account{
				Type:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service",
			},
			code: codes.InvalidArgument,
		},
		"update_username_empty_user": {
			initial: &gen.Account{
				Type:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					Username: "",
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"username"}},
			},
			expected: &gen.Account{
				Type:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
			code: codes.InvalidArgument,
		},
		"update_username_empty_service": {
			initial: &gen.Account{
				Type:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					Username: "",
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"username"}},
			},
			expected: &gen.Account{
				Type:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service",
			},
		},
		"invalid_update_mask": {
			initial: &gen.Account{
				Type:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
			update: &gen.UpdateAccountRequest{
				Account:    &gen.Account{},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"foo"}},
			},
			expected: &gen.Account{
				Type:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
			code: codes.InvalidArgument,
		},
		"wildcard_update_mask": {
			initial: &gen.Account{
				Type:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 1 MODIFIED",
					Username:    "user1-modified",
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"*"}},
			},
			expected: &gen.Account{
				Type:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1 MODIFIED",
				Username:    "user1-modified",
			},
		},
		"wildcard_update_mask_zero": {
			initial: &gen.Account{
				Type:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 1 MODIFIED",
					Username:    "",
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"*"}},
			},
			expected: &gen.Account{
				Type:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1", // not changed, because updates are all-or-nothing
				Username:    "user1",
			},
			code: codes.InvalidArgument, // because username is required, and we have tried to clear it
		},
		"wildcard_update_mask_service": {
			initial: &gen.Account{
				Type:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_SERVICE_ACCOUNT,
					DisplayName: "Service MODIFIED",
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"*"}},
			},
			expected: &gen.Account{
				Type:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service MODIFIED",
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			store := NewMemoryStore(zap.NewNop())
			server := NewServer(store, zap.NewNop())

			account, err := server.CreateAccount(ctx, &gen.CreateAccountRequest{Account: tc.initial})
			checkNilIfErrored(t, account, err)
			if err != nil {
				t.Fatalf("failed to create account: %v", err)
			}
			// inject ID, which is now known, into update and expected
			id := account.Id
			update := proto.Clone(tc.update).(*gen.UpdateAccountRequest)
			update.Account.Id = id
			expected := proto.Clone(tc.expected).(*gen.Account)
			expected.Id = id
			expected.CreateTime = account.CreateTime

			updated, err := server.UpdateAccount(ctx, update)
			checkNilIfErrored(t, updated, err)
			if status.Code(err) != tc.code {
				t.Errorf("expected error with code %v, got %v", tc.code, err)
			}
			if updated != nil {
				diff := cmp.Diff(expected, updated, protocmp.Transform())
				if diff != "" {
					t.Errorf("unexpected updated account value (-got +want):\n%s", diff)
				}
			}

			// fetch again to check that the update was persisted
			account, err = server.GetAccount(ctx, &gen.GetAccountRequest{Id: id})
			checkNilIfErrored(t, account, err)
			if err != nil {
				t.Fatalf("failed to get account: %v", err)
			}
			diff := cmp.Diff(expected, account, protocmp.Transform())
			if diff != "" {
				t.Errorf("unexpected retrieved account value (-got +want):\n%s", diff)
			}
		})
	}
}

func TestServer_DeleteAccount(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore(zap.NewNop())
	server := NewServer(store, zap.NewNop())

	createAccount := func(ty gen.Account_Type, username, displayName, password string) *gen.Account {
		t.Helper()
		res, err := server.CreateAccount(ctx, &gen.CreateAccountRequest{
			Account: &gen.Account{
				Type:        ty,
				Username:    username,
				DisplayName: displayName,
			},
			Password: password,
		})
		checkNilIfErrored(t, res, err)
		if err != nil {
			t.Fatalf("failed to create account: %v", err)
		}
		return res
	}
	user := createAccount(gen.Account_USER_ACCOUNT, "user1", "User 1", "user1Password")
	service := createAccount(gen.Account_SERVICE_ACCOUNT, "", "Service", "")

	// add a service credential to the service account, to check that
	//  - existence of service credentials does not prevent deletion of the account
	//  - service credentials are deleted when the account is deleted
	cred, err := server.CreateServiceCredential(ctx, &gen.CreateServiceCredentialRequest{
		ServiceCredential: &gen.ServiceCredential{
			AccountId:   service.Id,
			DisplayName: "Credential",
		},
	})
	if err != nil {
		t.Fatalf("failed to create service credential: %v", err)
	}

	// assign a role to the accounts
	// to check that role assignments
	//   - do not prevent deletion of the account
	//   - are deleted when the account is deleted
	role, err := server.CreateRole(ctx, &gen.CreateRoleRequest{
		Role: &gen.Role{
			DisplayName: "Foo Role",
			Permissions: []string{"foo"},
		},
	})
	if err != nil {
		t.Fatalf("failed to create role: %v", err)
	}
	ra1, err := server.CreateRoleAssignment(ctx, &gen.CreateRoleAssignmentRequest{
		RoleAssignment: &gen.RoleAssignment{
			AccountId: user.Id,
			RoleId:    role.Id,
		},
	})
	if err != nil {
		t.Fatalf("failed to create role assignment: %v", err)
	}
	ra2, err := server.CreateRoleAssignment(ctx, &gen.CreateRoleAssignmentRequest{
		RoleAssignment: &gen.RoleAssignment{
			AccountId: service.Id,
			RoleId:    role.Id,
		},
	})

	// delete the accounts
	_, err = server.DeleteAccount(ctx, &gen.DeleteAccountRequest{Id: user.Id})
	if err != nil {
		t.Errorf("failed to delete user account: %v", err)
	}
	_, err = server.DeleteAccount(ctx, &gen.DeleteAccountRequest{Id: service.Id})
	if err != nil {
		t.Errorf("failed to delete service account: %v", err)
	}

	// check that the accounts are actually gone
	_, err = server.GetAccount(ctx, &gen.GetAccountRequest{Id: user.Id})
	if status.Code(err) != codes.NotFound {
		t.Errorf("expected NotFound error for user account, got %v", err)
	}
	_, err = server.GetAccount(ctx, &gen.GetAccountRequest{Id: service.Id})
	if status.Code(err) != codes.NotFound {
		t.Errorf("expected NotFound error for service account, got %v", err)
	}

	// check that the service credential is gone
	_, err = server.GetServiceCredential(ctx, &gen.GetServiceCredentialRequest{Id: cred.Id})
	if status.Code(err) != codes.NotFound {
		t.Errorf("expected NotFound error for service credential, got %v", err)
	}

	// check that the role assignments are gone
	_, err = server.GetRoleAssignment(ctx, &gen.GetRoleAssignmentRequest{Id: ra1.Id})
	if status.Code(err) != codes.NotFound {
		t.Errorf("expected NotFound error for role assignment, got %v", err)
	}
	_, err = server.GetRoleAssignment(ctx, &gen.GetRoleAssignmentRequest{Id: ra2.Id})
	if status.Code(err) != codes.NotFound {
		t.Errorf("expected NotFound error for role assignment, got %v", err)
	}
}

// tests that uniqueness of usernames is enforced
func TestServer_Account_Username(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore(zap.NewNop())
	server := NewServer(store, zap.NewNop())

	_, err := server.CreateAccount(ctx, &gen.CreateAccountRequest{
		Account: &gen.Account{
			Type:        gen.Account_USER_ACCOUNT,
			DisplayName: "User 1",
			Username:    "user1",
		},
	})
	if err != nil {
		t.Fatalf("failed to create account1: %v", err)
	}
	user2, err := server.CreateAccount(ctx, &gen.CreateAccountRequest{
		Account: &gen.Account{
			Type:        gen.Account_USER_ACCOUNT,
			DisplayName: "User 2",
			Username:    "user2",
		},
	})
	if err != nil {
		t.Fatalf("failed to create account2: %v", err)
	}

	// try to create account that collides with user1
	_, err = server.CreateAccount(ctx, &gen.CreateAccountRequest{
		Account: &gen.Account{
			Type:        gen.Account_USER_ACCOUNT,
			DisplayName: "User 1A",
			Username:    "user1",
		},
	})
	if status.Code(err) != codes.AlreadyExists {
		t.Errorf("expected AlreadyExists error, got %v", err)
	}

	// try to change user2 to user1
	_, err = server.UpdateAccount(ctx, &gen.UpdateAccountRequest{
		Account: &gen.Account{
			Id:       user2.Id,
			Username: "user1",
		},
	})
	if status.Code(err) != codes.AlreadyExists {
		t.Errorf("expected AlreadyExists error, got %v", err)
	}
}

func TestServer_ServiceCredentials(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore(zap.NewNop())
	server := NewServer(store, zap.NewNop())

	account, err := server.CreateAccount(ctx, &gen.CreateAccountRequest{
		Account: &gen.Account{
			Type:        gen.Account_SERVICE_ACCOUNT,
			DisplayName: "Service",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	checkCreds := func(expect []*gen.ServiceCredential) {
		res, err := server.ListServiceCredentials(ctx, &gen.ListServiceCredentialsRequest{
			AccountId: account.Id,
		})
		checkNilIfErrored(t, res, err)
		if err != nil {
			t.Fatalf("failed to list service credentials: %v", err)
		}

		diff := cmp.Diff(expect, res.ServiceCredentials,
			protocmp.Transform(),
			protocmp.IgnoreFields(&gen.ServiceCredential{}, "secret"),
		)
		if diff != "" {
			t.Errorf("unexpected service credentials (-got +want):\n%s", diff)
		}

		// fetch them individually, checking we get the same result as in the list
		for _, cred := range res.ServiceCredentials {
			// a secret should never be returned except at creation time
			if cred.Secret != "" {
				t.Errorf("service credential %q has a secret", cred.Id)
			}

			fetched, err := server.GetServiceCredential(ctx, &gen.GetServiceCredentialRequest{Id: cred.Id})
			checkNilIfErrored(t, fetched, err)
			if err != nil {
				t.Fatalf("failed to get service credential %q: %v", cred.Id, err)
			}
			diff = cmp.Diff(cred, fetched, protocmp.Transform())
			if diff != "" {
				t.Errorf("unexpected fetched service credential %q (-got +want):\n%s", cred.Id, diff)
			}
		}
	}

	var creds []*gen.ServiceCredential
	expiry := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := range maxServiceCredentialsPerAccount {
		cred, err := server.CreateServiceCredential(ctx, &gen.CreateServiceCredentialRequest{
			ServiceCredential: &gen.ServiceCredential{
				AccountId:   account.Id,
				DisplayName: fmt.Sprintf("Credential %d", i),
				ExpireTime:  timestamppb.New(expiry),
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		creds = append(creds, cred)
	}

	checkCreds(creds)

	// creating a new credential should fail because the limit is reached
	_, err = server.CreateServiceCredential(ctx, &gen.CreateServiceCredentialRequest{
		ServiceCredential: &gen.ServiceCredential{
			AccountId:   account.Id,
			DisplayName: "Credential",
		},
	})
	if status.Code(err) != codes.ResourceExhausted {
		t.Errorf("expected ResourceExhausted error, got %v", err)
	}

	checkCreds(creds)

	// delete a credential, check it's gone
	_, err = server.DeleteServiceCredential(ctx, &gen.DeleteServiceCredentialRequest{
		Id: creds[0].Id,
	})
	if err != nil {
		t.Fatalf("failed to delete service credential: %v", err)
	}
	creds = creds[1:]
	checkCreds(creds)

	// delete the account, check that credential access fails
	_, err = server.DeleteAccount(ctx, &gen.DeleteAccountRequest{
		Id: account.Id,
	})
	if err != nil {
		t.Fatalf("failed to delete account: %v", err)
	}
	_, err = server.ListServiceCredentials(ctx, &gen.ListServiceCredentialsRequest{
		AccountId: account.Id,
	})
	if status.Code(err) != codes.NotFound {
		t.Errorf("expected NotFound error querying for credentials of deleted account, got %v", err)
	}

}

func TestServer_UpdateAccountPassword(t *testing.T) {
	ctx := context.Background()

	type testCase struct {
		create           *gen.CreateAccountRequest
		update           *gen.UpdateAccountPasswordRequest
		validPassword    string   // password to check after update
		invalidPasswords []string // passwords to check for invalidity
		code             codes.Code
	}

	cases := map[string]testCase{
		"add_password": {
			create: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 1",
					Username:    "user1",
				},
			},
			update: &gen.UpdateAccountPasswordRequest{
				NewPassword: "user1Password",
			},
			validPassword: "user1Password",
		},
		"change_password": {
			create: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 1",
					Username:    "user1",
				},
				Password: "thepassword1",
			},
			update: &gen.UpdateAccountPasswordRequest{
				NewPassword: "thepassword2",
			},
			validPassword:    "thepassword2",
			invalidPasswords: []string{"thepassword1"},
		},
		"change_password_valid_old_password": {
			create: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 1",
					Username:    "user1",
				},
				Password: "thepassword1",
			},
			update: &gen.UpdateAccountPasswordRequest{
				OldPassword: "thepassword1",
				NewPassword: "thepassword2",
			},
			validPassword:    "thepassword2",
			invalidPasswords: []string{"thepassword1"},
		},
		"change_password_invalid_old_password": {
			create: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 1",
					Username:    "user1",
				},
				Password: "thepassword1",
			},
			update: &gen.UpdateAccountPasswordRequest{
				OldPassword: "wrongPassword",
				NewPassword: "thepassword2",
			},
			validPassword:    "thepassword1",
			invalidPasswords: []string{"thepassword2", "wrongPassword"},
			code:             codes.FailedPrecondition,
		},
		"add_password_old_password_supplied": {
			create: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 1",
					Username:    "user1",
				},
			},
			update: &gen.UpdateAccountPasswordRequest{
				OldPassword: "thepassword1",
				NewPassword: "thepassword2",
			},
			invalidPasswords: []string{"thepassword1", "thepassword2"},
			code:             codes.FailedPrecondition,
		},
		"account_id_not_found": {
			update: &gen.UpdateAccountPasswordRequest{
				Id:          "12345",
				NewPassword: "thepassword",
			},
			code: codes.NotFound,
		},
		"account_id_invalid": {
			update: &gen.UpdateAccountPasswordRequest{
				Id:          "invalid",
				NewPassword: "thepassword",
			},
			code: codes.NotFound,
		},
		"account_id_empty": {
			update: &gen.UpdateAccountPasswordRequest{
				Id:          "",
				NewPassword: "thepassword",
			},
			code: codes.InvalidArgument,
		},
		"add_password_empty": {
			create: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 1",
					Username:    "user1",
				},
			},
			update: &gen.UpdateAccountPasswordRequest{
				NewPassword: "",
			},
			code: codes.InvalidArgument,
		},
		"add_password_too_short": {
			create: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 1",
					Username:    "user1",
				},
			},
			update: &gen.UpdateAccountPasswordRequest{
				NewPassword: "123456789",
			},
			code: codes.InvalidArgument,
		},
		"add_password_long": {
			create: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 1",
					Username:    "user1",
				},
			},
			update: &gen.UpdateAccountPasswordRequest{
				NewPassword: strings.Repeat("a", maxPasswordLength),
			},
		},
		"add_password_too_long": {
			create: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 1",
					Username:    "user1",
				},
			},
			update: &gen.UpdateAccountPasswordRequest{
				NewPassword: strings.Repeat("a", maxPasswordLength+1),
			},
			code: codes.InvalidArgument,
		},
		"password_ignores_leading_whitespace": {
			create: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 1",
					Username:    "user1",
				},
			},
			update: &gen.UpdateAccountPasswordRequest{
				NewPassword: "  thepassword",
			},
			validPassword: "thepassword",
		},
		"password_ignores_trailing_whitespace": {
			create: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 1",
					Username:    "user1",
				},
			},
			update: &gen.UpdateAccountPasswordRequest{
				NewPassword: "thepassword  ",
			},
			validPassword: "thepassword",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			store := NewMemoryStore(zap.NewNop())
			server := NewServer(store, zap.NewNop())

			var account *gen.Account
			var err error
			if tc.create != nil {
				account, err = server.CreateAccount(ctx, tc.create)
				if err != nil {
					t.Fatalf("failed to create account: %v", err)
				}
			}

			req := proto.Clone(tc.update).(*gen.UpdateAccountPasswordRequest)
			if account != nil {
				req.Id = account.Id
			}
			_, err = server.UpdateAccountPassword(ctx, req)
			if status.Code(err) != tc.code {
				t.Errorf("expected error with code %v, got %v", tc.code, err)
			}

			check := func(password string) error {
				return store.Read(ctx, func(tx *Tx) error {
					id, ok := parseID(account.Id)
					if !ok {
						t.Fatalf("failed to parse account ID %q", account.Id)
					}
					return tx.CheckAccountPassword(ctx, id, password)
				})
			}

			if tc.validPassword != "" {
				err = check(tc.validPassword)
				if err != nil {
					t.Errorf("password check for %q failed: %v", tc.validPassword, err)
				}
			}

			for _, password := range tc.invalidPasswords {
				err = check(password)
				if err == nil {
					t.Errorf("expected error for password %q, got nil", password)
				}
			}

		})
	}

}

func TestServer_Role(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore(zap.NewNop())
	server := NewServer(store, zap.NewNop())

	const numRoles = 200
	const numPermissions = 10

	var roles []*gen.Role
	t.Log("CreateRole:")
	for i := range numRoles {
		displayName := fmt.Sprintf("Role %d", i)
		role, err := server.CreateRole(ctx, &gen.CreateRoleRequest{
			Role: &gen.Role{
				DisplayName: displayName,
				// supply the permissions shuffled to check they are returned in order instead
				Permissions: shuffledPermissions(numPermissions),
			},
		})
		if err != nil {
			t.Fatalf("failed to create role %q: %v", displayName, err)
		}
		expect := &gen.Role{
			Id:          role.Id,
			DisplayName: displayName,
			Permissions: orderedPermissions(numPermissions),
		}
		diff := cmp.Diff(expect, role, protocmp.Transform())
		if diff != "" {
			t.Errorf("unexpected role value (-got +want):\n%s", diff)
		}

		roles = append(roles, role)
	}

	checkAll := func(expect []*gen.Role) {
		var pages [][]*gen.Role
		const pageSize = 42
		var nextPageToken string
		for {
			res, err := server.ListRoles(ctx, &gen.ListRolesRequest{
				PageToken: nextPageToken,
				PageSize:  pageSize,
			})
			checkNilIfErrored(t, res, err)
			if err != nil {
				t.Fatalf("failed to list roles: %v", err)
			}
			t.Logf("fetched page with token %q, returned %d results", nextPageToken, len(res.Roles))

			if res.NextPageToken != "" && len(res.Roles) < pageSize {
				t.Errorf("fewer results (%d) returned than expected (%d), but got a page token", len(res.Roles), pageSize)
			}

			pages = append(pages, res.Roles)
			nextPageToken = res.NextPageToken
			if nextPageToken == "" {
				break
			}
		}
		comparePages(t, pageSize, roles, pages)
	}

	checkAll(roles)

	// test that roles can be updated
	t.Log("UpdateRole:")
	role := roles[0]
	role.Permissions = append(role.Permissions, "000-new-permission") // should go at the beginning
	role.DisplayName += " MODIFIED"
	updated, err := server.UpdateRole(ctx, &gen.UpdateRoleRequest{
		Role: role,
	})
	slices.Sort(role.Permissions)
	checkNilIfErrored(t, role, err)
	if err != nil {
		t.Fatalf("failed to update role: %v", err)
	}
	diff := cmp.Diff(role, updated, protocmp.Transform())
	if diff != "" {
		t.Errorf("unexpected updated role value (-got +want):\n%s", diff)
	}
	// test that update is persisted
	updated, err = server.GetRole(ctx, &gen.GetRoleRequest{Id: role.Id})
	checkNilIfErrored(t, updated, err)
	if err != nil {
		t.Fatalf("failed to get updated role: %v", err)
	}
	diff = cmp.Diff(role, updated, protocmp.Transform())
	if diff != "" {
		t.Errorf("unexpected retrieved role value (-got +want):\n%s", diff)
	}

	// test that a role can't be deleted if it is assigned
	t.Log("DeleteRole:")
	account, err := server.CreateAccount(ctx, &gen.CreateAccountRequest{
		Account: &gen.Account{
			Type:        gen.Account_SERVICE_ACCOUNT,
			DisplayName: "foo account",
		},
	})
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}
	assignment, err := server.CreateRoleAssignment(ctx, &gen.CreateRoleAssignmentRequest{
		RoleAssignment: &gen.RoleAssignment{
			AccountId: account.Id,
			RoleId:    role.Id,
		},
	})
	if err != nil {
		t.Fatalf("failed to create role assignment: %v", err)
	}
	_, err = server.DeleteRole(ctx, &gen.DeleteRoleRequest{Id: role.Id})
	if status.Code(err) != codes.FailedPrecondition {
		t.Errorf("expected FailedPrecondition error when deleting role with assignments, got %v", err)
	}

	// delete the assignment and try again
	_, err = server.DeleteRoleAssignment(ctx, &gen.DeleteRoleAssignmentRequest{Id: assignment.Id})
	if err != nil {
		t.Fatalf("failed to delete role assignment: %v", err)
	}

	// test that roles can be deleted
	_, err = server.DeleteRole(ctx, &gen.DeleteRoleRequest{Id: role.Id})
	if err != nil {
		t.Fatalf("failed to delete role: %v", err)
	}

	// check that the role is actually gone
	_, err = server.GetRole(ctx, &gen.GetRoleRequest{Id: role.Id})
	if status.Code(err) != codes.NotFound {
		t.Errorf("expected NotFound error for role, got %v", err)
	}
	// deleting it again should fail
	_, err = server.DeleteRole(ctx, &gen.DeleteRoleRequest{Id: role.Id})
	if status.Code(err) != codes.NotFound {
		t.Errorf("expected NotFound error for role, got %v", err)
	}
}

func TestServer_UpdateRole(t *testing.T) {
	ctx := context.Background()

	type testCase struct {
		initial  *gen.Role
		update   *gen.UpdateRoleRequest
		expected *gen.Role
		code     codes.Code
	}
	cases := map[string]testCase{
		"empty update": {
			initial: &gen.Role{
				DisplayName: "Role 1",
				Permissions: []string{"foo", "bar"},
			},
			update: &gen.UpdateRoleRequest{
				Role: &gen.Role{},
			},
			expected: &gen.Role{
				DisplayName: "Role 1",
				Permissions: []string{"bar", "foo"},
			},
		},
		"update_display_name_implicit": {
			initial: &gen.Role{
				DisplayName: "Role 1",
				Permissions: []string{"foo", "bar"},
			},
			update: &gen.UpdateRoleRequest{
				Role: &gen.Role{
					DisplayName: "Role 1 MODIFIED",
				},
			},
			expected: &gen.Role{
				DisplayName: "Role 1 MODIFIED",
				Permissions: []string{"bar", "foo"},
			},
		},
		"update_display_name_explicit": {
			initial: &gen.Role{
				DisplayName: "Role 1",
				Permissions: []string{"foo", "bar"},
			},
			update: &gen.UpdateRoleRequest{
				Role: &gen.Role{
					DisplayName: "Role 1 MODIFIED",
					Permissions: []string{"foo2", "bar2"},
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"display_name"}},
			},
			expected: &gen.Role{
				DisplayName: "Role 1 MODIFIED",
				Permissions: []string{"bar", "foo"},
			},
		},
		"update_permissions_implicit": {
			initial: &gen.Role{
				DisplayName: "Role 1",
				Permissions: []string{"foo", "bar"},
			},
			update: &gen.UpdateRoleRequest{
				Role: &gen.Role{
					Permissions: []string{"foo2", "bar2"},
				},
			},
			expected: &gen.Role{
				DisplayName: "Role 1",
				Permissions: []string{"bar2", "foo2"},
			},
		},
		"update_permissions_explicit": {
			initial: &gen.Role{
				DisplayName: "Role 1",
				Permissions: []string{"foo", "bar"},
			},
			update: &gen.UpdateRoleRequest{
				Role: &gen.Role{
					DisplayName: "Role 1 MODIFIED",
					Permissions: []string{"foo2", "bar2"},
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"permissions"}},
			},
			expected: &gen.Role{
				DisplayName: "Role 1",
				Permissions: []string{"bar2", "foo2"},
			},
		},
		"update_permissions_empty": {
			initial: &gen.Role{
				DisplayName: "Role 1",
				Permissions: []string{"foo", "bar"},
			},
			update: &gen.UpdateRoleRequest{
				Role: &gen.Role{
					Permissions: nil,
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"permissions"}},
			},
			expected: &gen.Role{
				DisplayName: "Role 1",
				Permissions: nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			store := NewMemoryStore(zap.NewNop())
			server := NewServer(store, zap.NewNop())

			role, err := server.CreateRole(ctx, &gen.CreateRoleRequest{
				Role: tc.initial,
			})
			checkNilIfErrored(t, role, err)
			if err != nil {
				t.Fatalf("failed to create role: %v", err)
			}

			// inject ID, which is now known, into update and expected
			id := role.Id
			update := proto.Clone(tc.update).(*gen.UpdateRoleRequest)
			update.Role.Id = id
			expected := proto.Clone(tc.expected).(*gen.Role)
			expected.Id = id

			updated, err := server.UpdateRole(ctx, update)
			checkNilIfErrored(t, updated, err)
			if status.Code(err) != tc.code {
				t.Errorf("expected error with code %v, got %v", tc.code, err)
			}

			if updated != nil {
				diff := cmp.Diff(expected, updated, protocmp.Transform())
				if diff != "" {
					t.Errorf("unexpected updated role value (-got +want):\n%s", diff)
				}
			}

			// fetch again to check that the update was persisted
			role, err = server.GetRole(ctx, &gen.GetRoleRequest{Id: id})
			checkNilIfErrored(t, role, err)
			if err != nil {
				t.Fatalf("failed to get role: %v", err)
			}
			diff := cmp.Diff(expected, role, protocmp.Transform())
			if diff != "" {
				t.Errorf("unexpected retrieved role value (-got +want):\n%s", diff)
			}
		})
	}
}

func comparePages[T messageWithID](t *testing.T, pageSize int32, expect []T, gotPages [][]T, ignoreFields ...protoreflect.Name) {
	var zero T
	t.Helper()
	for i, page := range gotPages {
		var allIDs []string
		for _, a := range page {
			allIDs = append(allIDs, a.GetId())
		}
		t.Logf("page %d contains %d items: %s", i, len(page), strings.Join(allIDs, ", "))
		if len(page) > int(pageSize) {
			t.Errorf("page %d has more than %d items: %d", i, pageSize, len(page))
		}
		if len(page) > len(expect) {
			t.Errorf("page %d has more items (%d) than remaining expected items (%d)", i, len(page), len(expect))
			return
		}

		expectPage := expect[:len(page)]
		expect = expect[len(page):]

		diff := cmp.Diff(expectPage, page,
			protocmp.Transform(),
			protocmp.IgnoreFields(zero, ignoreFields...),
		)
		if diff != "" {
			t.Errorf("unexpected page %d contents (-got +want):\n%s", i, diff)
		}
	}

	if len(expect) > 0 {
		t.Errorf("expected %d more items than received", len(expect))
	}
}

func TestServer_RoleAssignments(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore(zap.NewNop())
	server := NewServer(store, zap.NewNop())

	createAccount := func(displayName string) *gen.Account {
		t.Helper()
		res, err := server.CreateAccount(ctx, &gen.CreateAccountRequest{
			Account: &gen.Account{
				Type:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: displayName,
			},
		})
		checkNilIfErrored(t, res, err)
		if err != nil {
			t.Fatalf("failed to create account: %v", err)
		}
		return res
	}
	createRole := func(displayName string, permissions ...string) *gen.Role {
		t.Helper()
		res, err := server.CreateRole(ctx, &gen.CreateRoleRequest{
			Role: &gen.Role{
				DisplayName: displayName,
				Permissions: permissions,
			},
		})
		checkNilIfErrored(t, res, err)
		if err != nil {
			t.Fatalf("failed to create role: %v", err)
		}
		return res
	}

	var accounts []*gen.Account
	var roles []*gen.Role
	const numAccounts = 50
	const numRoles = 50

	for i := range numAccounts {
		accounts = append(accounts, createAccount(fmt.Sprintf("Account %d", i)))
	}
	for i := range numRoles {
		roles = append(roles, createRole(fmt.Sprintf("Role %d", i)))
	}

	// create role assignments randomly
	var assignments []*gen.RoleAssignment
	for _, account := range accounts {
		for _, role := range roles {
			// don't assign all roles to all accounts, just some
			if rand.IntN(2) == 0 {
				continue
			}

			assignment, err := server.CreateRoleAssignment(ctx, &gen.CreateRoleAssignmentRequest{
				RoleAssignment: &gen.RoleAssignment{
					AccountId: account.Id,
					RoleId:    role.Id,
				},
			})
			checkNilIfErrored(t, assignment, err)
			if err != nil {
				t.Errorf("failed to create role assignment between account=%s and role=%s: %v", account.Id, role.Id, err)
			}

			assignments = append(assignments, assignment)
		}
	}

	type filter struct {
		name   string
		filter string
		expect []*gen.RoleAssignment
	}
	var filters []filter
	// no filter
	filters = append(filters, filter{
		name:   "no filter",
		filter: "",
		expect: assignments,
	})
	{
		// filter by account
		account := accounts[rand.IntN(len(accounts))]
		var expect []*gen.RoleAssignment
		for _, a := range assignments {
			if a.AccountId == account.Id {
				expect = append(expect, a)
			}
		}
		filters = append(filters, filter{
			name:   "account",
			filter: fmt.Sprintf("account_id = %s", account.Id),
			expect: expect,
		})
	}
	{
		// filter by role
		role := roles[rand.IntN(len(roles))]
		var expect []*gen.RoleAssignment
		for _, a := range assignments {
			if a.RoleId == role.Id {
				expect = append(expect, a)
			}
		}
		filters = append(filters, filter{
			name:   "role",
			filter: fmt.Sprintf("role_id = %s", role.Id),
			expect: expect,
		})
	}

	for _, f := range filters {
		t.Run(f.name, func(t *testing.T) {
			var got [][]*gen.RoleAssignment
			const pageSize = 15
			var nextPageToken string
			for {
				res, err := server.ListRoleAssignments(ctx, &gen.ListRoleAssignmentsRequest{
					Filter:    f.filter,
					PageToken: nextPageToken,
					PageSize:  pageSize,
				})
				checkNilIfErrored(t, res, err)
				if err != nil {
					t.Errorf("failed to list role assignments with filter %q: %v", f.filter, err)
				}
				t.Logf("fetched page with token %q, returned %d results", nextPageToken, len(res.RoleAssignments))

				if res.NextPageToken != "" && len(res.RoleAssignments) < pageSize {
					t.Errorf("fewer results (%d) returned than expected (%d), but got a page token", len(res.RoleAssignments), pageSize)
				}

				got = append(got, res.RoleAssignments)
				nextPageToken = res.NextPageToken
				if nextPageToken == "" {
					break
				}
			}

			comparePages(t, pageSize, f.expect, got)
		})
	}

	// delete a role assignment, check it's gone
	_, err := server.DeleteRoleAssignment(ctx, &gen.DeleteRoleAssignmentRequest{
		Id: assignments[0].Id,
	})
	if err != nil {
		t.Fatalf("failed to delete role assignment: %v", err)
	}
	// check that the role assignment is actually gone
	_, err = server.GetRoleAssignment(ctx, &gen.GetRoleAssignmentRequest{Id: assignments[0].Id})
	if status.Code(err) != codes.NotFound {
		t.Errorf("expected NotFound error for get role assignment, got %v", err)
	}
	// deleting it again should fail
	_, err = server.DeleteRoleAssignment(ctx, &gen.DeleteRoleAssignmentRequest{Id: assignments[0].Id})
	if status.Code(err) != codes.NotFound {
		t.Errorf("expected NotFound error for delete role assignment, got %v", err)
	}
}

type messageWithID interface {
	proto.Message
	GetId() string
}

func checkNilIfErrored[V any](t *testing.T, v *V, err error) {
	t.Helper()
	if err != nil && v != nil {
		t.Errorf("expected nil return because of error %v, found value %v", err, v)
	}
	if err == nil && v == nil {
		t.Error("expected non-nil return but got nil")
	}
}

func shuffledPermissions(n int) []string {
	source := rand.Perm(n)
	perms := make([]string, n)
	for i, p := range source {
		perms[i] = fmt.Sprintf("perm-%02d", p)
	}
	return perms
}

func orderedPermissions(n int) []string {
	perms := make([]string, n)
	for i := range n {
		perms[i] = fmt.Sprintf("perm-%02d", i)
	}
	return perms
}
