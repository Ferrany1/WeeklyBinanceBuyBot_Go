package Dirs

import (
	"WeeklyBinanceBuyBot_Go/lib/Utils"
	"bytes"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Creds struct {
	Binance      Binance
	SpereedSheet SpereedSheet
	Telegram     Telegram
}

type Binance struct {
	Key    string
	Secret string
}

type SpereedSheet struct {
	ID string
}

type Telegram struct {
	API    string
	ChatID int64
}

func GetFile(newFile string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.Buffer{}
	buf.WriteString(dir)
	buf.WriteString(newFile)
	result := buf.String()

	return result
}

func ReadFile(newFile string) Creds {
	f, err := os.Open(GetFile(newFile))
	Utils.Fatal(err)

	defer func() {
		err := f.Close()
		Utils.Fatal(err)
	}()

	Utils.Println(err)

	Credis := Creds{}
	err = json.NewDecoder(f).Decode(&Credis)

	Utils.Println(err)

	return Credis
}
