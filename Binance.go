package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/adshao/go-binance"
)

func binanceClient() *binance.Client {

	var (
		keys [2]string
	)

	f, err := os.Open(getFile("/Secret.txt"))
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for n := 0; scanner.Scan(); n++ {
		text := strings.SplitAfter(scanner.Text(), " = ")
		keys[n] = text[1]

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	client := binance.NewClient(keys[0], keys[1])

	return client
}

func getUsdtBalanceToTrade() float64 {

	client := binanceClient()

	res, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}

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
	if amountToTrade <= 10. {
		amountToTrade = 0
	}
	return amountToTrade
}

func getLastTrade() []string {

	client := binanceClient()

	orders, err := client.NewListOrdersService().Symbol("ETHUSDT").
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	var (
		lastOrder = []int{}
		max       int
		orderID   string
		dateTime  string
		ETH       string
		USDT      string
		exchRate  string
		totalETH  string
		totalUSDT string
		lastTrade = []string{}
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
	if err != nil {
		fmt.Println(err)
	}

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

func marketOrder(amountToTrade float64) {

	amountToTradeI := fmt.Sprintf("%f", amountToTrade)

	client := binanceClient()

	_, err := client.NewCreateOrderService().Symbol("ETHUSDT").
		Side(binance.SideTypeBuy).Type(binance.OrderTypeMarket).QuoteOrderQty(amountToTradeI).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
