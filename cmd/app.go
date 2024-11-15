package cmd

import (
	"fmt"
	"hot-coffee/internal/flag"
	"hot-coffee/internal/repositories/jsonrepository"
	"hot-coffee/internal/services/serviceinstance"
	"log"
	"log/slog"
	"net/http"
	"os"
)

// Main function
func Run() {
	err := flag.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	// Initialize storages
	jsonrepository.Init()
	// Initialize services
	serviceinstance.Init()

	// Router
	mux := routes()

	slog.Info(fmt.Sprintf("Listening on port: %d", flag.Port))
	err = http.ListenAndServe(fmt.Sprintf(":%d", flag.Port), mux)
	log.Fatal(err)

}
