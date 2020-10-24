package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"service.com/property"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var dynamoDbSession (*dynamodb.DynamoDB)

// init Create dynamo db session
func init() {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	dynamoDbSession = dynamodb.New(sess)
}

func main() {
	// properties = append(properties, exampleProperty)
	lambda.Start(HandleRequest)
}

// HandleRequest Handles REST routing
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Body size = %d. \n", len(request.Body))
	fmt.Println("Headers:")
	for key, value := range request.Headers {
		fmt.Printf("  %s: %s\n", key, value)
	}
	if request.HTTPMethod == "GET" {
		fmt.Printf("GET METHOD\n")
		b, err := json.Marshal(property.ListProperties(dynamoDbSession))
		if err != nil {
			return events.APIGatewayProxyResponse{Body: "JSON Transformation Error", StatusCode: 500}, err
		}
		return events.APIGatewayProxyResponse{Body: string(b), StatusCode: 200}, nil
	} else if request.HTTPMethod == "POST" {
		fmt.Printf("POST METHOD\n")
		return events.APIGatewayProxyResponse{Body: "POST", StatusCode: 200}, nil
	} else {
		var errMessage = "HTTP Method Not Supported"
		fmt.Printf(errMessage + "\n")
		return events.APIGatewayProxyResponse{Body: errMessage, StatusCode: 502}, errors.New(errMessage)
	}
}
