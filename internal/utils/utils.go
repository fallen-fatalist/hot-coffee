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

func JSONErrorRespond(w http.ResponseWriter, err error, statusCode int) {
	errJSON := errorEnveloper{err.Error()}
	if statusCode == 0 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(statusCode)
	}
	jsonError, err := json.MarshalIndent(errJSON, "", "   ")
	if err != nil {
		slog.Error(err.Error())
	}
	w.Write(jsonError)
}

func In(str string, arr []string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
