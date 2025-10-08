-- name: GetAccount :one
SELECT *
FROM accounts
WHERE id = :id;

-- name: GetAccountDetails :one
SELECT * FROM account_details
WHERE id = :id;

-- name: GetAccountByUsername :one
SELECT * FROM user_accounts
WHERE username = :username;

-- name: ListAccounts :many
SELECT *
FROM accounts
WHERE id > :after_id
ORDER BY id
LIMIT :limit;

-- name: ListAccountDetails :many
SELECT * FROM account_details
WHERE id > :after_id
ORDER BY id
LIMIT :limit;

-- name: CountAccounts :one
SELECT COUNT(*) AS count
FROM accounts;

-- name: CreateAccount :one
INSERT INTO accounts (display_name, description, type, create_time)
VALUES (:display_name, :description, :type, datetime('now', 'subsec'))
RETURNING *;

-- name: CreateUserAccount :one
INSERT INTO user_accounts (account_id, username, password_hash)
VALUES (:account_id, :username, nullif(:password_hash, x''))
RETURNING *;

-- name: CreateServiceAccount :one
INSERT INTO service_accounts (account_id, primary_secret_hash)
VALUES (:account_id, :primary_secret_hash)
RETURNING *;

-- name: UpdateAccountPasswordHash :exec
UPDATE user_accounts
SET password_hash = nullif(:password_hash, x'')
WHERE account_id = :account_id;

-- name: RotateServiceAccountSecret :exec
UPDATE service_accounts
SET primary_secret_hash = :primary_secret_hash,
    -- if no secondary_secret_expire_time is supplied, then we don't want a secondary secret
    secondary_secret_hash = CASE
        WHEN :secondary_secret_expire_time IS NULL THEN NULL
        ELSE primary_secret_hash
    END,
    secondary_secret_expire_time = :secondary_secret_expire_time
WHERE account_id = :account_id;

-- name: UpdateAccountDisplayName :exec
UPDATE accounts
SET display_name = :display_name
WHERE id = :id;

-- name: UpdateAccountUsername :exec
UPDATE user_accounts
SET username = :username
WHERE account_id = :account_id;

-- name: UpdateAccountDescription :exec
UPDATE accounts
SET description = :description
WHERE id = :id;

-- name: DeleteAccount :execrows
DELETE FROM accounts
WHERE id = :id;

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

-- name: ListRolesWithLegacyRole :many
SELECT *
FROM roles
WHERE legacy_role = :legacy_role
ORDER BY id;

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
-- refuse to update a role which is protected
UPDATE roles
SET display_name = :display_name
WHERE id = :id AND NOT protected;

-- name: UpdateRoleDescription :execrows
-- refuse to update a role which is protected
UPDATE roles
SET description = :description
WHERE id = :id AND NOT protected;

-- name: DeleteRole :execrows
-- refuse to update a role which is protected
DELETE FROM roles
WHERE id = :id AND NOT protected;

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

-- name: ListPermissionsForAccount :many
SELECT DISTINCT rp.permission, ra.scope_type, ra.scope_resource
FROM role_assignments ra
INNER JOIN role_permissions rp ON ra.role_id = rp.role_id
WHERE ra.account_id = :account_id
  AND rp.permission IS NOT NULL
ORDER BY rp.permission, ra.scope_type, ra.scope_resource;

-- name: ListLegacyRolesForAccount :many
SELECT DISTINCT r.legacy_role
FROM role_assignments ra
INNER JOIN roles r ON ra.role_id = r.id
WHERE ra.account_id = :account_id
  AND r.legacy_role IS NOT NULL
ORDER BY r.legacy_role;

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