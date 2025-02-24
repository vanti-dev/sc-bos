package account

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"
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
	store := NewMemoryStore(zap.NewNop())
	server := NewServer(store, zap.NewNop())

	createAccount := func(req *gen.CreateAccountRequest) *gen.Account {
		t.Helper()
		res, err := server.CreateAccount(context.Background(), req)
		checkNilIfErrored(t, res, err)
		if err != nil {
			t.Fatalf("failed to create account: %v", err)
		}
		return res
	}
	failCreateAccount := func(req *gen.CreateAccountRequest, code codes.Code) {
		t.Helper()
		res, err := server.CreateAccount(context.Background(), req)
		checkNilIfErrored(t, res, err)
		if status.Code(err) != code {
			t.Errorf("expected error code %v, got %v", code, status.Code(err))
		}
	}
	checkAccount := func(account, expect *gen.Account) {
		t.Helper()
		diff := cmp.Diff(expect, account,
			protocmp.Transform(),
			protocmp.IgnoreFields(&gen.Account{}, "id", "create_time"),
		)
		if diff != "" {
			t.Errorf("unexpected provided account value (-got +want):\n%s", diff)
		}

		// also retrieve using GetAccount and check it matches
		id := account.Id
		expect = proto.Clone(expect).(*gen.Account)
		expect.Id = id
		expect.CreateTime = account.CreateTime
		account, err := server.GetAccount(context.Background(), &gen.GetAccountRequest{Id: id})
		checkNilIfErrored(t, account, err)
		if err != nil {
			t.Fatalf("failed to get account %q: %v", id, err)
		}
		diff = cmp.Diff(expect, account, protocmp.Transform())
		if diff != "" {
			t.Errorf("unexpected retrieved account value (-got +want):\n%s", diff)
		}
	}

	user1 := createAccount(&gen.CreateAccountRequest{
		Account: &gen.Account{
			Kind:        gen.Account_USER_ACCOUNT,
			DisplayName: "User 1",
			Username:    "user1",
		},
	})
	checkAccount(user1, &gen.Account{
		Kind:        gen.Account_USER_ACCOUNT,
		DisplayName: "User 1",
		Username:    "user1",
	})
	user2 := createAccount(&gen.CreateAccountRequest{
		Account: &gen.Account{
			Kind:        gen.Account_USER_ACCOUNT,
			DisplayName: "User 2",
			Username:    "user2",
		},
		Password: "user2Password",
	})
	checkAccount(user2, &gen.Account{
		Kind:        gen.Account_USER_ACCOUNT,
		DisplayName: "User 2",
		Username:    "user2",
	})
	service := createAccount(&gen.CreateAccountRequest{
		Account: &gen.Account{
			Kind:        gen.Account_SERVICE_ACCOUNT,
			DisplayName: "Service",
		},
	})
	checkAccount(service, &gen.Account{
		Kind:        gen.Account_SERVICE_ACCOUNT,
		DisplayName: "Service",
	})

	// missing account kind
	failCreateAccount(&gen.CreateAccountRequest{
		Account: &gen.Account{
			DisplayName: "Missing Kind",
		},
	}, codes.InvalidArgument)
	// missing display name
	failCreateAccount(&gen.CreateAccountRequest{
		Account: &gen.Account{
			Kind:     gen.Account_USER_ACCOUNT,
			Username: "foo",
		},
	}, codes.InvalidArgument)
	// missing username for user account
	failCreateAccount(&gen.CreateAccountRequest{
		Account: &gen.Account{
			Kind:        gen.Account_USER_ACCOUNT,
			DisplayName: "Missing Username",
		},
	}, codes.InvalidArgument)
	// username supplied for service account
	failCreateAccount(&gen.CreateAccountRequest{
		Account: &gen.Account{
			Kind:        gen.Account_SERVICE_ACCOUNT,
			DisplayName: "Service Account",
			Username:    "service",
		},
	}, codes.InvalidArgument)

}

