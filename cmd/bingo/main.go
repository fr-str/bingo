package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/fr-str/bingo/pkg/api"
	"github.com/fr-str/bingo/pkg/bingo"
	_ "github.com/fr-str/bingo/pkg/config"
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

	api.ListenAndServe(os.Getenv("BINGO_ADDR"))
}
