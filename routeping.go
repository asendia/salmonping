package main

import (
	"fmt"
	"net/http"
	"os"

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

	// Prepare DB connection
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

	// Scrape all listings
	for _, ol := range listings {
		logJson(map[string]interface{}{
			"level":      "info",
			"message":    "Scraping a restaurant",
			"restaurant": ol.Name,
			"url":        ol.Url,
		})
		var status string
		var err error
		if ol.Platform == "gofood" {
			status, err = getGofoodStatus(ol.Url)
		} else if ol.Platform == "grabfood" {
			status, err = getGrabfoodStatus(ol.Url)
		} else {
			logJson(map[string]interface{}{
				"level":      "error",
				"message":    "Unsupported restaurant url",
				"restaurant": ol.Name,
				"url":        ol.Url,
			})
			continue
		}
		if err != nil {
			logJson(map[string]interface{}{
				"level":      "error",
				"message":    "Error scraping a restaurant",
				"restaurant": ol.Name,
				"error":      err.Error(),
			})
			continue
		}

		logJson(map[string]interface{}{
			"level":      "info",
			"restaurant": ol.Name,
			"status":     status,
		})

		// Log to database
		_, err = queries.InsertPing(ctx, db.InsertPingParams{
			OnlineListingID: ol.ID,
			Status:          status,
		})
		if err != nil {
			logJson(map[string]interface{}{
				"level":      "error",
				"message":    "Error inserting a ping",
				"restaurant": ol.Name,
				"error":      err.Error(),
			})
			continue
		}
		// currentTime := time.Now().UTC().Add(time.Hour * 7)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "ok"}`)
}
