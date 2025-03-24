package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type customWriter struct {
	http.ResponseWriter
	io.Writer
	statusCode string
}

func (c *customWriter) WriteHeader(status int) {
	c.statusCode = strconv.Itoa(status)
	c.ResponseWriter.WriteHeader(status)
}

func (c *customWriter) Write(b []byte) (int, error) {
	if c.statusCode == "" {
		c.statusCode = "200"
	}
	return c.Writer.Write(b)
}

func Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enc := "none"
		cw := &customWriter{ResponseWriter: w, Writer: w}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			enc = "gzip"

			if r.RequestURI != "/metrics" {
				w.Header().Set("Content-Encoding", enc)
				gz := gzip.NewWriter(w)
				defer gz.Close()
				cw.Writer = gz
			}
		}
		next.ServeHTTP(cw, r)
	})
}
