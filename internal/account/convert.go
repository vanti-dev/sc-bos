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

func accountToProto(account queries.Account) *gen.Account {
	converted := &gen.Account{
		Id:          formatID(account.ID),
		DisplayName: account.DisplayName,
		Type:        gen.Account_Type(gen.Account_Type_value[account.Type]), // default to ACCOUNT_TYPE_UNSPECIFIED
		CreateTime:  timestamppb.New(account.CreateTime),
	}
	if account.Username.Valid {
		converted.Username = account.Username.String
	}
	if account.Description.Valid {
		converted.Description = account.Description.String
	}
	return converted
}

func roleToProto(role queries.Role, permissions []string) *gen.Role {
	protoRole := &gen.Role{
		Id:            formatID(role.ID),
		DisplayName:   role.DisplayName,
		PermissionIds: permissions,
	}
	if role.Description.Valid {
		protoRole.Description = role.Description.String
	}
	return protoRole
}

func serviceCredentialToProto(cred queries.ServiceCredential, secret string) *gen.ServiceCredential {
	converted := &gen.ServiceCredential{
		Id:          formatID(cred.ID),
		DisplayName: cred.DisplayName,
		CreateTime:  timestamppb.New(cred.CreateTime),
		Secret:      secret,
	}
	if cred.ExpireTime.Valid {
		converted.ExpireTime = timestamppb.New(cred.ExpireTime.Time)
	}
	if cred.Description.Valid {
		converted.Description = cred.Description.String
	}
	return converted
}

func roleAssignmentToProto(assignment queries.RoleAssignment) *gen.RoleAssignment {
	ra := &gen.RoleAssignment{
		Id:        formatID(assignment.ID),
		AccountId: formatID(assignment.AccountID),
		RoleId:    formatID(assignment.RoleID),
	}
	if assignment.ScopeType.Valid && assignment.ScopeResource.Valid {
		ra.Scope = &gen.RoleAssignment_Scope{
			// defaults to RESOURCE_KIND_UNSPECIFIED
			ResourceType: gen.RoleAssignment_ResourceType(gen.RoleAssignment_ResourceType_value[assignment.ScopeType.String]),
			Resource:     assignment.ScopeResource.String,
		}
	}
	return ra
}

// in SQL queries that return a list of permissions per row, they are joined comma-separated
func splitPermissions(permissions string) []string {
	return strings.Split(permissions, ",")
}
