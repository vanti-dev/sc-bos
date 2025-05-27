CREATE TABLE history_sources (
    id              INTEGER PRIMARY KEY AUTOINCREMENT, -- use AUTOINCREMENT to ensure IDs are never reused
    source          TEXT NOT NULL UNIQUE
);

CREATE TABLE history (
    id              INTEGER PRIMARY KEY, -- encodes both the timestamp and a sequence number
    source_id       INTEGER NOT NULL,
    payload         BLOB,

    FOREIGN KEY (source_id) REFERENCES history_sources(id)
);

CREATE INDEX history_source_idx ON history (source_id, id);

-- expanded view of history data, useful when manually exploring data
CREATE VIEW history_denorm (id, timestamp, source, payload) AS
    SELECT h.id,
           datetime(CAST((h.id / 1000000) AS REAL) / 1000.0, 'unixepoch', 'subsec'),
           hs.source,
           h.payload
    FROM history h
    LEFT OUTER JOIN history_sources hs ON h.source_id = hs.id;
