package property

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// DbSession Wrapper pof a DynamoDB connector. Example of assignment:
//	svc := dynamodb.DynamoDB(sess)
//	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
type DbSession struct {
	DynamoDB dynamodbiface.DynamoDBAPI
}

// Property model
type Property struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Rent float32 `json:"rent"`
}

const (
	// QueryErrorMessageFormat Error message format for failed queries.
	QueryErrorMessageFormat = "Query API call failed, %v"

	// UnmarshalErrorMessageFormat Error message format for failed conversions from json to type.
	UnmarshalErrorMessageFormat = "Failed to unmarshal item, %v"

	// NotFoundErrorMessageFormat Error message format to use when an item has not been found
	NotFoundErrorMessageFormat = "Could not find '%s'"
)

// tableName DynamoDb table name to query against
const tableName = "props"

// GetPropertyList Lists all properties
func GetPropertyList(ig *DbSession) []Property {
	fmt.Printf("List Properties\n")
	propertyList := []Property{}

	queryParams := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := ig.DynamoDB.Scan(queryParams)
	if err != nil {
		panic(fmt.Sprintf(QueryErrorMessageFormat, err))
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &propertyList)
	if err != nil {
		panic(fmt.Sprintf(UnmarshalErrorMessageFormat, err))
	}

	return propertyList
}

// GetProperty Property related to provided id.
func GetProperty(id string, ig *DbSession) Property {
	fmt.Printf("Get Property with id: %s\n", id)
	property := Property{}

	result, err := ig.DynamoDB.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		panic(fmt.Sprintf(QueryErrorMessageFormat, err))
	}

	if result.Item == nil {
		panic(fmt.Sprintf(NotFoundErrorMessageFormat, id))
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &property)

	if err != nil {
		panic(fmt.Sprintf(UnmarshalErrorMessageFormat, err))
	}

	return property
}
