package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestServeDocumentation(t *testing.T) {
	// Create temp spec file
	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "test-spec.json")

	spec := `{
		"openapi": "3.0.0",
		"info": {
			"title": "Test API",
			"version": "1.0.0"
		},
		"servers": [{"url": "https://api.test.com"}],
		"paths": {
			"/test": {
				"get": {
					"summary": "Test endpoint"
				}
			}
		}
	}`

	if err := os.WriteFile(specPath, []byte(spec), 0644); err != nil {
		t.Fatalf("failed to create test spec: %v", err)
	}

	// Create server
	server, err := NewDocServer(specPath)
	if err != nil {
		t.Fatalf("NewDocServer failed: %v", err)
	}

	// Test request
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !contains(body, "Test API") {
		t.Error("response should contain API title")
	}

	if !contains(body, "Test endpoint") {
		t.Error("response should contain endpoint summary")
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || len(s) >= len(substr) &&
		(s[0:len(substr)] == substr || contains(s[1:], substr)))
}

func TestFullIntegration(t *testing.T) {
	// Use the real spec file
	specPath := "./spec/api.json"
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		t.Skip("Skipping integration test: spec file not found")
	}

	server, err := NewDocServer(specPath)
	if err != nil {
		t.Fatalf("NewDocServer failed: %v", err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	body := w.Body.String()

	// Verify key content
	checks := []string{
		"Webmail Info API",
		"/public_api/v1/get_webmail_info",
		"POST",
		"bearerAuth",
		"200",
		"400",
		"401",
	}

	for _, check := range checks {
		if !contains(body, check) {
			t.Errorf("response missing expected content: %s", check)
		}
	}
}
