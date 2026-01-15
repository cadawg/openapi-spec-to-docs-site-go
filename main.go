package main

import (
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cadawg/openapi-spec-to-docs-site-go/parser"
	"github.com/cadawg/openapi-spec-to-docs-site-go/templates"
)

//go:embed templates/index.html
var templateContent string

type DocServer struct {
	tmpl *template.Template
	data *templates.TemplateData
}

func NewDocServer(specPath string) (*DocServer, error) {
	// Load and parse spec
	spec, err := parser.LoadSpecFromFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load spec: %w", err)
	}

	// Build template data
	data := templates.BuildTemplateData(spec)

	// Parse template with custom functions
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}

	tmpl, err := template.New("index").Funcs(funcMap).Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &DocServer{
		tmpl: tmpl,
		data: data,
	}, nil
}

func (s *DocServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.Execute(w, s.data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func main() {
	specPath := flag.String("spec", "./spec/api.json", "Path to OpenAPI spec file")
	port := flag.Int("port", 8080, "Server port")
	flag.Parse()

	// Check if spec file exists
	if _, err := os.Stat(*specPath); os.IsNotExist(err) {
		log.Fatalf("Spec file not found: %s", *specPath)
	}

	// Create server
	server, err := NewDocServer(*specPath)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Documentation server running at http://localhost%s", addr)
	log.Printf("Using spec file: %s", *specPath)

	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
