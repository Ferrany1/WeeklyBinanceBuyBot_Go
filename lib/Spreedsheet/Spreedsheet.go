package Spreedsheet

import (
	"context"
	"io/ioutil"
	"log"
	"strconv"

	"WeeklyBinanceBuyBot_Go/lib/Dirs"
	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

func lastCellReturn() (int, int, int) {

	sheet := callSheet()

	lastRow, lastColumn := len(sheet.Rows), len(sheet.Columns)
	lastIndex, _ := strconv.Atoi(sheet.Rows[lastRow-1][0].Value)
	return lastRow, lastColumn, lastIndex
}

func callSheet() *spreadsheet.Sheet {

	data, _ := ioutil.ReadFile(Dirs.GetFile("/client_secret.json"))

	conf, _ := google.JWTConfigFromJSON(data, spreadsheet.Scope)

	client := conf.Client(context.TODO())

	service := spreadsheet.NewServiceWithClient(client)
	ss, _ := service.FetchSpreadsheet("1U0bu2wRjlMBiX5XdGlcSZqXWh7LDwBz-k1c_8kQzWOw")

	sheet, _ := ss.SheetByIndex(0)

	return sheet
}

func EditingSheet(lastOrder []string) {

	sheet := callSheet()

	lr, lc, li := lastCellReturn()

	newLi := strconv.Itoa(li + 1)

	updatedLastOrder := []string{newLi}

	updatedLastOrder = append(updatedLastOrder, lastOrder...)

	order, _ := strconv.ParseInt(updatedLastOrder[2], 10, 64)
	sheetOrder, _ := strconv.ParseInt(sheet.Rows[lr-1][2].Value, 10, 64)

	if order != sheetOrder {
		for clc := 0; clc < lc; clc++ {
			sheet.Update(lr, clc, updatedLastOrder[clc])
		}
	} else {
		log.Fatalln("No new orders was made")
	}

	sheet.Synchronize()

}
