CREATE TABLE IF NOT EXISTS alerts
(
    id               UUID        NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    description      TEXT        NULL,
    severity         int         NOT NULL             DEFAULT 13, /* WARN */
    create_time      TIMESTAMPTZ NOT NULL             DEFAULT now(),
    resolve_time     TIMESTAMPTZ NULL,
    floor            TEXT        NULL,
    zone             TEXT        NULL,
    source           TEXT        NULL,
    subsystem        TEXT        NULL,
    federation       TEXT        NULL,

    ack_time         TIMESTAMPTZ NULL,
    ack_author_name  TEXT        NULL,
    ack_author_email TEXT        NULL,
    ack_author_id    TEXT        NULL
);

ALTER TABLE alerts
    ADD COLUMN IF NOT EXISTS federation TEXT NULL;
ALTER TABLE alerts
    ADD COLUMN IF NOT EXISTS resolve_time TIMESTAMPTZ NULL;
ALTER TABLE alerts
    ADD COLUMN IF NOT EXISTS subsystem TEXT NULL;
