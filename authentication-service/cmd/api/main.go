package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/tepavcevic/microservices-golang/authentication/data"
)

const webPort = "8080"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting auth service")
	app := Config{}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
