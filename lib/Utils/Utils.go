package Utils

import (
	"log"
	"math"
)

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Println(err error) {
	if err != nil {
		log.Println(err)
	}
}

func Round(num float64, r int) float64 {
	fr := math.Pow10(r)
	num = math.Round(num*fr) / fr
	return num
}
