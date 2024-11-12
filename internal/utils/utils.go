package utils

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func FatalError(context string, err error) {
	if err != nil {
		log.Fatalf(context+": %s\n", err)
		os.Exit(1)
	}
}

type errorEnveloper struct {
	Err string `json:"error"`
}

func JSONErrorRespond(w http.ResponseWriter, err error) {
	errJSON := errorEnveloper{err.Error()}
	w.WriteHeader(http.StatusBadRequest)
	jsonError, err := json.MarshalIndent(errJSON, "", "   ")
	if err != nil {
		slog.Error(err.Error())
	}
	w.Write(jsonError)
}
