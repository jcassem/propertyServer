package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"service.com/property"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

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
		b, err := json.Marshal(property.ListProperties())
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
