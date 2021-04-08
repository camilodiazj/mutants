package repository

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/camilodiazj/mutants/application/dna"
	"log"
	"strconv"
)

type dynamoRepository struct {
	table  string
	dynamo *dynamodb.DynamoDB
}

func NewDynamoRepository(tableName string) dna.Repository {
	return &dynamoRepository{
		table:  tableName,
		dynamo: connectDynamoDb()}
}

func (r *dynamoRepository) Save(dna *dna.Entity) error {
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

func (r *dynamoRepository) CountMutants() (*dna.Counter, error) {
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
		return &dna.Counter{}, err
	}

	return &dna.Counter{
		CountResult: uint64(*result.Count),
		TotalCount:  uint64(*result.ScannedCount),
	}, nil
}

func connectDynamoDb() *dynamodb.DynamoDB {
	return dynamodb.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	})))
}
