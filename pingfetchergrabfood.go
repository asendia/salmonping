package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"

	"github.com/andybalholm/brotli"
)

func getGrabfoodStatus(url string) (string, http.Header, int, []byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nil, 0, nil, err
	}
	emulateBrowser(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", resp.Header, 0, nil, err
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
		return "", resp.Header, resp.StatusCode, body, err
	}
	if resp.StatusCode >= 400 {
		return "unknown", resp.Header, resp.StatusCode, body, nil
	}

	if bytes.Contains(body, []byte("Tutup</div></div>")) {
		return "closed", resp.Header, resp.StatusCode, body, nil
	} else if bytes.Contains(body, []byte("Jam Buka</label>")) {
		return "open", resp.Header, resp.StatusCode, body, nil
	}
	return "unknown", resp.Header, resp.StatusCode, body, nil
}
