CREATE TABLE accounts (
    id              INTEGER PRIMARY KEY,
    username        TEXT,
    display_name    TEXT NOT NULL,
    description     TEXT,
    type            TEXT NOT NULL,
    create_time     DATETIME NOT NULL,

    CONSTRAINT create_time_format CHECK ( create_time IS datetime(create_time, 'subsec') )
);

CREATE UNIQUE INDEX accounts_username ON accounts (username);

CREATE TABLE service_credentials (
    id                  INTEGER PRIMARY KEY,
    account_id          INTEGER NOT NULL,
    display_name        TEXT NOT NULL,
    description         TEXT,
    secret_hash         BLOB NOT NULL,
    create_time         DATETIME NOT NULL,
    expire_time         DATETIME,

    FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE,
    CONSTRAINT create_time_format CHECK ( create_time IS datetime(create_time, 'subsec') ),
    CONSTRAINT expire_time_format CHECK ( expire_time IS datetime(expire_time, 'subsec') )
);

CREATE TABLE password_credentials (
    account_id          INTEGER PRIMARY KEY, -- at most one password credential per account
    password_hash       BLOB NOT NULL,

    FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE
);

CREATE TABLE roles (
    id           INTEGER PRIMARY KEY,
    display_name TEXT NOT NULL,
    description  TEXT
);

CREATE UNIQUE INDEX roles_display_name ON roles (display_name);

CREATE TABLE role_permissions (
    role_id INTEGER NOT NULL,
    permission TEXT NOT NULL,

    FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX role_permissions_unique ON role_permissions (role_id, permission);

CREATE TABLE role_assignments (
    id              INTEGER PRIMARY KEY,
    account_id      INTEGER NOT NULL,
    role_id         INTEGER NOT NULL,
    scope_type      TEXT,
    scope_resource  TEXT,

    FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles (id)
);

CREATE UNIQUE INDEX role_assignments_unique ON role_assignments (account_id, role_id);
