CREATE TABLE IF NOT EXISTS enrollment
(
    address     TEXT  NOT NULL PRIMARY KEY,
    name        TEXT,
    description TEXT,
    cert        BYTEA NOT NULL
);
