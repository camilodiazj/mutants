package service

import (
	"encoding/json"
	dnaRepository "github.com/camilodiazj/mutants/application/repository"
	"github.com/camilodiazj/mutants/domain/mutant"
	"github.com/google/uuid"
	"log"
	"math"
)

type Stats struct {
	CountMutantDna uint64  `json:"count_mutant_dna"`
	CountHumanDna  uint64  `json:"count_human_dna"`
	Ratio          float64 `json:"ratio"`
}

type Dna struct {
	Sequence []string `json:"dna"`
}

type dnaProcessor struct {
	verifier   mutant.MutanVerifier
	repository dnaRepository.DnaRepository
}

type Processor interface {
	GetStats() (*Stats, error)
	ProcessDna(dna *Dna) (bool, error)
}

func NewDnaProcessor(verifier mutant.MutanVerifier, repository dnaRepository.DnaRepository) Processor {
	return &dnaProcessor{
		verifier:   verifier,
		repository: repository,
	}
}

func (p *dnaProcessor) ProcessDna(dna *Dna) (bool, error) {
	verifier := p.verifier
	isMutant, err := verifier.IsMutant(dna.Sequence)
	if err != nil {
		return false, err
	}
	go p.saveDna(dna.Sequence, isMutant)
	log.Println("Is mutant ?", isMutant)
	return isMutant, nil
}

func (p *dnaProcessor) GetStats() (*Stats, error) {
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

func (p *dnaProcessor) saveDna(sequence []string, isMutant bool) {
	r := p.repository
	bytes, _ := json.Marshal(sequence)
	entity := &dnaRepository.DnaEntity{
		Dna:      string(bytes),
		Id:       uuid.New().String(),
		IsMutant: isMutant,
	}
	r.Save(entity)
	log.Println("Dna Persisted")
}
