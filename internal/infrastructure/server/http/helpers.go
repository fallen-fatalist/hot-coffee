package httpserver

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type errorEnveloper struct {
	Err string `json:"error"`
}

// TODO: Move error handling from this helper to handlers
// TODO: Move error logging to another place
func jsonErrorRespond(w http.ResponseWriter, err error, statusCode int) {
	if statusCode == 0 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(statusCode)
	}

	if err != nil {
		slog.Error(err.Error())
	}

	// Hide error if related to server
	if statusCode >= 500 {
		err = fmt.Errorf("internal server error occured")
	}

	json, _ := json.Marshal(errorEnveloper{Err: err.Error()})

	w.Write(json)
}

type messageEnveloper struct {
	Msg string `json:"message"`
}

func jsonMessageRespond(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json, _ := json.Marshal(messageEnveloper{message})
	w.Write(json)

}
