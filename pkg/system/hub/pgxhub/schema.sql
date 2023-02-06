CREATE TABLE IF NOT EXISTS enrollment
(
    name        TEXT  NOT NULL PRIMARY KEY,
    description TEXT,
    address     TEXT  NOT NULL,
    cert        BYTEA NOT NULL
);
