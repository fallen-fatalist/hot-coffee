package cmd

import (
	"log"
	"net/http"
)

// Main function
func Run() {
	// Routes
	mux := routes()

	log.Print("Listening on port: 4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}
