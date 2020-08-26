package main

import (
	Bin "WeeklyBinanceBuyBot_Go/lib/Binance"
	SS "WeeklyBinanceBuyBot_Go/lib/Spreedsheet"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
)

func HandleRequest() {
	Bin.MarketOrder(Bin.GetUsdtBalanceToTrade())
	SS.EditingSheet(Bin.GetLastTrade())
	log.Println("Order was placed and recorded")
}

func main() {
	switch OS := os.Getenv("OS"); {
	case OS == "Lambda":
		lambda.Start(HandleRequest)
	default:
		HandleRequest()
	}
}
