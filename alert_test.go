package main

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/asendia/salmonping/db"
	"github.com/joho/godotenv"
)

func TestCreateTextMessage(t *testing.T) {
	godotenv.Load()
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		fmt.Printf("TELEGRAM_BOT_TOKEN is empty, optional and skipped")
		return
	}
	chatID, err := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	if err != nil {
		t.Errorf("Error parsing TELEGRAM_CHAT_ID")
		return
	}
	anomalies := []db.SelectOnlineListingPingsRow{
		{Status: "closed", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
		{Status: "unknown", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
		{Status: "closed", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
	}
	text := createTextMessage(anomalies)
	err = sendTelegramMessage(botToken, chatID, text)
	if err != nil {
		panic(err)
	}
	anomaly := []db.SelectOnlineListingPingsRow{
		{Status: "closed", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
	}
	text = createTextMessage(anomaly)
	err = sendTelegramMessage(botToken, chatID, text)
	if err != nil {
		panic(err)
	}
}