// tests ordering and pagination of ListAccounts
func TestServer_ListAccounts(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore(zap.NewNop())
	server := NewServer(store, zap.NewNop())

	createAccount := func(kind gen.Account_Kind, username, displayName string) string {
		t.Helper()
		res, err := server.CreateAccount(ctx, &gen.CreateAccountRequest{
			Account: &gen.Account{
				Kind:        kind,
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
				Kind:        gen.Account_USER_ACCOUNT,
				Username:    username,
				DisplayName: displayName,
			})
		} else {
			// make it a service account
			id := createAccount(gen.Account_SERVICE_ACCOUNT, "", displayName)
			expected = append(expected, &gen.Account{
				Id:          id,
				Kind:        gen.Account_SERVICE_ACCOUNT,
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
				Kind:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{},
			},
			expected: &gen.Account{
				Kind:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
		},
		"kind_change_prohibited": {
			initial: &gen.Account{
				Kind:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					Kind: gen.Account_SERVICE_ACCOUNT,
				},
			},
			expected: &gen.Account{
				Kind:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
			code: codes.InvalidArgument,
		},
		"same_kind_allowed": {
			initial: &gen.Account{
				Kind:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					Kind: gen.Account_USER_ACCOUNT,
				},
			},
			expected: &gen.Account{
				Kind:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
		},
		"update_display_name": {
			initial: &gen.Account{
				Kind:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					DisplayName: "Service MODIFIED",
				},
			},
			expected: &gen.Account{
				Kind:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service MODIFIED",
			},
		},
		"update_username": {
			initial: &gen.Account{
				Kind:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					Username: "user1-modified",
				},
			},
			expected: &gen.Account{
				Kind:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1-modified",
			},
		},
		"update_username_service_account": {
			initial: &gen.Account{
				Kind:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					Username: "username",
				},
			},
			expected: &gen.Account{
				Kind:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service",
			},
			code: codes.FailedPrecondition,
		},
		"update_display_name_empty": {
			initial: &gen.Account{
				Kind:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					DisplayName: "",
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"display_name"}},
			},
			expected: &gen.Account{
				Kind:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service",
			},
			code: codes.InvalidArgument,
		},
		"update_username_empty": {
			initial: &gen.Account{
				Kind:        gen.Account_USER_ACCOUNT,
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
				Kind:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
			code: codes.InvalidArgument,
		},
		"invalid_update_mask": {
			initial: &gen.Account{
				Kind:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
			update: &gen.UpdateAccountRequest{
				Account:    &gen.Account{},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"foo"}},
			},
			expected: &gen.Account{
				Kind:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
			code: codes.InvalidArgument,
		},
		"wildcard_update_mask": {
			initial: &gen.Account{
				Kind:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1",
				Username:    "user1",
			},
			update: &gen.UpdateAccountRequest{
				Account: &gen.Account{
					Kind:        gen.Account_USER_ACCOUNT,
					DisplayName: "User 1 MODIFIED",
					Username:    "user1-modified",
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"*"}},
			},
			expected: &gen.Account{
				Kind:        gen.Account_USER_ACCOUNT,
				DisplayName: "User 1 MODIFIED",
				Username:    "user1-modified",
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

// tests that uniqueness of usernames is enforced
func TestServer_Account_Username(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore(zap.NewNop())
	server := NewServer(store, zap.NewNop())

	_, err := server.CreateAccount(ctx, &gen.CreateAccountRequest{
		Account: &gen.Account{
			Kind:        gen.Account_USER_ACCOUNT,
			DisplayName: "User 1",
			Username:    "user1",
		},
	})
	if err != nil {
		t.Fatalf("failed to create account1: %v", err)
	}
	user2, err := server.CreateAccount(ctx, &gen.CreateAccountRequest{
		Account: &gen.Account{
			Kind:        gen.Account_USER_ACCOUNT,
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
			Kind:        gen.Account_USER_ACCOUNT,
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
			Kind:        gen.Account_SERVICE_ACCOUNT,
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
	}

	var creds []*gen.ServiceCredential
	expiry := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := range maxServiceCredentialsPerAccount {
		cred, err := server.CreateServiceCredential(ctx, &gen.CreateServiceCredentialRequest{
			ServiceCredential: &gen.ServiceCredential{
				AccountId:  account.Id,
				Title:      fmt.Sprintf("Credential %d", i),
				ExpireTime: timestamppb.New(expiry),
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
			AccountId: account.Id,
			Title:     "Credential",
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
