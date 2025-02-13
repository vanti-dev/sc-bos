package account

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/sc-bos/internal/database"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var migrations = database.MustLoadVersionedSchema(migrationsFS, "migrations")

const appID = 0x5C0501

var (
	ErrAccountNotFound        = status.Error(codes.NotFound, "account not found")
	ErrRoleNotFound           = status.Error(codes.NotFound, "role not found")
	ErrInvalidAccountKind     = status.Error(codes.InvalidArgument, "invalid account kind")
	ErrMissingUsername        = status.Error(codes.InvalidArgument, "user account requires username")
	ErrUnexpectedUsername     = status.Error(codes.InvalidArgument, "service account cannot have username")
	ErrUnexpectedServiceCreds = status.Error(codes.InvalidArgument, "user account cannot have service credentials")
	ErrUnexpectedPassword     = status.Error(codes.FailedPrecondition, "service account cannot have password")
	ErrInvalidPassword        = status.Error(codes.InvalidArgument, "password does not comply with policy")
	ErrInvalidPageSize        = status.Error(codes.InvalidArgument, "page size must be positive")
)

type Store struct {
	db *database.Database
}

func OpenStore(ctx context.Context, path string, logger *zap.Logger) (*Store, error) {
	db, err := database.Open(ctx, path,
		database.WithLogger(logger),
		database.WithApplicationID(appID),
	)
	if err != nil {
		return nil, err
	}

	err = db.Migrate(ctx, migrations)
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func NewMemoryStore(logger *zap.Logger) *Store {
	db := database.OpenMemory(
		database.WithLogger(logger),
		database.WithApplicationID(appID),
	)

	err := db.Migrate(context.Background(), migrations)
	if err != nil {
		// this can only happen if the migrations are broken
		panic(err)
	}

	return &Store{db: db}
}

func (s *Store) Close() error {
	return s.db.Close()
}

type ReadTx struct {
	readOps
}

type WriteTx struct {
	readOps
	writeOps
}

type Page[T any] struct {
	Items    []T
	NextPage string
}

type readOps struct {
	tx *sql.Tx
}

type writeOps struct {
	tx *sql.Tx
}

func (s *Store) Read(ctx context.Context, f func(*ReadTx) error) error {
	return s.db.ReadTx(ctx, func(tx *sql.Tx) error {
		readTx := &ReadTx{readOps{tx}}
		return f(readTx)
	})
}

func (s *Store) Write(ctx context.Context, f func(*WriteTx) error) error {
	return s.db.WriteTx(ctx, func(tx *sql.Tx) error {
		writeTx := &WriteTx{readOps{tx}, writeOps{tx}}
		return f(writeTx)
	})
}

func (r *readOps) GetAccount(ctx context.Context, id string) (*gen.Account, error) {
	idInt, err := parseAccountID(id)
	if err != nil {
		return nil, err
	}

	page, err := r.selectAccounts(ctx, filter{mode: matchExact, value: idInt}, 1)
	if err != nil {
		return nil, err
	}
	if len(page.Items) == 0 {
		return nil, ErrAccountNotFound
	}
	return page.Items[0], nil
}

func (r *readOps) AccountByUsername(ctx context.Context, username string) (*gen.Account, PasswordHash, error) {
	// TODO: could do all this with a single query instead of fetching the ID first
	const query = `
		SELECT a.id, p.password_hash
		FROM accounts a
		LEFT OUTER JOIN password_credentials p ON a.id = p.account_id
		WHERE a.username = ?;
	`

	var (
		id   accountID
		hash []byte
	)
	row := r.tx.QueryRowContext(ctx, query, username)
	err := row.Scan(&id, &hash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, ErrAccountNotFound
	} else if err != nil {
		return nil, nil, err
	}

	acc, err := r.GetAccount(ctx, id.String())
	if err != nil {
		return nil, nil, err
	}
	return acc, hash, nil
}

func (r *readOps) ListAccounts(ctx context.Context, pageToken string, pageSize int64) (Page[*gen.Account], error) {
	empty := Page[*gen.Account]{}
	if pageSize <= 0 {
		return empty, ErrInvalidPageSize
	}

	var f filter
	if pageToken != "" {
		idInt, err := parseAccountID(pageToken)
		if err != nil {
			return empty, err
		}
		f = filter{mode: matchAfter, value: idInt}
	}

	return r.selectAccounts(ctx, f, pageSize)
}

func (r *readOps) selectAccounts(ctx context.Context, filter filter, limit int64) (Page[*gen.Account], error) {
	empty := Page[*gen.Account]{}

	var whereClause string
	switch filter.mode {
	case matchAny:
		whereClause = ""
	case matchExact:
		whereClause = "WHERE id = ?1"
	case matchAfter:
		whereClause = "WHERE id > ?1"
	}
	query := fmt.Sprintf(`
		SELECT id, display_name, kind, create_time, username
		FROM accounts
		%s
		ORDER BY id
		LIMIT ?2;
    `, whereClause)

	rows, err := r.tx.QueryContext(ctx, query, filter.value, limit)
	if err != nil {
		return empty, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var accounts []*gen.Account
	for rows.Next() {
		var (
			id          accountID
			displayName string
			kind        string
			username    sql.NullString
			createTime  database.Timestamp
		)
		err = rows.Scan(&id, &displayName, &kind, &createTime, &username)
		if err != nil {
			return empty, err
		}

		account := &gen.Account{
			Id:          id.String(),
			CreateTime:  timestamppb.New(time.Time(createTime)),
			Kind:        gen.Account_Kind(gen.Account_Kind_value[kind]), // default to zero value ACCOUNT_KIND_UNSPECIFIED
			DisplayName: displayName,
		}
		if username.Valid {
			account.Username = username.String
		}
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		return empty, err
	}

	var nextPage string
	if len(accounts) > 0 {
		nextPage = accounts[len(accounts)-1].Id
	}
	return Page[*gen.Account]{accounts, nextPage}, nil
}

func (w *writeOps) CreateUserAccount(ctx context.Context, username, displayName string) (*gen.Account, error) {
	if username == "" {
		return nil, ErrMissingUsername
	}

	acct, err := createAccount(ctx, w.tx, gen.Account_USER_ACCOUNT, username, displayName)
	if err != nil {
		return nil, err
	}

	return acct, nil
}

func (w *writeOps) CreateServiceAccount(ctx context.Context, displayName string) (*gen.Account, error) {
	acct, err := createAccount(ctx, w.tx, gen.Account_SERVICE_ACCOUNT, "", displayName)
	if err != nil {
		return nil, err
	}

	return acct, nil
}

func (w *writeOps) UpdateAccountPasswordHash(ctx context.Context, id string, hash []byte) error {
	parsedID, err := parseAccountID(id)
	if err != nil {
		return err
	}
	return updateAccountPasswordHash(ctx, w.tx, parsedID, hash)
}

func (r *readOps) GetRole(ctx context.Context, id string) (*gen.Role, error) {
	parsedID, err := parseRoleID(id)
	if err != nil {
		return nil, err
	}

	page, err := r.selectRoles(ctx, filter{mode: matchExact, value: parsedID}, 1)
	if err != nil {
		return nil, err
	}
	if len(page.Items) == 0 {
		return nil, ErrRoleNotFound
	}
	return page.Items[0], nil
}

func (r *readOps) ListRoles(ctx context.Context, pageToken string, pageSize int64) (Page[*gen.Role], error) {
	empty := Page[*gen.Role]{}
	if pageSize <= 0 {
		return empty, ErrInvalidPageSize
	}

	var f filter
	if pageToken != "" {
		idInt, err := parseRoleID(pageToken)
		if err != nil {
			return empty, err
		}
		f = filter{mode: matchAfter, value: idInt}
	}

	return r.selectRoles(ctx, f, pageSize)
}

func (r *readOps) selectRoles(ctx context.Context, filter filter, limit int64) (Page[*gen.Role], error) {
	empty := Page[*gen.Role]{}

	var whereClause string
	switch filter.mode {
	case matchAny:
		whereClause = ""
	case matchExact:
		whereClause = "WHERE id = ?1"
	case matchAfter:
		whereClause = "WHERE id > ?1"
	}
	query := fmt.Sprintf(`
		SELECT r.id, r.name, coalesce(group_concat(p.permission, x'00'), '')
		FROM roles r
		LEFT OUTER JOIN role_permissions p ON r.id = p.role_id	
		%s
		GROUP BY r.id
		ORDER BY r.id
		LIMIT ?2;
	`, whereClause)

	rows, err := r.tx.QueryContext(ctx, query, filter.value, limit)
	if err != nil {
		return empty, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var roles []*gen.Role
	for rows.Next() {
		var (
			id          roleID
			name        string
			permissions string
		)
		err = rows.Scan(&id, &name)
		if err != nil {
			return empty, err
		}
		roles = append(roles, &gen.Role{
			Id:          id.String(),
			Title:       name,
			Permissions: strings.Split(permissions, "\x00"),
		})
	}
	if err := rows.Err(); err != nil {
		return empty, err
	}

	var nextPage string
	if len(roles) > 0 {
		nextPage = roles[len(roles)-1].Id
	}
	return Page[*gen.Role]{roles, nextPage}, nil
}

func (w *writeOps) CreateRole(ctx context.Context, title string) (*gen.Role, error) {
	const queryCreateRole = `
		INSERT INTO roles (name)
		VALUES (?)
		RETURNING id;
	`

	var id roleID
	row := w.tx.QueryRowContext(ctx, queryCreateRole, title)
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}

	return &gen.Role{
		Id:    id.String(),
		Title: title,
	}, nil
}

func (w *writeOps) UpdateRoleName(ctx context.Context, id string, name string) error {
	parsedID, err := parseRoleID(id)
	if err != nil {
		return err
	}

	err = checkRoleExists(ctx, w.tx, parsedID)
	if err != nil {
		return err
	}
	const query = `
		UPDATE roles
		SET name = ?
		WHERE id = ?;
	`

	_, err2 := w.tx.ExecContext(ctx, query, name, parsedID)
	return err2
}

func (w *writeOps) UpdateRolePermissions(ctx context.Context, id string, permissions []string) error {
	parsedID, err := parseRoleID(id)
	if err != nil {
		return err
	}

	err = checkRoleExists(ctx, w.tx, parsedID)
	if err != nil {
		return err
	}
	return replaceRolePermissions(ctx, w.tx, parsedID, permissions)
}

func (w *writeOps) DeleteRole(ctx context.Context, id string) error {
	parsedID, err := parseRoleID(id)
	if err != nil {
		return err
	}

	const queryDeleteRole = `
		DELETE FROM roles
		WHERE id = ?;
	`

	err = checkRoleExists(ctx, w.tx, parsedID)
	if err != nil {
		return err
	}

	_, err = w.tx.ExecContext(ctx, queryDeleteRole, parsedID)
	return err
}

func createAccount(ctx context.Context, tx *sql.Tx, kind gen.Account_Kind, username, displayName string) (*gen.Account, error) {
	const query = `
		INSERT INTO accounts (username, display_name, kind, create_time)
		VALUES (?, ?, ?, datetime('now', 'subsec'))
		RETURNING id, create_time;
	`

	var usernameValue sql.NullString
	if username != "" {
		usernameValue.String = username
		usernameValue.Valid = true
	}
	kindStr := gen.Account_Kind_name[int32(kind)]
	row := tx.QueryRowContext(ctx, query, usernameValue, displayName, kindStr)

	var (
		id      accountID
		created database.Timestamp
	)
	err := row.Scan(&id, &created)
	if err != nil {
		return nil, err
	}
	return &gen.Account{
		Id:          id.String(),
		CreateTime:  timestamppb.New(time.Time(created)),
		Kind:        kind,
		DisplayName: displayName,
		Username:    username,
	}, nil
}

func checkAccountKind(ctx context.Context, tx *sql.Tx, id accountID, expect gen.Account_Kind, mismatch error) error {
	const query = `
		SELECT kind
		FROM accounts
		WHERE id = ?;
	`

	var kind string
	err := tx.QueryRowContext(ctx, query, id).Scan(&kind)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrAccountNotFound
	} else if err != nil {
		return err
	}
	if kind != gen.Account_Kind_name[int32(expect)] {
		return mismatch
	}
	return nil
}

func updateAccountPasswordHash(ctx context.Context, tx *sql.Tx, id accountID, hash []byte) error {
	const query = `
		INSERT INTO password_credentials (account_id, password_hash)
		VALUES (?, ?)
		ON CONFLICT (account_id) DO UPDATE
		SET password_hash = excluded.password_hash;
	`

	err := checkAccountKind(ctx, tx, id, gen.Account_USER_ACCOUNT, ErrUnexpectedPassword)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, hash, id)
	return err
}

