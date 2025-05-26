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
	Count int
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

func (b Bingo) SaveBingoCell(ctx context.Context, session, field string) error {
	day := dayStamp(time.Now())
	entry, err := b.DB.GetEntry(ctx, store.GetEntryParams{
		Field: field, Session: session, Day: day,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	err = b.DB.SaveBingoEntry(ctx, store.SaveBingoEntryParams{
		Field:     field,
		Session:   session,
		Day:       day,
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
	alreadySet, err := b.DB.GetEntries(ctx, store.GetEntriesParams{
		Session: session, Day: dayStamp(time.Now()),
	})
	if err != nil {
		return nil, err
	}

	m := make(map[string]struct {
		IsSet bool
		Count int
	})
	for _, entry := range alreadySet {
		m[entry.Field] = struct {
			IsSet bool
			Count int
		}{
			IsSet: true,
			// subtract 1 because we are counted too
			// and count is supposed to be noumber of other people who marked this square
			Count: int(entry.DailyFieldCount) - 1,
		}
	}

	data := make([]BingoCell, len(fields))
	for i, cell := range fields {
		data[i] = BingoCell{
			Field: cell,
			IsSet: m[cell].IsSet,
			Count: m[cell].Count,
		}
	}
	return data, nil
}
