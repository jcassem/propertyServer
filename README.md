# Property Server Go

REST API for Property management written in Go to be deployed as AWS Lambdas.

## Build/Deploy

```bash
cd property
GOOS=linux go build main.go
zip function.zip main
```

Upload zip to AWS Lambda using AWS console (auto-deployemnt coming soon).