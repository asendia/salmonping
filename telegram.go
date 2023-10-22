package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SendMessageRequest struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

func sendTelegramMessage(token string, chatID int64, text string) error {
	const telegramAPIURL = "https://api.telegram.org/bot"
	url := fmt.Sprintf("%s%s/sendMessage", telegramAPIURL, token)

	requestBody, err := json.Marshal(&SendMessageRequest{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "Markdown", // or "HTML", if you want
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed sending message: %s", body)
	}
	return nil
}

func storeStatusToEmoji(status string) string {
	switch status {
	case "open":
		return "✅"
	case "closed":
		return "🚫"
	case "unknown":
		return "❓"
	default:
		return "❓"
	}
}
