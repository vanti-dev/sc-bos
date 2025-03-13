-- name: GetAccount :one
SELECT *
FROM accounts
WHERE id = :id;

-- name: ListAccounts :many
SELECT *
FROM accounts
WHERE id > :after_id
ORDER BY id
LIMIT :limit;

-- name: CountAccounts :one
SELECT COUNT(*) AS count
FROM accounts;

-- name: GetAccountByUsername :one
SELECT *
FROM accounts
WHERE username = :username;

-- name: CreateUserAccount :one
INSERT INTO accounts (username, display_name, description, type, create_time)
VALUES (:username, :display_name, :description, 'USER_ACCOUNT', datetime('now', 'subsec'))
RETURNING *;

-- name: CreateServiceAccount :one
INSERT INTO accounts (display_name, description, type, create_time)
VALUES (:display_name, :description, 'SERVICE_ACCOUNT', datetime('now', 'subsec'))
RETURNING *;

-- name: GetAccountPasswordHash :one
SELECT password_hash
FROM password_credentials
WHERE account_id = :account_id;

-- name: UpdateAccountPasswordHash :exec
INSERT INTO password_credentials (account_id, password_hash)
VALUES (:account_id, :password_hash)
ON CONFLICT (account_id) DO UPDATE
SET password_hash = excluded.password_hash;

-- name: UpdateAccountDisplayName :exec
UPDATE accounts
SET display_name = :display_name
WHERE id = :id;

-- name: UpdateAccountUsername :exec
UPDATE accounts
SET username = :username
WHERE id = :id;

-- name: UpdateAccountDescription :exec
UPDATE accounts
SET description = :description
WHERE id = :id;

-- name: DeleteAccount :execrows
DELETE FROM accounts
WHERE id = :id;

-- name: CreateServiceCredential :one
INSERT INTO service_credentials (account_id, display_name, description, secret_hash, create_time, expire_time)
VALUES (:account_id, :display_name, :description, :secret_hash, datetime('now', 'subsec'), :expire_time)
RETURNING *;

-- name: GetServiceCredential :one
SELECT *
FROM service_credentials
WHERE id = :id;

-- name: ListAccountServiceCredentials :many
SELECT *
FROM service_credentials
WHERE account_id = :account_id
ORDER BY id;

-- name: DeleteServiceCredential :execrows
DELETE FROM service_credentials
WHERE id = :id;

-- name: CountServiceCredentialsForAccount :one
SELECT COUNT(*) AS count
FROM service_credentials
WHERE account_id = :account_id;

-- name: GetRole :one
SELECT *
FROM roles
WHERE id = :id;

-- name: ListRoles :many
SELECT *
FROM roles
WHERE id > :after_id
ORDER BY id
LIMIT :limit;

-- name: CountRoles :one
SELECT COUNT(*) AS count
FROM roles;

-- name: ListRolesAndPermissions :many
SELECT sqlc.embed(roles), group_concat(coalesce(role_permissions.permission, ''), ',') AS permissions
FROM roles
LEFT OUTER JOIN role_permissions ON roles.id = role_permissions.role_id
WHERE roles.id > :after_id
GROUP BY roles.id
ORDER BY roles.id
LIMIT :limit;

-- name: CreateRole :one
INSERT INTO roles (display_name, description)
VALUES (:display_name, :description)
RETURNING *;

-- name: UpdateRoleDisplayName :execrows
UPDATE roles
SET display_name = :display_name
WHERE id = :id;

-- name: UpdateRoleDescription :execrows
UPDATE roles
SET description = :description
WHERE id = :id;

-- name: DeleteRole :execrows
DELETE FROM roles
WHERE id = :id;

-- name: ListRolePermissions :many
SELECT permission
FROM role_permissions
WHERE role_id = :role_id
ORDER BY permission;

-- name: AddRolePermission :exec
INSERT INTO role_permissions (role_id, permission)
VALUES (:role_id, :permission)
ON CONFLICT (role_id, permission) DO NOTHING;

-- name: DeleteRolePermission :execrows
DELETE FROM role_permissions
WHERE role_id = :role_id AND permission = :permission;

-- name: ClearRolePermissions :execrows
DELETE FROM role_permissions
WHERE role_id = :role_id;

-- name: GetRoleAssignment :one
SELECT *
FROM role_assignments
WHERE id = :id;

-- name: ListRoleAssignments :many
SELECT *
FROM role_assignments
WHERE id > :after_id
ORDER BY id
LIMIT :limit;

-- name: CountRoleAssignments :one
SELECT COUNT(*)
FROM role_assignments;

-- name: ListRoleAssignmentsForAccount :many
SELECT *
FROM role_assignments
WHERE account_id = :account_id
  AND id > :after_id
ORDER BY id
LIMIT :limit;

-- name: CountRoleAssignmentsForAccount :one
SELECT COUNT(*) AS count
FROM role_assignments
WHERE account_id = :account_id;

-- name: ListRoleAssignmentsForRole :many
SELECT *
FROM role_assignments
WHERE role_id = :role_id
  AND id > :after_id
ORDER BY id
LIMIT :limit;

-- name: CountRoleAssignmentsForRole :one
SELECT COUNT(*) AS count
FROM role_assignments
WHERE role_id = :role_id;

-- name: CreateRoleAssignment :one
INSERT INTO role_assignments (account_id, role_id, scope_type, scope_resource)
VALUES (:account_id, :role_id, :scope_kind, :scope_resource)
RETURNING *;

-- name: DeleteRoleAssignment :execrows
DELETE FROM role_assignments
WHERE id = :id;