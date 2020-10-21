# Property Server Go

REST API for Property management written in Go to be deployed as AWS Lambdas.

## Build/Deploy

To build the code:
```bash
cd property
go build main.go
```

To build and update an AWS lambda:
```bash
cd property
GOOS=linux go build main.go
zip function.zip main
aws lambda update-function-code --function-name property-server-go --zip-file fileb://function.zip
```

Upload zip to AWS Lambda using AWS console (auto-deployemnt coming soon).