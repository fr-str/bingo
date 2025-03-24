package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func ManageCookie(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			// create a new session
			cookie = &http.Cookie{}
			cookie.Name = "session"
			cookie.Value = uuid.New().String()
			http.SetCookie(w, cookie)
		}
		ctx := context.WithValue(r.Context(), "session", cookie.Value)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
