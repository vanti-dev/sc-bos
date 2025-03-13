package account

import (
	"context"
	"errors"
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
		"display_name_long": {
			req: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_SERVICE_ACCOUNT,
					DisplayName: strings.Repeat("a", maxDisplayNameLength),
				},
			},
			expect: &gen.Account{
				Type:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: strings.Repeat("a", maxDisplayNameLength),
			},
		},
		"display_name_too_long": {
			req: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_SERVICE_ACCOUNT,
					DisplayName: strings.Repeat("a", maxDisplayNameLength+1),
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
		"username_long": {
			req: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 1",
					Username:    strings.Repeat("a", maxUsernameLength),
				},
			},
			expect: &gen.Account{
				Type:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    strings.Repeat("a", maxUsernameLength),
			},
		},
		"username_too_long": {
			req: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 1",
					Username:    strings.Repeat("a", maxUsernameLength+1),
				},
			},
			code: codes.InvalidArgument,
		},
		"description_short": {
			req: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_SERVICE_ACCOUNT,
					DisplayName: "Service",
					Description: "a description",
				},
			},
			expect: &gen.Account{
				Type:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service",
				Description: "a description",
			},
		},
		"description_long": {
			req: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_SERVICE_ACCOUNT,
					DisplayName: "Service",
					Description: strings.Repeat("a", maxDescriptionLength),
				},
			},
			expect: &gen.Account{
				Type:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service",
				Description: strings.Repeat("a", maxDescriptionLength),
			},
		},
		"description_too_long": {
			req: &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_SERVICE_ACCOUNT,
					DisplayName: "Service",
					Description: strings.Repeat("a", maxDescriptionLength+1),
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
		"nil": {
			req: &gen.CreateAccountRequest{
				Account: nil,
			},
			code: codes.InvalidArgument,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			logger := testLogger(t)
			store := NewMemoryStore(logger)
			server := NewServer(store, logger)

			for _, other := range tc.others {
				_, err := server.CreateAccount(context.Background(), &gen.CreateAccountRequest{Account: other})
				if err != nil {
					t.Fatalf("failed to create other account: %v", err)
				}
			}

			res, err := server.CreateAccount(context.Background(), tc.req)
			checkNilIfErrored(t, res, err)
			if status.Code(err) != tc.code {
				t.Errorf("expected error with code %v, got %v", tc.code, err)
			}
			diff := cmp.Diff(tc.expect, res,
				protocmp.Transform(),
				protocmp.IgnoreFields(&gen.Account{}, "id", "create_time"),
			)
			if diff != "" {
				t.Errorf("unexpected provided account value (-want +got):\n%s", diff)
			}

			// also retrieve using GetAccount and check it matches
			if res != nil {
				id := res.Id
				var expect *gen.Account
				if tc.expect != nil {
					expect = proto.Clone(tc.expect).(*gen.Account)
					expect.Id = id
					expect.CreateTime = res.CreateTime
				}
				account, err := server.GetAccount(context.Background(), &gen.GetAccountRequest{Id: id})
				checkNilIfErrored(t, account, err)
				if err != nil {
					t.Fatalf("failed to get account %q: %v", id, err)
				}
				diff = cmp.Diff(expect, account, protocmp.Transform())
				if diff != "" {
					t.Errorf("unexpected retrieved account value (-want +got):\n%s", diff)
				}
			}
		})
	}
}

