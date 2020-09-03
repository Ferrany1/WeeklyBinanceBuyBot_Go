package main

import (
	"WeeklyBinanceBuyBot_Go/lib/LRequest"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

func main() {

	switch OS := os.Getenv("OS"); {
	case OS == "Lambda":
		lambda.Start(LRequest.HandleRequest)
	default:
		LRequest.HandleNonApiRequest()
	}

}
