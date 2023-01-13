CREATE TABLE IF NOT EXISTS alerts
(
    id               UUID        NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    description      TEXT        NULL,
    severity         int         NOT NULL             DEFAULT 13, /* WARN */
    create_time      TIMESTAMPTZ NOT NULL             DEFAULT now(),
    floor            TEXT        NULL,
    zone             TEXT        NULL,
    source           TEXT        NULL,

    ack_time         TIMESTAMPTZ NULL,
    ack_author_name  TEXT        NULL,
    ack_author_email TEXT        NULL,
    ack_author_id    TEXT        NULL
);
