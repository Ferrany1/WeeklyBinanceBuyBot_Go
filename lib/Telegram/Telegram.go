package Telegram

import (
	"WeeklyBinanceBuyBot_Go/lib/Dirs"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

var (
	Config      = Dirs.ReadFile("/Config.json")
	API         = Config.Telegram.API
	URLTelegram = fmt.Sprintf("https://api.telegram.org/bot%s/", API)
)

type WebhookReqBody struct {
	Message struct {
		MessageID int `json:"message_id"`

		From struct {
			Username string `json:"username"`
		} `json:"from"`

		Date int `json:"date"`

		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`

		Text string `json:"text"`
	} `json:"message"`
}

type sendMessageReqBody struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	MessageID int    `json:"message_id"`
}

func SendMessage(b *WebhookReqBody, Text string) {
	reqBody := &sendMessageReqBody{
		ChatID: b.Message.Chat.ID,
		Text:   Text,
	}

	Send("sendMessage", reqBody)

}

func Send(SendMethod string, reqBody *sendMessageReqBody) {

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(reqBody, reqBytes)

	res, err := http.Post(URLTelegram+SendMethod, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Fatalln(err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalln(errors.New("unexpected status" + res.Status))
	}
}
