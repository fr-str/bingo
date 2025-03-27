package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/fr-str/bingo/pkg/api"
	"github.com/fr-str/bingo/pkg/bingo"
	_ "github.com/fr-str/bingo/pkg/config"
	"github.com/fr-str/bingo/pkg/db"
	"github.com/fr-str/env"
)

func main() {
	dbDir := env.Get("BINGO_DB_DIR", "./data")
	var dbPath string
	if dbDir != ":memory:" {
		dbPath = filepath.Join(dbDir, "bingo.db")
	}
	st, err := db.ConnectStore(context.Background(), dbPath)
	if err != nil {
		panic(err)
	}
	api := api.API{
		Bingo: bingo.Bingo{DB: st},
	}

	api.ListenAndServe(os.Getenv("BINGO_ADDR"))
}
