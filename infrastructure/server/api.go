package server

import (
	"encoding/json"
	"github.com/camilodiazj/mutants/application/service"
	"github.com/gorilla/mux"
	"net/http"
)

type api struct {
	router         http.Handler
	mutantVerifier service.Processor
}

type Server interface {
	Router() http.Handler
}

func New() Server {
	a := &api{}
	r := mux.NewRouter()
	r.HandleFunc("/mutant", a.processDna).Methods(http.MethodPost)
	r.HandleFunc("/stats", a.getStats).Methods(http.MethodGet)

	a.router = r
	a.mutantVerifier = service.NewDnaProcessor()
	return a
}

func (a *api) Router() http.Handler {
	return a.router
}

func (a *api) getStats(w http.ResponseWriter, _ *http.Request) {

}

func (a *api) processDna(w http.ResponseWriter, r *http.Request) {
	var dnaSequence service.Dna
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
