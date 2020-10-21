# Property Server Go

REST API for Property management written in Go to be deployed as AWS Lambdas.

## Build/Deploy

To build the code:
```bash
cd property
go build main.go
```

Run locally using [AWS Serverless Application Model (AWS SAM)](https://amzn.to/37uQjEa):
```bash
go build main.go
sam local start-api
```

To build and update an AWS lambda:
```bash
cd property
GOOS=linux go build main.go
zip function.zip main
aws lambda update-function-code --function-name property-server-go --zip-file fileb://function.zip
```

Upload zip to AWS Lambda using AWS console (auto-deployemnt coming soon).