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
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
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
	case "log":
		app.logItem(w, r, requestPayload.Log)
	case "mail":
		app.mail(w, r, requestPayload.Mail)
	default:
		app.errorJSON(w, r, errors.New("unknown action"))
	}
}

func (app *Config) logItem(w http.ResponseWriter, r *http.Request, l LogPayload) {
	jsonData, _ := json.MarshalIndent(l, "", "\t")

	logServiceURL := "http://logger-service:8080/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, r, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, r, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, r, errors.New("error calling logger service"), http.StatusInternalServerError)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "log created"

	app.writeJSON(w, r, http.StatusAccepted, payload)
}

func (app *Config) authenticate(w http.ResponseWriter, r *http.Request, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", "http://authentication-service:8080/authenticate", bytes.NewBuffer(jsonData))
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
	if response.StatusCode != http.StatusAccepted {
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

func (app *Config) mail(w http.ResponseWriter, r *http.Request, m MailPayload) {
	jsonData, _ := json.MarshalIndent(m, "", "\t")

	mailServiceURL := "http://mail-service:8080/send"

	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, r, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, r, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, r, errors.New("error calling mail service"), http.StatusInternalServerError)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "mail sent to " + m.To

	app.writeJSON(w, r, http.StatusAccepted, payload)
}
