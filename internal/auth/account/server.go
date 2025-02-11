package account

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type Server struct {
	gen.UnimplementedAccountApiServiceServer
	store *Store
}

func NewServer(store *Store) *Server {
	return &Server{store: store}
}

// GetAccount returns a single account by ID.
//
// TODO: implement read_mask
func (s *Server) GetAccount(ctx context.Context, req *gen.GetAccountRequest) (*gen.GetAccountResponse, error) {
	var account *gen.Account
	err := s.store.Read(ctx, func(tx *ReadTx) error {
		var err error
		account, err = tx.GetAccount(ctx, req.Id)
		return err
	})

	if err != nil {
		return nil, err
	}
	return &gen.GetAccountResponse{Account: account}, nil
}

func (s *Server) ListAccounts(ctx context.Context, req *gen.ListAccountsRequest) (*gen.ListAccountsResponse, error) {
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = defaultPageSize
	} else if pageSize < minPageSize {
		pageSize = minPageSize
	} else if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	res := &gen.ListAccountsResponse{}
	err := s.store.Read(ctx, func(tx *ReadTx) error {
		var err error
		res.Accounts, res.NextPageToken, err = tx.ListAccounts(ctx, req.PageToken, int64(pageSize))
		return err
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Server) CreateAccount(ctx context.Context, req *gen.CreateAccountRequest) (*gen.CreateAccountResponse, error) {
	account := req.Account
	if account == nil {
		return nil, status.Error(codes.InvalidArgument, "account is required")
	}
	switch account.Kind {
	case gen.Account_USER_ACCOUNT:
		if account.Username == "" {
			return nil, ErrMissingUsername
		}
	case gen.Account_SERVICE_ACCOUNT:
		if account.Username != "" {
			return nil, status.Error(codes.InvalidArgument, "service accounts cannot have a username")
		}
	default:
		return nil, ErrInvalidAccountKind
	}

	var created *gen.Account
	err := s.store.Write(ctx, func(tx *WriteTx) error {
		var err error
		switch req.Account.Kind {
		case gen.Account_USER_ACCOUNT:
			created, err = tx.CreateUserAccount(ctx, account.Username, account.DisplayName)
		case gen.Account_SERVICE_ACCOUNT:
			created, err = tx.CreateServiceAccount(ctx, account.DisplayName)
		default:
			panic("already validated account kind")
		}
		if err != nil {
			return err
		}

		if len(account.RoleAssignments) > 0 {
			err = tx.UpdateRoleAssignments(ctx, created.Id, account.RoleAssignments)
			if err != nil {
				return err
			}
		}

		created, err = tx.GetAccount(ctx, created.Id)
		return err
	})
	if err != nil {
		return nil, err
	}

	return &gen.CreateAccountResponse{Account: account}, nil
}

func (s *Server) GetRole(ctx context.Context, req *gen.GetRoleRequest) (*gen.GetRoleResponse, error) {
	var role *gen.Role
	err := s.store.Read(ctx, func(tx *ReadTx) error {
		var err error
		role, err = tx.GetRole(ctx, req.Id)
		return err
	})
	if err != nil {
		return nil, err
	}

	return &gen.GetRoleResponse{Role: role}, nil
}

func (s *Server) CreateRole(ctx context.Context, req *gen.CreateRoleRequest) (*gen.CreateRoleResponse, error) {
	if req.Role == nil {
		return nil, status.Error(codes.InvalidArgument, "role is required")
	}

	var created *gen.Role
	err := s.store.Write(ctx, func(tx *WriteTx) error {
		var err error
		created, err = tx.CreateRole(ctx, req.Role.Title)
		if err != nil {
			return err
		}

		if len(req.Role.Permissions) > 0 {
			err = tx.UpdateRolePermissions(ctx, created.Id, req.Role.Permissions)
			if err != nil {
				return err
			}
		}
		created, err = tx.GetRole(ctx, created.Id)
		return err
	})
	if err != nil {
		return nil, err
	}

	return &gen.CreateRoleResponse{Role: created}, nil
}

const (
	minPageSize     = 1
	maxPageSize     = 100
	defaultPageSize = 30
)
