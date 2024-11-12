package cmd

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
)

// Main function
func Run() {
	// Router
	mux := routes()

	slog.Info("Listening on port: 4000")
	err := http.ListenAndServe(fmt.Sprintf(":%d", Port), mux)
	log.Fatal(err)

}
