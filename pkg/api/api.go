package api

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/fr-str/bingo/pkg/api/middleware"
	"github.com/fr-str/bingo/pkg/bingo"
	"github.com/fr-str/bingo/pkg/web"
	"github.com/fr-str/log"
)

type API struct {
	Bingo bingo.Bingo
	mux   *http.ServeMux
}

func (api *API) ListenAndServe(addr string) {
	mux := http.NewServeMux()
	api.mux = mux
	api.RegisterAll()

	handler := middleware.Use(mux, middleware.Panic,
		middleware.Compress,
		middleware.ManageCookie,
		log.HTTPHandler,
	)

	log.Info("ListenAndServe", log.String("addr", addr))
	err := http.ListenAndServe(addr, handler)
	if err != nil {
		panic(err)
	}
}

func (api *API) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request) error) {
	api.mux.HandleFunc(pattern, handlerErrors(handler))
}

func handlerErrors(handler func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			log.InfoCtx(r.Context(), "error handling request", log.String("error", strings.ReplaceAll(err.Error(), "\n", " ")))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (api *API) RegisterAll() {
	api.HandleFunc("GET /", api.index)
	api.HandleFunc("GET /api/square/click", api.handleSquareClick)
	api.HandleFunc("GET /api/stats", api.handleStats)
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

	return web.Index(data).Render(r.Context(), w)
}

func (api *API) handleSquareClick(w http.ResponseWriter, r *http.Request) error {
	field := r.URL.Query().Get("field")
	if field == "" {
		return errors.New("field is required")
	}

	session, ok := r.Context().Value("session").(string)
	if !ok {
		return errors.New("nie wiem jak ale nie ma sesji ¯\\_(ツ)_/¯")
	}
	err := api.Bingo.SaveBingoCell(r.Context(), session, field)
	if err != nil {
		return err
	}

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)

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
