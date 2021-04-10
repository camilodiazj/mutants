package main

import (
	"github.com/camilodiazj/mutants/infrastructure/configuration"
	"github.com/camilodiazj/mutants/infrastructure/server"
	"log"
	"net/http"
)

const port = ":8000"

func main() {
	s := server.New(configuration.GetInjections())
	log.Println("Init mutants API REST")
	log.Println("Listening on Port: ", port)
	log.Fatal(http.ListenAndServe(port, s.Router()))
}
