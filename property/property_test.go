package property

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const (
	expectedPropertyID   = "fefd41e1-a66f-4e18-bb81-710587ac4574"
	expectedPropertyName = "Test Property"
	expectedPropertyRent = 999.99
)

func getExpectedProperty() Property {
	return Property{
		ID:   expectedPropertyID,
		Name: expectedPropertyName,
		Rent: expectedPropertyRent,
	}
}

// A fakeDynamoDB instance to use in testing
type fakeDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	// payload map[string]string // Store expected return values
	err error
}

// Mock GetItem such that the output returned carries values identical to input.
func (fd *fakeDynamoDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	output := new(dynamodb.GetItemOutput)
	output.Item = make(map[string]*dynamodb.AttributeValue)

	output.Item["id"] = &dynamodb.AttributeValue{
		S: aws.String(expectedPropertyID),
	}
	output.Item["name"] = &dynamodb.AttributeValue{
		S: aws.String(expectedPropertyName),
	}
	output.Item["rent"] = &dynamodb.AttributeValue{
		N: aws.String(fmt.Sprintf("%f", expectedPropertyRent)),
	}

	return output, fd.err
}

func TestGetProperty(t *testing.T) {

	var expectedProperty = getExpectedProperty()

	getter := new(DbSession)
	getter.DynamoDB = &fakeDynamoDB{}

	if actualProperty := GetProperty(expectedPropertyID, getter); actualProperty != expectedProperty {
		t.Errorf("Expected %v but got %v", expectedProperty, actualProperty)
	}
}
