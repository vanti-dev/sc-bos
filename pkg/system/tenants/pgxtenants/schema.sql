CREATE TABLE IF NOT EXISTS tenant
(
    id          UUID        NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    title       TEXT        NOT NULL,
    create_time TIMESTAMPTZ NOT NULL             DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tenant_zone
(
    tenant    UUID NOT NULL REFERENCES tenant (id) ON DELETE CASCADE,
    zone_name TEXT NOT NULL,

    PRIMARY KEY (tenant, zone_name),
    CONSTRAINT zone_not_empty CHECK (zone_name <> '')
);

CREATE TABLE IF NOT EXISTS tenant_secret
(
    id             UUID        NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    note           TEXT        NOT NULL,
    tenant         UUID        NOT NULL REFERENCES tenant (id) ON DELETE CASCADE,
    secret_hash    BYTEA       NOT NULL,
    create_time    TIMESTAMPTZ NOT NULL             DEFAULT now(),
    expire_time    TIMESTAMPTZ,
    first_use_time TIMESTAMPTZ,
    last_use_time  TIMESTAMPTZ
);
