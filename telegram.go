package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/asendia/salmonping/db"
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

func sendTelegramAlert(rows []db.SelectOnlineListingPingsRow) error {
	if !isWithinTimeCriteria() {
		return nil
	}

	chatID, err := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	if err != nil {
		return err
	}

	if len(rows) == 0 {
		return nil
	}
	ctr := 0
	text := "🚨🚨 Anomaly detected 🚨🚨\n\n"
	for _, row := range rows {
		if row.Status == "open" {
			continue
		}
		ctr++
		text += fmt.Sprintf("%d. %s [%s](%s)\n", ctr, storeStatusToEmoji(row.Status), row.Name, row.Url)
	}

	if ctr == 0 {
		return nil
	}

	// Send message
	err = sendTelegramMessage(os.Getenv("TELEGRAM_BOT_TOKEN"), chatID, text)
	if err != nil {
		return err
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

func isWithinTimeCriteria() bool {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// handle error
		return false
	}

	now := time.Now().In(loc)
	dayOfWeek := now.Weekday()
	hour := now.Hour()

	// Check if it's within Monday - Saturday
	if dayOfWeek == time.Sunday {
		return false
	}

	// Check if it's within 10AM - 8PM
	if hour < 10 || hour >= 20 {
		return false
	}

	return true
}
