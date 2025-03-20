PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
PRAGMA application_id = 0x5C0501;
PRAGMA user_version = 1;
CREATE TABLE accounts (
    id              INTEGER PRIMARY KEY,
    username        TEXT,
    display_name    TEXT NOT NULL,
    description     TEXT,
    type            TEXT NOT NULL,
    create_time     DATETIME NOT NULL,

    CONSTRAINT create_time_format CHECK ( create_time IS datetime(create_time, 'subsec') )
);
INSERT INTO accounts VALUES(2,NULL,'My Service 0','A service with 0 service credentials','SERVICE_ACCOUNT','2025-03-19 12:20:00.000');
INSERT INTO accounts VALUES(3,NULL,'My Service 1','A service with 1 service credential','SERVICE_ACCOUNT','2025-03-19 12:21:00.000');
INSERT INTO accounts VALUES(4,NULL,'My Service 2','A service with 2 service credentials','SERVICE_ACCOUNT','2025-03-19 12:22:00.000');
INSERT INTO accounts VALUES(6,'userpassword','User With Password',NULL,'USER_ACCOUNT','2025-03-19 12:23:00.000');
INSERT INTO accounts VALUES(7,'usernopassword','User Without Password',NULL,'USER_ACCOUNT','2025-03-19 12:24:00.000');
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
INSERT INTO service_credentials VALUES(1,3,'Service Credential 1',NULL,X'9c65c105be9cd6cff33aa0ee81c86d86b67283fc13bb4a05c60131e4285b1ca8','2025-03-19 12:21:00.000',NULL);
INSERT INTO service_credentials VALUES(3,4,'Service Credential 2',NULL,X'336041e20250eed54119d84b4da659298036eeb88e9e316caedae227c5e43589','2025-03-19 12:22:00.000','2025-03-19 12:30:00.000');
INSERT INTO service_credentials VALUES(4,4,'Service Credential 1',NULL,X'a303fa3f5c73236056e0238b1e3962cea5a6abaecb053400e66d8ca82578ddf0','2025-03-19 12:23:00.000','2025-03-19 17:00:00.000');
INSERT INTO service_credentials VALUES(5,5,'Service Credential 1',NULL,X'c64f3cbd5fb06eee2b474f96be9c96f2cbcc1d4b362d54aa2616c15743961a6c','2025-03-19 12:24:00.000','2025-03-19 17:00:00.000');
INSERT INTO service_credentials VALUES(6,5,'Service Credential 2',NULL,X'ec05238ebd9e2d27c0e0beecf3367f200061b935de9440ad5f643c87460438ab','2025-03-19 12:25:00.000','2025-03-19 17:01:00.000');
CREATE TABLE password_credentials (
    account_id          INTEGER PRIMARY KEY, -- at most one password credential per account
    password_hash       BLOB NOT NULL,

    FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE
);
INSERT INTO password_credentials VALUES(6,X'2432612431302475356166476c6f7a305646496974363961644e32582e584f644852676e664e70544a652e755232517074504c366835446843503447');
CREATE TABLE roles (
    id           INTEGER PRIMARY KEY,
    display_name TEXT NOT NULL,
    description  TEXT
);
INSERT INTO roles VALUES(1,'My Role',NULL);
CREATE TABLE role_permissions (
    role_id INTEGER NOT NULL,
    permission TEXT NOT NULL,

    FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE
);
INSERT INTO role_permissions VALUES(1,'account:read');
INSERT INTO role_permissions VALUES(1,'account:write');
CREATE TABLE role_assignments (
    id              INTEGER PRIMARY KEY,
    account_id      INTEGER NOT NULL,
    role_id         INTEGER NOT NULL,
    scope_type      TEXT,
    scope_resource  TEXT,

    FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles (id)
);
INSERT INTO role_assignments VALUES(1,6,1,'ZONE','foo');
CREATE UNIQUE INDEX accounts_username ON accounts (username);
CREATE UNIQUE INDEX roles_display_name ON roles (display_name);
CREATE UNIQUE INDEX role_permissions_unique ON role_permissions (role_id, permission);
CREATE UNIQUE INDEX role_assignments_unique ON role_assignments (account_id, role_id);
COMMIT;
