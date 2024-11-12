package cmd

import (
	"fmt"
	"hot-coffee/internal/flag"
	"hot-coffee/internal/repositories/jsonrepository"
	"hot-coffee/internal/services/serviceinstance"
	"log"
	"log/slog"
	"net/http"
)

// Main function
func Run() {
	// Initialize storages
	jsonrepository.Init()
	// Initialize services
	serviceinstance.Init()

	// Router
	mux := routes()

	slog.Info("Listening on port: 4000")
	err := http.ListenAndServe(fmt.Sprintf(":%d", flag.Port), mux)
	log.Fatal(err)

}
