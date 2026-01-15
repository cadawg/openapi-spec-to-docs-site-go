package parser

import (
	"encoding/json"
	"fmt"
	"os"
)

// OpenAPISpec represents the parsed OpenAPI specification
type OpenAPISpec struct {
	OpenAPI    string              `json:"openapi"`
	Info       Info                `json:"info"`
	Servers    []Server            `json:"servers"`
	Paths      map[string]PathItem `json:"paths"`
	Components Components          `json:"components,omitempty"`
}

type Info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

type Server struct {
	URL         string               `json:"url"`
	Description string               `json:"description,omitempty"`
	Variables   map[string]ServerVar `json:"variables,omitempty"`
}

type ServerVar struct {
	Enum    []string `json:"enum,omitempty"`
	Default string   `json:"default"`
}

type PathItem struct {
	Get    *Operation `json:"get,omitempty"`
	Post   *Operation `json:"post,omitempty"`
	Put    *Operation `json:"put,omitempty"`
	Delete *Operation `json:"delete,omitempty"`
	Patch  *Operation `json:"patch,omitempty"`
}

type Operation struct {
	Summary     string                `json:"summary,omitempty"`
	Description string                `json:"description,omitempty"`
	Security    []map[string][]string `json:"security,omitempty"`
	RequestBody *RequestBody          `json:"requestBody,omitempty"`
	Responses   map[string]Response   `json:"responses"`
}

type RequestBody struct {
	Required bool                 `json:"required,omitempty"`
	Content  map[string]MediaType `json:"content"`
}

type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

type MediaType struct {
	Schema   *Schema            `json:"schema,omitempty"`
	Examples map[string]Example `json:"examples,omitempty"`
	Example  interface{}        `json:"example,omitempty"`
}

type Schema struct {
	Type       string              `json:"type,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
	Nullable   bool                `json:"nullable,omitempty"`
}

type Property struct {
	Type        string      `json:"type,omitempty"`
	Description string      `json:"description,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Nullable    bool        `json:"nullable,omitempty"`
}

type Example struct {
	Value       interface{} `json:"value"`
	Summary     string      `json:"summary,omitempty"`
	Description string      `json:"description,omitempty"`
}

type Components struct {
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`
}

type SecurityScheme struct {
	Type   string `json:"type"`
	Scheme string `json:"scheme,omitempty"`
	Name   string `json:"name,omitempty"`
	In     string `json:"in,omitempty"`
}

// ParseSpec parses an OpenAPI specification from JSON bytes
func ParseSpec(data []byte) (*OpenAPISpec, error) {
	var spec OpenAPISpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}
	return &spec, nil
}

// LoadSpecFromFile loads and parses an OpenAPI spec from a file
func LoadSpecFromFile(path string) (*OpenAPISpec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec file: %w", err)
	}
	return ParseSpec(data)
}
