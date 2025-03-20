CREATE TABLE user_accounts (
    account_id INTEGER PRIMARY KEY,
    username TEXT NOT NULL,
    password_hash BLOB,

    FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE
);

CREATE TABLE service_accounts (
    account_id INTEGER PRIMARY KEY,
    primary_secret_hash BLOB NOT NULL,
    secondary_secret_hash BLOB,
    secondary_secret_expire_time DATETIME,

    CONSTRAINT secondary_secret_expire_time_format CHECK ( service_accounts.secondary_secret_expire_time IS datetime(secondary_secret_expire_time, 'subsec') ),
    FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE
);

INSERT INTO user_accounts (account_id, username, password_hash)
SELECT id, username, password_hash
FROM accounts
LEFT OUTER JOIN password_credentials ON accounts.id = password_credentials.account_id
WHERE type = 'USER_ACCOUNT';

CREATE UNIQUE INDEX user_accounts_username ON user_accounts (username);

INSERT INTO service_accounts (account_id, primary_secret_hash, secondary_secret_hash, secondary_secret_expire_time)
SELECT
    id,
    -- the old schema required zero-or-more secrets per service account
    -- the new schema requires one-or-two secrets per service account
    -- if the old schema had zero secrets, we'll add a zero hash in place of the primary secret which won't match anything
    coalesce((SELECT secret_hash FROM service_credentials WHERE account_id = accounts.id ORDER BY create_time DESC, id LIMIT 1), zeroblob(32)),
    (SELECT secret_hash FROM service_credentials WHERE account_id = accounts.id ORDER BY create_time DESC, id LIMIT 1 OFFSET 1),
    (SELECT expire_time FROM service_credentials WHERE account_id = accounts.id ORDER BY create_time DESC, id LIMIT 1 OFFSET 1)
FROM accounts
WHERE type = 'SERVICE_ACCOUNT';

DROP TABLE password_credentials;
DROP TABLE service_credentials;
DROP INDEX accounts_username;
ALTER TABLE accounts DROP COLUMN username;

CREATE VIEW account_details AS
SELECT accounts.*, username, password_hash, primary_secret_hash, secondary_secret_hash, secondary_secret_expire_time
FROM accounts
LEFT OUTER JOIN user_accounts ON accounts.id = user_accounts.account_id
LEFT OUTER JOIN service_accounts ON accounts.id = service_accounts.account_id;

