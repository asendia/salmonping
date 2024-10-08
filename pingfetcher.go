package main

import (
	"context"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/asendia/salmonping/db"
)

func fetchListings(ctx context.Context, queries *db.Queries) error {
	listings, err := queries.SelectListings(ctx, db.SelectListingsParams{
		EnablePing: []bool{true},
		Names:      []string{},
		Platforms:  []string{},
		Statuses:   []string{},
	})
	if err != nil {
		return err
	}
	for _, ol := range listings {
		if !ol.EnablePing {
			logJson(map[string]interface{}{
				"level":   "info",
				"message": "Skipping a listing because ping is disabled",
				"listing": ol.Name,
				"url":     ol.Url,
			})
			continue
		}
		logJson(map[string]interface{}{
			"level":   "info",
			"message": "Scraping a listing",
			"listing": ol.Name,
			"url":     ol.Url,
		})
		var status string
		var header http.Header
		var code int
		var body []byte
		var err error
		if ol.Platform == "gofood" {
			status, header, code, body, err = getGofoodStatus(ol.Url)
		} else if ol.Platform == "grabfood" {
			status, header, code, body, err = getGrabfoodStatus(ol.Url)
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
			"header":  header,
			"code":    code,
			"level":   "info",
			"listing": ol.Name,
			"status":  status,
		})

		if status == "unknown" {
			// Store the body in Cloud Storage
			// Create objectname with this format dump/2023/12/25/grabfood_21_05.html
			now := time.Now()
			objectName := path.Join("dump", strconv.Itoa(now.Year()), strconv.Itoa(int(now.Month())), strconv.Itoa(now.Day()), ol.Platform+"_"+strconv.Itoa(now.Hour())+"_"+strconv.Itoa(now.Minute())+".html")
			writeToCloudStorage("salmonping", objectName, body)
		}

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
	return nil
}
