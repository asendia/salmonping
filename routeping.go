package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/asendia/salmonping/db"
)

func routePing(w http.ResponseWriter, r *http.Request) {
	// Check API key
	apiKey := r.Header.Get("X-API-Key")
	if apiKey != os.Getenv("API_KEY") {
		j, _ := logJson(map[string]interface{}{
			"level":   "warning",
			"message": "Invalid API key",
		})
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(j))
		return
	}

	ctx := r.Context()
	tx, conn, _, message, err := prepareDBConn(ctx)
	if conn != nil {
		defer conn.Close(ctx)
	}
	if tx != nil {
		// Commit everything
		defer tx.Commit(ctx)
	}
	if err != nil {
		j, _ := logJson(map[string]interface{}{
			"level":   "error",
			"error":   err.Error(),
			"message": message,
		})
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(j))
		return
	}
	queries := db.New(tx)

	listings, err := queries.SelectListings(ctx)
	if err != nil {
		j, _ := logJson(map[string]interface{}{
			"level":   "error",
			"error":   err.Error(),
			"message": "Error selecting listings",
		})
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(j))
		return
	}
	scrapeListings(ctx, queries, listings)

	err = verifyPingAndSendAlert(ctx, queries)
	if err != nil {
		logJson(map[string]interface{}{
			"level":   "warning",
			"error":   err.Error(),
			"message": "Failed to send Telegram alert",
		})
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "ok"}`)
}

func scrapeListings(ctx context.Context, queries *db.Queries, listings []db.OnlineListing) {
	for _, ol := range listings {
		logJson(map[string]interface{}{
			"level":   "info",
			"message": "Scraping a listing",
			"listing": ol.Name,
			"url":     ol.Url,
		})
		var status string
		var err error
		if ol.Platform == "gofood" {
			status, err = getGofoodStatus(ol.Url)
		} else if ol.Platform == "grabfood" {
			status, err = getGrabfoodStatus(ol.Url)
		} else {
			logJson(map[string]interface{}{
				"level":   "error",
				"message": "Unsupported url",
				"listing": ol.Name,
				"url":     ol.Url,
			})
			continue
		}
		if err != nil {
			logJson(map[string]interface{}{
				"level":   "error",
				"message": "Error scraping a listing",
				"listing": ol.Name,
				"error":   err.Error(),
			})
			continue
		}

		logJson(map[string]interface{}{
			"level":   "info",
			"listing": ol.Name,
			"status":  status,
		})

		// Log to database
		_, err = queries.InsertPing(ctx, db.InsertPingParams{
			OnlineListingID: ol.ID,
			Status:          status,
		})
		if err != nil {
			logJson(map[string]interface{}{
				"level":   "error",
				"message": "Error inserting a ping",
				"listing": ol.Name,
				"error":   err.Error(),
			})
			continue
		}
	}
}

func verifyPingAndSendAlert(ctx context.Context, queries *db.Queries) error {
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
	if len(anomalies) == 0 {
		return nil
	}
	text := "🚨🚨 Anomaly detected 🚨🚨\n\n"
	// Check if current ping status is different from previous ping status
	// If different, send message to Telegram
	for num, row := range anomalies {
		// Check if row.Name key exists in currentListingPingMap
		text += fmt.Sprintf("%d. %s [%s](%s)\n", num+1, storeStatusToEmoji(row.Status), row.Name, row.Url)
	}

	// Send message
	err = sendTelegramMessage(botToken, chatID, text)
	if err != nil {
		return err
	}
	return nil
}
