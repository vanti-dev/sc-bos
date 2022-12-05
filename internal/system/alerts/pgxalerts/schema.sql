CREATE TABLE IF NOT EXISTS alerts
(
    id          UUID        NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    description TEXT        NULL,
    severity    int         NOT NULL             DEFAULT 13, /* WARN */
    create_time TIMESTAMPTZ NOT NULL             DEFAULT now(),
    ack_time    TIMESTAMPTZ NULL,
    floor       TEXT        NULL,
    zone        TEXT        NULL,
    source      TEXT        NULL
);
