package property

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Property model
type Property struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Rent float32 `json:"rent"`
}

// tableName DynamoDb table name to query against
const tableName = "props"

// ListProperties Lists all properties
func ListProperties(dynamoDbSession *dynamodb.DynamoDB) []Property {
	fmt.Printf("List Properties\n")
	propertyList := []Property{}

	queryParams := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynamoDbSession.Scan(queryParams)
	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		os.Exit(1)
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &propertyList)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return propertyList
}
