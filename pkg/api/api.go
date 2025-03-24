package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fr-str/bingo/pkg/api/middleware"
	"github.com/fr-str/bingo/pkg/bingo"
	"github.com/fr-str/bingo/pkg/store"
	"github.com/fr-str/bingo/pkg/web"
	"github.com/fr-str/log"
)

type API struct {
	DB  *store.Queries
	mux *http.ServeMux
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
	fields := []string{
		"you've forgot charger", "Feature only works on prod", "new jFrog token expired", "challenges and oportunities",
		"bug or a new feature", "unjustified PD call", "workstation WOL crash", "timesheets crash",
		"slack connection problems", "SVPN won't connect", "bugged migration", "Random exception expired",
		"VS code removed", "Feature works everywhere but prod", "work planned without details", "sandwiches gut already gone",
	}

	session, ok := r.Context().Value("session").(string)
	if !ok {
		session = ""
	}
	alreadySet, err := api.DB.GetEntries(r.Context(), session)
	if err != nil {
		return err
	}
	m := make(map[string]bool)
	for _, entry := range alreadySet {
		m[entry.Field] = true
	}

	data := make([]bingo.BingoCell, len(fields))
	for i, cell := range fields {
		data[i] = bingo.BingoCell{
			Field: cell,
			IsSet: m[cell],
		}
	}

	return web.Index(data).Render(r.Context(), w)
}

func (api *API) handleSquareClick(w http.ResponseWriter, r *http.Request) error {
	session := r.Context().Value("session").(string)
	field := r.URL.Query().Get("field")
	if field == "" {
		return errors.New("field is required")
	}

	entry, err := api.DB.GetEntry(r.Context(), bingoEntryID(session, field))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	err = api.DB.SaveBingoEntry(r.Context(), store.SaveBingoEntryParams{
		ID:        bingoEntryID(session, field),
		Field:     field,
		Session:   session,
		IsSet:     sql.NullInt64{Int64: 1, Valid: !entry.IsSet.Valid},
		CreatedAt: entry.CreatedAt,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)

	return nil
}

func bingoEntryID(session string, field string) string {
	// timestamp unix with only year, month and day

	timestamp := time.Now().Unix()
	timestamp = timestamp / 86400
	timestamp = timestamp * 86400
	return fmt.Sprintf("%s/%s/%d", session, field, timestamp)
}
