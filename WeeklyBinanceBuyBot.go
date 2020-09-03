package main

import (
	"WeeklyBinanceBuyBot_Go/lib/LRequest"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(LRequest.HandleRequest)
}
