package main

import (
	"encoding/json"
	"github.com/camilodiazj/mutants/domain/mutant"
	"github.com/camilodiazj/mutants/infrastructure"
	"github.com/gorilla/mux"
	"log"
	"math"
	"net/http"
	"sync"
)

var mutantService = mutant.NewMutanService()
var wg sync.WaitGroup

type Stats struct {
	CountMutantDna uint64  `json:"count_mutant_dna"`
	CountHumanDna  uint64  `json:"count_human_dna"`
	Ratio          float64 `json:"ratio"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/mutant", processDna).Methods("POST")
	router.HandleFunc("/stats", getStats).Methods("GET")
	log.Println("Init Server")
	wg.Wait()
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getStats(w http.ResponseWriter, _ *http.Request) {
	result, err := infrastructure.GetCountOf(true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		ratioRounded := 0.0
		if result.ItemsCount != 0 {
			ratio := float32(result.Counter) / float32(result.ItemsCount)
			ratioRounded = math.Round(float64(ratio*100)) / 100
		}

		stats := &Stats{
			CountHumanDna:  result.ItemsCount - result.Counter,
			CountMutantDna: result.Counter,
			Ratio:          ratioRounded,
		}
		res, _ := json.Marshal(stats)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(res)
	}
}

func processDna(w http.ResponseWriter, r *http.Request) {
	var dna mutant.Dna
	_ = json.NewDecoder(r.Body).Decode(&dna)

	isMutant := validateDna(dna.Sequence)
	log.Println("Dna processed, is mutant? ", isMutant)
	if isMutant {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
	wg.Add(1)
	go persistDnaResult(dna, isMutant)
}

func validateDna(dna []string) bool {
	return mutantService.IsMutant(dna)
}

func persistDnaResult(dna mutant.Dna, isMutant bool) {
	bytes, err := json.Marshal(dna.Sequence)
	infrastructure.PutItem(&infrastructure.Dna{
		IsMutant: isMutant,
		Sequence: string(bytes),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Dna persisted")
	wg.Done()
}
