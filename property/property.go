package property

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
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

	// NotFoundErrorMessageFormat Error message format to use when an item has not been found.
	NotFoundErrorMessageFormat = "Could not find '%s'"

	// InvalidPropertErrorMessageFormat Error message format to use when an item is missing key information.
	InvalidPropertErrorMessageFormat = "Property is not valid: %v"

	// PersistenceErrorMessageFormat Error message format to use when an error occurred during persisting/saving an item.
	PersistenceErrorMessageFormat = "An error occured while saving: %v"
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

// CreateProperty Persists the provided property item.
func CreateProperty(property Property, ig *DbSession) Property {
	if property.Name == "" {
		panic(fmt.Sprintf(InvalidPropertErrorMessageFormat, property))
	}

	property.ID = fmt.Sprintf("%v", uuid.Must(uuid.NewRandom()))
	fmt.Printf("Create property with id: %s\n", property.ID)

	return persistProperty(property, ig)
}

// UpdateProperty Updates and persists the provided property details against the provided id.
func UpdateProperty(id string, property Property, ig *DbSession) Property {
	if id == "" || id != property.ID {
		panic(fmt.Sprintf(InvalidPropertErrorMessageFormat, property))
	}

	fmt.Printf("Update property with id: %s\n", property.ID)
	return persistProperty(property, ig)
}

// persistProperty Saves the provided property against the provided DB session.
func persistProperty(property Property, ig *DbSession) Property {

	attributeValue, err := dynamodbattribute.MarshalMap(property)
	if err != nil {
		panic(fmt.Sprintf(UnmarshalErrorMessageFormat, err))
	}

	putItemInput := &dynamodb.PutItemInput{
		Item:      attributeValue,
		TableName: aws.String(tableName),
	}

	_, err = ig.DynamoDB.PutItem(putItemInput)
	if err != nil {
		panic(fmt.Sprintf(PersistenceErrorMessageFormat, err))
	}

	// Updated data not retuned so we send back the request item (or can perform a GetItem request)
	return property
}
