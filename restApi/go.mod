module property

go 1.15

replace service.com/property => ./service

require (
	github.com/aws/aws-lambda-go v1.19.1
	service.com/property v0.0.0-00010101000000-000000000000
)
