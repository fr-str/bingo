package api

import (
	"errors"
	"net/http"
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

	session := r.Context().Value("session").(string)
	err := api.Bingo.SaveBingoCell(r.Context(), session, field)
	if err != nil {
		return err
	}

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)

	return nil
}
