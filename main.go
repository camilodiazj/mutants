package main

import (
	"github.com/camilodiazj/mutants/infrastructure/configuration"
	"github.com/camilodiazj/mutants/infrastructure/server"
	"log"
	"net/http"
)

func main() {
	s := server.New(configuration.GetInjections())
	log.Fatal(http.ListenAndServe(":8000", s.Router()))
}
