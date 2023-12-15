package main

import (
	"testing"
)

func TestGetGofoodStatusIntegration(t *testing.T) {
	url := "https://gofood.co.id/jakarta/restaurant/salmon-fit-apartemen-menara-kebon-jeruk-06f0dcc6-14f4-4092-810f-2bcc81214d23"

	status, _, code, _, err := getGofoodStatus(url)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if status == "unknown" {
		t.Errorf("Unknown status, code: %v", code)
	}
}
