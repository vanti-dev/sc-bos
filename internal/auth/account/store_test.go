package account

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

func TestServer(t *testing.T) {
	store := NewMemoryStore(zap.NewNop())
	server := NewServer(store, zap.NewNop())

	type testCase struct {
		input        *gen.Account
		expect       *gen.Account
		expectStatus codes.Code
	}
	cases := map[string]testCase{
		"basic_user": {
			input: &gen.Account{
				Kind:        gen.Account_USER_ACCOUNT,
				DisplayName: "Foo User",
				Username:    "foo@example.com",
			},
		},
		"user_with_roles": {
			input: &gen.Account{
				Kind:        gen.Account_USER_ACCOUNT,
				Username:    "user_with_roles@example.com",
				DisplayName: "Foo User with roles",
			},
		},
		"service_account": {
			input: &gen.Account{
				Kind:        gen.Account_SERVICE_ACCOUNT,
				DisplayName: "Service Account",
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			created, err := server.CreateAccount(ctx, &gen.CreateAccountRequest{Account: tc.input})
			checkNilIfErrored(t, created, err)
			if code := status.Code(err); code != tc.expectStatus {
				t.Errorf("expected status %v, got error %v", tc.expectStatus, err)
			}
			if err != nil {
				return
			}
			if created.GetId() == "" {
				t.Error("CreateAccount returned an empty ID")
			}

			// ID is assigned by the store so we don't check it
			expect := tc.expect
			if expect == nil {
				// common case where we expect the input to be returned
				expect = tc.input
			}
			diff := cmp.Diff(expect, created,
				protocmp.Transform(),
				protocmp.IgnoreFields(new(gen.Account), "id", "create_time"),
			)
			if diff != "" {
				t.Errorf("CreateAccount returned incorrect Account (-want +got): %s", diff)
			}

			// Check that the account was actually stored
			got, err := server.GetAccount(ctx, &gen.GetAccountRequest{Id: created.Id})
			checkNilIfErrored(t, got, err)
			if err != nil {
				t.Errorf("error getting account %q: %v", created.Id, err)
			}
			expect = proto.Clone(expect).(*gen.Account)
			expect.Id = created.Id
			expect.CreateTime = created.CreateTime
			if diff := cmp.Diff(expect, got, protocmp.Transform()); diff != "" {
				t.Errorf("GetAccount returned incorrect Account (-want +got): %s", diff)
			}
		})
	}
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

func createTestRole(t *testing.T, server *Server, title string, permissions ...string) (id string) {
	t.Helper()
	role := &gen.Role{
		Title:       title,
		Permissions: permissions,
	}
	res, err := server.CreateRole(context.Background(), &gen.CreateRoleRequest{Role: role})
	if err != nil {
		t.Fatalf("failed to create role %q: %v", title, err)
	}
	return res.Id
}
