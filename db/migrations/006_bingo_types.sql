-- +goose up
ALTER TABLE bingo_history ADD COLUMN type INTEGER NOT NULL DEFAULT 1;


