package property

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// Property model
type Property struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Rent float32 `json:"rent"`
}

const tableName = "props"

// ListProperties Lists all properties
func ListProperties(dynamoDbSession *dynamodb.DynamoDB) []Property {
	fmt.Printf("List Properties\n")

	// Create the Expression to fill the input struct with.
	// Get all properties with name
	filt := expression.Name("name").Equal(expression.Value("123 Fake Street"))

	// Get specific fields
	proj := expression.NamesList(expression.Name("ID"), expression.Name("Name"), expression.Name("Rent"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		fmt.Println("Got error building expression:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	propertyList := []Property{}

	// Make the DynamoDB Query API call
	result, err := dynamoDbSession.Scan(params)
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
