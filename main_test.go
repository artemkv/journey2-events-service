package main

import (
	"fmt"
	"net/http"
	"testing"
)

var port string

func init() {
	LoadDotEnv()
	port = GetPort(":8600")
}

func TestHealthCheckIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	url := fmt.Sprintf("http://127.0.0.1%s/health", port)
	statusCode := request(t, url)

	if statusCode != 200 {
		t.Errorf("Expected 200, actual: %d", statusCode)
	}
}

func TestLivenessCheckIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	url := fmt.Sprintf("http://127.0.0.1%s/liveness", port)
	statusCode := request(t, url)

	if statusCode != 200 {
		t.Errorf("Expected 200, actual: %d", statusCode)
	}
}

func TestReadinessCheckIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	url := fmt.Sprintf("http://127.0.0.1%s/readiness", port)
	statusCode := request(t, url)

	if statusCode != 200 {
		t.Errorf("Expected 200, actual: %d", statusCode)
	}
}

func request(t *testing.T, url string) int {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	return resp.StatusCode
}
