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
