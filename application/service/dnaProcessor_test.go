package service

import (
	"errors"
	"github.com/camilodiazj/mutants/application/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

var processorTester Processor

func init() {
	processorTester = NewDnaProcessor(&mutanServiceMock{}, &dynamoRepositoryMock{})
}

func TestProcessDna(t *testing.T) {
	_, err := processorTester.ProcessDna(&Dna{[]string{"", ""}})
	assert.Nil(t, err)
}

func TestProcessDnaShouldFailDueInvalidDna(t *testing.T) {
	processorTester = NewDnaProcessor(&mutanServiceMock{true}, &dynamoRepositoryMock{})
	_, err := processorTester.ProcessDna(&Dna{[]string{"", ""}})
	assert.NotNil(t, err)
}

func TestGetStats(t *testing.T) {
	stats, _ := processorTester.GetStats()
	assert.Equal(t, uint64(10), stats.CountMutantDna)
}

func TestGetStatsShouldFailDueDbError(t *testing.T) {
	processorTester = NewDnaProcessor(&mutanServiceMock{}, &dynamoRepositoryMock{true})
	_, err := processorTester.GetStats()
	assert.NotNil(t, err)
}

func (s mutanServiceMock) IsMutant(dna []string) (bool, error) {
	if s.isInvalidDna {
		return false, errors.New("Invalid Dna")
	}
	return true, nil
}

type mutanServiceMock struct{
	isInvalidDna bool
}

type dynamoRepositoryMock struct{
	isDeadService bool
}

func (r *dynamoRepositoryMock) Save(dna *repository.DnaEntity) error {
	return nil
}

func (r *dynamoRepositoryMock) CountMutants() (*repository.Count, error) {
	if r.isDeadService {
		return nil, errors.New("dead")
	}
	return &repository.Count{
		CountResult: 10,
		TotalCount:  15,
	}, nil
}
