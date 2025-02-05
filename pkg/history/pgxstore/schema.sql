CREATE TABLE IF NOT EXISTS history
(
    id          BIGSERIAL PRIMARY KEY,
    source      TEXT        NOT NULL,
    create_time TIMESTAMPTZ NOT NULL,
    payload     BYTEA       NULL
);

DROP INDEX IF EXISTS history_create_time_idx; -- replaced by history_id_source_idx;
DROP INDEX IF EXISTS history_source_idx; -- replaced by history_source_create_time_idx
CREATE INDEX IF NOT EXISTS history_id_source_idx ON history (id, source);
CREATE INDEX IF NOT EXISTS history_source_create_time_idx ON history (source, create_time);
