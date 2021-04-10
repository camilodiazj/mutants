package service

import (
	"github.com/camilodiazj/mutants/application/repository"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

var processor Processor

func init() {
	wg := sync.WaitGroup{} //TODO: Mock
	processor = NewDnaProcessor(&wg, &mutanServiceMock{}, &dynamoRepositoryMock{})
}

func TestProcessDna(t *testing.T) {
	_, err := processor.ProcessDna(&Dna{[]string{"", ""}})
	assert.Nil(t, err)
}

func TestGetStats(t *testing.T) {
	stats, _ := processor.GetStats()
	assert.Equal(t, uint64(10), stats.CountMutantDna)
}

type mutanServiceMock struct{}

func (mutanServiceMock) IsMutant(dna []string) bool {
	return true
}

type dynamoRepositoryMock struct{}

func (r *dynamoRepositoryMock) Save(dna *repository.DnaEntity) error {
	return nil
}

func (r *dynamoRepositoryMock) CountMutants() (*repository.Counter, error) {
	return &repository.Counter{
		CountResult: 10,
		TotalCount:  15,
	}, nil
}
