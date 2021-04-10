package server

import (
	"bytes"
	"encoding/json"
	"github.com/camilodiazj/mutants/application/service"
	"github.com/camilodiazj/mutants/infrastructure/configuration"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

var httpRecorder *httptest.ResponseRecorder
var router *mux.Router
var injections *configuration.Injections
var server Server

func init() {
	httpRecorder = httptest.NewRecorder()
	router = mux.NewRouter()
	injections = &configuration.Injections{
		Processor: &processMock{},
		Router:    router,
	}
	server = New(injections)
}

func TestRouterGETStatsEndpoint(t *testing.T) {
	server.Router()
	router.ServeHTTP(httpRecorder, httptest.NewRequest(http.MethodGet, "/stats", nil))

	if httpRecorder.Code != http.StatusOK {
		t.Error("Did not get expected HTTP status code, got", httpRecorder.Code)
	}
}

func TestRouterPOSTProcessDNA(t *testing.T) {
	server.Router()
	input := []string{"ATGCGA"}
	res, _ := json.Marshal(input)
	router.ServeHTTP(httpRecorder, httptest.NewRequest(http.MethodPost, "/mutant", bytes.NewReader(res)))

	if httpRecorder.Code != http.StatusOK {
		t.Error("Did not get expected HTTP status code, got", httpRecorder.Code)
	}
}

type processMock struct{}

func (*processMock) GetStats() (*service.Stats, error) {
	return &service.Stats{
		CountMutantDna: 10,
		CountHumanDna:  15,
		Ratio:          0,
	}, nil
}

func (*processMock) ProcessDna(dna *service.Dna) (bool, error) {
	return true, nil
}
