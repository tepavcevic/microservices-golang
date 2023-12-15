package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct {
}

const webPort = "8080"

func main() {
	app := Config{}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Println("Starting mail service on port", webPort)

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
