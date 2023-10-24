package main

import (
	"fmt"
	"io"
	"net/http"
)

func routeWebhookGofood(w http.ResponseWriter, r *http.Request) {
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
			"level":   "error",
			"message": "Error validating signature",
		})
		// w.WriteHeader(http.StatusUnauthorized)
		// w.Header().Set("Content-Type", "application/json")
		// fmt.Fprintf(w, `{"message": "error"}`)
		// return
	}

	// Convert body to json
	log := map[string]interface{}{
		"body":    writeGofoodWebhookToCloudStorage(bodyBytes),
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
