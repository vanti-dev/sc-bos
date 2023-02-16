CREATE TABLE IF NOT EXISTS history
(
    id          BIGSERIAL PRIMARY KEY,
    source      TEXT        NOT NULL,
    create_time TIMESTAMPTZ NOT NULL,
    payload     BYTEA       NULL
);

CREATE INDEX ON history (source);
CREATE INDEX ON history (create_time);
