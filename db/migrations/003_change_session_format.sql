-- +goose up
-- update sessions to use new format, <session>/<timestamp>
UPDATE bingo_history
SET session = session || '/' || strftime('%s', date(substr(updated_at, 1, 10)))
WHERE session NOT LIKE '%/%'; -- Avoid double appending
