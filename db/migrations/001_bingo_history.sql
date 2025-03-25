-- +goose up
CREATE TABLE IF NOT EXISTS bingo_history (
    id TEXT PRIMARY KEY,
    field TEXT NOT NULL,
    is_set INTEGER,
    session TEXT NOT NULL,
    created_at RFC3339 NOT NULL,
    updated_at RFC3339 NOT NULL
);
