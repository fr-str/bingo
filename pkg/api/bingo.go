package api

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/fr-str/bingo/pkg/bingo"
	"github.com/fr-str/bingo/pkg/web"
)

func (api *API) RegisterBingo() {
	api.HandleFunc("GET /", api.index)
	api.HandleFunc("GET /api/square/click", api.handleSquareClick)
	api.HandleFunc("GET /api/stats", api.handleStats)
	api.HandleFunc("GET /bingo-board", api.handleBoard)
}

func (api *API) index(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return nil
	}
	session, ok := r.Context().Value("session").(string)
	if !ok {
		return errors.New("nie wiem jak ale nie ma sesji ¯\\_(ツ)_/¯")
	}
	data, err := api.Bingo.GetBingoCells(r.Context(), session)
	if err != nil {
		return err
	}

	return web.Index(bingo.BingoBoard{Cells: data, Type: bingo.Regular}).Render(r.Context(), w)
}

func (api *API) handleBoard(w http.ResponseWriter, r *http.Request) error {
	session, ok := r.Context().Value("session").(string)
	if !ok {
		return errors.New("nie wiem jak ale nie ma sesji ¯\\_(ツ)_/¯")
	}

	data, err := api.Bingo.GetBingoCells(r.Context(), session)
	if err != nil {
		return err
	}

	return web.BingoBoard(bingo.BingoBoard{Cells: data, Type: bingo.Regular}).Render(r.Context(), w)
}

func (api *API) handleSquareClick(w http.ResponseWriter, r *http.Request) error {
	session, ok := r.Context().Value("session").(string)
	if !ok {
		return errors.New("nie wiem jak ale nie ma sesji ¯\\_(ツ)_/¯")
	}

	field := r.URL.Query().Get("field")
	if field == "" {
		return errors.New("field is required")
	}
	bType := r.URL.Query().Get("type")
	if bType == "" {
		return errors.New("type is required")
	}

	if len(bType) > 1 {
		return errors.New("type can only be 1 character")
	}

	typeValue := int64(bType[0]) - '0'
	if !bingo.TypeExists(typeValue) {
		return errors.New("invalid type")
	}

	err := api.Bingo.SaveBingoCell(r.Context(), session, field, typeValue)
	if err != nil {
		return err
	}

	w.Header().Set("HX-Trigger", "force-load")
	w.WriteHeader(http.StatusOK)
	// switch typeValue {
	// case bingo.Regular:
	// 	w.Header().Set("Location", "/")
	// case bingo.AllHands:
	// 	w.Header().Set("Location", "/all-hands")
	// }
	// w.WriteHeader(http.StatusFound)

	return nil
}

func (api *API) handleStats(w http.ResponseWriter, r *http.Request) error {
	data, err := api.Bingo.DB.BingoStats(r.Context())
	if err != nil {
		return err
	}

	// return in format requested by client
	accept := r.Header.Get("Accept")
	switch {
	case strings.Contains(accept, "text/csv"):
		sdata := [][]string{{"field", "count", "date"}}
		for _, d := range data {
			sdata = append(sdata, []string{d.Field, strconv.FormatInt(d.Count, 10), d.Date.(string)})
		}
		w.Header().Add("content-type", "text/csv")
		return csv.NewWriter(w).WriteAll(sdata)

	case strings.Contains(accept, "application/json"):
		fallthrough
	default:
		w.Header().Add("content-type", "application/json")
		return json.NewEncoder(w).Encode(data)
	}
}
