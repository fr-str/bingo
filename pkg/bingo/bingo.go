package bingo

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/fr-str/bingo/pkg/db/types"
	"github.com/fr-str/bingo/pkg/store"
)

type BingoCell struct {
	Field string
	Count int
	IsSet bool
}

type BingoBoard struct {
	Cells []BingoCell
	Type  int64
}

func dayStamp(t time.Time) int64 {
	timestamp := t.Unix()
	timestamp = timestamp / 86400
	return timestamp * 86400
}

type Bingo struct {
	DB *store.Queries
}

func (b Bingo) SaveBingoCell(ctx context.Context, session, field string, bType int64) error {
	day := dayStamp(time.Now())
	entry, err := b.DB.GetEntry(ctx, store.GetEntryParams{
		Field: field, Session: session, Day: day, Type: bType,
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
		Type:      bType,
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
		Type: Regular,
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

	// shuffle the data to make it look random
	// using session + day as a seed
	// this way seed will be different for each user
	// but stay the same during single day
	seed := session + strconv.FormatInt(dayStamp(time.Now()), 10)
	ShuffleSliceWithSeed(data, seed)

	return data, nil
}

func stringToSeed(s string) int64 {
	hasher := sha256.New()
	hasher.Write([]byte(s))
	hashBytes := hasher.Sum(nil)

	// Take the first 8 bytes of the hash for the seed.
	// SHA256 produces a 32-byte hash, so we have plenty.
	if len(hashBytes) < 8 {
		// This should not happen with SHA256, but as a fallback
		var seed int64
		for i, b := range hashBytes {
			seed += int64(b) << (i * 8) // Simple combination if too short
		}
		return seed
	}
	return int64(binary.BigEndian.Uint64(hashBytes[:8]))
}

func ShuffleSliceWithSeed(slice []BingoCell, seedString string) {
	seed := stringToSeed(seedString)
	// Create a new random source and a new rand.Rand instance
	source := rand.NewSource(seed)
	r := rand.New(source)

	r.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}
