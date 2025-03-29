package httpserver

import (
	"fmt"
	"log/slog"
	"net/http"
)

func RequestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info(fmt.Sprintf("%s request with URL: %s", r.Method, r.URL.String()))
		next.ServeHTTP(w, r)
	})
}

func HeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