func checkRoleExists(ctx context.Context, tx *sql.Tx, id roleID) error {
	const query = `
		SELECT 1
		FROM roles
		WHERE id = ?;
	`

	var unused int
	err := tx.QueryRowContext(ctx, query, id).Scan(&unused)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrRoleNotFound
	}
	return err
}

func replaceRolePermissions(ctx context.Context, tx *sql.Tx, id roleID, permissions []string) error {
	const queryDeletePermissions = `
		DELETE FROM role_permissions
		WHERE role_id = ?;
	`

	_, err := tx.ExecContext(ctx, queryDeletePermissions, id)
	if err != nil {
		return err
	}

	const queryInsertPermission = `
		INSERT INTO role_permissions (role_id, permission)
		VALUES (?, ?);
	`

	stmt, err := tx.PrepareContext(ctx, queryInsertPermission)
	if err != nil {
		return err
	}
	for _, perm := range permissions {
		_, err = stmt.ExecContext(ctx, id, perm)
		if err != nil {
			return err
		}
	}
	return nil
}

type accountID int64

func (id accountID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

func parseAccountID(idStr string) (accountID, error) {
	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		// an account ID must be a valid integer
		// so any other account ID can't possibly exist in the DB
		return 0, ErrAccountNotFound
	}
	return accountID(idInt), nil
}

type roleID int64

func (id roleID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

func parseRoleID(idStr string) (roleID, error) {
	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		// a role ID must be a valid integer
		// so any other role ID can't possibly exist in the DB
		return 0, ErrRoleNotFound
	}
	return roleID(idInt), nil
}

type filter struct {
	mode  matchMode
	value any
}

type matchMode int

const (
	matchAny   matchMode = iota // matches everything, ignoring the match value
	matchExact                  // matches the exact value
	matchAfter                  // matches values that sort after the given value
)
