package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/fr-str/log"
)

type ContextKey string

const (
	ClientKey ContextKey = "client"
)

type Middleware func(http.Handler) http.Handler

func Use(mux http.Handler, midlewares ...Middleware) http.Handler {
	for i := len(midlewares) - 1; i >= 0; i-- {
		mux = midlewares[i](mux)
	}
	return mux
}

func Panic(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				log.ErrorCtx(r.Context(), "HTTP Panic", log.Err(err.(error)))
				fmt.Println(string(debug.Stack()))
			}
		}()
		h.ServeHTTP(w, r)
	})
}
