package main

import (
	"context"
	"path/filepath"

	"github.com/fr-str/bingo/pkg/api"
	"github.com/fr-str/bingo/pkg/bingo"
	"github.com/fr-str/bingo/pkg/db"
)

const (
	dbDir = "./data"
)

func main() {
	st, err := db.ConnectStore(context.Background(), filepath.Join(dbDir, "bingo.db"))
	if err != nil {
		panic(err)
	}
	api := api.API{
		Bingo: bingo.Bingo{DB: st},
	}

	api.ListenAndServe(":2137")
}