// tests ordering and pagination of ListAccounts
func TestServer_ListAccounts(t *testing.T) {
	ctx := context.Background()
	logger := testLogger(t)
	store := NewMemoryStore(logger)
	server := NewServer(store, logger)

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

		if res.TotalSize != numAccounts {
			t.Errorf("expected total size %d, got %d", numAccounts, res.TotalSize)
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
		others   []*gen.Account // other accounts that should be created before the test account
		initial  *gen.Account   // initial account to create
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
		"update_description_implicit": {
			initial: &gen.Account{
				Type:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					Description: "Description",
				},
			},
			expected: &gen.Account{
				Type:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service",
				Description: "Description",
			},
		},
		"update_description_explicit": {
			initial: &gen.Account{
				Type:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service",
				Description: "Before",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					Description: "After",
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"description"}},
			},
			expected: &gen.Account{
				Type:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service",
				Description: "After",
			},
		},
		"update_description_empty": {
			initial: &gen.Account{
				Type:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service",
				Description: "Before",
			},
			update: &gen.UpdateAccountRequest{
				Account:    &gen.Account{},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"description"}},
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
				Description: "Description",
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
		"username_conflict": {
			others: []*gen.Account{
				{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "Foo",
					Username:    "foo",
				},
			},
			initial: &gen.Account{
				Type:        gen.Account_USER_ACCOUNT,
				DisplayName: "Bar",
				Username:    "bar",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					Username: "foo",
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"username"}},
			},
			expected: &gen.Account{
				Type:        gen.Account_USER_ACCOUNT,
				DisplayName: "Bar",
				Username:    "bar",
			},
			code: codes.AlreadyExists,
		},
		"not_exist": {
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					Id:       "123",
					Username: "bar",
				},
			},
			code: codes.NotFound,
		},
		"invalid_id": {
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					Id:          "invalid",
					DisplayName: "Barbar",
				},
			},
			code: codes.NotFound,
		},
		"nil": {
			update: &gen.UpdateAccountRequest{
				Account: nil,
			},
			code: codes.InvalidArgument,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			logger := testLogger(t)
			store := NewMemoryStore(logger)
			server := NewServer(store, logger)

			for _, other := range tc.others {
				_, err := server.CreateAccount(ctx, &gen.CreateAccountRequest{Account: other})
				if err != nil {
					t.Fatalf("failed to create other account: %v", err)
				}
			}

			var (
				account  *gen.Account
				expected *gen.Account
			)
			update := proto.Clone(tc.update).(*gen.UpdateAccountRequest)
			if tc.initial != nil {
				var err error
				account, err = server.CreateAccount(ctx, &gen.CreateAccountRequest{Account: tc.initial})
				checkNilIfErrored(t, account, err)
				if err != nil {
					t.Fatalf("failed to create account: %v", err)
				}
				// inject ID, which is now known, into update and expected
				update.Account.Id = account.Id
			}
			id := update.GetAccount().GetId()
			if tc.expected != nil {
				expected = proto.Clone(tc.expected).(*gen.Account)
				expected.Id = id
				expected.CreateTime = account.CreateTime
			}

			updated, err := server.UpdateAccount(ctx, update)
			checkNilIfErrored(t, updated, err)
			if status.Code(err) != tc.code {
				t.Errorf("expected error with code %v, got %v", tc.code, err)
			}
			if updated != nil {
				diff := cmp.Diff(expected, updated, protocmp.Transform())
				if diff != "" {
					t.Errorf("unexpected updated account value (-want +got):\n%s", diff)
				}
			}

			if tc.expected != nil {
				// fetch again to check that the update was persisted
				account, err = server.GetAccount(ctx, &gen.GetAccountRequest{Id: id})
				checkNilIfErrored(t, account, err)
				if err != nil {
					t.Fatalf("failed to get account: %v", err)
				}
				diff := cmp.Diff(expected, account, protocmp.Transform())
				if diff != "" {
					t.Errorf("unexpected retrieved account value (-want +got):\n%s", diff)
				}
			}
		})
	}
}

var validUsernames = []string{
	"foo",
	"foo.bar",
	"foo-bar",
	"foo_bar",
	"foo123",
	"foo@example.com",
	strings.Repeat("a", maxUsernameLength),
}

var invalidUsernames = []string{
	"ab",
	"foo bar",
	"foo\nbar",
	"foo😀bar",
	"foo/bar",
	"foo\\bar",
	"foo:bar",
}

func TestServer_CreateAccount_Usernames(t *testing.T) {
	template := &gen.CreateAccountRequest{
		Account: &gen.Account{
			Type:        gen.Account_USER_ACCOUNT,
			DisplayName: "User",
		},
	}

	ctx := context.Background()
	logger := testLogger(t)
	store := NewMemoryStore(logger)
	server := NewServer(store, logger)
	for _, username := range validUsernames {
		req := proto.Clone(template).(*gen.CreateAccountRequest)
		req.Account.Username = username
		_, err := server.CreateAccount(ctx, req)
		if err != nil {
			t.Errorf("valid username %q failed: %v", username, err)
		}
	}
	for _, username := range invalidUsernames {
		req := proto.Clone(template).(*gen.CreateAccountRequest)
		req.Account.Username = username
		_, err := server.CreateAccount(ctx, req)
		if !errors.Is(err, ErrInvalidUsername) {
			t.Errorf("invalid username %q returned %v", username, err)
		}
	}
}

