package repository

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/camilodiazj/mutants/application/repository"
	"github.com/gusaul/go-dynamock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MyDynamo struct {
	Db dynamodbiface.DynamoDBAPI
}

var mock *dynamock.DynaMock
var Dyna *MyDynamo
var r repository.DnaRepository

func init() {
	Dyna = new(MyDynamo)
	Dyna.Db, mock = dynamock.New()
	r = NewDynamoRepository("tableName", Dyna.Db)
}

func TestSave(t *testing.T) {
	mock.ExpectPutItem().ToTable("tableName")
	err := r.Save(&repository.DnaEntity{
		Dna:      "dna_sequence",
		Id:       "uuid",
		IsMutant: false,
	})
	assert.Nil(t, err)
}

func TestSaveShouldFail(t *testing.T) {
	err := r.Save(&repository.DnaEntity{
		Dna:      "dna_sequence",
		Id:       "uuid",
		IsMutant: false,
	})
	assert.NotNil(t, err)
}

func TestCountMutants(t *testing.T) {
	var count, scannedC int64
	count = 5
	scannedC = 10
	mock.ExpectScan().Table("tableName").WillReturns(dynamodb.ScanOutput{
		Count:        &count,
		ScannedCount: &scannedC,
	})
	countResult, _ := r.CountMutants()
	assert.Equal(t, uint64(5), countResult.CountResult)
	assert.Equal(t, uint64(10), countResult.TotalCount)
}

func TestCountMutantsShouldFail(t *testing.T) {
	_, err := r.CountMutants()
	assert.NotNil(t, err)
}
