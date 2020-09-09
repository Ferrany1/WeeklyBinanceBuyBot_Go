package LRequest

import (
	"WeeklyBinanceBuyBot_Go/lib/Binance"
	"WeeklyBinanceBuyBot_Go/lib/Dirs"
	"WeeklyBinanceBuyBot_Go/lib/Spreedsheet"
	"WeeklyBinanceBuyBot_Go/lib/Telegram"
	"WeeklyBinanceBuyBot_Go/lib/Utils"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"strings"
	"time"
)

var (
	Config = Dirs.ReadFile("/Config.json")
)

func HandleRequest(event events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	eventReader := strings.NewReader(event.Body)
	body := &Telegram.WebhookReqBody{}

	switch err := json.NewDecoder(eventReader).Decode(body); {

	case err != nil:
		HandleNonApiRequest()

	default:
		HandleApiRequest(body)
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func HandleNonApiRequest() {
	switch err := Binance.MarketOrder(Binance.GetUsdtBalanceToTrade()); {
	case err == nil:
		time.Sleep(10 * time.Second)
		Spreedsheet.EditingSheet(Binance.GetLastTrade())
		log.Println("Order was placed and recorded")
	default:
		Utils.Fatal(err)
	}
}

func HandleApiRequest(body *Telegram.WebhookReqBody) {
	var (
		TChatID  = body.Message.Chat.ID
		TMessage = body.Message.Text
	)

	switch {

	case TChatID == Config.Telegram.ChatID:

		switch {

		case TMessage == "/profit":
			var (
				BlastPrice             = Binance.LastPrice()
				ETH, USDT, RateETHUSDT = Spreedsheet.GetAveragePriceHistory()
				PL                     = ETH*BlastPrice - USDT
				PLProcent              = (1 - RateETHUSDT/BlastPrice) * 100
				Text                   = fmt.Sprintf(
					"Total:\nETH: %0.6f\nUSDT: %0.2f$\nBuy Rate: %0.2f$\nCurrent Rate: %0.2f$\n\nP/L: %0.2f$ (%0.2f%%)",
					ETH, USDT, RateETHUSDT, BlastPrice, PL, PLProcent,
				)
			)
			Telegram.DeleteCommand(body)
			Telegram.SendMessage(body, Text)
			log.Println("Telegram message sent")

		default:
			var Text = "Wrong command, only /profit available"

			Telegram.SendMessage(body, Text)
			log.Println("Telegram error sent")
		}
	default:
		log.Println(body)
	}
}
