package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/andybalholm/brotli"
)

func getGofoodStatus(url string) (string, http.Header, []byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nil, nil, err
	}
	emulateBrowser(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, nil, err
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
		return "", resp.Header, body, err
	}

	if bytes.Contains(body, []byte(">Buka")) {
		return "open", resp.Header, body, nil
	} else if bytes.Contains(body, []byte(">Tutup")) {
		return "closed", resp.Header, body, nil
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
		return "unknown", resp.Header, body, nil
	}
}
