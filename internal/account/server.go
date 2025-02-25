package account

import (
	"context"
	"database/sql"
	"errors"
	"regexp"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vanti-dev/sc-bos/internal/account/queries"
	"github.com/vanti-dev/sc-bos/internal/database"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

var (
	ErrDatabase                  = status.Error(codes.Internal, "database internal error")
	ErrAccountNotFound           = status.Error(codes.NotFound, "account not found")
	ErrRoleNotFound              = status.Error(codes.NotFound, "role not found")
	ErrRoleAssignmentNotFound    = status.Error(codes.NotFound, "role assignment not found")
	ErrServiceCredentialNotFound = status.Error(codes.NotFound, "service credential not found")
	ErrInvalidAccountKind        = status.Error(codes.InvalidArgument, "invalid account kind")
	ErrMissingUsername           = status.Error(codes.InvalidArgument, "user account requires username")
	ErrMissingDisplayName        = status.Error(codes.InvalidArgument, "account requires display name")
	ErrUnexpectedUsernameCreate  = status.Error(codes.InvalidArgument, "service account cannot have username")
	ErrUnexpectedUsernameUpdate  = status.Error(codes.FailedPrecondition, "service account cannot have username")
	ErrUnexpectedServiceCreds    = status.Error(codes.FailedPrecondition, "user account cannot have service credentials")
	ErrServiceCredentialLimit    = status.Error(codes.ResourceExhausted, "too many service credentials")
	ErrUsernameExists            = status.Error(codes.AlreadyExists, "username already exists")
	ErrUnexpectedPassword        = status.Error(codes.FailedPrecondition, "only user account can have password")
	ErrInvalidPassword           = status.Error(codes.InvalidArgument, "password does not comply with policy")
	ErrInvalidPageToken          = status.Error(codes.InvalidArgument, "invalid page token")
	ErrInvalidFilter             = status.Error(codes.InvalidArgument, "invalid filter")
	ErrRoleInUse                 = status.Error(codes.FailedPrecondition, "role is in use")
)

type Server struct {
	gen.UnimplementedAccountApiServiceServer
	store  *Store
	logger *zap.Logger
}

func NewServer(store *Store, logger *zap.Logger) *Server {
	return &Server{store: store, logger: logger}
}

// GetAccount returns a single account by ID.
func (s *Server) GetAccount(ctx context.Context, req *gen.GetAccountRequest) (*gen.Account, error) {
	id, ok := parseID(req.Id)
	if !ok {
		return nil, ErrAccountNotFound
	}

	var dbAccount queries.Account
	err := s.store.Read(ctx, func(tx *Tx) error {
		var err error
		dbAccount, err = tx.GetAccount(ctx, id)
		return err
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAccountNotFound
	} else if err != nil {
		s.logger.Error("failed to get account", zap.Error(err), zap.String("id", req.Id))
		return nil, ErrDatabase
	}

	return accountToProto(dbAccount), nil
}

func (s *Server) ListAccounts(ctx context.Context, req *gen.ListAccountsRequest) (*gen.ListAccountsResponse, error) {
	pageSize := resolvePageSize(req.PageSize)

	afterID, ok := parsePageToken(req.PageToken)
	if !ok {
		return nil, ErrInvalidPageToken
	}

	res := &gen.ListAccountsResponse{}
	err := s.store.Read(ctx, func(tx *Tx) error {
		page, err := tx.ListAccounts(ctx, queries.ListAccountsParams{
			AfterID: afterID,
			Limit:   pageSize + 1, // fetch one extra to determine if there are more
		})
		if err != nil {
			return err
		}

		if int64(len(page)) > pageSize {
			last := page[pageSize-1] // last element that we are going to send
			res.NextPageToken = formatPageToken(last.ID)
			page = page[:pageSize]
		}
		for _, dbAccount := range page {
			res.Accounts = append(res.Accounts, accountToProto(dbAccount))
		}
		return nil
	})
	if err != nil {
		s.logger.Error("failed to list accounts",
			zap.Error(err),
			zap.String("pageToken", req.PageToken),
			zap.Int64("pageSize", pageSize),
		)
		return nil, ErrDatabase
	}

	return res, nil
}

func (s *Server) CreateAccount(ctx context.Context, req *gen.CreateAccountRequest) (*gen.Account, error) {
	account := req.Account
	if account == nil {
		return nil, status.Error(codes.InvalidArgument, "account is required")
	}
	if account.DisplayName == "" {
		return nil, ErrMissingDisplayName
	}
	switch account.Kind {
	case gen.Account_USER_ACCOUNT:
		if account.Username == "" {
			return nil, ErrMissingUsername
		}
		if req.Password != "" && !permitPassword(req.Password) {
			return nil, ErrInvalidPassword
		}
	case gen.Account_SERVICE_ACCOUNT:
		if account.Username != "" {
			return nil, ErrUnexpectedUsernameCreate
		}
		if req.Password != "" {
			return nil, ErrUnexpectedPassword
		}
	default:
		return nil, ErrInvalidAccountKind
	}

	var created queries.Account
	err := s.store.Write(ctx, func(tx *Tx) error {
		var err error
		switch req.Account.Kind {
		case gen.Account_USER_ACCOUNT:
			created, err = tx.CreateUserAccount(ctx, queries.CreateUserAccountParams{
				Username:    sql.NullString{Valid: true, String: account.Username},
				DisplayName: account.DisplayName,
			})
		case gen.Account_SERVICE_ACCOUNT:
			created, err = tx.CreateServiceAccount(ctx, account.DisplayName)
		default:
			panic("already validated account kind")
		}
		if err != nil {
			return err
		}

		if req.Password != "" {
			err = tx.UpdateAccountPassword(ctx, created.ID, req.Password)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if database.IsUniqueConstraintError(err) {
		return nil, ErrUsernameExists
	} else if err != nil {
		s.logger.Error("failed to create account",
			zap.Error(err),
			zap.String("kind", account.Kind.String()),
			zap.String("username", account.Username),
		)
		return nil, ErrDatabase
	}

	return accountToProto(created), nil
}

func (s *Server) UpdateAccount(ctx context.Context, req *gen.UpdateAccountRequest) (*gen.Account, error) {
	if req.Account == nil {
		return nil, status.Error(codes.InvalidArgument, "account is required")
	}

	id, ok := parseID(req.Account.Id)
	if !ok {
		return nil, ErrAccountNotFound
	}

	const (
		fieldDisplayName = "display_name"
		fieldUsername    = "username"
	)

	mask := maskOrDefault(req.Account, req.UpdateMask)
	// validate updated field values
	if maskContains(mask, fieldDisplayName) && req.Account.DisplayName == "" {
		return nil, status.Error(codes.InvalidArgument, "display_name must be non-empty")
	}
	if maskContains(mask, fieldUsername) && req.Account.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username must be non-empty")
	}

	var account queries.Account
	err := s.store.Write(ctx, func(tx *Tx) error {
		var err error
		account, err = tx.GetAccount(ctx, id)
		if err != nil {
			s.logger.Error("failed to get account", zap.Error(err), zap.String("id", req.Account.Id))
			return ErrDatabase
		}

		var (
			updateUsername    bool
			updateDisplayName bool
		)
		fields, err := fieldsToUpdate(accountToProto(account), req.Account, mask)
		if err != nil {
			return err
		}
		for _, field := range fields {
			switch field {
			case fieldDisplayName:
				updateDisplayName = true
			case fieldUsername:
				updateUsername = true
			default:
				return status.Errorf(codes.InvalidArgument, "field %q unsupported for update", field)
			}
		}

		if updateDisplayName {
			err = tx.UpdateAccountDisplayName(ctx, queries.UpdateAccountDisplayNameParams{
				ID:          id,
				DisplayName: req.Account.DisplayName,
			})
			if err != nil {
				s.logger.Error("failed to update account display name", zap.Error(err),
					zap.String("id", req.Account.Id), zap.String("displayName", req.Account.DisplayName))
				return ErrDatabase
			}
			account.DisplayName = req.Account.DisplayName
		}

		// only user accounts can have usernames
		usernameAllowed := account.Kind == gen.Account_USER_ACCOUNT.String()
		if updateUsername && !usernameAllowed {
			return ErrUnexpectedUsernameUpdate
		}

		if updateUsername {
			err = tx.UpdateAccountUsername(ctx, queries.UpdateAccountUsernameParams{
				ID:       id,
				Username: sql.NullString{Valid: true, String: req.Account.Username},
			})
			if database.IsUniqueConstraintError(err) {
				return ErrUsernameExists
			} else if err != nil {
				s.logger.Error("failed to update account username", zap.Error(err),
					zap.String("id", req.Account.Id),
					zap.String("username", req.Account.Username))
				return ErrDatabase
			}
			account.Username = sql.NullString{Valid: true, String: req.Account.Username}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return accountToProto(account), nil
}

func (s *Server) DeleteAccount(ctx context.Context, req *gen.DeleteAccountRequest) (*gen.DeleteAccountResponse, error) {
	id, ok := parseID(req.Id)
	if !ok {
		return nil, ErrAccountNotFound
	}

	var deleted bool
	err := s.store.Write(ctx, func(tx *Tx) error {
		rowsDeleted, err := tx.DeleteAccount(ctx, id)
		if err != nil {
			return err
		}
		deleted = rowsDeleted > 0
		return nil
	})
	if err != nil {
		s.logger.Error("failed to delete account", zap.Error(err), zap.String("id", req.Id))
		return nil, ErrDatabase
	}
	if !deleted {
		return nil, ErrAccountNotFound
	}
	return &gen.DeleteAccountResponse{}, nil
}

func (s *Server) GetServiceCredential(ctx context.Context, req *gen.GetServiceCredentialRequest) (*gen.ServiceCredential, error) {
	id, ok := parseID(req.Id)
	if !ok {
		return nil, ErrServiceCredentialNotFound
	}

	var cred queries.ServiceCredential
	err := s.store.Read(ctx, func(tx *Tx) error {
		var err error
		cred, err = tx.GetServiceCredential(ctx, id)
		return err
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrServiceCredentialNotFound
	} else if err != nil {
		s.logger.Error("failed to get service credential", zap.Error(err), zap.String("id", req.Id))
		return nil, ErrDatabase
	}

	return serviceCredentialToProto(cred, ""), nil
}

func (s *Server) ListServiceCredentials(ctx context.Context, req *gen.ListServiceCredentialsRequest) (*gen.ListServiceCredentialsResponse, error) {
	accountID, ok := parseID(req.AccountId)
	if !ok {
		return nil, ErrAccountNotFound
	}

	res := &gen.ListServiceCredentialsResponse{}
	err := s.store.Read(ctx, func(tx *Tx) error {
		account, err := tx.GetAccount(ctx, accountID)
		if err != nil {
			return err
		}
		if account.Kind != gen.Account_SERVICE_ACCOUNT.String() {
			return ErrUnexpectedServiceCreds
		}

		page, err := tx.ListAccountServiceCredentials(ctx, accountID)
		if err != nil {
			return err
		}

		for _, cred := range page {
			res.ServiceCredentials = append(res.ServiceCredentials, serviceCredentialToProto(cred, ""))
		}
		return nil
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAccountNotFound
	} else if err != nil {
		s.logger.Error("failed to list service credentials", zap.Error(err), zap.String("accountId", req.AccountId))
		return nil, ErrDatabase
	}

	return res, nil
}

func (s *Server) CreateServiceCredential(ctx context.Context, req *gen.CreateServiceCredentialRequest) (*gen.ServiceCredential, error) {
	if req.ServiceCredential == nil {
		return nil, status.Error(codes.InvalidArgument, "service_credential is required")
	}

	if !validateTitle(req.ServiceCredential.Title) {
		return nil, status.Error(codes.InvalidArgument, "invalid title")
	}

	accountID, ok := parseID(req.ServiceCredential.AccountId)
	if !ok {
		return nil, ErrAccountNotFound
	}

	var generated GeneratedServiceCredential
	err := s.store.Write(ctx, func(tx *Tx) error {
		var expiry sql.NullTime
		if req.ServiceCredential.ExpireTime != nil {
			expiry = sql.NullTime{Valid: true, Time: req.ServiceCredential.ExpireTime.AsTime()}
		}
		var err error
		generated, err = tx.GenerateServiceCredential(ctx, accountID, req.ServiceCredential.Title, expiry)
		return err
	})
	if err != nil {
		return nil, err
	}

	return serviceCredentialToProto(generated.ServiceCredential, generated.Secret), nil
}

func (s *Server) DeleteServiceCredential(ctx context.Context, req *gen.DeleteServiceCredentialRequest) (*gen.DeleteServiceCredentialResponse, error) {
	credID, ok := parseID(req.Id)
	if !ok {
		return nil, ErrServiceCredentialNotFound
	}

	var deleted bool
	err := s.store.Write(ctx, func(tx *Tx) error {
		rowsDeleted, err := tx.DeleteServiceCredential(ctx, credID)
		if err != nil {
			return err
		}
		deleted = rowsDeleted > 0
		return nil
	})
	if err != nil {
		s.logger.Error("failed to delete service credential", zap.Error(err), zap.String("id", req.Id))
		return nil, ErrDatabase
	}
	if !deleted {
		return nil, ErrServiceCredentialNotFound
	}

	return &gen.DeleteServiceCredentialResponse{}, nil
}

func (s *Server) UpdateAccountPassword(ctx context.Context, req *gen.UpdateAccountPasswordRequest) (*gen.UpdateAccountPasswordResponse, error) {
	if !permitPassword(req.NewPassword) {
		return nil, ErrInvalidPassword
	}

	id, ok := parseID(req.Id)
	if !ok {
		return nil, ErrAccountNotFound
	}

	err := s.store.Write(ctx, func(tx *Tx) error {
		if req.OldPassword != "" {
			err := tx.CheckAccountPassword(ctx, id, req.OldPassword)
			if err != nil {
				return err
			}
		}

		return tx.UpdateAccountPassword(ctx, id, req.NewPassword)
	})
	if err != nil {
		return nil, err
	}

	return &gen.UpdateAccountPasswordResponse{}, nil
}

func (s *Server) GetRole(ctx context.Context, req *gen.GetRoleRequest) (*gen.Role, error) {
	id, ok := parseID(req.Id)
	if !ok {
		return nil, ErrRoleNotFound
	}

	var (
		role        queries.Role
		permissions []string
	)
	err := s.store.Read(ctx, func(tx *Tx) error {
		var err error
		role, err = tx.GetRole(ctx, id)
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRoleNotFound
		} else if err != nil {
			return err
		}
		permissions, err = tx.ListRolePermissions(ctx, id)
		return err
	})
	if err != nil {
		return nil, err
	}

	return roleToProto(role, permissions), nil
}

func (s *Server) ListRoles(ctx context.Context, req *gen.ListRolesRequest) (*gen.ListRolesResponse, error) {
	pageSize := resolvePageSize(req.PageSize)

	afterID, ok := parsePageToken(req.PageToken)
	if !ok {
		return nil, ErrInvalidPageToken
	}

	res := &gen.ListRolesResponse{}
	err := s.store.Read(ctx, func(tx *Tx) error {
		page, err := tx.ListRolesWithPermissions(ctx, queries.ListRolesWithPermissionsParams{
			AfterID: afterID,
			Limit:   pageSize + 1, // fetch one extra to determine if there are more
		})
		if err != nil {
			return err
		}

		if int64(len(page)) > pageSize {
			last := page[pageSize-1] // last element that we are going to send
			res.NextPageToken = formatPageToken(last.Role.ID)
			page = page[:pageSize]
		}

		for _, role := range page {
			permissions := splitPermissions(role.Permissions)
			res.Roles = append(res.Roles, roleToProto(role.Role, permissions))
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Server) CreateRole(ctx context.Context, req *gen.CreateRoleRequest) (*gen.Role, error) {
	if req.Role == nil {
		return nil, status.Error(codes.InvalidArgument, "role is required")
	}

	var (
		role        queries.Role
		permissions []string
	)
	err := s.store.Write(ctx, func(tx *Tx) error {
		var err error
		role, err = tx.CreateRole(ctx, req.Role.Title)
		if err != nil {
			return err
		}

		for _, perm := range req.Role.Permissions {
			err = tx.AddRolePermission(ctx, queries.AddRolePermissionParams{
				RoleID:     role.ID,
				Permission: perm,
			})
			if err != nil {
				return err
			}
		}

		// re-fetch the permissions, as they have been reordered and deduplicated
		// when added to the database
		permissions, err = tx.ListRolePermissions(ctx, role.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return roleToProto(role, permissions), nil
}

func (s *Server) UpdateRole(ctx context.Context, req *gen.UpdateRoleRequest) (*gen.Role, error) {
	if req.Role == nil {
		return nil, status.Error(codes.InvalidArgument, "role is required")
	}
	id, ok := parseID(req.Role.Id)
	if !ok {
		return nil, ErrRoleNotFound
	}

	var (
		updateTitle       bool
		updatePermissions bool
	)
	if req.UpdateMask == nil {
		// if no update mask is provided, update all fields that are set
		updateTitle = req.Role.Title != ""
		updatePermissions = len(req.Role.Permissions) > 0
	} else {
		for _, path := range req.UpdateMask.Paths {
			switch path {
			case "title":
				updateTitle = true
			case "permissions":
				updatePermissions = true
			default:
				return nil, status.Errorf(codes.InvalidArgument, "unsupported field %q in update mask", path)
			}
		}
	}

	var (
		role        queries.Role
		permissions []string
	)
	err := s.store.Write(ctx, func(tx *Tx) error {
		var err error
		role, err = tx.GetRole(ctx, id)
		if err != nil {
			return err
		}

		if updateTitle {
			_, err = tx.UpdateRoleName(ctx, queries.UpdateRoleNameParams{
				ID:   id,
				Name: req.Role.Title,
			})
			if err != nil {
				return err
			}
			role.Name = req.Role.Title
		}

		if updatePermissions {
			// clear existing permissions
			_, err = tx.ClearRolePermissions(ctx, id)
			if err != nil {
				return err
			}

			// add new permissions
			for _, perm := range req.Role.Permissions {
				err = tx.AddRolePermission(ctx, queries.AddRolePermissionParams{
					RoleID:     id,
					Permission: perm,
				})
				if err != nil {
					return err
				}
			}
		}

		permissions, err = tx.ListRolePermissions(ctx, id)
		return err
	})
	if err != nil {
		return nil, err
	}

	return roleToProto(role, permissions), nil
}

func (s *Server) DeleteRole(ctx context.Context, req *gen.DeleteRoleRequest) (*gen.DeleteRoleResponse, error) {
	id, ok := parseID(req.Id)
	if !ok {
		return nil, ErrRoleNotFound
	}

	var deleted bool
	err := s.store.Write(ctx, func(tx *Tx) error {
		rowsDeleted, err := tx.DeleteRole(ctx, id)
		if database.IsForeignKeyError(err) {
			return ErrRoleInUse
		} else if err != nil {
			s.logger.Error("failed to delete role", zap.Error(err), zap.String("id", req.Id))
			return ErrDatabase
		}
		deleted = rowsDeleted > 0
		return nil
	})
	if err != nil {
		return nil, err
	}
	if !deleted {
		return nil, ErrRoleNotFound
	}

	return &gen.DeleteRoleResponse{}, nil
}

func (s *Server) GetRoleAssignment(ctx context.Context, req *gen.GetRoleAssignmentRequest) (*gen.RoleAssignment, error) {
	id, ok := parseID(req.Id)
	if !ok {
		return nil, ErrRoleAssignmentNotFound
	}

	var assignment queries.RoleAssignment
	err := s.store.Read(ctx, func(tx *Tx) error {
		var err error
		assignment, err = tx.GetRoleAssignment(ctx, id)
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRoleAssignmentNotFound
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	return roleAssignmentToProto(assignment), nil
}

func (s *Server) ListRoleAssignments(ctx context.Context, req *gen.ListRoleAssignmentsRequest) (*gen.ListRoleAssignmentsResponse, error) {
	pageSize := resolvePageSize(req.PageSize)

	afterID, ok := parsePageToken(req.PageToken)
	if !ok {
		return nil, ErrInvalidPageToken
	}

	filterField, filterID, ok := parseRoleAssignmentFilter(req.Filter)
	if !ok {
		return nil, ErrInvalidFilter
	}

	var (
		res = &gen.ListRoleAssignmentsResponse{}
		err error
	)
	err = s.store.Read(ctx, func(tx *Tx) error {
		var page []queries.RoleAssignment
		switch filterField {
		case roleAssignmentAccountID:
			page, err = tx.ListRoleAssignmentsForAccount(ctx, queries.ListRoleAssignmentsForAccountParams{
				AfterID:   afterID,
				Limit:     pageSize + 1, // fetch one extra to determine if there are more
				AccountID: filterID,
			})
		case roleAssignmentRoleID:
			page, err = tx.ListRoleAssignmentsForRole(ctx, queries.ListRoleAssignmentsForRoleParams{
				AfterID: afterID,
				Limit:   pageSize + 1,
				RoleID:  filterID,
			})
		case roleAssignmentUnfiltered:
			page, err = tx.ListRoleAssignments(ctx, queries.ListRoleAssignmentsParams{
				AfterID: afterID,
				Limit:   pageSize + 1,
			})
		default:
			// unreachable because parseRoleAssignmentFilter only allows account_id and role_id
			panic("unreachable")
		}
		if err != nil {
			return err
		}

		if int64(len(page)) > pageSize {
			last := page[pageSize-1] // last element that we are going to send
			res.NextPageToken = formatPageToken(last.ID)
			page = page[:pageSize]
		}

		for _, assignment := range page {
			res.RoleAssignments = append(res.RoleAssignments, roleAssignmentToProto(assignment))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Server) CreateRoleAssignment(ctx context.Context, req *gen.CreateRoleAssignmentRequest) (*gen.RoleAssignment, error) {
	accountID, ok := parseID(req.RoleAssignment.AccountId)
	if !ok {
		return nil, ErrAccountNotFound
	}
	roleID, ok := parseID(req.RoleAssignment.RoleId)
	if !ok {
		return nil, ErrRoleNotFound
	}

	var (
		scopeKind, scopeResource sql.NullString
	)
	if req.RoleAssignment.Scope != nil {
		scopeKind = sql.NullString{Valid: true, String: req.RoleAssignment.Scope.ResourceKind.String()}
		scopeResource = sql.NullString{Valid: true, String: req.RoleAssignment.Scope.Resource}
	}

	var assignment queries.RoleAssignment
	err := s.store.Write(ctx, func(tx *Tx) error {
		var err error
		assignment, err = tx.CreateRoleAssignment(ctx, queries.CreateRoleAssignmentParams{
			AccountID:     accountID,
			RoleID:        roleID,
			ScopeKind:     scopeKind,
			ScopeResource: scopeResource,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return roleAssignmentToProto(assignment), nil
}

func (s *Server) DeleteRoleAssignment(ctx context.Context, req *gen.DeleteRoleAssignmentRequest) (*gen.DeleteRoleAssignmentResponse, error) {
	id, ok := parseID(req.Id)
	if !ok {
		return nil, ErrRoleAssignmentNotFound
	}

	var deleted bool
	err := s.store.Write(ctx, func(tx *Tx) error {
		rowsDeleted, err := tx.DeleteRoleAssignment(ctx, id)
		if err != nil {
			return err
		}
		deleted = rowsDeleted > 0
		return nil
	})
	if err != nil {
		s.logger.Error("failed to delete role assignment", zap.Error(err), zap.String("id", req.Id))
		return nil, ErrDatabase
	}
	if !deleted {
		return nil, ErrRoleAssignmentNotFound
	}

	return &gen.DeleteRoleAssignmentResponse{}, nil
}

func resolvePageSize(pageSize int32) int64 {
	if pageSize == 0 {
		return defaultPageSize
	}
	if pageSize < minPageSize {
		return minPageSize
	}
	if pageSize > maxPageSize {
		return maxPageSize
	}
	return int64(pageSize)
}

const (
	minPageSize     = 1
	maxPageSize     = 100
	defaultPageSize = 30
)

var filterRoleAssignmentsRegexp = regexp.MustCompile(`^ *(account_id|role_id) *= *(\d+) *$`)

type roleAssignmentField string

const (
	roleAssignmentAccountID  roleAssignmentField = "account_id"
	roleAssignmentRoleID     roleAssignmentField = "role_id"
	roleAssignmentUnfiltered roleAssignmentField = ""
)

// parses a filter string like:
//   - "account_id=123" - matches role assignment where account_id=123
//   - "role_id=456"    - matches role assignment where role_id=456
//   - ""               - matches all role assignments (the empty string is a valid filter)
func parseRoleAssignmentFilter(filter string) (field roleAssignmentField, id int64, ok bool) {
	if filter == "" {
		return roleAssignmentUnfiltered, 0, true
	}
	m := filterRoleAssignmentsRegexp.FindStringSubmatch(filter)
	if m == nil {
		return "", 0, false
	}
	id, ok = parseID(m[2])
	if !ok {
		return "", 0, false
	}
	return roleAssignmentField(m[1]), id, true
}
