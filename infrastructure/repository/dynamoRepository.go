package repository

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/camilodiazj/mutants/application/repository"
	"log"
	"strconv"
)

type dynamoRepository struct {
	table  string
	dynamo dynamodbiface.DynamoDBAPI
}

func NewDynamoRepository(tableName string, dynamo dynamodbiface.DynamoDBAPI) repository.DnaRepository {
	return &dynamoRepository{
		table:  tableName,
		dynamo: dynamo}
}

func (r *dynamoRepository) Save(dna *repository.DnaEntity) error {
	dynamo := r.dynamo
	_, err := dynamo.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(dna.Id),
			},
			"sequence": {
				S: aws.String(dna.Dna),
			},
			"isMutant": {
				S: aws.String(strconv.FormatBool(dna.IsMutant)),
			},
		},
		TableName: aws.String(r.table),
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *dynamoRepository) CountMutants() (*repository.Counter, error) {
	dynamo := r.dynamo
	params := &dynamodb.ScanInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":isMutant": {
				S: aws.String(strconv.FormatBool(true)),
			},
		},
		FilterExpression: aws.String("isMutant = :isMutant"),
		Select:           aws.String("COUNT"),
		TableName:        aws.String(r.table),
		TotalSegments:    nil,
	}
	result, err := dynamo.Scan(params)

	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
		return &repository.Counter{}, err
	}

	return &repository.Counter{
		CountResult: uint64(*result.Count),
		TotalCount:  uint64(*result.ScannedCount),
	}, nil
}
