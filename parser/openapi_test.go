package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseSpec(t *testing.T) {
	spec := `{
		"openapi": "3.0.0",
		"info": {
			"title": "Test API",
			"version": "1.0.0"
		},
		"paths": {}
	}`

	result, err := ParseSpec([]byte(spec))
	if err != nil {
		t.Fatalf("ParseSpec failed: %v", err)
	}

	if result.Info.Title != "Test API" {
		t.Errorf("expected title 'Test API', got '%s'", result.Info.Title)
	}

	if result.Info.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", result.Info.Version)
	}
}

func TestLoadSpecFromFile(t *testing.T) {
	// Create temp file
	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "spec.json")

	specContent := `{
		"openapi": "3.0.0",
		"info": {
			"title": "File API",
			"version": "2.0.0"
		},
		"paths": {}
	}`

	if err := os.WriteFile(specPath, []byte(specContent), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	result, err := LoadSpecFromFile(specPath)
	if err != nil {
		t.Fatalf("LoadSpecFromFile failed: %v", err)
	}

	if result.Info.Title != "File API" {
		t.Errorf("expected title 'File API', got '%s'", result.Info.Title)
	}
}

func TestLoadSpecFromFile_NotFound(t *testing.T) {
	_, err := LoadSpecFromFile("/nonexistent/spec.json")
	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
}
