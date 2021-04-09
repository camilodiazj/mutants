package main

import (
	"github.com/camilodiazj/mutants/infrastructure/server"
	"log"
	"net/http"
)


func main() {
	s := server.New()
	log.Fatal(http.ListenAndServe(":8000", s.Router()))
}





