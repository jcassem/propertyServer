package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Property model
type Property struct {
	Name string  `json:"name"`
	Rent float32 `json:"rent"`
}

// Example slice
var properties = []Property{
	Property{
		Name: "123 Fake Street",
		Rent: 1200.00,
	},
	Property{
		Name: "2 Main Road",
		Rent: 899.50,
	},
	Property{
		Name: "Flat A 120 Regents Street",
		Rent: 14060.66,
	},
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
		b, err := json.Marshal(properties)
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
