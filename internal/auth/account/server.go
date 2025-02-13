package account

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vanti-dev/sc-bos/pkg/gen"
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
	var account *gen.Account
	err := s.store.Read(ctx, func(tx *ReadTx) error {
		var err error
		account, err = tx.GetAccount(ctx, req.Id)
		return err
	})

	if err != nil {
		return nil, err
	}
	return account, nil
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

	var res *gen.ListAccountsResponse
	err := s.store.Read(ctx, func(tx *ReadTx) error {
		page, err := tx.ListAccounts(ctx, req.PageToken, int64(pageSize))
		if err != nil {
			return err
		}

		res = &gen.ListAccountsResponse{
			Accounts:      page.Items,
			NextPageToken: page.NextPage,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Server) CreateAccount(ctx context.Context, req *gen.CreateAccountRequest) (*gen.Account, error) {
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
			return nil, ErrUnexpectedUsername
		}
		if req.Password != "" {
			return nil, ErrUnexpectedPassword
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

		if req.Password != "" {
			err = s.setPassword(ctx, tx, created.Id, req.Password)
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

	return created, nil
}

func (s *Server) GetRole(ctx context.Context, req *gen.GetRoleRequest) (*gen.Role, error) {
	var role *gen.Role
	err := s.store.Read(ctx, func(tx *ReadTx) error {
		var err error
		role, err = tx.GetRole(ctx, req.Id)
		return err
	})
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *Server) CreateRole(ctx context.Context, req *gen.CreateRoleRequest) (*gen.Role, error) {
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

	return created, nil
}

func (s *Server) setPassword(ctx context.Context, tx *WriteTx, id, password string) error {
	if !permitPassword(password) {
		return ErrInvalidPassword
	}

	hash, err := HashPassword(password)
	if err != nil {
		s.logger.Error("failed to hash password", zap.String("userID", id), zap.Error(err))
		return status.Error(codes.Internal, "failed to set password")
	}

	return tx.UpdateAccountPasswordHash(ctx, id, hash)
}

const (
	minPageSize     = 1
	maxPageSize     = 100
	defaultPageSize = 30
)
