package server

import (
	"encoding/json"
	"github.com/camilodiazj/mutants/application/dna"
	"github.com/gorilla/mux"
	"net/http"
)

type api struct {
	router         http.Handler
	mutantVerifier dna.Processor
}

type Server interface {
	Router() http.Handler
}

func (a *api) Router() http.Handler {
	return a.router
}

func New() Server {
	a := &api{}
	r := mux.NewRouter()
	r.HandleFunc("/mutant", a.processDna).Methods(http.MethodPost)
	r.HandleFunc("/stats", a.getStats).Methods(http.MethodGet)

	a.router = r
	a.mutantVerifier = dna.NewDnaProcessor()
	return &api{}
}

func (a *api) getStats(w http.ResponseWriter, _ *http.Request) {

}

func (a *api) processDna(w http.ResponseWriter, r *http.Request) {
	var dnaSequence dna.Dna
	_ = json.NewDecoder(r.Body).Decode(&dnaSequence)

	isMutant, err := a.mutantVerifier.ProcessDna(&dnaSequence)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
	}

	if isMutant {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}
