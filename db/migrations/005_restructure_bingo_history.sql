-- +goose Up
ALTER TABLE bingo_history RENAME TO bingo_history_old;

CREATE TABLE bingo_history (
    field TEXT NOT NULL,
    session TEXT NOT NULL,
    day INTEGER NOT NULL,
    is_set INTEGER,
    created_at RFC3339 NOT NULL,
    updated_at RFC3339 NOT NULL,
    PRIMARY KEY (session,day,field)
);

INSERT INTO bingo_history (field, session,day, is_set, created_at, updated_at)
SELECT 
    T.field,
    SUBSTR(T.session,1,  INSTR(T.session, '/')-1) as session, 
    SUBSTR(T.session,  INSTR(T.session, '/')+1) AS day,
    T.is_set,
    T.created_at,
    T.updated_at
FROM
    bingo_history_old AS T;

DROP TABLE bingo_history_old;
