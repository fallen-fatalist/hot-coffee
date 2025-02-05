package httpserver

import (
	"log/slog"
	"net/http"
)

type errorEnveloper struct {
	Err string `json:"error"`
}

func jsonErrorRespond(w http.ResponseWriter, err error, statusCode int) {
	slog.Error(err.Error())
	if statusCode == 0 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(statusCode)
	}

	if err != nil {
		slog.Error(err.Error())
	}
	message := "Undefined server error"

	w.Write([]byte(message))
}
