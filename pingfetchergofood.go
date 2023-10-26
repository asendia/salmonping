package main

import (
	"bytes"
	"io"
	"net/http"
	"path"
	"strconv"
	"time"
)

func getGofoodStatus(url string) (string, http.Header, []byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp.Header, nil, err
	}
	if bytes.Contains(body, []byte("<span><p>Buka</p></span>")) {
		return "open", resp.Header, body, nil
	} else if bytes.Contains(body, []byte("<span><p>Tutup</p></span>")) {
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
