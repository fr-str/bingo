// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package store

import (
	"database/sql"

	"github.com/fr-str/bingo/pkg/db/types"
)

type BingoHistory struct {
	ID        string
	Field     string
	IsSet     sql.NullInt64
	Session   string
	CreatedAt types.RFC3339
	UpdatedAt types.RFC3339
}
