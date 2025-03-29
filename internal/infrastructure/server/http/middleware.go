package cmd

import (
	"fmt"
	"log/slog"
	"net/http"
)

func requestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info(fmt.Sprintf("%s request with URL: %s", r.Method, r.URL.String()))
		next.ServeHTTP(w, r)
	})
}
