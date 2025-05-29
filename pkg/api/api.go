package api

import (
	"net/http"
	"strings"

	"github.com/fr-str/bingo/pkg/api/middleware"
	"github.com/fr-str/bingo/pkg/bingo"
	"github.com/fr-str/log"
)

type API struct {
	Bingo bingo.Bingo
	mux   *http.ServeMux
}

func (api *API) ListenAndServe(addr string) {
	mux := http.NewServeMux()
	api.mux = mux
	api.RegisterBingo()
	api.RegisterAllHandsBingo()

	// serve static
	fs := http.FileServer(http.Dir("static"))
	api.mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

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
