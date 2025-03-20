DROP INDEX role_assignments_unique;
-- nulls are considered distinct in unique indices, so we need to use coalesce to treat them as equal
CREATE UNIQUE INDEX role_assignments_unique ON role_assignments (account_id, role_id, coalesce(scope_resource, ''), coalesce(scope_type, ''));