package main

import (
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
)

func TestTelegramSendMessage(*testing.T) {
	godotenv.Load()
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	// Parse chatID from env
	chatID, err := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	if err != nil {
		panic(err)
	}

	err = sendTelegramMessage(token, chatID, "Test message from SalmonPing")
	if err != nil {
		panic(err)
	}
}
