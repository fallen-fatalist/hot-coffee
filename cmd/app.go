package cmd

import (
	"fmt"
	"hot-coffee/internal/002Repositories/jsonrepository"
	"hot-coffee/internal/flag"
	"log"
	"log/slog"
	"net/http"
)

// Main function
func Run() {
	// Initialize stroage
	jsonrepository.Init()

	// Router
	mux := routes()

	slog.Info("Listening on port: 4000")
	err := http.ListenAndServe(fmt.Sprintf(":%d", flag.Port), mux)
	log.Fatal(err)

}
