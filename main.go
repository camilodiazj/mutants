package main

import (
	"github.com/camilodiazj/mutants/infrastructure/server"
	"log"
	"net/http"
)

//var wg sync.WaitGroup

type Stats struct {
	CountMutantDna uint64  `json:"count_mutant_dna"`
	CountHumanDna  uint64  `json:"count_human_dna"`
	Ratio          float64 `json:"ratio"`
}

func main() {
	s := server.New()
	//wg.Wait()
	log.Fatal(http.ListenAndServe(":8000", s.Router()))
}





