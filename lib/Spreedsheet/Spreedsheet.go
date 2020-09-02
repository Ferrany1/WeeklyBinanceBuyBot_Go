package Spreedsheet

import (
	"WeeklyBinanceBuyBot_Go/lib/Utils"
	"context"
	"io/ioutil"
	"log"
	"strconv"

	"WeeklyBinanceBuyBot_Go/lib/Dirs"
	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

var SSID = Dirs.ReadFile("/Config.json").SpereedSheet.ID

func CallSheet() *spreadsheet.Sheet {

	data, _ := ioutil.ReadFile(Dirs.GetFile("/client_secret.json"))

	conf, _ := google.JWTConfigFromJSON(data, spreadsheet.Scope)

	client := conf.Client(context.TODO())

	service := spreadsheet.NewServiceWithClient(client)
	ss, _ := service.FetchSpreadsheet(SSID)

	sheet, _ := ss.SheetByIndex(0)

	return sheet
}

func LastCellReturn() (int, int, int) {

	sheet := CallSheet()

	lastRow, lastColumn := len(sheet.Rows), len(sheet.Columns)
	lastIndex, _ := strconv.Atoi(sheet.Rows[lastRow-1][0].Value)

	return lastRow, lastColumn, lastIndex
}

func EditingSheet(lastOrder []string) {

	sheet := CallSheet()

	lr, lc, li := LastCellReturn()

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

	err := sheet.Synchronize()
	Utils.Println(err)

}
