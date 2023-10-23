package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/asendia/salmonping/db"
)

func fetchListings(ctx context.Context, queries *db.Queries) error {
	listings, err := queries.SelectListings(ctx)
	if err != nil {
		return err
	}
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
	return nil
}

func getGofoodStatus(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if bytes.Contains(body, []byte("<span><p>Buka</p></span>")) {
		return "open", nil
	} else if bytes.Contains(body, []byte("<span><p>Tutup</p></span>")) {
		return "closed", nil
	} else {
		logJson(map[string]interface{}{
			"level":  "error",
			"error":  "Unknown status",
			"header": resp.Header,
			"url":    url,
		})
		// Store the body in Cloud Storage
		// Create objectname with this format dump/2023/12/25/grabfood_21_05.html
		now := time.Now()
		objectName := path.Join("dump", strconv.Itoa(now.Year()), strconv.Itoa(int(now.Month())), strconv.Itoa(now.Day()), "gofood_"+strconv.Itoa(now.Hour())+"_"+strconv.Itoa(now.Minute())+".html")
		writeToCloudStorage("salmonping", objectName, body)
		return "unknown", nil
	}
}

func getGrabfoodStatus(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:89.0) Gecko/20100101 Firefox/89.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "id,en-US;q=0.7,en;q=0.3")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Cache-Control", "max-age=0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var body []byte
	encoding := resp.Header.Get("Content-Encoding")
	switch encoding {
	case "br":
		bodyReader := brotli.NewReader(resp.Body)
		body, err = io.ReadAll(bodyReader)
	case "gzip":
		bodyReader, readerErr := gzip.NewReader(resp.Body)
		if readerErr == nil {
			body, err = io.ReadAll(bodyReader)
		} else {
			err = readerErr
		}
	default:
		body, err = io.ReadAll(resp.Body)
	}
	if err != nil {
		return "", err
	}

	if bytes.Contains(body, []byte("Tutup</div></div>")) {
		return "closed", nil
	} else if bytes.Contains(body, []byte("Jam Buka</label>")) {
		return "open", nil
	} else {
		// Store the body in Cloud Storage
		// Create objectname with this format dump/2023/12/25/grabfood_21_05.html
		now := time.Now()
		objectName := path.Join("dump", strconv.Itoa(now.Year()), strconv.Itoa(int(now.Month())), strconv.Itoa(now.Day()), "grabfood_"+strconv.Itoa(now.Hour())+"_"+strconv.Itoa(now.Minute())+".html")
		writeToCloudStorage("salmonping", objectName, body)
		return "unknown", nil
	}
}
