CREATE TABLE IF NOT EXISTS publication
(
    id            TEXT NOT NULL PRIMARY KEY,
    audience_name TEXT
);

CREATE TABLE IF NOT EXISTS publication_version
(
    id             UUID        NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    publication_id TEXT        NOT NULL REFERENCES publication (id) ON DELETE CASCADE,
    publish_time   TIMESTAMPTZ NOT NULL,
    body           BYTEA       NOT NULL,
    media_type     TEXT,
    changelog      TEXT
);

CREATE TABLE IF NOT EXISTS publication_ack
(
    id              UUID        NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    -- at the moment it's only possible to acknowledge a version once, hence the uniqueness constraint
    version_id      UUID        NOT NULL UNIQUE REFERENCES publication_version (id) ON DELETE CASCADE,
    accepted        BOOLEAN     NOT NULL,
    rejected_reason TEXT,
    receipt_time    TIMESTAMPTZ NOT NULL
);
