package Binance

import (
	"WeeklyBinanceBuyBot_Go/lib/Dirs"
	"WeeklyBinanceBuyBot_Go/lib/Utils"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/adshao/go-binance"
)

var (
	Config = Dirs.ReadFile("/Config.json")
	Key    = Config.Binance.Key
	Secret = Config.Binance.Secret
)

type CurrentPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func binanceClient() *binance.Client {
	client := binance.NewClient(Key, Secret)

	return client
}

func GetUsdtBalanceToTrade() float64 {
	client := binanceClient()

	res, err := client.NewGetAccountService().Do(context.Background())
	Utils.Println(err)

	var (
		USDT                  float64
		remainingWeeksDivider float64
	)

	for _, asset := range res.Balances {
		if asset.Asset == "USDT" {
			USDT, _ = strconv.ParseFloat(asset.Free, 64)
		}
	}

	switch todayDay := time.Now().Day(); {
	case todayDay >= 28 || todayDay%7 == 0:
		remainingWeeksDivider = float64(5 - (todayDay / 7))
	default:
		remainingWeeksDivider = float64(5 - (todayDay/7 + 1))
	}

	amountToTrade := math.Round((USDT/remainingWeeksDivider)*100) / 100
	switch {
	case amountToTrade <= 10.:
		amountToTrade = 0
	case amountToTrade >= 10. && amountToTrade <= 11.:
		amountToTrade = float64(int(amountToTrade)) + 0.01
	default:
		amountToTrade = float64(int(amountToTrade))
	}
	return amountToTrade
}

func GetLastTrade() []string {
	client := binanceClient()

	orders, err := client.NewListOrdersService().Symbol("ETHUSDT").
		Do(context.Background())
	Utils.Println(err)

	var (
		lastOrder []int
		max       int
		orderID   string
		dateTime  string
		ETH       string
		USDT      string
		exchRate  string
		totalETH  string
		totalUSDT string
		lastTrade []string
	)

	for _, o := range orders {
		lastOrder = append(lastOrder, int(o.Time))

	}

	for _, lO := range lastOrder {
		if lO > max {
			max = lO
		}
	}

	for _, o := range orders {
		if int(o.Time) == int(max) {

			rate1, _ := strconv.ParseFloat(o.CummulativeQuoteQuantity, 64)
			rate2, _ := strconv.ParseFloat(o.ExecutedQuantity, 64)

			dateTime = time.Unix(o.Time/1000, 0).Format("2006-01-02 15:04:05")
			orderID = strconv.FormatInt(o.OrderID, 10)
			ETH = strings.Replace(o.ExecutedQuantity, ".", ",", 1)
			USDT = strings.Replace(o.CummulativeQuoteQuantity, ".", ",", 1)
			exchRate = strings.Replace(strconv.FormatFloat(rate1/rate2, 'g', -1, 32), ".", ",", 1)

		}
	}

	res, err := client.NewGetAccountService().Do(context.Background())
	Utils.Println(err)

	for _, asset := range res.Balances {
		if asset.Asset == "USDT" {
			totalUSDT = strings.Replace(asset.Free, ".", ",", 1)
		} else if asset.Asset == "ETH" {
			totalETH = strings.Replace(asset.Free, ".", ",", 1)
		}
	}
	lastTrade = []string{dateTime, orderID, ETH, USDT, exchRate, totalETH, totalUSDT}
	return lastTrade
}

func MarketOrder(amountToTrade float64) error {
	amountToTradeI := fmt.Sprintf("%f", amountToTrade)

	client := binanceClient()

	_, err := client.NewCreateOrderService().Symbol("ETHUSDT").
		Side(binance.SideTypeBuy).Type(binance.OrderTypeMarket).QuoteOrderQty(amountToTradeI).Do(context.Background())

	return err
}

func LastPrice() float64 {
	var (
		symbol = "ETHUSDT"
		url    = fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbol)
	)

	resp, err := http.Get(url)
	Utils.Println(err)

	var CP CurrentPrice
	err = json.NewDecoder(resp.Body).Decode(&CP)
	Utils.Println(err)

	price, err := strconv.ParseFloat(CP.Price, 64)
	Utils.Println(err)

	return price
}
