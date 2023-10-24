package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func routeWebhookGofood(w http.ResponseWriter, r *http.Request) {
	// Print all request headers and body
	bodyBytes := make([]byte, r.ContentLength)
	r.Body.Read(bodyBytes)
	// Convert body to json
	bodyJson := make(map[string]interface{})
	log := map[string]interface{}{
		"level":   "info",
		"message": "Gofood Webhook received",
		"headers": r.Header,
	}
	err := json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		log["body"] = string(bodyBytes)
	} else {
		log["body"] = bodyJson
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
