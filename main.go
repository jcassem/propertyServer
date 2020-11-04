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

// IDParameterName Name of the expected Id path parameter for GET requests
const IDParameterName = "id"

const (
	httpMethodNotSupportedErrorMessage = "HTTP Method Not Supported"

	requestBodyCannotBeEmptyErrorMessage = "Request Body Cannot Be Empty"

	jsonTransformationErrorMessage = "JSON Transformation Error"
)

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
	lambda.Start(HandleRequest)
}

// HandleRequest Handles REST routing
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Headers:")
	for key, value := range request.Headers {
		fmt.Printf("  %s: %s\n", key, value)
	}

	if request.HTTPMethod == "GET" {
		if id, ok := request.PathParameters[IDParameterName]; ok {
			return handleGet(id)
		}

		return handleGetList()
	} else if request.HTTPMethod == "POST" {
		if request.Body == "" {
			return events.APIGatewayProxyResponse{Body: requestBodyCannotBeEmptyErrorMessage, StatusCode: 502}, errors.New(requestBodyCannotBeEmptyErrorMessage)
		}

		return handlePost(request.Body)
	} else {
		return events.APIGatewayProxyResponse{Body: httpMethodNotSupportedErrorMessage, StatusCode: 502}, errors.New(httpMethodNotSupportedErrorMessage)
	}
}

func handleGetList() (events.APIGatewayProxyResponse, error) {
	fmt.Printf("(GET) LIST\n")

	b, err := json.Marshal(property.GetPropertyList(dbSession))
	if err != nil {
		return events.APIGatewayProxyResponse{Body: jsonTransformationErrorMessage, StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{Body: string(b), StatusCode: 200}, nil
}

func handleGet(id string) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("(GET) ITEM\n")

	b, err := json.Marshal(property.GetProperty(id, dbSession))
	if err != nil {
		return events.APIGatewayProxyResponse{Body: jsonTransformationErrorMessage, StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{Body: string(b), StatusCode: 200}, nil
}

func handlePost(requestBody string) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("(POST) ITEM\n")

	var propertyToCreate property.Property
	json.Unmarshal([]byte(requestBody), &propertyToCreate)

	b, err := json.Marshal(property.CreateProperty(propertyToCreate, dbSession))
	if err != nil {
		return events.APIGatewayProxyResponse{Body: jsonTransformationErrorMessage, StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{Body: string(b), StatusCode: 200}, nil
}
