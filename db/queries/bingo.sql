-- name: SaveBingoEntry :exec
INSERT INTO bingo_history (id, field,session, is_set, created_at, updated_at)
VALUES (:id, :field,:session, :is_set, :created_at, :updated_at)
ON CONFLICT (id) DO UPDATE SET
    is_set = excluded.is_set,
    updated_at = excluded.updated_at;

-- name: GetEntry :one
SELECT * FROM bingo_history WHERE id = :id;

-- name: GetEntries :many
SELECT * FROM bingo_history WHERE session = :session AND is_set IS NOT NULL;
