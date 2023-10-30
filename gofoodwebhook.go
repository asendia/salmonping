package main

import (
	"fmt"
)

func verifyGofoodSignature(msg []byte, signature string, secretKey string) error {
	if secretKey == "" {
		return nil
	}
	key := []byte(secretKey)
	valid, err := hmacVerifyMessage(msg, key, signature)
	if err != nil {
		return err
	}
	if !valid {
		return fmt.Errorf("invalid signature")
	}
	return nil
}

// Generated using ChatGPT-4
type GofoodWebhookPayload struct {
	Header GofoodWebhookHeader `json:"header"`
	Body   GofoodWebhookBody   `json:"body"`
}

type GofoodWebhookHeader struct {
	EventName string `json:"event_name"`
	Version   int    `json:"version"`
	Timestamp string `json:"timestamp"`
	EventID   string `json:"event_id"`
}

type GofoodWebhookBody struct {
	ServiceType string         `json:"service_type"`
	Order       GofoodOrder    `json:"order"`
	Driver      GofoodDriver   `json:"driver"`
	Outlet      GofoodOutlet   `json:"outlet"`
	Customer    GofoodCustomer `json:"customer"`
}

type GofoodOrder struct {
	OrderItems        []GofoodOrderItem `json:"order_items"`
	Currency          string            `json:"currency"`
	OrderNumber       string            `json:"order_number"`
	CreatedAt         string            `json:"created_at"`
	TakeawayCharges   float32           `json:"takeaway_charges"`
	OrderTotal        float32           `json:"order_total"`
	Pin               string            `json:"pin"`
	CutleryRequested  *bool             `json:"cutlery_requested"` // Use pointer for nullable bool
	Status            string            `json:"status"`
	AppliedPromotions []string          `json:"applied_promotions"` // Assuming promotions are strings; adjust if needed
}

type GofoodOrderItem struct {
	Quantity   int             `json:"quantity"`
	Name       string          `json:"name"`
	ExternalID string          `json:"external_id"`
	ID         string          `json:"id"`
	Notes      string          `json:"notes"`
	Price      float32         `json:"price"`
	Variants   []GofoodVariant `json:"variants"`
}

type GofoodVariant struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ExternalID string `json:"external_id"`
}

type GofoodDriver struct {
	Name string `json:"name"`
}

type GofoodOutlet struct {
	ExternalOutletID *string `json:"external_outlet_id"` // Use pointer for nullable string
	ID               string  `json:"id"`
}

type GofoodCustomer struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}
