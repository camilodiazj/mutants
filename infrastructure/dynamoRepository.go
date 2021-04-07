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

type Stats struct {
	CountMutantDna float32
	CountHumanDna  float32
	Ratio          float32
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

func GetItem() {
	params := &dynamodb.ScanInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":isMutant": {
				S: aws.String("false"),
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
	}

	count := float32(*result.Count)
	fulCount := float32(*result.ScannedCount)
	fmt.Sprintf("%.1f\n", count/fulCount)
	var stats = &Stats{
		CountHumanDna:  count,
		CountMutantDna: fulCount - count,
		Ratio:          count / fulCount,
	}

	fmt.Println(stats)
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

func PutItem(dna Dna) {
	output, err := dynamo.PutItem(&dynamodb.PutItemInput{
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
	} else {
		fmt.Println(output.ConsumedCapacity)
	}
}
