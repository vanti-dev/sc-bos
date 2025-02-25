package account

import (
	"strconv"
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/sc-bos/internal/account/queries"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

func parseID(id string) (int64, bool) {
	parsed, err := strconv.ParseInt(id, 10, 64)
	if err != nil || parsed <= 0 {
		return 0, false
	}
	return parsed, true
}

func formatID(id int64) string {
	return strconv.FormatInt(id, 10)
}

func formatPageToken(nextID int64) string {
	if nextID == 0 {
		return ""
	}
	return formatID(nextID)
}

func parsePageToken(token string) (int64, bool) {
	if token == "" {
		return 0, true
	}
	return parseID(token)
}

func accountToProto(account queries.Account) *gen.Account {
	var (
		username string
		// default to ACCOUNT_KIND_UNSPECIFIED
		kind = gen.Account_Kind(gen.Account_Kind_value[account.Kind])
	)
	if account.Username.Valid {
		username = account.Username.String
	}

	return &gen.Account{
		Id:          formatID(account.ID),
		Username:    username,
		DisplayName: account.DisplayName,
		Kind:        kind,
		CreateTime:  timestamppb.New(account.CreateTime),
	}
}

func roleToProto(role queries.Role, permissions []string) *gen.Role {
	return &gen.Role{
		Id:          formatID(role.ID),
		Title:       role.Name,
		Permissions: permissions,
	}
}

func serviceCredentialToProto(cred queries.ServiceCredential, secret string) *gen.ServiceCredential {
	return &gen.ServiceCredential{
		Id:         formatID(cred.ID),
		Title:      cred.Title,
		CreateTime: timestamppb.New(cred.CreateTime),
		ExpireTime: timestamppb.New(cred.ExpireTime.Time),
		Secret:     secret,
	}
}

func roleAssignmentToProto(assignment queries.RoleAssignment) *gen.RoleAssignment {
	ra := &gen.RoleAssignment{
		Id:        formatID(assignment.ID),
		AccountId: formatID(assignment.AccountID),
		RoleId:    formatID(assignment.RoleID),
	}
	if assignment.ScopeKind.Valid && assignment.ScopeResource.Valid {
		ra.Scope = &gen.RoleAssignment_Scope{
			// defaults to RESOURCE_KIND_UNSPECIFIED
			ResourceKind: gen.RoleAssignment_ResourceKind(gen.RoleAssignment_ResourceKind_value[assignment.ScopeKind.String]),
			Resource:     assignment.ScopeResource.String,
		}
	}
	return ra
}

// in SQL queries that return a list of permissions per row, they are joined comma-separated
func splitPermissions(permissions string) []string {
	return strings.Split(permissions, ",")
}
