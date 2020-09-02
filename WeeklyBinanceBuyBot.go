package main

import (
	Bin "WeeklyBinanceBuyBot_Go/lib/Binance"
	SS "WeeklyBinanceBuyBot_Go/lib/Spreedsheet"
	"WeeklyBinanceBuyBot_Go/lib/Utils"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
)

func HandleRequest() {
	switch err := Bin.MarketOrder(Bin.GetUsdtBalanceToTrade()); {
	case err == nil:
		SS.EditingSheet(Bin.GetLastTrade())
		log.Println("Order was placed and recorded")
	default:
		Utils.Fatal(err)
	}
}

func main() {
	switch OS := os.Getenv("OS"); {
	case OS == "Lambda":
		lambda.Start(HandleRequest)
	default:
		HandleRequest()
	}
}
