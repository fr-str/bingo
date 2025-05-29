package api

import (
	"errors"
	"net/http"

	"github.com/fr-str/bingo/pkg/bingo"
	"github.com/fr-str/bingo/pkg/web"
)

func (api *API) RegisterAllHandsBingo() {
	api.HandleFunc("GET /all-hands", api.handleAllHandsBingo)
	api.HandleFunc("GET /all-hands-bingo-board", api.handleAllHandsBoard)
}

func (api *API) handleAllHandsBingo(w http.ResponseWriter, r *http.Request) error {
	session, ok := r.Context().Value("session").(string)
	if !ok {
		return errors.New("nie wiem jak ale nie ma sesji ¯\\_(ツ)_/¯")
	}
	data, err := api.Bingo.GetAllHandsBingoCells(r.Context(), session)
	if err != nil {
		return err
	}

	return web.AllHandsIndex(bingo.BingoBoard{Cells: data, Type: bingo.Regular}).Render(r.Context(), w)
}

func (api *API) handleAllHandsBoard(w http.ResponseWriter, r *http.Request) error {
	session, ok := r.Context().Value("session").(string)
	if !ok {
		return errors.New("nie wiem jak ale nie ma sesji ¯\\_(ツ)_/¯")
	}
	data, err := api.Bingo.GetAllHandsBingoCells(r.Context(), session)
	if err != nil {
		return err
	}

	return web.BingoBoard(bingo.BingoBoard{Cells: data, Type: bingo.AllHands}).Render(r.Context(), w)
}
