package cmd

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"hot-coffee/internal/flag"
	"hot-coffee/internal/service/serviceinstance"
)

// Main function
func Run() {
	err := flag.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	// Initialize services
	serviceinstance.Init()

	// Router
	mux := routes()

	slog.Info(fmt.Sprintf("Listening on port: %d", flag.Port))
	err = http.ListenAndServe(fmt.Sprintf(":%d", flag.Port), mux)
	log.Fatal(err)
}
