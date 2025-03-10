package account

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"regexp"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vanti-dev/sc-bos/internal/account/queries"
	"github.com/vanti-dev/sc-bos/internal/sqlite"
	"github.com/vanti-dev/sc-bos/internal/util/pass"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

var (
	ErrDatabase                  = status.Error(codes.Internal, "database internal error")
	ErrAccountNotFound           = status.Error(codes.NotFound, "account not found")
	ErrRoleNotFound              = status.Error(codes.NotFound, "role not found")
	ErrRoleAssignmentNotFound    = status.Error(codes.NotFound, "role assignment not found")
	ErrServiceCredentialNotFound = status.Error(codes.NotFound, "service credential not found")
	ErrPermissionNotFound        = status.Error(codes.NotFound, "permission not found")
	ErrInvalidAccountType        = status.Error(codes.InvalidArgument, "invalid account type")
	ErrMissingUsername           = status.Error(codes.InvalidArgument, "user account requires username")
	ErrMissingDisplayName        = status.Error(codes.InvalidArgument, "account requires display name")
	ErrUnexpectedUsernameCreate  = status.Error(codes.InvalidArgument, "service account cannot have username")
	ErrUnexpectedUsernameUpdate  = status.Error(codes.FailedPrecondition, "service account cannot have username")
	ErrUnexpectedServiceCreds    = status.Error(codes.FailedPrecondition, "user account cannot have service credentials")
	ErrServiceCredentialLimit    = status.Error(codes.ResourceExhausted, "too many service credentials")
	ErrUsernameExists            = status.Error(codes.AlreadyExists, "username already exists")
	ErrRoleAssignmentExists      = status.Error(codes.AlreadyExists, "role assignment already exists")
	ErrUnexpectedPasswordCreate  = status.Error(codes.InvalidArgument, "only user account can have password")
	ErrUnexpectedPasswordUpdate  = status.Error(codes.FailedPrecondition, "only user account can have password")
	ErrInvalidUsername           = status.Error(codes.InvalidArgument, "invalid username")
	ErrInvalidDisplayName        = status.Error(codes.InvalidArgument, "invalid display name")
	ErrInvalidDescription        = status.Error(codes.InvalidArgument, "invalid description")
	ErrInvalidPassword           = status.Error(codes.InvalidArgument, "password does not comply with policy")
	ErrInvalidResourceType       = status.Error(codes.InvalidArgument, "invalid scope resource type")
	ErrInvalidResource           = status.Error(codes.InvalidArgument, "invalid scope resource")
	ErrIncorrectPassword         = status.Error(codes.FailedPrecondition, "incorrect password")
	ErrInvalidPageToken          = status.Error(codes.InvalidArgument, "invalid page token")
	ErrInvalidFilter             = status.Error(codes.InvalidArgument, "invalid filter")
	ErrRoleInUse                 = status.Error(codes.FailedPrecondition, "role is in use")
	ErrResourceMissing           = status.Error(codes.InvalidArgument, "resource to create/update not supplied")
)

type Server struct {
	gen.UnimplementedAccountApiServer
	gen.UnimplementedAccountInfoServer
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
		return nil, s.processError(err, zap.String("rpc", "GetAccount"), zap.String("id", req.Id))
	}

	return accountToProto(dbAccount), nil
}

