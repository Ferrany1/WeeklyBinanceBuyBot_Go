package main

import "github.com/aws/aws-lambda-go/lambda"

func HandleRequest() {

	marketOrder(getUsdtBalanceToTrade())
	editingSheet(getLastTrade())

}

func main() {

	lambda.Start(HandleRequest)

}
