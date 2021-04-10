package configuration

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/camilodiazj/mutants/application/service"
	"github.com/camilodiazj/mutants/domain/mutant"
	"github.com/camilodiazj/mutants/infrastructure/repository"
	"github.com/gorilla/mux"
	"sync"
)

var wg sync.WaitGroup

type Injections struct {
	Processor service.Processor
	Router    *mux.Router
}

func GetInjections() *Injections {
	verifier := mutant.NewMutanVerifier()
	dynamoRepository := repository.NewDynamoRepository("DNA", configureDynamoDB())
	return &Injections{
		Processor: service.NewDnaProcessor(&wg, verifier, dynamoRepository),
		Router:    mux.NewRouter(),
	}
}

func configureDynamoDB() dynamodbiface.DynamoDBAPI {
	awsSession, _ := session.NewSession(&aws.Config{Region: aws.String("us-east-2")})
	svc := dynamodb.New(awsSession)
	return dynamodbiface.DynamoDBAPI(svc)
}
