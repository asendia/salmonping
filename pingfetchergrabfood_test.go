package main

import (
	"testing"
)

func TestGetGrabfoodStatusIntegration(t *testing.T) {
	url := "https://food.grab.com/id/id/restaurant/salmon-fit-apartemen-menara-kebun-jeruk-delivery/6-C2XUWAX3PEU1JT"

	unknownCombo := 0
	for i := 0; i < 3; i++ {
		status, _, code, _, err := getGrabfoodStatus(url)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if status == "unknown" {
			t.Logf("Unknown status, code: %v", code)
			unknownCombo++
		}
	}
	if unknownCombo >= 3 {
		t.Errorf("Failed fetching grabfood 3x in a row")
	}
}
