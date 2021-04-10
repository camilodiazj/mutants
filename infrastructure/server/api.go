package server

import (
	"encoding/json"
	"github.com/camilodiazj/mutants/application/service"
	"github.com/camilodiazj/mutants/infrastructure/configuration"
	"net/http"
)

type api struct {
	router    http.Handler
	processor service.Processor
}

type Server interface {
	Router() http.Handler
}

func New(injections *configuration.Injections) Server {
	a := &api{}
	r := injections.Router
	a.processor = injections.Processor
	r.HandleFunc("/mutant", a.processDna).Methods(http.MethodPost)
	r.HandleFunc("/stats", a.getStats).Methods(http.MethodGet)
	a.router = r
	return a
}

func (a *api) Router() http.Handler {
	return a.router
}

func (a *api) getStats(w http.ResponseWriter, _ *http.Request) {
	stats, err := a.processor.GetStats()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		res, _ := json.Marshal(stats)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(res)
	}
}

func (a *api) processDna(w http.ResponseWriter, r *http.Request) {
	var dnaSequence service.Dna
	_ = json.NewDecoder(r.Body).Decode(&dnaSequence)

	isMutant, err := a.processor.ProcessDna(&dnaSequence)

	if err != nil || !isMutant {
		w.WriteHeader(http.StatusForbidden)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
