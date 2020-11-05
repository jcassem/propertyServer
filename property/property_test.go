package property

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func getExpectedProperty() Property {
	return Property{
		ID:   "fefd41e1-a66f-4e18-bb81-710587ac4574",
		Name: "Test Property",
		Rent: 999.99,
	}
}

func getExpectedPropertyAsDynamoDbOutput() map[string]*dynamodb.AttributeValue {
	var expectedProperty = getExpectedProperty()
	outputItem := make(map[string]*dynamodb.AttributeValue)

	outputItem["id"] = &dynamodb.AttributeValue{
		S: aws.String(expectedProperty.ID),
	}
	outputItem["name"] = &dynamodb.AttributeValue{
		S: aws.String(expectedProperty.Name),
	}
	outputItem["rent"] = &dynamodb.AttributeValue{
		N: aws.String(fmt.Sprintf("%f", expectedProperty.Rent)),
	}

	return outputItem
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
	output.Item = getExpectedPropertyAsDynamoDbOutput()

	return output, fd.err
}

// Mock Scan such that the output returned carries values identical to input.
func (fd *fakeDynamoDB) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	output := new(dynamodb.ScanOutput)
	outputItems := []map[string]*dynamodb.AttributeValue{getExpectedPropertyAsDynamoDbOutput()}
	output.SetItems(outputItems)

	return output, fd.err
}

// Mock Scan such that the output returned carries values identical to input.
func (fd *fakeDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	output := new(dynamodb.PutItemOutput)
	return output, fd.err
}

func TestGetProperty(t *testing.T) {
	var expected = getExpectedProperty()

	getter := new(DbSession)
	getter.DynamoDB = &fakeDynamoDB{}

	if actual := GetProperty(expected.ID, getter); actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func TestGetPropertyList(t *testing.T) {
	var expected = []Property{getExpectedProperty()}

	getter := new(DbSession)
	getter.DynamoDB = &fakeDynamoDB{}

	if actual := GetPropertyList(getter); !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func TestCreateProperty(t *testing.T) {
	var expected = getExpectedProperty()

	getter := new(DbSession)
	getter.DynamoDB = &fakeDynamoDB{}

	actual := CreateProperty(expected, getter)

	if expected.Name != actual.Name || expected.Rent != actual.Rent {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func TestUpdateProperty(t *testing.T) {
	var expected = getExpectedProperty()

	getter := new(DbSession)
	getter.DynamoDB = &fakeDynamoDB{}

	if actual := UpdateProperty(expected.ID, expected, getter); !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}
