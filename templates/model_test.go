package templates

import (
	"testing"

	"github.com/cadawg/openapi-spec-to-docs-site-go/parser"
)

func TestBuildTemplateData(t *testing.T) {
	spec := &parser.OpenAPISpec{
		Info: parser.Info{
			Title:   "My API",
			Version: "1.0.0",
		},
		Servers: []parser.Server{
			{URL: "https://api.example.com"},
		},
		Paths: map[string]parser.PathItem{
			"/test": {
				Post: &parser.Operation{
					Summary: "Test endpoint",
				},
			},
		},
	}

	data := BuildTemplateData(spec)

	if data.APITitle != "My API" {
		t.Errorf("expected APITitle 'My API', got '%s'", data.APITitle)
	}

	if len(data.Endpoints) != 1 {
		t.Fatalf("expected 1 endpoint, got %d", len(data.Endpoints))
	}

	if data.Endpoints[0].Path != "/test" {
		t.Errorf("expected path '/test', got '%s'", data.Endpoints[0].Path)
	}

	if data.Endpoints[0].Method != "POST" {
		t.Errorf("expected method 'POST', got '%s'", data.Endpoints[0].Method)
	}
}
