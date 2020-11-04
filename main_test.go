package main

import (
	"encoding/json"
	"testing"

	"github.com/jcassem/propertyServer/property"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestGetPropertyListHandler(t *testing.T) {
	request := events.APIGatewayProxyRequest{HTTPMethod: "GET"}
	expected := []property.Property{}
	actual := []property.Property{}

	response, err := HandleRequest(nil, request)

	assert.IsType(t, nil, err)
	assert.IsType(t, "string", response.Body)

	err = json.Unmarshal([]byte(response.Body), &actual)
	if err != nil {
		t.Errorf("Error un-marshalling response, %v", err)
	}
	assert.IsType(t, expected, actual)
}

// func TestGetPropertyHandler(t *testing.T) {
// }
