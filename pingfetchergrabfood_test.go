package main

import (
	"math/rand"
	"testing"
)

func TestGetGrabfoodStatusIntegration(t *testing.T) {
	url := "https://food.grab.com/id/id/restaurant/salmon-fit-apartemen-menara-kebun-jeruk-delivery/6-C2XUWAX3PEU1JT"
	rand.Shuffle(len(userAgents), func(i, j int) { userAgents[i], userAgents[j] = userAgents[j], userAgents[i] })

	for _, ua := range userAgents {
		status, header, _, err := getGrabfoodStatus(url, ua)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if status == "unknown" {
			t.Errorf("Unexpected status: %v %v", status, header)
		}
	}
}
