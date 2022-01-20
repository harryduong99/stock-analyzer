package middleware

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func AdminMiddleware() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("x-api-key") == os.Getenv("X-API-KEY") {
				h.ServeHTTP(w, r)
			}
		})
	}
}
