-- legacy_role column specifies which fixed legacy roles this role maps to,
-- while we are still using the old authentication method.
ALTER TABLE roles
ADD COLUMN legacy_role TEXT;

-- protected can be set to true which will prevent the role from being modified or deleted using the API.
-- intended for built-in roles that are required for the system to function properly.
ALTER TABLE roles
ADD COLUMN protected BOOLEAN NOT NULL DEFAULT FALSE;

-- Insert the built-in roles so that legacy roles can be referenced.
-- Could fail if roles with these names already exist - will require manual intervention in that case to rename the roles
-- before the migration can proceed.
INSERT INTO roles (display_name, description, legacy_role, protected) VALUES
   ('Admin', 'Full system access (built-in role)', 'admin', TRUE),
   ('Super Admin', 'Full system access (built-in role)', 'super-admin', TRUE),
   ('Commissioner', 'Alter configurations (built-in role)', 'commissioner', TRUE),
   ('Operator', 'View data and control devices (built-in role)', 'operator', TRUE),
   ('Viewer', 'View data (built-in role)', 'viewer', TRUE);