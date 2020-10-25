package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jcassem/propertyServer/property"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const idParameterName = "id"

// DynamoDb session
var dbSession = new(property.DbSession)

// init Initialize a db session that the SDK will use to load
// 		credentials from the shared credentials file ~/.aws/credentials
// 		and region from the shared configuration file ~/.aws/config.
func init() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	var svc *dynamodb.DynamoDB = dynamodb.New(sess)
	dbSession.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
}

func main() {
	// properties = append(properties, exampleProperty)
	lambda.Start(HandleRequest)
}

// HandleRequest Handles REST routing
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Headers:")
	for key, value := range request.Headers {
		fmt.Printf("  %s: %s\n", key, value)
	}

	if request.HTTPMethod == "GET" {
		if id, ok := request.PathParameters[idParameterName]; ok {
			return handleGet(id)
		}

		return handleGetList()
	} else if request.HTTPMethod == "POST" {
		fmt.Printf("POST METHOD\n")
		return events.APIGatewayProxyResponse{Body: "POST", StatusCode: 200}, nil
	} else {
		var errMessage = "HTTP Method Not Supported"
		fmt.Printf(errMessage + "\n")
		return events.APIGatewayProxyResponse{Body: errMessage, StatusCode: 502}, errors.New(errMessage)
	}
}

func handleGetList() (events.APIGatewayProxyResponse, error) {
	fmt.Printf("(GET) LIST\n")
	b, err := json.Marshal(property.List(dbSession))
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "JSON Transformation Error", StatusCode: 500}, err
	}
	return events.APIGatewayProxyResponse{Body: string(b), StatusCode: 200}, nil
}

func handleGet(id string) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("(GET) ITEM\n")
	b, err := json.Marshal(property.Get(id, dbSession))
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "JSON Transformation Error", StatusCode: 500}, err
	}
	return events.APIGatewayProxyResponse{Body: string(b), StatusCode: 200}, nil
}
