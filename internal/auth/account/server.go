package account

import (
	"context"

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
	account, err := s.store.GetAccount(ctx, req.Id)
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

	accounts, nextPage, err := s.store.ListAccounts(ctx, req.PageToken, int64(pageSize))
	if err != nil {
		return nil, err
	}

	return &gen.ListAccountsResponse{
		Accounts:      accounts,
		NextPageToken: nextPage,
	}, nil
}

func (s *Server) CreateAccount(ctx context.Context, req *gen.CreateAccountRequest) (*gen.CreateAccountResponse, error) {
	account, err := s.store.CreateAccount(ctx, req.Account)
	if err != nil {
		return nil, err
	}

	return &gen.CreateAccountResponse{Account: account}, nil
}

func (s *Server) GetRole(ctx context.Context, req *gen.GetRoleRequest) (*gen.GetRoleResponse, error) {
	role, err := s.store.GetRole(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &gen.GetRoleResponse{Role: role}, nil
}

func (s *Server) CreateRole(ctx context.Context, req *gen.CreateRoleRequest) (*gen.CreateRoleResponse, error) {
	role, err := s.store.CreateRole(ctx, req.Role)
	if err != nil {
		return nil, err
	}

	return &gen.CreateRoleResponse{Role: role}, nil
}

const (
	minPageSize     = 1
	maxPageSize     = 100
	defaultPageSize = 30
)