func TestServer_UpdateAccount_Usernames(t *testing.T) {
	ctx := context.Background()
	logger := testLogger(t)
	store := NewMemoryStore(logger)
	server := NewServer(store, logger)

	account, err := server.CreateAccount(ctx, &gen.CreateAccountRequest{
		Account: &gen.Account{
			Type:        gen.Account_USER_ACCOUNT,
			DisplayName: "User",
			Username:    "user",
		},
	})
	if err != nil {
		t.Fatalf("failed to set up test account: %v", err)
	}

	template := &gen.UpdateAccountRequest{
		Account: &gen.Account{
			Id: account.Id,
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"username"}},
	}
	for _, username := range validUsernames {
		req := proto.Clone(template).(*gen.UpdateAccountRequest)
		req.Account.Username = username
		_, err := server.UpdateAccount(ctx, req)
		if err != nil {
			t.Errorf("valid username %q failed: %v", username, err)
		}
	}
	for _, username := range invalidUsernames {
		req := proto.Clone(template).(*gen.UpdateAccountRequest)
		req.Account.Username = username
		_, err := server.UpdateAccount(ctx, req)
		if !errors.Is(err, ErrInvalidUsername) {
			t.Errorf("invalid username %q returned %v", username, err)
		}
	}
}

func TestServer_DeleteAccount(t *testing.T) {
	ctx := context.Background()
	logger := testLogger(t)
	store := NewMemoryStore(logger)
	server := NewServer(store, logger)

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
			DisplayName:   "Foo Role",
			PermissionIds: []string{"foo"},
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
	// deleting an account again should fail
	_, err = server.DeleteAccount(ctx, &gen.DeleteAccountRequest{Id: user.Id})
	if status.Code(err) != codes.NotFound {
		t.Errorf("expected NotFound error when deleting user account twice, got %v", err)
	}
	// deleting an account again with allow_missing should succeed
	_, err = server.DeleteAccount(ctx, &gen.DeleteAccountRequest{Id: user.Id, AllowMissing: true})
	if err != nil {
		t.Errorf("failed to delete user account with allow_missing: %v", err)
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
	logger := testLogger(t)
	store := NewMemoryStore(logger)
	server := NewServer(store, logger)

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
	logger := testLogger(t)
	store := NewMemoryStore(logger)
	server := NewServer(store, logger)

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
			t.Errorf("unexpected service credentials (-want +got):\n%s", diff)
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
				t.Errorf("unexpected fetched service credential %q (-want +got):\n%s", cred.Id, diff)
			}
		}
	}

	var creds []*gen.ServiceCredential
	expiry := time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)
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

