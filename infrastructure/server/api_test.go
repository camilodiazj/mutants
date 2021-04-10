package server

import (
	"bytes"
	"encoding/json"
	"errors"
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

func TestRouterGETStatsEndpointShouldResponseInternalServerError(t *testing.T) {
	newHttpRecorder := httptest.NewRecorder()
	newRouter := mux.NewRouter()
	injections = &configuration.Injections{
		Processor: &processMock{true},
		Router:    newRouter,
	}
	newServer := New(injections)
	newServer.Router()
	newRouter.ServeHTTP(newHttpRecorder, httptest.NewRequest(http.MethodGet, "/stats", nil))

	if newHttpRecorder.Code != http.StatusInternalServerError {
		t.Error("Did not get expected HTTP status code, got", newHttpRecorder.Code)
	}
}

func TestRouterPOSTProcessDNA(t *testing.T) {
	server.Router()
	input := []string{"ATGCGA"}
	res, _ := json.Marshal(input)
	router.ServeHTTP(httpRecorder, httptest.NewRequest(http.MethodPost, "/mutant/", bytes.NewReader(res)))

	if httpRecorder.Code != http.StatusOK {
		t.Error("Did not get expected HTTP status code, got", httpRecorder.Code)
	}
}

func TestRouterPOSTProcessDNAShouldFailDueError(t *testing.T) {
	newHttpRecorder := httptest.NewRecorder()
	newRouter := mux.NewRouter()
	injections = &configuration.Injections{
		Processor: &processMock{true},
		Router:    newRouter,
	}
	newServer := New(injections)
	newServer.Router()
	input := []string{"ATGCGA"}
	res, _ := json.Marshal(input)
	newRouter.ServeHTTP(newHttpRecorder, httptest.NewRequest(http.MethodPost, "/mutant/", bytes.NewReader(res)))

	if newHttpRecorder.Code != http.StatusForbidden {
		t.Error("Did not get expected HTTP status code, got", httpRecorder.Code)
	}
}

type processMock struct {
	shouldFail bool
}

func (p *processMock) GetStats() (*service.Stats, error) {
	if p.shouldFail {
		return &service.Stats{}, errors.New("Fail Get Stats")
	}
	return &service.Stats{
		CountMutantDna: 10,
		CountHumanDna:  15,
		Ratio:          0,
	}, nil
}

func (p *processMock) ProcessDna(dna *service.Dna) (bool, error) {
	if p.shouldFail {
		return false, errors.New("Fail Get Stats")
	}
	return true, nil
}
