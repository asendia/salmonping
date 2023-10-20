package main

import "testing"

func TestGetGrabfoodStatusIntegration(t *testing.T) {
	url := "https://food.grab.com/id/id/restaurant/salmon-fit-apartemen-menara-kebun-jeruk-delivery/6-C2XUWAX3PEU1JT"
	status, err := getGrabfoodStatus(url)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if status == "unknown" {
		t.Errorf("Unexpected status: %v", status)
	}
}
