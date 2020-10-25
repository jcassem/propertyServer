module propertyServer

go 1.15

replace github.com/jcassem/propertyServer/property => ./property

require (
	github.com/aws/aws-lambda-go v1.19.1
	github.com/aws/aws-sdk-go v1.35.12
	github.com/google/uuid v1.1.2 // indirect
	github.com/jcassem/propertyServer/property v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.6.1
)
