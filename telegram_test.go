package main

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
)

func TestTelegramSendMessage(t *testing.T) {
	godotenv.Load()
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		fmt.Printf("TELEGRAM_BOT_TOKEN is empty, optional and skipped")
		return
	}
	// Parse chatID from env
	chatID, err := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	if err != nil {
		t.Errorf("Error parsing TELEGRAM_CHAT_ID")
		return
	}

	err = sendTelegramMessage(token, chatID, "Test message from SalmonPing")
	if err != nil {
		panic(err)
	}
}
