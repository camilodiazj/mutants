package infrastructure

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"log"
	"strconv"
)

var dynamo *dynamodb.DynamoDB

type Dna struct {
	IsMutant bool
	Sequence string
}

type Count struct {
	Counter    uint64
	ItemsCount uint64
}

const TableName = "DNA"

func init() {
	dynamo = connectDynamoDb()
}

func connectDynamoDb() *dynamodb.DynamoDB {
	return dynamodb.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	})))
}

//Maybe should receive as params TableName
func GetCountOf(isMutant bool) (Count, error) {
	params := &dynamodb.ScanInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":isMutant": {
				S: aws.String(strconv.FormatBool(isMutant)),
			},
		},
		FilterExpression: aws.String("isMutant = :isMutant"),
		Select:           aws.String("COUNT"),
		TableName:        aws.String(TableName),
		TotalSegments:    nil,
	}

	result, err := dynamo.Scan(params)

	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
		return Count{}, err
	}

	return Count{
		Counter:    uint64(*result.Count),
		ItemsCount: uint64(*result.ScannedCount),
	}, nil
}

func Query() {
	//Item Count updates every six hours; SAD :(
	tableDescription, err := dynamo.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(TableName),
	})

	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
	}
	params := &dynamodb.QueryInput{
		TableName: aws.String(TableName),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":isMutant": {
				S: aws.String("false"),
			},
		},
		Select:                 aws.String("COUNT"),
		KeyConditionExpression: aws.String("isMutant = :isMutant"),
	}
	result, err := dynamo.Query(params)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
	}
	count := result.Count
	fulCount := tableDescription
	fmt.Println(count, fulCount)
}

func PutItem(dna *Dna) {
	_, err := dynamo.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(uuid.New().String()),
			},
			"sequence": {
				S: aws.String(dna.Sequence),
			},
			"isMutant": {
				S: aws.String(strconv.FormatBool(dna.IsMutant)),
			},
		},
		TableName: aws.String(TableName),
	})

	if err != nil {
		fmt.Print(err)
	}
}
