-- +goose up
-- bacfill created_at using updated_at
UPDATE bingo_history SET created_at = updated_at;
