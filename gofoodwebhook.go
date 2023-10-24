package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

func verifyGofoodSignature(message []byte, signature string) error {
	secretKey := os.Getenv("GOFOOD_NOTIFICATION_SECRET_KEY")
	valid := checkHMAC(message, []byte(signature), []byte(secretKey))
	if secretKey != "" && !valid {
		return fmt.Errorf("invalid signature")
	}
	return nil
}

func checkHMAC(message, messageHMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	// encode expectedHMAC to hex
	expectedMACHex := []byte(hex.EncodeToString(expectedMAC))
	return hmac.Equal(messageHMAC, expectedMACHex)
}

func generateHMAC(message, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hex.EncodeToString(expectedMAC)
}

func writeGofoodWebhookToCloudStorage(bodyBytes []byte) interface{} {
	var body interface{} = nil
	bodyJson := make(map[string]interface{})
	filename := "gofoodwebhook/unknown"
	err := json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		body = string(bodyBytes)
	} else {
		body = bodyJson
		header, ok := bodyJson["header"].(map[string]interface{})
		if ok {
			eventName, ok := header["event_name"].(string)
			if ok {
				filename = fmt.Sprintf("gofoodwebhook/%s", eventName)
			}
		}
	}
	writeToCloudStorage("salmonping", filename, bodyBytes)
	return body
}
