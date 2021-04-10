package server

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/camilodiazj/mutants/application/service"
	"github.com/camilodiazj/mutants/domain/mutant"
	"github.com/camilodiazj/mutants/infrastructure/repository"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
)

type api struct {
	router         http.Handler
	mutantVerifier service.Processor
}

var wg sync.WaitGroup

type Server interface {
	Router() http.Handler
}

func New() Server {
	a := &api{}
	r := mux.NewRouter()
	r.HandleFunc("/mutant", a.processDna).Methods(http.MethodPost)
	r.HandleFunc("/stats", a.getStats).Methods(http.MethodGet)

	a.router = r
	a.mutantVerifier = service.NewDnaProcessor(&wg, mutant.NewMutanVerifier(), repository.NewDynamoRepository("DNA", ConfigureDynamoDB()))
	wg.Wait()
	return a
}

func (a *api) Router() http.Handler {
	return a.router
}

func (a *api) getStats(w http.ResponseWriter, _ *http.Request) {
	stats, err := a.mutantVerifier.GetStats()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	res, _ := json.Marshal(stats)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
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

//TODO: Move to configurations file or something like that
func ConfigureDynamoDB() dynamodbiface.DynamoDBAPI {
	awsSession, _ := session.NewSession(&aws.Config{Region: aws.String("us-east-2")})
	svc := dynamodb.New(awsSession)
	return dynamodbiface.DynamoDBAPI(svc)
}
