package bingo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/fr-str/bingo/pkg/db/types"
	"github.com/fr-str/bingo/pkg/store"
)

type BingoCell struct {
	Field string
	IsSet bool
}

func bingoSession(session string) string {
	return fmt.Sprintf("%s/%d", session, dayStamp(time.Now()))
}

func dayStamp(t time.Time) int64 {
	timestamp := t.Unix()
	timestamp = timestamp / 86400
	return timestamp * 86400
}

type Bingo struct {
	DB *store.Queries
}

func bingoIDFormat(field, session string) string {
	return fmt.Sprintf("%s/%s/%d", session, field, dayStamp(time.Now()))
}

func (b Bingo) SaveBingoCell(ctx context.Context, session, field string) error {
	session = bingoSession(session)

	entry, err := b.DB.GetEntry(ctx, store.GetEntryParams{ID: field, Session: session})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	err = b.DB.SaveBingoEntry(ctx, store.SaveBingoEntryParams{
		ID:        bingoIDFormat(field, session),
		Field:     field,
		Session:   session,
		IsSet:     sql.NullInt64{Int64: 1, Valid: !entry.IsSet.Valid},
		CreatedAt: types.RFC3339{Time: time.Now()},
		UpdatedAt: types.RFC3339{Time: time.Now()},
	})
	if err != nil {
		return err
	}

	return nil
}

var fields = []string{
	"you've forgot charger", "Feature only works on prod", "new jFrog token expired", "challenges and oportunities",
	"bug or a new feature", "unjustified PD call", "workstation WOL crash", "timesheets crash",
	"slack connection problems", "SVPN won't connect", "bugged migration", "Random exception expired",
	"VS code removed", "Feature works everywhere but prod", "work planned without details", "sandwiches guy already gone",
}

func (b Bingo) GetBingoCells(ctx context.Context, session string) ([]BingoCell, error) {
	session = bingoSession(session)

	alreadySet, err := b.DB.GetEntries(ctx, session)
	if err != nil {
		return nil, err
	}

	m := make(map[string]bool)
	for _, entry := range alreadySet {
		m[entry.Field] = true
	}

	data := make([]BingoCell, len(fields))
	for i, cell := range fields {
		data[i] = BingoCell{
			Field: cell,
			IsSet: m[cell],
		}
	}
	return data, nil
}
