package property

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// DbSession can be assigned a DynamoDB connector like:
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

// tableName DynamoDb table name to query against
const tableName = "props"

// List Lists all properties
func List(ig *DbSession) []Property {
	fmt.Printf("List Properties\n")
	propertyList := []Property{}

	queryParams := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := ig.DynamoDB.Scan(queryParams)
	if err != nil {
		panic(fmt.Sprintf("Query API call failed, %v", err))
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &propertyList)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return propertyList
}

// Get Property related to provided id.
func Get(id string, ig *DbSession) Property {
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
		panic(fmt.Sprintf("Query API call failed, %v", err))
	}

	if result.Item == nil {
		panic(fmt.Sprintf("Could not find '%s'", id))
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &property)

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return property
}
