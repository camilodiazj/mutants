package main

import (
	"encoding/json"
	"github.com/camilodiazj/mutants/domain/mutant"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var mutantService = mutant.NewMutanService()

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/mutant", processDna).Methods("POST")
	log.Println("Init Server")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func processDna(w http.ResponseWriter, r *http.Request) {
	var dna mutant.Dna
	_ = json.NewDecoder(r.Body).Decode(&dna)
	isMutant := mutantService.IsMutant(dna.Sequence)
	//bytes, _ := json.Marshal(dna.Sequence)
	if isMutant {
		w.WriteHeader(http.StatusOK)
	}
	//infrastructure.PutItem(infrastructure.Dna{
	//	IsMutant: isMutant,
	//	Sequence: string(bytes),
	//})
	//infrastructure.GetItem()
	w.WriteHeader(http.StatusForbidden)
}
