package account

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"math"
	"strconv"
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

	var account *gen.Account
	account, err = getAccountOnly(ctx, r.tx, idInt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAccountNotFound
	} else if err != nil {
		return nil, err
	}

	account.ServiceCredentials, err = getAccountServiceCreds(ctx, r.tx, idInt)
	if err != nil {
		return nil, err
	}

	account.RoleAssignments, err = getAccountRoleAssignments(ctx, r.tx, idInt)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (r *readOps) ListAccounts(ctx context.Context, pageToken string, pageSize int64) (accounts []*gen.Account, nextPage string, err error) {
	if pageSize <= 0 {
		return nil, "", ErrInvalidPageSize
	}

	var minAcctID accountID
	if pageToken != "" {
		idInt, err := parseAccountID(pageToken)
		if err != nil {
			return nil, "", err
		}
		minAcctID = idInt
	}

	var foundRange accountRange
	accounts, foundRange, err = listAccountsOnly(ctx, r.tx, minAcctID, pageSize)
	if err != nil {
		return nil, "", err
	}

	creds, err := listAccountServiceCreds(ctx, r.tx, foundRange)
	if err != nil {
		return nil, "", err
	}

	roleAssignments, err := listRoleAssignments(ctx, r.tx, foundRange)
	if err != nil {
		return nil, "", err
	}

	for _, acct := range accounts {
		acct.ServiceCredentials = creds[acct.Id]
		acct.RoleAssignments = roleAssignments[acct.Id]
	}

	if len(accounts) > 0 {
		nextPage = accounts[len(accounts)-1].Id
	}

	return accounts, nextPage, nil
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

func (w *writeOps) UpdateRoleAssignments(ctx context.Context, accountID string, roleAssignments []*gen.RoleAssignment) error {
	id, err := parseAccountID(accountID)
	if err != nil {
		return err
	}
	return updateAccountRoleAssignments(ctx, w.tx, int64(id), true, roleAssignments)
}

func (r *readOps) GetRole(ctx context.Context, id string) (*gen.Role, error) {
	parsedID, err := parseRoleID(id)
	if err != nil {
		return nil, err
	}

	const querySelectRole = `
		SELECT name
		FROM roles
		WHERE id = ?;
	`

	const querySelectPermissions = `
		SELECT permission
		FROM role_permissions
		WHERE role_id = ?
		ORDER BY permission;
	`

	var (
		name  string
		perms []string
	)
	row := r.tx.QueryRowContext(ctx, querySelectRole, parsedID)
	err = row.Scan(&name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRoleNotFound
	} else if err != nil {
		return nil, err
	}

	rows, err := r.tx.QueryContext(ctx, querySelectPermissions, parsedID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var perm string
		err = rows.Scan(&perm)
		if err != nil {
			return nil, err
		}
		perms = append(perms, perm)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return &gen.Role{
		Id:          parsedID.String(),
		Title:       name,
		Permissions: perms,
	}, nil
}

func (r *readOps) ListRoles(ctx context.Context, pageToken string, pageSize int64) (roles []*gen.Role, nextPage string, err error) {
	if pageSize <= 0 {
		return nil, "", ErrInvalidPageSize
	}

	var minRoleID roleID
	if pageToken != "" {
		idInt, err := parseRoleID(pageToken)
		if err != nil {
			return nil, "", err
		}
		minRoleID = idInt
	}

	const queryRoles = `
		SELECT r.id, r.name, p.permission
		FROM roles r 
		LEFT OUTER JOIN role_permissions p ON r.id = p.role_id
		WHERE r.id > ?
		ORDER BY r.id
		LIMIT ?;
	`

	rows, err := r.tx.QueryContext(ctx, queryRoles, minRoleID, pageSize)
	if err != nil {
		return nil, "", err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var current *gen.Role
	for rows.Next() {
		var (
			id         roleID
			name       string
			permission sql.NullString
		)
		err = rows.Scan(&id, &name, &permission)
		if err != nil {
			return nil, "", err
		}

		// flush the current role if the rows iterator has moved to a new role
		if current == nil {
			current = &gen.Role{Id: id.String(), Title: name}
		} else if current.Id != id.String() {
			roles = append(roles, current)
			current = &gen.Role{Id: id.String(), Title: name}
		}

		if permission.Valid {
			current.Permissions = append(current.Permissions, permission.String)
		}
	}
	if current != nil {
		roles = append(roles, current)
	}
	if err := rows.Err(); err != nil {
		return nil, "", err
	}

	if len(roles) > 0 {
		nextPage = roles[len(roles)-1].Id
	}

	return roles, nextPage, nil
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
	return updateRoleName(ctx, w.tx, parsedID, name)
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

func getAccountOnly(ctx context.Context, tx *sql.Tx, id accountID) (*gen.Account, error) {
	const querySelectAccount = `
		SELECT display_name, kind, create_time, username
		FROM accounts 
		WHERE id = ?;
	`

	var (
		displayName string
		kind        string
		username    sql.NullString
		createTime  database.Timestamp
	)
	row := tx.QueryRowContext(ctx, querySelectAccount, id)
	err := row.Scan(&displayName, &kind, &createTime, &username)
	if err != nil {
		return nil, err
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
	return account, nil
}

// lists accounts without populating service credentials or role assignments
func listAccountsOnly(ctx context.Context, tx *sql.Tx, startAfter accountID, limit int64) ([]*gen.Account, accountRange, error) {
	const query = `
		SELECT id, display_name, kind, create_time, username
		FROM accounts
		WHERE id > ?
		ORDER BY id
		LIMIT ?;
    `

	rows, err := tx.QueryContext(ctx, query, startAfter, limit)
	if err != nil {
		return nil, accountRange{}, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var accounts []*gen.Account
	foundRange := accountRange{Min: math.MaxInt64}
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
			return nil, accountRange{}, err
		}
		if id < foundRange.Min {
			foundRange.Min = id
		}
		if id > foundRange.Max {
			foundRange.Max = id
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
		return nil, accountRange{}, err
	}
	return accounts, foundRange, nil
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

func getAccountServiceCreds(ctx context.Context, tx *sql.Tx, id accountID) ([]*gen.ServiceCredential, error) {
	const querySelectServiceCreds = `
		SELECT id, title, create_time, expire_time
		FROM service_credentials
		WHERE account_id = ?
		ORDER BY id;
	`

	var serviceCreds []*gen.ServiceCredential
	rows, err := tx.QueryContext(ctx, querySelectServiceCreds, id)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var (
			id         int64
			title      string
			createTime database.Timestamp
			expireTime sql.Null[database.Timestamp]
		)

		err = rows.Scan(&id, &title, &createTime, &expireTime)
		if err != nil {
			return nil, err
		}

		cred := &gen.ServiceCredential{
			Id:         strconv.FormatInt(id, 10),
			AccountId:  strconv.FormatInt(id, 10),
			Title:      title,
			CreateTime: timestamppb.New(time.Time(createTime)),
		}
		if expireTime.Valid {
			cred.ExpireTime = timestamppb.New(time.Time(expireTime.V))
		}
		serviceCreds = append(serviceCreds, cred)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return serviceCreds, nil
}

func listAccountServiceCreds(ctx context.Context, tx *sql.Tx, accounts accountRange) (map[string][]*gen.ServiceCredential, error) {
	const query = `
		SELECT account_id, id, title, create_time, expire_time
		FROM service_credentials
		WHERE account_id >= ? AND account_id <= ?
		ORDER BY account_id, id;
    `

	rows, err := tx.QueryContext(ctx, query, accounts.Min, accounts.Max)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	creds := make(map[string][]*gen.ServiceCredential)
	for rows.Next() {
		var (
			account    accountID
			id         int64
			title      string
			createTime database.Timestamp
			expireTime sql.Null[database.Timestamp]
		)
		err = rows.Scan(&account, &id, &title, &createTime, &expireTime)
		if err != nil {
			return nil, err
		}
		accountIDStr := account.String()
		cred := &gen.ServiceCredential{
			Id:         strconv.FormatInt(id, 10),
			AccountId:  accountIDStr,
			Title:      title,
			CreateTime: timestamppb.New(time.Time(createTime)),
		}
		if expireTime.Valid {
			cred.ExpireTime = timestamppb.New(time.Time(expireTime.V))
		}
		creds[accountIDStr] = append(creds[accountIDStr], cred)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return creds, nil
}

func getAccountRoleAssignments(ctx context.Context, tx *sql.Tx, id accountID) ([]*gen.RoleAssignment, error) {
	const queryRoleAssignments = `
		SELECT role_id, scope_kind, scope_resource
		FROM role_assignments
		WHERE account_id = ?
		ORDER BY role_id;
	`

	var roleAssignments []*gen.RoleAssignment
	rows, err := tx.QueryContext(ctx, queryRoleAssignments, id)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var (
			roleID        int64
			scopeKind     sql.NullString
			scopeResource sql.NullString
		)

		err = rows.Scan(&roleID, &scopeKind, &scopeResource)
		if err != nil {
			return nil, err
		}

		var scope *gen.RoleAssignment_Scope
		if scopeKind.Valid {
			scope = &gen.RoleAssignment_Scope{
				// default to zero value RESOURCE_KIND_UNSPECIFIED if the value is not recognized
				ResourceKind: gen.RoleAssignment_ResourceKind(gen.RoleAssignment_ResourceKind_value[scopeKind.String]),
				Resource:     scopeResource.String,
			}
		}
		ra := &gen.RoleAssignment{
			Role:  strconv.FormatInt(roleID, 10),
			Scope: scope,
		}
		roleAssignments = append(roleAssignments, ra)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return roleAssignments, nil
}

func listRoleAssignments(ctx context.Context, tx *sql.Tx, accounts accountRange) (map[string][]*gen.RoleAssignment, error) {
	const query = `
		SELECT account_id, role_id, scope_kind, scope_resource
		FROM role_assignments
		WHERE account_id >= ? AND account_id <= ?
		ORDER BY account_id, role_id;
	`

	rows, err := tx.QueryContext(ctx, query, accounts.Min, accounts.Max)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	roleAssignments := make(map[string][]*gen.RoleAssignment)
	for rows.Next() {
		var (
			account       accountID
			roleID        int64
			scopeKind     sql.NullString
			scopeResource sql.NullString
		)

		err = rows.Scan(&account, &roleID, &scopeKind, &scopeResource)
		if err != nil {
			return nil, err
		}

		var scope *gen.RoleAssignment_Scope
		if scopeKind.Valid {
			scope = &gen.RoleAssignment_Scope{
				// default to zero value RESOURCE_KIND_UNSPECIFIED if the value is not recognized
				ResourceKind: gen.RoleAssignment_ResourceKind(gen.RoleAssignment_ResourceKind_value[scopeKind.String]),
				Resource:     scopeResource.String,
			}
		}
		ra := &gen.RoleAssignment{
			Role:  strconv.FormatInt(roleID, 10),
			Scope: scope,
		}
		accountIDStr := account.String()
		roleAssignments[accountIDStr] = append(roleAssignments[accountIDStr], ra)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return roleAssignments, nil
}

// updateAccountRoleAssignments updates the role assignments for an account.
// If replace is true, all existing role assignments are replaced with the supplied role assignments.
// Otherwise, the supplied role assignments are merged with the existing role assignments:
//   - If the account does not have that role ID, a new role assignment is created.
//   - If the account already has that role ID, that assignment is replaced (with the supplied scope).
func updateAccountRoleAssignments(ctx context.Context, tx *sql.Tx, accountID int64, replace bool, roleAssignments []*gen.RoleAssignment) error {
	const queryDeleteRoleAssignments = `
		DELETE FROM role_assignments
		WHERE account_id = ?;
	`

	if replace {
		_, err := tx.ExecContext(ctx, queryDeleteRoleAssignments, accountID)
		if err != nil {
			return err
		}
	}

	const queryInsertRoleAssignment = `
		INSERT INTO role_assignments (account_id, role_id, scope_kind, scope_resource)
		VALUES (?, ?, ?, ?)
		ON CONFLICT (account_id, role_id) DO UPDATE
		SET scope_kind = excluded.scope_kind, scope_resource = excluded.scope_resource;
	`

	stmt, err := tx.PrepareContext(ctx, queryInsertRoleAssignment)
	if err != nil {
		return err
	}
	for _, ra := range roleAssignments {
		var scopeKind, scopeResource sql.NullString
		if ra.Scope != nil {
			scopeKind.String = ra.Scope.ResourceKind.String()
			scopeKind.Valid = true
			scopeResource.String = ra.Scope.Resource
			scopeResource.Valid = true
		}
		_, err = stmt.ExecContext(ctx, accountID, ra.Role, scopeKind, scopeResource)
		if err != nil {
			return err
		}
	}
	return nil
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

func updateRoleName(ctx context.Context, tx *sql.Tx, id roleID, name string) error {
	const query = `
		UPDATE roles
		SET name = ?
		WHERE id = ?;
	`

	_, err := tx.ExecContext(ctx, query, name, id)
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

type accountRange struct {
	Min, Max accountID
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
