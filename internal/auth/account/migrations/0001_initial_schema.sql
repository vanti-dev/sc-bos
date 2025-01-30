CREATE TABLE accounts (
    id              INTEGER PRIMARY KEY,
    username        TEXT,
    display_name    TEXT NOT NULL,
    kind            TEXT NOT NULL,
    create_time     TEXT NOT NULL,

    CONSTRAINT create_time_format CHECK ( create_time IS datetime(create_time, 'subsec') )
) STRICT;

CREATE UNIQUE INDEX accounts_username ON accounts (username);

CREATE TABLE service_credentials (
    id                  INTEGER PRIMARY KEY,
    account_id          INTEGER,
    title               TEXT NOT NULL,
    secret_hash         TEXT NOT NULL,
    create_time         TEXT NOT NULL,
    expire_time         TEXT,

    FOREIGN KEY (account_id) REFERENCES accounts (id),
    CONSTRAINT create_time_format CHECK ( create_time IS datetime(create_time, 'subsec') ),
    CONSTRAINT expire_time_format CHECK ( expire_time IS datetime(expire_time, 'subsec') )
) STRICT;

CREATE TABLE password_credentials (
    account_id          INTEGER PRIMARY KEY, -- at most one password credential per account
    password_hash       TEXT NOT NULL,

    FOREIGN KEY (account_id) REFERENCES accounts (id)
) STRICT;

CREATE TABLE roles (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
) STRICT;

CREATE UNIQUE INDEX roles_name ON roles (name);

CREATE TABLE role_permissions (
    role_id INTEGER NOT NULL,
    permission TEXT NOT NULL,

    FOREIGN KEY (role_id) REFERENCES roles (id)
) STRICT;

CREATE UNIQUE INDEX role_permissions_unique ON role_permissions (role_id, permission);

CREATE TABLE role_assignments (
    account_id      INTEGER NOT NULL,
    role_id         INTEGER NOT NULL,
    scope_kind      TEXT,
    scope_resource  TEXT,

    FOREIGN KEY (account_id) REFERENCES accounts (id),
    FOREIGN KEY (role_id) REFERENCES roles (id)
) STRICT;

CREATE UNIQUE INDEX role_assignments_unique ON role_assignments (account_id, role_id);
