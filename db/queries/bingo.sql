-- name: SaveBingoEntry :exec
INSERT INTO bingo_history (field,session,day, is_set, created_at, updated_at,type)
VALUES (:field,:session,:day, :is_set, :created_at, :updated_at,:type)
ON CONFLICT (session,day,field) DO UPDATE SET
    is_set = excluded.is_set,
    updated_at = excluded.updated_at;

-- name: GetEntry :one
SELECT * FROM bingo_history WHERE field = :field AND session = :session AND day = :day AND type = :type LIMIT 1;

-- name: GetEntries :many
WITH DailyFieldCounts AS (
    SELECT 
        field,
        COUNT(*) as daily_field_count
    FROM bingo_history WHERE bingo_history.day = :day AND bingo_history.type = :type AND is_set IS NOT NULL
    GROUP BY field
) SELECT 
    bh_session.*,
    dfc.daily_field_count 
FROM bingo_history bh_session
JOIN DailyFieldCounts dfc ON bh_session.field = dfc.field
WHERE 
    bh_session.session = :session
AND 
    bh_session.type = :type
AND 
    bh_session.day = :day
AND
    bh_session.is_set IS NOT NULL;

-- aggregates by field and date
-- name: BingoStats :many
SELECT
    field,
    count(*) as count,
    date(created_at) as date
FROM bingo_history
WHERE is_set IS NOT NULL
GROUP BY field, date(created_at)
ORDER BY date(created_at) DESC;

