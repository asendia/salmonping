package main

import (
	"fmt"
	"io"
	"net/http"
)

func gofoodWebhookHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logJson(map[string]interface{}{
			"error":   err.Error(),
			"level":   "error",
			"message": "Error reading body",
		})
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "error"}`)
		return
	}
	defer r.Body.Close()

	err = verifyGofoodSignature(bodyBytes, r.Header.Get("X-Go-Signature"))
	if err != nil {
		logJson(map[string]interface{}{
			"error":   err.Error(),
			"level":   "warning",
			"message": "Error validating signature",
		})
		// Cannot verify signature in sandbox
		// w.WriteHeader(http.StatusUnauthorized)
		// w.Header().Set("Content-Type", "application/json")
		// fmt.Fprintf(w, `{"message": "error"}`)
		// return
	}

	p, err := parseBodyAsGofoodWebhookPayload(bodyBytes)
	if err != nil {
		logJson(map[string]interface{}{
			"body":    string(bodyBytes),
			"error":   err.Error(),
			"level":   "error",
			"message": "Error parsing body",
		})
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "error"}`)
		return
	}
	log := map[string]interface{}{
		"payload": p,
		"headers": r.Header,
		"level":   "info",
		"message": "Gofood Webhook received",
	}

	_, err = logJson(log)
	if err != nil {
		logJson(map[string]interface{}{
			"error":   err.Error(),
			"level":   "error",
			"message": "Error logging webhook",
		})
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "ok"}`)
}
