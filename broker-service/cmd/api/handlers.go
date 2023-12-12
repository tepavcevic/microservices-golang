package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	if err := app.readJSON(w, r, &requestPayload); err != nil {
		app.errorJSON(w, r, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, r, requestPayload.Auth)
	default:
		app.errorJSON(w, r, errors.New("unknown actiont"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, r *http.Request, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, r, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, r, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, r, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}
	if response.StatusCode != http.StatusOK {
		app.errorJSON(w, r, errors.New("error calling auth service"), http.StatusInternalServerError)
		return
	}

	var jsonFromService jsonResponse
	if err := json.NewDecoder(response.Body).Decode(&jsonFromService); err != nil {
		app.errorJSON(w, r, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, r, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, r, http.StatusAccepted, payload)
}
