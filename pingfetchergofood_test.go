package main

import (
	"testing"
)

func TestGetGofoodStatusIntegration(t *testing.T) {
	url := "https://gofood.co.id/jakarta/restaurant/salmon-fit-apartemen-menara-kebon-jeruk-06f0dcc6-14f4-4092-810f-2bcc81214d23"

	status, _, code, body, err := getGofoodStatus(url)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	t.Logf("Status: %s, Code: %d, Body length: %d", status, code, len(body))
	// GoFood now returns a WAF challenge page (probe.js) which should be detected as "unknown"
	if status != "open" && status != "closed" && status != "unknown" {
		t.Errorf("Unexpected status: %s", status)
	}
}

func TestGetGofoodStatusWAFDetection(t *testing.T) {
	// Simulate a WAF challenge page response
	body := []byte(`<!DOCTYPE html><html>
<head>
  <meta charset="UTF-8">
  <script>
    var buid = "fffffffffffffffffff"
  </script>
  <script src="/C2WF946J0/probe.js?v=vc1jasc"></script>
</head>
<body></body>
</html>`)

	if containsProbeJS(body) != true {
		t.Error("Should detect probe.js in WAF challenge page")
	}

	normalBody := []byte(`<html><body><div>Salmon Fit</div><div>Buka</div></body></html>`)
	if containsProbeJS(normalBody) != false {
		t.Error("Should not detect probe.js in normal page")
	}
}
