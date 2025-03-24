-- +goose up
CREATE TABLE bingo_history (
    id TEXT PRIMARY KEY,
    field TEXT NOT NULL,
    is_set INTEGER,
    session TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
