package dna

import (
	"encoding/json"
	"github.com/camilodiazj/mutants/domain/mutant"
	"github.com/camilodiazj/mutants/infrastructure/repository"
	"github.com/google/uuid"
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

type MutantProcessor struct {
	verifier   mutant.MutanVerifier
	repository Repository
}

type Processor interface {
	GetStats() (*Stats, error)
	ProcessDna(dna *Dna) (bool, error)
}

func NewDnaProcessor() Processor {
	return &MutantProcessor{
		verifier:   mutant.NewMutanVerifier(),
		repository: repository.NewDynamoRepository("DNA"),
	}
}

func (p *MutantProcessor) ProcessDna(dna *Dna) (bool, error) {
	verifier := p.verifier
	isMutant := verifier.IsMutant(dna.Sequence)
	p.saveDna(dna.Sequence, isMutant)
	return isMutant, nil
}

func (p *MutantProcessor) saveDna(sequence []string, isMutant bool) {
	r := p.repository
	bytes, _ := json.Marshal(sequence)
	entity := &Entity{
		Dna:      string(bytes),
		Id:       uuid.New().String(),
		IsMutant: isMutant,
	}
	err := r.Save(entity)
	if err != nil {
		return
	}
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
	if result.TotalCount != 0 {
		ratio := float32(mutantCount) / float32(humanCount)
		ratioR = math.Round(float64(ratio*100)) / 100
	}
	return &Stats{
		CountMutantDna: mutantCount,
		CountHumanDna:  humanCount,
		Ratio:          ratioR,
	}, nil
}