func (s *Server) ListAccounts(ctx context.Context, req *gen.ListAccountsRequest) (*gen.ListAccountsResponse, error) {
	pageSize := resolvePageSize(req.PageSize)

	res := &gen.ListAccountsResponse{
		TotalSize: -1, // sentinel, indicates we need to calculate this
	}
	var afterID int64 = 0
	if req.PageToken != "" {
		// this RPC does not support filtering
		pageToken, err := parsePageToken(req.PageToken, "")
		if err != nil {
			return nil, ErrInvalidPageToken
		}
		afterID = pageToken.LastId
		res.TotalSize = pageToken.TotalSize
	}

	err := s.store.Read(ctx, func(tx *Tx) error {
		page, err := tx.ListAccounts(ctx, queries.ListAccountsParams{
			AfterID: afterID,
			Limit:   pageSize + 1, // fetch one extra to determine if there are more
		})
		if err != nil {
			return err
		}

		if res.TotalSize < 0 {
			count, err := tx.CountAccounts(ctx)
			if err != nil {
				return err
			}
			if count > math.MaxInt32 {
				res.TotalSize = 0
			} else {
				res.TotalSize = int32(count)
			}
		}

		if int64(len(page)) > pageSize {
			last := page[pageSize-1] // last element that we are going to send
			res.NextPageToken = encodePageToken(&PageToken{
				LastId:    last.ID,
				TotalSize: res.TotalSize,
			})
			page = page[:pageSize]
		}
		for _, dbAccount := range page {
			res.Accounts = append(res.Accounts, accountToProto(dbAccount))
		}
		return nil
	})
	if err != nil {
		return nil, s.processError(err,
			zap.String("rpc", "ListAccounts"),
			zap.String("pageToken", req.PageToken),
			zap.Int32("pageSize", req.PageSize),
		)
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
	} else if !validateDisplayName(account.DisplayName) {
		return nil, ErrInvalidDisplayName
	}
	if !validateDescription(account.Description) {
		return nil, ErrInvalidDescription
	}
	switch account.Type {
	case gen.Account_USER_ACCOUNT:
		if account.Username == "" {
			return nil, ErrMissingUsername
		}
		if !validateUsername(account.Username) {
			return nil, ErrInvalidUsername
		}
	case gen.Account_SERVICE_ACCOUNT:
		if account.Username != "" {
			return nil, ErrUnexpectedUsernameCreate
		}
		if req.Password != "" {
			return nil, ErrUnexpectedPasswordCreate
		}
	default:
		return nil, ErrInvalidAccountType
	}

	var created queries.Account
	err := s.store.Write(ctx, func(tx *Tx) error {
		var description sql.NullString
		if req.Account.Description != "" {
			description = sql.NullString{Valid: true, String: req.Account.Description}
		}

		var err error
		switch req.Account.Type {
		case gen.Account_USER_ACCOUNT:
			created, err = tx.CreateUserAccount(ctx, queries.CreateUserAccountParams{
				Username:    sql.NullString{Valid: true, String: account.Username},
				DisplayName: account.DisplayName,
				Description: description,
			})
		case gen.Account_SERVICE_ACCOUNT:
			created, err = tx.CreateServiceAccount(ctx, queries.CreateServiceAccountParams{
				DisplayName: account.DisplayName,
				Description: description,
			})
		default:
			return ErrInvalidAccountType
		}
		if sqlite.IsUniqueConstraintError(err) {
			return ErrUsernameExists
		} else if err != nil {
			return err
		}

		if req.Password != "" {
			err = tx.UpdateAccountPassword(ctx, created.ID, req.Password)
			if errors.Is(err, ErrInvalidPassword) {
				return err
			} else if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, s.processError(err, zap.String("rpc", "CreateAccount"))
	}

	return accountToProto(created), nil
}

func (s *Server) UpdateAccount(ctx context.Context, req *gen.UpdateAccountRequest) (*gen.Account, error) {
	const (
		fieldDisplayName = "display_name"
		fieldDescription = "description"
		fieldUsername    = "username"
		fieldCreateTime  = "create_time"
	)
	// ignore output-only fields in masks, as per AIP-203
	mask, err := resolveMask(req.Account, req.UpdateMask, fieldCreateTime)
	if err != nil {
		return nil, err
	}

	if req.Account == nil {
		return nil, status.Error(codes.InvalidArgument, "account is required")
	}

	id, ok := parseID(req.Account.Id)
	if !ok {
		return nil, ErrAccountNotFound
	}

	var account queries.Account
	err = s.store.Write(ctx, func(tx *Tx) error {
		var err error
		account, err = tx.GetAccount(ctx, id)
		if errors.Is(err, sql.ErrNoRows) {
			return ErrAccountNotFound
		} else if err != nil {
			return err
		}

		var (
			updateUsername    bool
			updateDisplayName bool
			updateDescription bool
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
			case fieldDescription:
				updateDescription = true
			default:
				return status.Errorf(codes.InvalidArgument, "field %q unsupported for update", field)
			}
		}

		if updateDisplayName {
			if !validateDisplayName(req.Account.DisplayName) {
				return ErrMissingDisplayName
			}
			err = tx.UpdateAccountDisplayName(ctx, queries.UpdateAccountDisplayNameParams{
				ID:          id,
				DisplayName: req.Account.DisplayName,
			})
			if err != nil {
				return err
			}
			account.DisplayName = req.Account.DisplayName
		}

		// only user accounts can have usernames
		usernameAllowed := account.Type == gen.Account_USER_ACCOUNT.String()
		if updateUsername && !usernameAllowed {
			return ErrUnexpectedUsernameUpdate
		}

		if updateUsername {
			if !validateUsername(req.Account.Username) {
				return ErrInvalidUsername
			}
			err = tx.UpdateAccountUsername(ctx, queries.UpdateAccountUsernameParams{
				ID:       id,
				Username: sql.NullString{Valid: true, String: req.Account.Username},
			})
			if sqlite.IsUniqueConstraintError(err) {
				return ErrUsernameExists
			} else if err != nil {
				return err
			}
			account.Username = sql.NullString{Valid: true, String: req.Account.Username}
		}

		if updateDescription {
			if !validateDescription(req.Account.Description) {
				return ErrInvalidResource
			}
			var description sql.NullString
			if req.Account.Description != "" {
				description = sql.NullString{Valid: true, String: req.Account.Description}
			}

			err = tx.UpdateAccountDescription(ctx, queries.UpdateAccountDescriptionParams{
				ID:          id,
				Description: description,
			})
			if err != nil {
				return err
			}
			account.Description = description
		}

		return nil
	})
	if err != nil {
		return nil, s.processError(err, zap.String("rpc", "UpdateAccount"), zap.String("id", req.Account.Id))
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
		return nil, s.processError(err, zap.String("rpc", "DeleteAccount"), zap.String("id", req.Id))
	}
	if !deleted && !req.AllowMissing {
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
		return nil, s.processError(err, zap.String("rpc", "GetServiceCredential"), zap.String("id", req.Id))
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
		if errors.Is(err, sql.ErrNoRows) {
			return ErrAccountNotFound
		} else if err != nil {
			return err
		}
		if account.Type != gen.Account_SERVICE_ACCOUNT.String() {
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
	if err != nil {
		return nil, s.processError(err, zap.String("rpc", "ListServiceCredentials"), zap.String("accountId", req.AccountId))
	}

	return res, nil
}

func (s *Server) CreateServiceCredential(ctx context.Context, req *gen.CreateServiceCredentialRequest) (*gen.ServiceCredential, error) {
	if req.ServiceCredential == nil {
		return nil, status.Error(codes.InvalidArgument, "service_credential is required")
	}

	if !validateDisplayName(req.ServiceCredential.DisplayName) {
		return nil, ErrInvalidDisplayName
	}
	if !validateDescription(req.ServiceCredential.Description) {
		return nil, ErrInvalidDescription
	}

	accountID, ok := parseID(req.ServiceCredential.AccountId)
	if !ok {
		return nil, ErrAccountNotFound
	}

	var generated GeneratedServiceCredential
	err := s.store.Write(ctx, func(tx *Tx) error {
		var (
			expireTime  sql.NullTime
			description sql.NullString
		)
		if req.ServiceCredential.ExpireTime != nil {
			expireTime = sql.NullTime{Valid: true, Time: req.ServiceCredential.ExpireTime.AsTime()}
		}
		if req.ServiceCredential.Description != "" {
			description = sql.NullString{Valid: true, String: req.ServiceCredential.Description}
		}
		var err error
		generated, err = tx.GenerateServiceCredential(ctx, queries.ServiceCredential{
			AccountID:   accountID,
			DisplayName: req.ServiceCredential.DisplayName,
			Description: description,
			ExpireTime:  expireTime,
		})
		return err
	})
	if err != nil {
		return nil, s.processError(err, zap.String("rpc", "CreateServiceCredential"), zap.String("accountId", req.ServiceCredential.AccountId))
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
	if !deleted && !req.AllowMissing {
		return nil, ErrServiceCredentialNotFound
	}

	return &gen.DeleteServiceCredentialResponse{}, nil
}

func (s *Server) UpdateAccountPassword(ctx context.Context, req *gen.UpdateAccountPasswordRequest) (*gen.UpdateAccountPasswordResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if !permitPassword(req.NewPassword) {
		return nil, ErrInvalidPassword
	}

	id, ok := parseID(req.Id)
	if !ok {
		return nil, ErrAccountNotFound
	}

	err := s.store.Write(ctx, func(tx *Tx) error {
		if req.OldPassword != "" {
			_, err := tx.GetAccount(ctx, id)
			if errors.Is(err, sql.ErrNoRows) {
				return ErrAccountNotFound
			} else if err != nil {
				return err
			}

			err = tx.CheckAccountPassword(ctx, id, req.OldPassword)
			if errors.Is(err, sql.ErrNoRows) {
				// account has no password saved
				return ErrIncorrectPassword
			} else if errors.Is(err, pass.ErrMismatchedHashAndPassword) {
				return ErrIncorrectPassword
			} else if err != nil {
				return err
			}
		}

		return tx.UpdateAccountPassword(ctx, id, req.NewPassword)
	})
	if err != nil {
		return nil, s.processError(err, zap.String("rpc", "UpdateAccountPassword"), zap.String("id", req.Id))
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
		return nil, s.processError(err, zap.String("rpc", "GetRole"), zap.String("id", req.Id))
	}

	return roleToProto(role, permissions), nil
}

func (s *Server) ListRoles(ctx context.Context, req *gen.ListRolesRequest) (*gen.ListRolesResponse, error) {
	pageSize := resolvePageSize(req.PageSize)

	res := &gen.ListRolesResponse{
		TotalSize: -1, // sentinel, indicates we need to calculate this
	}
	var afterID int64
	if req.PageToken != "" {
		token, err := parsePageToken(req.PageToken, "")
		if err != nil {
			return nil, ErrInvalidPageToken
		}
		afterID = token.LastId
		res.TotalSize = token.TotalSize
	}

	err := s.store.Read(ctx, func(tx *Tx) error {
		page, err := tx.ListRolesAndPermissions(ctx, queries.ListRolesAndPermissionsParams{
			AfterID: afterID,
			Limit:   pageSize + 1, // fetch one extra to determine if there are more
		})
		if err != nil {
			return err
		}

		if res.TotalSize < 0 {
			count, err := tx.CountRoles(ctx)
			if err != nil {
				return err
			}
			if count > math.MaxInt32 {
				res.TotalSize = 0
			} else {
				res.TotalSize = int32(count)
			}
		}

		if int64(len(page)) > pageSize {
			last := page[pageSize-1] // last element that we are going to send
			res.NextPageToken = encodePageToken(&PageToken{
				LastId:    last.Role.ID,
				TotalSize: res.TotalSize,
			})
			page = page[:pageSize]
		}

		for _, role := range page {
			permissions := splitPermissions(role.Permissions)
			res.Roles = append(res.Roles, roleToProto(role.Role, permissions))
		}

		return nil
	})
	if err != nil {
		return nil, s.processError(err, zap.String("rpc", "ListRoles"), zap.String("pageToken", req.PageToken))
	}

	return res, nil
}

func (s *Server) CreateRole(ctx context.Context, req *gen.CreateRoleRequest) (*gen.Role, error) {
	if req.Role == nil {
		return nil, ErrResourceMissing
	}

	if !validateDisplayName(req.Role.DisplayName) {
		return nil, ErrInvalidDisplayName
	}
	if !validateDescription(req.Role.Description) {
		return nil, ErrInvalidDescription
	}

	var (
		role        queries.Role
		permissions []string
	)
	err := s.store.Write(ctx, func(tx *Tx) error {
		var err error
		params := queries.CreateRoleParams{
			DisplayName: req.Role.DisplayName,
		}
		if req.Role.Description != "" {
			params.Description = sql.NullString{Valid: true, String: req.Role.Description}
		}
		role, err = tx.CreateRole(ctx, params)
		if err != nil {
			return err
		}

		for _, perm := range req.Role.PermissionIds {
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
		return nil, s.processError(err, zap.String("rpc", "CreateRole"))
	}

	return roleToProto(role, permissions), nil
}

func (s *Server) UpdateRole(ctx context.Context, req *gen.UpdateRoleRequest) (*gen.Role, error) {
	if req.Role == nil {
		return nil, status.Error(codes.InvalidArgument, "role is required")
	}

	const (
		fieldDisplayName   = "display_name"
		fieldPermissionIDs = "permission_ids"
		fieldDescription   = "description"
	)
	mask, err := resolveMask(req.Role, req.UpdateMask)
	if err != nil {
		return nil, err
	}

	id, ok := parseID(req.Role.Id)
	if !ok {
		return nil, ErrRoleNotFound
	}

	var (
		role        queries.Role
		permissions []string
	)
	err = s.store.Write(ctx, func(tx *Tx) error {
		var err error
		role, err = tx.GetRole(ctx, id)
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRoleNotFound
		} else if err != nil {
			return err
		}
		permissions, err = tx.ListRolePermissions(ctx, id)
		if err != nil {
			return err
		}

		// decide which fields to update
		var (
			updateDisplayName bool
			updatePermissions bool
			updateDescription bool
		)
		fields, err := fieldsToUpdate(roleToProto(role, permissions), req.Role, mask)
		if err != nil {
			return err
		}
		for _, field := range fields {
			switch field {
			case fieldDisplayName:
				updateDisplayName = true
			case fieldPermissionIDs:
				updatePermissions = true
			case fieldDescription:
				updateDescription = true
			default:
				return status.Errorf(codes.InvalidArgument, "field %q unsupported for update", field)
			}
		}

		if updateDisplayName {
			if !validateDisplayName(req.Role.DisplayName) {
				return ErrInvalidDisplayName
			}

			_, err = tx.UpdateRoleDisplayName(ctx, queries.UpdateRoleDisplayNameParams{
				ID:          id,
				DisplayName: req.Role.DisplayName,
			})
			if err != nil {
				return err
			}
			role.DisplayName = req.Role.DisplayName
		}

		if updateDescription {
			if !validateDescription(req.Role.Description) {
				return ErrInvalidDescription
			}

			var value sql.NullString
			if req.Role.Description != "" {
				value = sql.NullString{String: req.Role.Description, Valid: true}
			}
			_, err = tx.UpdateRoleDescription(ctx, queries.UpdateRoleDescriptionParams{
				ID:          id,
				Description: value,
			})
			if err != nil {
				return err
			}
			role.Description = value
		}

		if updatePermissions {
			// clear existing permissions
			_, err = tx.ClearRolePermissions(ctx, id)
			if err != nil {
				return err
			}

			// add new permissions
			for _, perm := range req.Role.PermissionIds {
				err = tx.AddRolePermission(ctx, queries.AddRolePermissionParams{
					RoleID:     id,
					Permission: perm,
				})
				if err != nil {
					return err
				}
			}

			permissions, err = tx.ListRolePermissions(ctx, id)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, s.processError(err, zap.String("rpc", "UpdateRole"), zap.String("id", req.Role.Id))
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
		if sqlite.IsForeignKeyError(err) {
			return ErrRoleInUse
		} else if err != nil {
			s.logger.Error("failed to delete role", zap.Error(err), zap.String("id", req.Id))
			return ErrDatabase
		}
		deleted = rowsDeleted > 0
		return nil
	})
	if err != nil {
		return nil, s.processError(err, zap.String("rpc", "DeleteRole"), zap.String("id", req.Id))
	}
	if !deleted && !req.AllowMissing {
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
		return nil, s.processError(err, zap.String("rpc", "GetRoleAssignment"), zap.String("id", req.Id))
	}

	return roleAssignmentToProto(assignment), nil
}

func (s *Server) ListRoleAssignments(ctx context.Context, req *gen.ListRoleAssignmentsRequest) (*gen.ListRoleAssignmentsResponse, error) {
	pageSize := resolvePageSize(req.PageSize)

	filterField, filterID, ok := parseRoleAssignmentFilter(req.Filter)
	if !ok {
		return nil, ErrInvalidFilter
	}

	var token *PageToken
	if req.PageToken != "" {
		var err error
		token, err = parsePageToken(req.PageToken, req.Filter)
		if err != nil {
			return nil, err
		}
	}

	var (
		res = &gen.ListRoleAssignmentsResponse{}
		err error
	)
	err = s.store.Read(ctx, func(tx *Tx) error {
		page, err := tx.ListRoleAssignmentsFiltered(ctx, filterField, filterID, token, pageSize)
		if err != nil {
			return err
		}
		res.TotalSize = page.TotalSize
		if page.More {
			res.NextPageToken = encodePageToken(&PageToken{
				LastId:    page.LastID,
				TotalSize: page.TotalSize,
				Filter:    req.Filter,
			})
		}

		for _, assignment := range page.RoleAssignments {
			res.RoleAssignments = append(res.RoleAssignments, roleAssignmentToProto(assignment))
		}
		return nil
	})
	if err != nil {
		return nil, s.processError(err, zap.String("rpc", "ListRoleAssignments"), zap.String("filter", req.Filter))
	}

	return res, nil
}

func (s *Server) CreateRoleAssignment(ctx context.Context, req *gen.CreateRoleAssignmentRequest) (*gen.RoleAssignment, error) {
	if req.RoleAssignment == nil {
		return nil, ErrResourceMissing
	}

	accountID, ok := parseID(req.RoleAssignment.AccountId)
	if !ok {
		return nil, ErrAccountNotFound
	}
	roleID, ok := parseID(req.RoleAssignment.RoleId)
	if !ok {
		return nil, ErrRoleNotFound
	}

	var (
		scopeType, scopeResource sql.NullString
	)
	if scope := req.RoleAssignment.Scope; scope != nil {
		if !validateResourceType(scope.ResourceType) {
			return nil, ErrInvalidResourceType
		}
		if !validateResource(scope.Resource) {
			return nil, ErrInvalidResource
		}

		scopeType = sql.NullString{Valid: true, String: scope.ResourceType.String()}
		scopeResource = sql.NullString{Valid: true, String: scope.Resource}
	}

	var assignment queries.RoleAssignment
	err := s.store.Write(ctx, func(tx *Tx) error {
		var err error
		assignment, err = tx.CreateRoleAssignment(ctx, queries.CreateRoleAssignmentParams{
			AccountID:     accountID,
			RoleID:        roleID,
			ScopeKind:     scopeType,
			ScopeResource: scopeResource,
		})
		if sqlite.IsUniqueConstraintError(err) {
			return ErrRoleAssignmentExists
		} else if sqlite.IsForeignKeyError(err) {
			// figure out which one is missing
			_, accountErr := tx.GetAccount(ctx, accountID)
			_, roleErr := tx.GetRole(ctx, roleID)
			if errors.Is(accountErr, sql.ErrNoRows) {
				return ErrAccountNotFound
			} else if errors.Is(roleErr, sql.ErrNoRows) {
				return ErrRoleNotFound
			} else {
				// unexpected database error
				return errors.Join(err, accountErr, roleErr)
			}
		}
		return err
	})
	if err != nil {
		return nil, s.processError(err,
			zap.String("rpc", "CreateRoleAssignment"),
			zap.String("accountId", req.RoleAssignment.AccountId),
			zap.String("roleId", req.RoleAssignment.RoleId),
		)
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
		return nil, s.processError(err, zap.String("rpc", "DeleteRoleAssignment"), zap.String("id", req.Id))
	}
	if !deleted && !req.AllowMissing {
		return nil, ErrRoleAssignmentNotFound
	}

	return &gen.DeleteRoleAssignmentResponse{}, nil
}

func (s *Server) GetPermission(ctx context.Context, req *gen.GetPermissionRequest) (*gen.Permission, error) {
	// TODO: add list of valid permissions and support fetching them
	return nil, ErrPermissionNotFound
}

func (s *Server) ListPermissions(ctx context.Context, req *gen.ListPermissionsRequest) (*gen.ListPermissionsResponse, error) {
	// TODO: add list of valid permissions and support fetching them
	return &gen.ListPermissionsResponse{}, nil
}

func (s *Server) GetAccountLimits(ctx context.Context, req *gen.GetAccountLimitsRequest) (*gen.AccountLimits, error) {
	return &gen.AccountLimits{
		Username: &gen.AccountLimits_Field{
			MinLength: minUsernameLength,
			MaxLength: maxUsernameLength,
		},
		Password: &gen.AccountLimits_Field{
			MinLength: minPasswordLength,
			MaxLength: maxPasswordLength,
		},
		DisplayName: &gen.AccountLimits_Field{
			MinLength: minDisplayNameLength,
			MaxLength: maxDisplayNameLength,
		},
		Description: &gen.AccountLimits_Field{
			MinLength: minDescriptionLength,
			MaxLength: maxDescriptionLength,
		},
		MaxServiceCredentialsPerAccount: maxServiceCredentialsPerAccount,
	}, nil
}

func (s *Server) processError(err error, fields ...zap.Field) error {
	logger := s.logger.With(fields...)
	switch {
	case errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded):
		logger.Debug("request cancelled or timed out", zap.Error(err))
		return err
	case status.Code(err) != codes.Unknown:
		return err
	default:
		s.logger.Error("unexpected account service internal error", zap.Error(err))
		return ErrDatabase
	}
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
