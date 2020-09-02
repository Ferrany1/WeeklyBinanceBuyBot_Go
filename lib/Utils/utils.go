package Utils

import (
	"log"
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
