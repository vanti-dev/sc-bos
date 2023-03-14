CREATE TABLE IF NOT EXISTS history
(
    id          BIGSERIAL PRIMARY KEY,
    source      TEXT        NOT NULL,
    create_time TIMESTAMPTZ NOT NULL,
    payload     BYTEA       NULL
);

CREATE INDEX IF NOT EXISTS history_source_idx ON history (source);
CREATE INDEX IF NOT EXISTS history_create_time_idx ON history (create_time);
