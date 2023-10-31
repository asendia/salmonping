package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/asendia/salmonping/db"
)

func sendAlerts(ctx context.Context, queries *db.Queries) error {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN is not set")
	}
	chatID, err := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	if err != nil {
		return err
	}
	schedules, err := getTodaySchedules(ctx, queries)
	if err != nil {
		return err
	}
	listingPings, err := getTodayPings(ctx, queries)
	if err != nil {
		return err
	}
	anomalies := getPingAnomalies(schedules, listingPings)
	totalAnomalies := len(anomalies)
	if totalAnomalies == 0 {
		return nil
	}

	text := createTextMessage(anomalies)
	if text == "" {
		return nil
	}
	err = sendTelegramMessage(botToken, chatID, text)
	if err != nil {
		return err
	}
	return nil
}

func createTextMessage(anomalies []db.SelectOnlineListingPingsRow) string {
	var text string
	if len(anomalies) == 0 {
		return ""
	}
	if len(anomalies) == 1 {
		row := anomalies[0]
		text = fmt.Sprintf("%s [%s](%s) is %s\n", storeStatusToEmoji(row.Status), row.Name, row.Url, row.Status)
	} else {
		text = "🚨🚨 Anomalies detected 🚨🚨\n\n"
		// Check if current ping status is different from previous ping status
		// If different, send message to Telegram
		for num, row := range anomalies {
			// Check if row.Name key exists in currentListingPingMap
			text += fmt.Sprintf("%d. %s [%s - %s](%s)\n", num+1, storeStatusToEmoji(row.Status), row.Name, row.Platform, row.Url)
		}
	}
	text += "\nsalmonfit.com/status"
	return text
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
