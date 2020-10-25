# Property Server Go

REST API for Property management written in Go to be deployed as AWS Lambdas.

This uses an AWS DynamoDb instance with a partition key of 'id', which is a string intended to house UUIDs.

Both the Lambda and DynamoDb are assumed to be in the same region.

## Build/Deploy

To build the code:
```bash
go build main.go
```

Run latest build locally using [AWS Serverless Application Model (AWS SAM)](https://amzn.to/37uQjEa):
```bash
sam local start-api
```

To build and update an AWS lambda:
```bash
GOOS=linux go build main.go
zip function.zip main
aws lambda update-function-code --function-name property-server-go --zip-file fileb://function.zip
```

Upload zip to AWS Lambda using AWS console (auto-deployemnt coming soon).