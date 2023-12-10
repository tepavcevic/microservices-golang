package main

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Broker working great!",
	}

	output, _ := json.MarshalIndent(payload, "", "/t")
	w.Header().Set("Content-Type", "appliaction/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(output)
}
