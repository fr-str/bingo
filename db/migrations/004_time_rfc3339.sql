-- +goose up
UPDATE bingo_history
SET created_at = REPLACE(SUBSTR(created_at, 1, 25), ' ', 'T') || 'Z';
UPDATE bingo_history
SET updated_at = REPLACE(SUBSTR(updated_at, 1, 25), ' ', 'T') || 'Z';



