package main

import (
	"testing"
)

func TestGetGrabfoodStatusIntegration(t *testing.T) {
	t.Skip("Temporarily skipping this since grab website blocks GCP IP")
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

func TestGetGrabfoodStatusSkeletonDetection(t *testing.T) {
	// GrabFood now serves skeleton-only HTML (client-side rendered)
	// The scraper should return "unknown" since actual data is loaded via JS
	skeletonBody := []byte(`<html><body><div class="Skeleton___2reg0 title___3X_DY"></div></body></html>`)
	if !isGrabfoodSkeleton(skeletonBody) {
		t.Error("Should detect skeleton-only page")
	}

	closedBody := []byte(`<html><body><div>Tutup</div></body></html>`)
	if isGrabfoodSkeleton(closedBody) {
		t.Error("Should not detect non-skeleton page as skeleton")
	}
}
