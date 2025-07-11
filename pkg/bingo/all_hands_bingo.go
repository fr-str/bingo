package bingo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/fr-str/bingo/pkg/store"
)

const (
	pepePoland      = `<img alt="peepoPoland, 1x" src="https://cdn.betterttv.net/emote/636794749013520589f5900c/1x.webp">`
	pepeIndia       = `<img src="https://cdn.betterttv.net/emote/6145dc0399488d19dad00332/1x.webp">`
	residentSleeper = `<img alt="residentSleeper, 1x" src="https://cdn.betterttv.net/emote/6097044739b5010444d0e0df/1x.webp">`
)

var allHandsFields = []string{
	"AI is great",
	fmt.Sprintf(`<div>%[1]s%[1]s</div><div>%[2]s%[2]s%[2]s%[2]s%[2]s%[2]s%[2]s%[2]s</div>`, pepePoland, pepeIndia),
	"audio issues",
	"AdManager has \"Clients\"",
	"Hackathon went great",
	"dad jokes",
	"POLAND STRONK",
	"AI = massive boost for productivity",
	"new AI product",
	"We are improving test coverage",
	"diversity and inclusion",
	"the fucked up music drowns out the person talking",
	"shity and confusing charts",
	"INCIDENTS",
	fmt.Sprintf(`<div>%[1]s%[1]s%[1]s%[1]s%[1]s</div>`, residentSleeper),
	"shoutout SRIB",
}

func (b *Bingo) GetAllHandsBingoCells(ctx context.Context, session string) ([]BingoCell, error) {
	alreadySet, err := b.DB.GetEntries(ctx, store.GetEntriesParams{
		Session: session, Day: dayStamp(time.Now()),
		Type: AllHands,
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

	data := make([]BingoCell, len(allHandsFields))
	for i, cell := range allHandsFields {
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
