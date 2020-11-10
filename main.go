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

	idParameterCannotBeEmptyErrorMessage = "Id parameter missing"

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
	// fmt.Println("Headers:")
	// for key, value := range request.Headers {
	// 	fmt.Printf("  %s: %s\n", key, value)
	// }

	switch request.HTTPMethod {
	case "GET":
		if id, ok := request.PathParameters[IDParameterName]; ok {
			return handleGet(id)
		}

		return handleGetList()
	case "POST":
		if request.Body == "" {
			return raiseError(400, requestBodyCannotBeEmptyErrorMessage)
		}

		return handlePost(request.Body)
	case "PUT":
		if id, ok := request.PathParameters[IDParameterName]; ok {
			if request.Body == "" {
				return raiseError(400, requestBodyCannotBeEmptyErrorMessage)
			}

			return handlePut(id, request.Body)
		}

		return raiseError(400, idParameterCannotBeEmptyErrorMessage)
	case "DELETE":
		if id, ok := request.PathParameters[IDParameterName]; ok {
			return handleDelete(id)
		}

		return raiseError(400, idParameterCannotBeEmptyErrorMessage)
	default:
		return raiseError(502, httpMethodNotSupportedErrorMessage)
	}
}

func raiseError(statusCode int, errorMessage string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body: errorMessage, StatusCode: statusCode}, errors.New(errorMessage)
}

func handleGetList() (events.APIGatewayProxyResponse, error) {
	fmt.Printf("(GET) LIST\n")

	prop, err := property.GetPropertyList(dbSession)
	propJSON, err := json.Marshal(prop)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: jsonTransformationErrorMessage, StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{Body: string(propJSON), StatusCode: 200}, nil
}

func handleGet(id string) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("(GET) ITEM\n")

	prop, err := property.GetProperty(id, dbSession)
	propJSON, err := json.Marshal(prop)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: jsonTransformationErrorMessage, StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{Body: string(propJSON), StatusCode: 200}, nil
}

func handlePost(requestBody string) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("(POST) ITEM\n")

	var propertyToCreate property.Property
	json.Unmarshal([]byte(requestBody), &propertyToCreate)

	prop, err := property.CreateProperty(propertyToCreate, dbSession)
	propJSON, err := json.Marshal(prop)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: jsonTransformationErrorMessage, StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{Body: string(propJSON), StatusCode: 200}, nil
}

func handlePut(id string, requestBody string) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("(PUT) ITEM\n")

	var propertyToUpdate property.Property
	json.Unmarshal([]byte(requestBody), &propertyToUpdate)

	prop, err := property.UpdateProperty(id, propertyToUpdate, dbSession)
	propJSON, err := json.Marshal(prop)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: jsonTransformationErrorMessage, StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{Body: string(propJSON), StatusCode: 200}, nil
}

func handleDelete(id string) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("(DELETE) ITEM\n")

	err := property.DeleteProperty(id, dbSession)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: jsonTransformationErrorMessage, StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}
