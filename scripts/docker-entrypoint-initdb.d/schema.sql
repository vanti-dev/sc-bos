CREATE TABLE publication
(
    id            TEXT NOT NULL PRIMARY KEY,
    audience_name TEXT
);

CREATE TABLE publication_version
(
    id             UUID        NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    publication_id TEXT        NOT NULL REFERENCES publication (id) ON DELETE CASCADE,
    publish_time   TIMESTAMPTZ NOT NULL,
    body           BYTEA       NOT NULL,
    media_type     TEXT,
    changelog      TEXT
);

CREATE TABLE acknowledgement
(
    id              UUID        NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    -- at the moment it's only possible to acknowledge a version once, hence the uniqueness constraint
    version_id      UUID        NOT NULL UNIQUE REFERENCES publication_version (id) ON DELETE CASCADE,
    accepted        BOOLEAN     NOT NULL,
    rejected_reason TEXT,
    receipt_time    TIMESTAMPTZ NOT NULL
);

CREATE TABLE enrollment
(
    name        TEXT  NOT NULL PRIMARY KEY,
    description TEXT,
    address     TEXT  NOT NULL,
    cert        BYTEA NOT NULL
);

CREATE TABLE tenant
(
    id          UUID        NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    title       TEXT        NOT NULL,
    create_time TIMESTAMPTZ NOT NULL             DEFAULT now()
);

CREATE TABLE tenant_zone
(
    tenant    UUID NOT NULL REFERENCES tenant (id) ON DELETE CASCADE,
    zone_name TEXT NOT NULL,

    PRIMARY KEY (tenant, zone_name),
    CONSTRAINT zone_not_empty CHECK (zone_name <> '')
);

CREATE TABLE tenant_secret
(
    id             UUID        NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    note           TEXT        NOT NULL,
    tenant         UUID        NOT NULL REFERENCES tenant (id) ON DELETE CASCADE,
    secret_hash    BYTEA       NOT NULL,
    expiration     TIMESTAMPTZ NOT NULL,
    first_use_time TIMESTAMPTZ,
    last_use_time  TIMESTAMPTZ
);