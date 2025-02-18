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

-- name: GetAccountByUsername :one
SELECT sqlc.embed(accounts), password_credentials.password_hash
FROM accounts
LEFT OUTER JOIN password_credentials ON accounts.id = password_credentials.account_id
WHERE username = :username;

-- name: CreateUserAccount :one
INSERT INTO accounts (username, display_name, kind, create_time)
VALUES (:username, :display_name, 'USER_ACCOUNT', datetime('now', 'subsec'))
RETURNING *;

-- name: CreateServiceAccount :one
INSERT INTO accounts (display_name, kind, create_time)
VALUES (:display_name, 'SERVICE_ACCOUNT', datetime('now', 'subsec'))
RETURNING *;

-- name: UpdateAccountPasswordHash :exec
INSERT INTO password_credentials (account_id, password_hash)
VALUES (:account_id, :password_hash)
ON CONFLICT (account_id) DO UPDATE
SET password_hash = :password_hash;

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

-- name: ListRolesWithPermissions :many
SELECT sqlc.embed(roles), group_concat(role_permissions.permission, ',') AS permissions
FROM roles
LEFT OUTER JOIN role_permissions ON roles.id = role_permissions.role_id
WHERE roles.id > :after_id
GROUP BY roles.id
ORDER BY roles.id
LIMIT :limit;

-- name: CreateRole :one
INSERT INTO roles (name)
VALUES (:name)
RETURNING *;

-- name: UpdateRoleName :execrows
UPDATE roles
SET name = :name
WHERE id = :id;

-- name: DeleteRole :execrows
DELETE FROM roles
WHERE id = :id;

-- name: ListPermissionsForRole :many
SELECT *
FROM role_permissions
WHERE role_id = :role_id
ORDER BY permission;
