package service

import (
	"encoding/json"
	repository3 "github.com/camilodiazj/mutants/application/repository"
	"github.com/camilodiazj/mutants/domain/mutant"
	"github.com/camilodiazj/mutants/infrastructure/repository"
	"github.com/google/uuid"
	"log"
	"math"
	"sync"
)

type Stats struct {
	CountMutantDna uint64  `json:"count_mutant_dna"`
	CountHumanDna  uint64  `json:"count_human_dna"`
	Ratio          float64 `json:"ratio"`
}

type Dna struct {
	Sequence []string `json:"dna"`
}

type MutantProcessor struct {
	wg         *sync.WaitGroup
	verifier   mutant.MutanVerifier
	repository repository3.DnaRepository
}

type Processor interface {
	GetStats() (*Stats, error)
	ProcessDna(dna *Dna) (bool, error)
}

func NewDnaProcessor(wg *sync.WaitGroup) Processor {
	return &MutantProcessor{
		verifier:   mutant.NewMutanVerifier(),
		repository: repository.NewDynamoRepository("DNA"),
		wg:         wg,
	}
}

func (p *MutantProcessor) ProcessDna(dna *Dna) (bool, error) {
	verifier := p.verifier
	isMutant := verifier.IsMutant(dna.Sequence)
	p.wg.Add(1)
	go p.saveDna(dna.Sequence, isMutant)
	log.Println("Is mutant ?", isMutant)
	return isMutant, nil
}

func (p *MutantProcessor) GetStats() (*Stats, error) {
	r := p.repository
	result, err := r.CountMutants()
	if err != nil {
		return &Stats{}, err
	}
	ratioR := 0.0
	mutantCount := result.CountResult
	humanCount := result.TotalCount - mutantCount
	if result.TotalCount != 0 && humanCount != 0 {
		ratio := float32(mutantCount) / float32(humanCount)
		ratioR = math.Round(float64(ratio*100)) / 100
	}
	return &Stats{
		CountMutantDna: mutantCount,
		CountHumanDna:  humanCount,
		Ratio:          ratioR,
	}, nil
}

func (p *MutantProcessor) saveDna(sequence []string, isMutant bool) {
	r := p.repository
	bytes, _ := json.Marshal(sequence)
	entity := &repository3.DnaEntity{
		Dna:      string(bytes),
		Id:       uuid.New().String(),
		IsMutant: isMutant,
	}
	err := r.Save(entity)
	if err != nil {
		return
	}
	log.Println("Dna Persisted")
	p.wg.Done()
}
