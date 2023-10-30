package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestVerifyGofoodSignature(t *testing.T) {
	godotenv.Load()
	msg, err := loadWebhookExamples("gofood.order.completed")
	if err != nil {
		t.Error(err)
	}
	signature := hmacSignMessage(msg, []byte(os.Getenv("GOFOOD_NOTIFICATION_SECRET_KEY")))
	err = verifyGofoodSignature(msg, signature, os.Getenv("GOFOOD_NOTIFICATION_SECRET_KEY"))
	if err != nil {
		t.Error(err)
	}
}

func TestParseBodyAsWebhookPayload(t *testing.T) {
	events := []string{
		"gofood.order.cancelled",
		"gofood.order.completed",
		"gofood.order.driver_arrived",
		"gofood.order.driver_otw_pickup",
		"gofood.order.merchant_accepted",
		"gofood.order.placed",
	}
	for _, event := range events {
		msg, err := loadWebhookExamples(event)
		if err != nil {
			t.Error(err)
		}
		var p GofoodWebhookPayload
		err = json.Unmarshal(msg, &p)
		if err != nil {
			t.Error(err)
		}
		if p.Header.EventName != event {
			t.Errorf("Expected event name %s, got %s", event, p.Header.EventName)
		}
		if p.Body.Order.OrderItems[0].Variants[0].ExternalID != "variant-external-id-test-1" {
			t.Errorf("Expected variant ID variant-external-id-test-1, got %s", p.Body.Order.OrderItems[0].Variants[0].ID)
		}
	}
}

func loadWebhookExamples(eventName string) ([]byte, error) {
	file, err := os.Open(fmt.Sprintf("./gofoodwebhookexamples/%s.json", eventName))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	msg, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