func TestServer_CreateServiceCredential(t *testing.T) {
	type testCase struct {
		req    *gen.CreateServiceCredentialRequest
		expect *gen.ServiceCredential
		code   codes.Code
	}

	cases := map[string]testCase{
		"no_expiry": {
			req: &gen.CreateServiceCredentialRequest{
				ServiceCredential: &gen.ServiceCredential{
					DisplayName: "Credential 1",
				},
			},
			expect: &gen.ServiceCredential{
				DisplayName: "Credential 1",
			},
		},
		"expiry": {
			req: &gen.CreateServiceCredentialRequest{
				ServiceCredential: &gen.ServiceCredential{
					DisplayName: "Credential 2",
					ExpireTime:  timestamppb.New(time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
			expect: &gen.ServiceCredential{
				DisplayName: "Credential 2",
				ExpireTime:  timestamppb.New(time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		"no_display_name": {
			req: &gen.CreateServiceCredentialRequest{
				ServiceCredential: &gen.ServiceCredential{},
			},
			code: codes.InvalidArgument,
		},
		"long_display_name": {
			req: &gen.CreateServiceCredentialRequest{
				ServiceCredential: &gen.ServiceCredential{
					DisplayName: strings.Repeat("a", maxDisplayNameLength),
				},
			},
			expect: &gen.ServiceCredential{
				DisplayName: strings.Repeat("a", maxDisplayNameLength),
			},
		},
		"display_name_too_long": {
			req: &gen.CreateServiceCredentialRequest{
				ServiceCredential: &gen.ServiceCredential{
					DisplayName: strings.Repeat("a", maxDisplayNameLength+1),
				},
			},
			code: codes.InvalidArgument,
		},
		"description": {
			req: &gen.CreateServiceCredentialRequest{
				ServiceCredential: &gen.ServiceCredential{
					DisplayName: "Credential 3",
					Description: "This is a description",
				},
			},
			expect: &gen.ServiceCredential{
				DisplayName: "Credential 3",
				Description: "This is a description",
			},
		},
		"nil": {
			req:  &gen.CreateServiceCredentialRequest{},
			code: codes.InvalidArgument,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			logger := testLogger(t)
			store := NewMemoryStore(logger)
			server := NewServer(store, logger)

			account, err := server.CreateAccount(context.Background(), &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_SERVICE_ACCOUNT,
					DisplayName: "Service",
				},
			})
			if err != nil {
				t.Fatalf("failed to create account: %v", err)
			}

			req := proto.Clone(tc.req).(*gen.CreateServiceCredentialRequest)
			if req.ServiceCredential != nil {
				req.ServiceCredential.AccountId = account.Id
			}
			created, err := server.CreateServiceCredential(context.Background(), req)
			checkNilIfErrored(t, created, err)
			if status.Code(err) != tc.code {
				t.Errorf("expected error with code %v, got %v", tc.code, err)
			}

			diff := cmp.Diff(tc.expect, created,
				protocmp.Transform(),
				protocmp.IgnoreFields(&gen.ServiceCredential{}, "id", "create_time", "secret"),
			)
			if diff != "" {
				t.Errorf("unexpected provided service credential value (-want +got):\n%s", diff)
			}
		})
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
			logger := testLogger(t)
			store := NewMemoryStore(logger)
			server := NewServer(store, logger)

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
	// tests a sequence of role operations, checking that role lifecycles are handled correctly

	ctx := context.Background()
	logger := testLogger(t)
	store := NewMemoryStore(logger)
	server := NewServer(store, logger)

	const numRoles = 200
	const numPermissions = 10

	const description = "A role for testing"
	var roles []*gen.Role
	t.Log("CreateRole:")
	for i := range numRoles {
		displayName := fmt.Sprintf("Role %d", i)
		numPermissions := i % numPermissions // test roles with different numbers of permissions
		role, err := server.CreateRole(ctx, &gen.CreateRoleRequest{
			Role: &gen.Role{
				DisplayName: displayName,
				// supply the permissions shuffled to check they are returned in order instead
				PermissionIds: shuffledPermissions(numPermissions),
				Description:   description,
			},
		})
		if err != nil {
			t.Fatalf("failed to create role %q: %v", displayName, err)
		}
		expect := &gen.Role{
			Id:            role.Id,
			DisplayName:   displayName,
			PermissionIds: orderedPermissions(numPermissions),
			Description:   description,
		}
		diff := cmp.Diff(expect, role, protocmp.Transform())
		if diff != "" {
			t.Errorf("unexpected role value (-want +got):\n%s", diff)
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
			if res.TotalSize != numRoles {
				t.Errorf("expected total size %d, got %d", numRoles, res.TotalSize)
			}

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
	role.PermissionIds = append(role.PermissionIds, "000-new-permission") // should go at the beginning
	role.DisplayName += " MODIFIED"
	updated, err := server.UpdateRole(ctx, &gen.UpdateRoleRequest{
		Role: role,
	})
	slices.Sort(role.PermissionIds)
	checkNilIfErrored(t, role, err)
	if err != nil {
		t.Fatalf("failed to update role: %v", err)
	}
	diff := cmp.Diff(role, updated, protocmp.Transform())
	if diff != "" {
		t.Errorf("unexpected updated role value (-want +got):\n%s", diff)
	}
	// test that update is persisted
	updated, err = server.GetRole(ctx, &gen.GetRoleRequest{Id: role.Id})
	checkNilIfErrored(t, updated, err)
	if err != nil {
		t.Fatalf("failed to get updated role: %v", err)
	}
	diff = cmp.Diff(role, updated, protocmp.Transform())
	if diff != "" {
		t.Errorf("unexpected retrieved role value (-want +got):\n%s", diff)
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
	// deleting with allow_missing should succeed
	_, err = server.DeleteRole(ctx, &gen.DeleteRoleRequest{Id: role.Id, AllowMissing: true})
	if err != nil {
		t.Errorf("failed to delete role with allow_missing: %v", err)
	}
}

func TestServer_CreateRole(t *testing.T) {
	ctx := context.Background()

	type testCase struct {
		existing []*gen.Role
		role     *gen.Role
		err      error
	}

	cases := map[string]testCase{
		"missing_display_name": {
			role: &gen.Role{},
			err:  ErrInvalidDisplayName,
		},
		"display_name_too_long": {
			role: &gen.Role{DisplayName: strings.Repeat("a", maxDisplayNameLength+1)},
			err:  ErrInvalidDisplayName,
		},
		"long_display_name": {
			role: &gen.Role{DisplayName: strings.Repeat("a", maxDisplayNameLength)},
		},
		"short_display_name": {
			role: &gen.Role{DisplayName: "Role"},
		},
		"short_description": {
			role: &gen.Role{
				DisplayName: "Role",
				Description: "Role Description",
			},
		},
		"description_too_long": {
			role: &gen.Role{
				DisplayName: "Role",
				Description: strings.Repeat("a", maxDescriptionLength+1),
			},
			err: ErrInvalidDescription,
		},
		"long_description": {
			role: &gen.Role{
				DisplayName: "Role",
				Description: strings.Repeat("a", maxDescriptionLength),
			},
		},
		"single_permission": {
			role: &gen.Role{
				DisplayName:   "Role",
				PermissionIds: []string{"foo"},
			},
		},
		"multiple_permissions_ordered": {
			role: &gen.Role{
				DisplayName:   "Role",
				PermissionIds: []string{"bar", "baz", "foo"},
			},
		},
		"multiple_permissions_unordered": {
			role: &gen.Role{
				DisplayName:   "Role",
				PermissionIds: []string{"foo", "bar", "baz"},
			},
		},
		"duplicate_permissions": {
			role: &gen.Role{
				DisplayName:   "Role",
				PermissionIds: []string{"foo", "bar", "foo"},
			},
		},
		"nil": {
			role: nil,
			err:  ErrResourceMissing,
		},
		"display_name_conflict": {
			existing: []*gen.Role{
				{DisplayName: "Role 1"},
			},
			role: &gen.Role{DisplayName: "Role 1"},
			err:  ErrRoleDisplayNameExists,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			logger := testLogger(t)
			store := NewMemoryStore(logger)
			server := NewServer(store, logger)

			for _, existing := range tc.existing {
				_, err := server.CreateRole(ctx, &gen.CreateRoleRequest{Role: existing})
				if err != nil {
					t.Fatalf("failed to create existing role: %v", err)
				}
			}

			created, err := server.CreateRole(ctx, &gen.CreateRoleRequest{Role: tc.role})
			if !errors.Is(err, tc.err) {
				t.Errorf("expected error %v, got %v", tc.err, err)
			}

			if created != nil {
				expect := proto.Clone(tc.role).(*gen.Role)
				expect.Id = created.Id
				// normalise the permission IDs
				slices.Sort(expect.PermissionIds)
				expect.PermissionIds = slices.Compact(expect.PermissionIds)

				diff := cmp.Diff(expect, created, protocmp.Transform())
				if diff != "" {
					t.Errorf("unexpected created role value (-want +got):\n%s", diff)
				}

				fetched, err := server.GetRole(ctx, &gen.GetRoleRequest{Id: created.Id})
				if err != nil {
					t.Errorf("failed to fetch created role: %v", err)
				}

				diff = cmp.Diff(expect, fetched, protocmp.Transform())
				if diff != "" {
					t.Errorf("unexpected fetched role value (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestServer_UpdateRole(t *testing.T) {
	ctx := context.Background()

	type testCase struct {
		initial  *gen.Role
		others   []*gen.Role // other roles to create before the update
		update   *gen.UpdateRoleRequest
		expected *gen.Role
		code     codes.Code
	}
	cases := map[string]testCase{
		"empty update": {
			initial: &gen.Role{
				DisplayName:   "Role 1",
				PermissionIds: []string{"foo", "bar"},
			},
			update: &gen.UpdateRoleRequest{
				Role: &gen.Role{},
			},
			expected: &gen.Role{
				DisplayName:   "Role 1",
				PermissionIds: []string{"bar", "foo"},
			},
		},
		"update_display_name_implicit": {
			initial: &gen.Role{
				DisplayName:   "Role 1",
				PermissionIds: []string{"foo", "bar"},
			},
			update: &gen.UpdateRoleRequest{
				Role: &gen.Role{
					DisplayName: "Role 1 MODIFIED",
				},
			},
			expected: &gen.Role{
				DisplayName:   "Role 1 MODIFIED",
				PermissionIds: []string{"bar", "foo"},
			},
		},
		"update_display_name_explicit": {
			initial: &gen.Role{
				DisplayName:   "Role 1",
				PermissionIds: []string{"foo", "bar"},
			},
			update: &gen.UpdateRoleRequest{
				Role: &gen.Role{
					DisplayName:   "Role 1 MODIFIED",
					PermissionIds: []string{"foo2", "bar2"},
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"display_name"}},
			},
			expected: &gen.Role{
				DisplayName:   "Role 1 MODIFIED",
				PermissionIds: []string{"bar", "foo"},
			},
		},
		"update_display_name_clear": {
			initial: &gen.Role{
				DisplayName: "Role 1",
			},
			update: &gen.UpdateRoleRequest{
				Role:       &gen.Role{},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"display_name"}},
			},
			expected: &gen.Role{
				DisplayName: "Role 1",
			},
			code: codes.InvalidArgument,
		},
		"update_display_name_too_long": {
			initial: &gen.Role{
				DisplayName: "Role 1",
			},
			update: &gen.UpdateRoleRequest{
				Role: &gen.Role{
					DisplayName: strings.Repeat("a", maxDisplayNameLength+1),
				},
			},
			expected: &gen.Role{
				DisplayName: "Role 1",
			},
			code: codes.InvalidArgument,
		},
		"update_permissions_implicit": {
			initial: &gen.Role{
				DisplayName:   "Role 1",
				PermissionIds: []string{"foo", "bar"},
			},
			update: &gen.UpdateRoleRequest{
				Role: &gen.Role{
					PermissionIds: []string{"foo2", "bar2"},
				},
			},
			expected: &gen.Role{
				DisplayName:   "Role 1",
				PermissionIds: []string{"bar2", "foo2"},
			},
		},
		"update_permissions_explicit": {
			initial: &gen.Role{
				DisplayName:   "Role 1",
				PermissionIds: []string{"foo", "bar"},
			},
			update: &gen.UpdateRoleRequest{
				Role: &gen.Role{
					DisplayName:   "Role 1 MODIFIED",
					PermissionIds: []string{"foo2", "bar2"},
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"permission_ids"}},
			},
			expected: &gen.Role{
				DisplayName:   "Role 1",
				PermissionIds: []string{"bar2", "foo2"},
			},
		},
		"update_permissions_empty": {
			initial: &gen.Role{
				DisplayName:   "Role 1",
				PermissionIds: []string{"foo", "bar"},
			},
			update: &gen.UpdateRoleRequest{
				Role: &gen.Role{
					PermissionIds: nil,
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"permission_ids"}},
			},
			expected: &gen.Role{
				DisplayName:   "Role 1",
				PermissionIds: nil,
			},
		},
		"update_description_implicit": {
			initial: &gen.Role{
				DisplayName: "Role 1",
			},
			update: &gen.UpdateRoleRequest{
				Role: &gen.Role{
					Description: "A role for testing",
				},
			},
			expected: &gen.Role{
				DisplayName: "Role 1",
				Description: "A role for testing",
			},
		},
		"update_description_explicit": {
			initial: &gen.Role{
				DisplayName: "Role 1",
			},
			update: &gen.UpdateRoleRequest{
				Role: &gen.Role{
					Description: "A role for testing",
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"description"}},
			},
			expected: &gen.Role{
				DisplayName: "Role 1",
				Description: "A role for testing",
			},
		},
		"update_description_clear": {
			initial: &gen.Role{
				DisplayName: "Role 1",
				Description: "A role for testing",
			},
			update: &gen.UpdateRoleRequest{
				Role:       &gen.Role{},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"description"}},
			},
			expected: &gen.Role{
				DisplayName: "Role 1",
			},
		},
		"update_description_invalid": {
			initial: &gen.Role{
				DisplayName: "Role 1",
			},
			update: &gen.UpdateRoleRequest{
				Role: &gen.Role{Description: strings.Repeat("a", maxDescriptionLength+1)},
			},
			expected: &gen.Role{
				DisplayName: "Role 1",
			},
			code: codes.InvalidArgument,
		},
		"nil": {
			initial:  &gen.Role{DisplayName: "Role 1"},
			update:   &gen.UpdateRoleRequest{Role: nil},
			expected: &gen.Role{DisplayName: "Role 1"},
			code:     codes.InvalidArgument,
		},
		"display_name_conflict": {
			initial: &gen.Role{DisplayName: "Role 1"},
			others:  []*gen.Role{{DisplayName: "Role 2"}},
			update: &gen.UpdateRoleRequest{
				Role: &gen.Role{DisplayName: "Role 2"},
			},
			expected: &gen.Role{DisplayName: "Role 1"},
			code:     codes.AlreadyExists,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			logger := testLogger(t)
			store := NewMemoryStore(logger)
			server := NewServer(store, logger)

			role, err := server.CreateRole(ctx, &gen.CreateRoleRequest{
				Role: tc.initial,
			})
			checkNilIfErrored(t, role, err)
			if err != nil {
				t.Fatalf("failed to create role: %v", err)
			}

			for _, other := range tc.others {
				_, err := server.CreateRole(ctx, &gen.CreateRoleRequest{
					Role: other,
				})
				if err != nil {
					t.Fatalf("failed to create other role: %v", err)
				}
			}

			// inject ID, which is now known, into update and expected
			id := role.Id
			update := proto.Clone(tc.update).(*gen.UpdateRoleRequest)
			if update.Role != nil {
				update.Role.Id = id
			}
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
					t.Errorf("unexpected updated role value (-want +got):\n%s", diff)
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
				t.Errorf("unexpected retrieved role value (-want +got):\n%s", diff)
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
			t.Errorf("unexpected page %d contents (-want +got):\n%s", i, diff)
		}
	}

	if len(expect) > 0 {
		t.Errorf("expected %d more items than received", len(expect))
	}
}

func TestServer_RoleAssignments(t *testing.T) {
	ctx := context.Background()
	logger := testLogger(t)
	store := NewMemoryStore(logger)
	server := NewServer(store, logger)

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
				DisplayName:   displayName,
				PermissionIds: permissions,
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

				if int(res.TotalSize) != len(f.expect) {
					t.Errorf("expected total size %d, got %d", len(f.expect), res.TotalSize)
				}

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
	// deleting with allow_missing should succeed
	_, err = server.DeleteRoleAssignment(ctx, &gen.DeleteRoleAssignmentRequest{Id: assignments[0].Id, AllowMissing: true})
	if err != nil {
		t.Errorf("failed to delete role assignment with allow_missing: %v", err)
	}
}

func TestServer_CreateRoleAssignment(t *testing.T) {
	ctx := context.Background()

	type testCase struct {
		existing []*gen.RoleAssignment
		req      *gen.CreateRoleAssignmentRequest
		code     codes.Code
	}

	// will be replaced with IDs created during the test run
	const (
		accountPlaceholder = "ACCOUNT PLACEHOLDER"
		rolePlaceholder    = "ROLE PLACEHOLDER"
	)

	cases := map[string]testCase{
		"role_does_not_exist": {
			req: &gen.CreateRoleAssignmentRequest{
				RoleAssignment: &gen.RoleAssignment{
					AccountId: accountPlaceholder,
					RoleId:    "999",
				},
			},
			code: codes.NotFound,
		},
		"account_does_not_exist": {
			req: &gen.CreateRoleAssignmentRequest{
				RoleAssignment: &gen.RoleAssignment{
					AccountId: "999",
					RoleId:    rolePlaceholder,
				},
			},
			code: codes.NotFound,
		},
		"role_id_invalid": {
			req: &gen.CreateRoleAssignmentRequest{
				RoleAssignment: &gen.RoleAssignment{
					AccountId: accountPlaceholder,
					RoleId:    "invalid",
				},
			},
			code: codes.NotFound,
		},
		"account_id_invalid": {
			req: &gen.CreateRoleAssignmentRequest{
				RoleAssignment: &gen.RoleAssignment{
					AccountId: "invalid",
					RoleId:    rolePlaceholder,
				},
			},
			code: codes.NotFound,
		},
		"scope_provided": {
			req: &gen.CreateRoleAssignmentRequest{
				RoleAssignment: &gen.RoleAssignment{
					AccountId: accountPlaceholder,
					RoleId:    rolePlaceholder,
					Scope: &gen.RoleAssignment_Scope{
						ResourceType: gen.RoleAssignment_NAMED_RESOURCE,
						Resource:     "named-resource",
					},
				},
			},
			code: codes.OK,
		},
		"scope_not_provided": {
			req: &gen.CreateRoleAssignmentRequest{
				RoleAssignment: &gen.RoleAssignment{
					AccountId: accountPlaceholder,
					RoleId:    rolePlaceholder,
				},
			},
			code: codes.OK,
		},
		"scope_invalid_resource_type": {
			req: &gen.CreateRoleAssignmentRequest{
				RoleAssignment: &gen.RoleAssignment{
					AccountId: accountPlaceholder,
					RoleId:    rolePlaceholder,
					Scope: &gen.RoleAssignment_Scope{
						ResourceType: 500,
						Resource:     "invalid resource type",
					},
				},
			},
			code: codes.InvalidArgument,
		},
		"scope_empty_resource_string": {
			req: &gen.CreateRoleAssignmentRequest{
				RoleAssignment: &gen.RoleAssignment{
					AccountId: accountPlaceholder,
					RoleId:    rolePlaceholder,
					Scope: &gen.RoleAssignment_Scope{
						ResourceType: gen.RoleAssignment_NAMED_RESOURCE,
						Resource:     "",
					},
				},
			},
			code: codes.InvalidArgument,
		},
		"conflict": {
			existing: []*gen.RoleAssignment{
				{
					AccountId: accountPlaceholder,
					RoleId:    rolePlaceholder,
				},
			},
			req: &gen.CreateRoleAssignmentRequest{
				RoleAssignment: &gen.RoleAssignment{
					AccountId: accountPlaceholder,
					RoleId:    rolePlaceholder,
					Scope: &gen.RoleAssignment_Scope{
						ResourceType: gen.RoleAssignment_NAMED_RESOURCE,
						Resource:     "resource name",
					},
				},
			},
			code: codes.AlreadyExists,
		},
		"nil": {
			req:  &gen.CreateRoleAssignmentRequest{},
			code: codes.InvalidArgument,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			logger := testLogger(t)
			store := NewMemoryStore(logger)
			server := NewServer(store, logger)

			// Create a valid account and role for positive test cases
			account, err := server.CreateAccount(ctx, &gen.CreateAccountRequest{
				Account: &gen.Account{
					Type:        gen.Account_USER_ACCOUNT,
					DisplayName: "User",
					Username:    "user",
				},
			})
			if err != nil {
				t.Fatalf("failed to create account: %v", err)
			}
			role, err := server.CreateRole(ctx, &gen.CreateRoleRequest{
				Role: &gen.Role{
					DisplayName:   "Role",
					PermissionIds: []string{"foo"},
				},
			})
			if err != nil {
				t.Fatalf("failed to create role: %v", err)
			}

			substitute := func(ra *gen.RoleAssignment) *gen.RoleAssignment {
				if ra == nil {
					return nil
				}
				// Inject valid IDs into the request if needed
				cloned := proto.Clone(ra).(*gen.RoleAssignment)
				if cloned.RoleId == rolePlaceholder {
					cloned.RoleId = role.Id
				}
				if cloned.AccountId == accountPlaceholder {
					cloned.AccountId = account.Id
				}
				return cloned
			}

			// Create any existing role assignments
			for _, existing := range tc.existing {
				existing = substitute(existing)

				_, err = server.CreateRoleAssignment(ctx, &gen.CreateRoleAssignmentRequest{
					RoleAssignment: existing,
				})
				if err != nil {
					t.Fatalf("failed to create existing role assignment: %v", err)
				}
			}

			req := proto.Clone(tc.req).(*gen.CreateRoleAssignmentRequest)
			req.RoleAssignment = substitute(req.RoleAssignment)

			created, err := server.CreateRoleAssignment(ctx, req)
			if status.Code(err) != tc.code {
				t.Errorf("expected error code %v, got %v", tc.code, err)
			}

			if created != nil {
				expected := proto.Clone(req.RoleAssignment).(*gen.RoleAssignment)
				expected.Id = created.Id

				diff := cmp.Diff(expected, created, protocmp.Transform())
				if diff != "" {
					t.Errorf("unexpected created role assignment value (-want +got):\n%s", diff)
				}
			}
		})
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

func testLogger(t *testing.T) *zap.Logger {
	t.Helper()
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	return logger
}
