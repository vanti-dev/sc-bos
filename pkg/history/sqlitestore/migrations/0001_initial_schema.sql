CREATE TABLE history_meta (
    key             TEXT PRIMARY KEY,
    value           ANY
);

INSERT INTO history_meta (key, value) VALUES
    ('epoch', datetime('subsec')); -- don't modify this once data is inserted, as it changes the meaning of the data

CREATE TABLE history_sources (
    id              INTEGER PRIMARY KEY AUTOINCREMENT, -- use AUTOINCREMENT to ensure IDs are never reused
    source          TEXT NOT NULL UNIQUE
);

CREATE TABLE history (
    id              INTEGER PRIMARY KEY,
    source_id       INTEGER NOT NULL,
    epoch_offset_ms INTEGER NOT NULL,
    payload         BLOB,

    FOREIGN KEY (source_id) REFERENCES history_sources(id)
);

CREATE INDEX history_source_create_time_idx ON history (source_id, epoch_offset_ms);
