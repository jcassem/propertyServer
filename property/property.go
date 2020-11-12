package property

import (
	"errors"
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
	// QueryErrorMessageType Error message format for failed queries.
	QueryErrorMessageType = "Query API call failed"

	// UnmarshalErrorMessageType Error message format for failed conversions from json to type.
	UnmarshalErrorMessageType = "Failed to unmarshal item"

	// NotFoundErrorMessageType Error message format to use when an item has not been found.
	NotFoundErrorMessageType = "Could not find property"

	// InvalidPropertErrorMessageType Error message format to use when an item is missing key information.
	InvalidPropertErrorMessageType = "Property is not valid"

	// PersistenceErrorMessageType Error message format to use when an error occurred during persisting/saving an item.
	PersistenceErrorMessageType = "An error occured while saving"
)

// ServiceError Error wrapper for service
type ServiceError struct {
	ErrorType string // Error handling constant
	Error     error  // Error
}

// tableName DynamoDb table name to query against
const tableName = "props"

// GetPropertyList Lists all properties
func GetPropertyList(ig *DbSession) (*[]Property, *ServiceError) {
	fmt.Printf("List Properties\n")
	propertyList := []Property{}

	queryParams := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := ig.DynamoDB.Scan(queryParams)
	if err != nil {
		serviceError := ServiceError{
			ErrorType: QueryErrorMessageType,
			Error:     err,
		}
		return nil, &serviceError
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &propertyList)
	if err != nil {
		serviceError := ServiceError{
			ErrorType: UnmarshalErrorMessageType,
			Error:     err,
		}
		return nil, &serviceError
	}

	return &propertyList, nil
}

// GetProperty Property related to provided id.
func GetProperty(id string, ig *DbSession) (*Property, *ServiceError) {
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
		serviceError := ServiceError{
			ErrorType: QueryErrorMessageType,
			Error:     err,
		}
		return nil, &serviceError
	}

	if result.Item == nil {
		serviceError := ServiceError{
			ErrorType: NotFoundErrorMessageType,
			Error:     err,
		}
		return nil, &serviceError
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &property)
	if err != nil {
		serviceError := ServiceError{
			ErrorType: UnmarshalErrorMessageType,
			Error:     err,
		}
		return nil, &serviceError
	}

	return &property, nil
}

// CreateProperty Persists the provided property item.
func CreateProperty(property Property, ig *DbSession) (*Property, *ServiceError) {
	if property.Name == "" {
		serviceError := ServiceError{
			ErrorType: InvalidPropertErrorMessageType,
			Error:     errors.New("Name missing"),
		}
		return nil, &serviceError
	}

	property.ID = fmt.Sprintf("%v", uuid.Must(uuid.NewRandom()))
	fmt.Printf("Create property with id: %s\n", property.ID)

	return persistProperty(property, ig)
}

// UpdateProperty Updates and persists the provided property details against the provided id.
func UpdateProperty(id string, property Property, ig *DbSession) (*Property, *ServiceError) {
	if id == "" || id != property.ID {
		serviceError := ServiceError{
			ErrorType: InvalidPropertErrorMessageType,
			Error:     errors.New("Id missing"),
		}
		return nil, &serviceError
	}

	fmt.Printf("Update property with id: %s\n", property.ID)
	return persistProperty(property, ig)
}

// DeleteProperty Deletes property associated to the provided id.
func DeleteProperty(id string, ig *DbSession) *ServiceError {
	if id == "" {
		serviceError := ServiceError{
			ErrorType: InvalidPropertErrorMessageType,
			Error:     errors.New("Id missing"),
		}
		return &serviceError
	}

	fmt.Printf("Delete property with id: %s\n", id)

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := ig.DynamoDB.DeleteItem(input)
	if err != nil {
		serviceError := ServiceError{
			ErrorType: PersistenceErrorMessageType,
			Error:     err,
		}
		return &serviceError
	}

	return nil
}

// persistProperty Saves the provided property against the provided DB session.
func persistProperty(property Property, ig *DbSession) (*Property, *ServiceError) {

	attributeValue, err := dynamodbattribute.MarshalMap(property)
	if err != nil {
		serviceError := ServiceError{
			ErrorType: UnmarshalErrorMessageType,
			Error:     err,
		}
		return nil, &serviceError
	}

	putItemInput := &dynamodb.PutItemInput{
		Item:      attributeValue,
		TableName: aws.String(tableName),
	}

	_, err = ig.DynamoDB.PutItem(putItemInput)
	if err != nil {
		serviceError := ServiceError{
			ErrorType: PersistenceErrorMessageType,
			Error:     err,
		}
		return nil, &serviceError
	}

	// Updated data not retuned so we send back the request item (or can perform a GetItem request)
	return &property, nil
}
