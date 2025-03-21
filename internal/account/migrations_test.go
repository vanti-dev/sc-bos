package account

import (
	"context"
	"embed"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/ncruces/go-sqlite3"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

//go:embed testdata/*
var testdata embed.FS

func TestMigrations(t *testing.T) {
	type testCase struct {
		dumpFile       string
		expectAccounts []*gen.Account
		expectRoles    []*gen.Role
		expectAssigns  []*gen.RoleAssignment
	}

	cases := map[string]testCase{
		"schema_0001": {
			dumpFile: "testdata/from_schema_0001.sql",
			expectAccounts: []*gen.Account{
				{
					Id:          "2",
					DisplayName: "My Service 0",
					Description: "A service with 0 service credentials",
					Type:        gen.Account_SERVICE_ACCOUNT,
					CreateTime:  timestamppb.New(time.Date(2025, 3, 19, 12, 20, 0, 0, time.UTC)),
					Details: &gen.Account_ServiceDetails{
						ServiceDetails: &gen.ServiceAccount{
							ClientId: "2",
						},
					},
				},
				{
					Id:          "3",
					DisplayName: "My Service 1",
					Description: "A service with 1 service credential",
					Type:        gen.Account_SERVICE_ACCOUNT,
					CreateTime:  timestamppb.New(time.Date(2025, 3, 19, 12, 21, 0, 0, time.UTC)),
					Details: &gen.Account_ServiceDetails{
						ServiceDetails: &gen.ServiceAccount{
							ClientId: "3",
						},
					},
				},
				{
					Id:          "4",
					DisplayName: "My Service 2",
					Description: "A service with 2 service credentials",
					Type:        gen.Account_SERVICE_ACCOUNT,
					CreateTime:  timestamppb.New(time.Date(2025, 3, 19, 12, 22, 0, 0, time.UTC)),
					Details: &gen.Account_ServiceDetails{
						ServiceDetails: &gen.ServiceAccount{
							ClientId:              "4",
							PreviousSecretExpires: timestamppb.New(time.Date(2025, 3, 19, 12, 30, 0, 0, time.UTC)),
						},
					},
				},
				{
					Id:          "6",
					DisplayName: "User With Password",
					Type:        gen.Account_USER_ACCOUNT,
					CreateTime:  timestamppb.New(time.Date(2025, 3, 19, 12, 23, 0, 0, time.UTC)),
					Details: &gen.Account_UserDetails{
						UserDetails: &gen.UserAccount{
							Username:    "userpassword",
							HasPassword: true,
						},
					},
				},
				{
					Id:          "7",
					DisplayName: "User Without Password",
					Type:        gen.Account_USER_ACCOUNT,
					CreateTime:  timestamppb.New(time.Date(2025, 3, 19, 12, 24, 0, 0, time.UTC)),
					Details: &gen.Account_UserDetails{
						UserDetails: &gen.UserAccount{
							Username: "usernopassword",
						},
					},
				},
			},
			expectRoles: []*gen.Role{
				{
					Id:            "1",
					DisplayName:   "My Role",
					PermissionIds: []string{"account:read", "account:write"},
				},
			},
			expectAssigns: []*gen.RoleAssignment{
				{
					Id:        "1",
					AccountId: "6",
					RoleId:    "1",
					Scope: &gen.RoleAssignment_Scope{
						ResourceType: gen.RoleAssignment_ZONE,
						Resource:     "foo",
					},
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			dump, err := testdata.ReadFile(tc.dumpFile)
			if err != nil {
				t.Fatalf("could not read dump file %q: %v", tc.dumpFile, err)
			}
			dbPath := restoreDump(t, string(dump))

			ctx := context.Background()
			logger := testLogger(t)
			store, err := OpenStore(ctx, dbPath, logger)
			if err != nil {
				t.Fatalf("OpenStore: %v", err)
			}
			server := NewServer(store, logger)

			accounts := listAccounts(t, server)
			roles := listRoles(t, server)
			roleAssignments := listRoleAssignments(t, server)

			diff := cmp.Diff(tc.expectAccounts, accounts, protocmp.Transform())
			if diff != "" {
				t.Errorf("accounts mismatch (-want +got):\n%s", diff)
			}
			diff = cmp.Diff(tc.expectRoles, roles, protocmp.Transform())
			if diff != "" {
				t.Errorf("roles mismatch (-want +got):\n%s", diff)
			}
			diff = cmp.Diff(tc.expectAssigns, roleAssignments, protocmp.Transform())
			if diff != "" {
				t.Errorf("role assignments mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func restoreDump(t *testing.T, dump string) (dbPath string) {
	t.Helper()
	dir := t.TempDir()
	dbPath = filepath.Join(dir, "db.sqlite")

	conn, err := sqlite3.Open(dbPath)
	if err != nil {
		t.Fatalf("sqlite3.Open: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Fatalf("conn.Close: %v", err)
		}
	}()

	err = conn.Exec(dump)
	if err != nil {
		t.Fatalf("conn.Exec: %v", err)
	}
	return dbPath
}

func listAccounts(t *testing.T, server *Server) []*gen.Account {
	t.Helper()
	ctx := context.Background()
	var nextPageToken string
	var accounts []*gen.Account
	for {
		resp, err := server.ListAccounts(ctx, &gen.ListAccountsRequest{PageToken: nextPageToken})
		if err != nil {
			t.Fatalf("ListAccounts: %v", err)
		}
		accounts = append(accounts, resp.Accounts...)
		nextPageToken = resp.NextPageToken
		if nextPageToken == "" {
			break
		}
	}
	return accounts
}

func listRoles(t *testing.T, server *Server) []*gen.Role {
	t.Helper()
	ctx := context.Background()
	var nextPageToken string
	var roles []*gen.Role
	for {
		resp, err := server.ListRoles(ctx, &gen.ListRolesRequest{PageToken: nextPageToken})
		if err != nil {
			t.Fatalf("ListRoles: %v", err)
		}
		roles = append(roles, resp.Roles...)
		nextPageToken = resp.NextPageToken
		if nextPageToken == "" {
			break
		}
	}
	return roles
}

func listRoleAssignments(t *testing.T, server *Server) []*gen.RoleAssignment {
	t.Helper()
	ctx := context.Background()
	var nextPageToken string
	var roleAssignments []*gen.RoleAssignment
	for {
		resp, err := server.ListRoleAssignments(ctx, &gen.ListRoleAssignmentsRequest{PageToken: nextPageToken})
		if err != nil {
			t.Fatalf("ListRoleAssignments: %v", err)
		}
		roleAssignments = append(roleAssignments, resp.RoleAssignments...)
		nextPageToken = resp.NextPageToken
		if nextPageToken == "" {
			break
		}
	}
	return roleAssignments
}
