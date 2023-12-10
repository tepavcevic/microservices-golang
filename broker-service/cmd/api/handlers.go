package main

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Broker working great!",
	}

	err := app.writeJSON(w, r, http.StatusOK, payload)

	if err != nil {
		_ = app.errorJSON(w, r, err)
	}
}
