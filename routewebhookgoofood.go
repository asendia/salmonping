package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	secretKey := os.Getenv("GOFOOD_NOTIFICATION_SECRET_KEY")
	clientSignature := r.Header.Get("X-Go-Signature")
	valid := checkHMAC(bodyBytes, []byte(clientSignature), []byte(secretKey))
	if secretKey != "" && !valid {
		logJson(map[string]interface{}{
			"error":   "Invalid signature",
			"level":   "error",
			"message": "Error validating signature",
		})
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "error"}`)
		return
	}

	// Convert body to json
	bodyJson := make(map[string]interface{})
	log := map[string]interface{}{
		"level":   "info",
		"message": "Gofood Webhook received",
		"headers": r.Header,
	}
	err = json.Unmarshal(bodyBytes, &bodyJson)
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

func checkHMAC(message, messageHMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageHMAC, expectedMAC)
}
