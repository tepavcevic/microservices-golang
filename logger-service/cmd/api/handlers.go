package main

import (
	"net/http"

	"github.com/tepavcevic/microservices-golang/logger/data"
)

type jsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload jsonPayload

	_ = app.readJSON(w, r, &requestPayload)

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	if err := app.Models.LogEntry.Insert(event); err != nil {
		app.errorJSON(w, r, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "log created",
	}

	app.writeJSON(w, r, http.StatusAccepted, resp)
}
