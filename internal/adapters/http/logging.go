package httpadapter

import (
	"log"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
)

// Logging prints METHOD PATH STATUS DURATION using the provided logger.
func Logging(l *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			dur := time.Since(start)
			l.Printf("%s %s %d %s", r.Method, r.URL.Path, ww.Status(), dur)
		})
	}
}
