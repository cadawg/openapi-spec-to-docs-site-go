package templates

import (
	"encoding/json"
	"strings"

	"github.com/cadawg/openapi-spec-to-docs-site-go/parser"
)

// TemplateData is the data passed to the HTML template
type TemplateData struct {
	APITitle   string
	APIVersion string
	Servers    []ServerData
	Security   []SecurityData
	Endpoints  []EndpointData
}

type ServerData struct {
	URL         string
	Description string
}

type SecurityData struct {
	Name   string
	Type   string
	Scheme string
}

type EndpointData struct {
	Path        string
	Method      string
	Summary     string
	Description string
	Security    []string
	RequestBody *RequestBodyData
	Responses   []ResponseData
}

type RequestBodyData struct {
	Required bool
	Schema   string // JSON formatted
	Examples []ExampleData
}

type ResponseData struct {
	StatusCode  string
	Description string
	Schema      string // JSON formatted
	Examples    []ExampleData
}

type ExampleData struct {
	Name    string
	Summary string
	Value   string // JSON formatted
}

// BuildTemplateData converts OpenAPI spec to template data
func BuildTemplateData(spec *parser.OpenAPISpec) *TemplateData {
	data := &TemplateData{
		APITitle:   spec.Info.Title,
		APIVersion: spec.Info.Version,
		Servers:    buildServers(spec.Servers),
		Security:   buildSecurity(spec.Components.SecuritySchemes),
		Endpoints:  buildEndpoints(spec.Paths),
	}
	return data
}

func buildServers(servers []parser.Server) []ServerData {
	result := make([]ServerData, len(servers))
	for i, s := range servers {
		result[i] = ServerData{
			URL:         s.URL,
			Description: s.Description,
		}
	}
	return result
}

func buildSecurity(schemes map[string]parser.SecurityScheme) []SecurityData {
	result := []SecurityData{}
	for name, scheme := range schemes {
		result = append(result, SecurityData{
			Name:   name,
			Type:   scheme.Type,
			Scheme: scheme.Scheme,
		})
	}
	return result
}

func buildEndpoints(paths map[string]parser.PathItem) []EndpointData {
	endpoints := []EndpointData{}

	for path, item := range paths {
		if item.Get != nil {
			endpoints = append(endpoints, buildEndpoint(path, "GET", item.Get))
		}
		if item.Post != nil {
			endpoints = append(endpoints, buildEndpoint(path, "POST", item.Post))
		}
		if item.Put != nil {
			endpoints = append(endpoints, buildEndpoint(path, "PUT", item.Put))
		}
		if item.Delete != nil {
			endpoints = append(endpoints, buildEndpoint(path, "DELETE", item.Delete))
		}
		if item.Patch != nil {
			endpoints = append(endpoints, buildEndpoint(path, "PATCH", item.Patch))
		}
	}

	return endpoints
}

func buildEndpoint(path, method string, op *parser.Operation) EndpointData {
	endpoint := EndpointData{
		Path:        path,
		Method:      method,
		Summary:     op.Summary,
		Description: op.Description,
		Security:    extractSecurityNames(op.Security),
		Responses:   buildResponses(op.Responses),
	}

	if op.RequestBody != nil {
		endpoint.RequestBody = buildRequestBody(op.RequestBody)
	}

	return endpoint
}

func extractSecurityNames(security []map[string][]string) []string {
	names := []string{}
	for _, sec := range security {
		for name := range sec {
			names = append(names, name)
		}
	}
	return names
}

func buildRequestBody(rb *parser.RequestBody) *RequestBodyData {
	data := &RequestBodyData{
		Required: rb.Required,
		Examples: []ExampleData{},
	}

	for contentType, media := range rb.Content {
		if strings.Contains(contentType, "json") && media.Schema != nil {
			schemaJSON, _ := json.MarshalIndent(media.Schema, "", "  ")
			data.Schema = string(schemaJSON)
		}

		for name, ex := range media.Examples {
			valueJSON, _ := json.MarshalIndent(ex.Value, "", "  ")
			data.Examples = append(data.Examples, ExampleData{
				Name:    name,
				Summary: ex.Summary,
				Value:   string(valueJSON),
			})
		}

		// Handle single example
		if media.Example != nil && len(media.Examples) == 0 {
			valueJSON, _ := json.MarshalIndent(media.Example, "", "  ")
			data.Examples = append(data.Examples, ExampleData{
				Name:  "example",
				Value: string(valueJSON),
			})
		}
	}

	return data
}

func buildResponses(responses map[string]parser.Response) []ResponseData {
	result := []ResponseData{}

	for code, resp := range responses {
		data := ResponseData{
			StatusCode:  code,
			Description: resp.Description,
			Examples:    []ExampleData{},
		}

		for contentType, media := range resp.Content {
			if strings.Contains(contentType, "json") && media.Schema != nil {
				schemaJSON, _ := json.MarshalIndent(media.Schema, "", "  ")
				data.Schema = string(schemaJSON)
			}

			for name, ex := range media.Examples {
				valueJSON, _ := json.MarshalIndent(ex.Value, "", "  ")
				data.Examples = append(data.Examples, ExampleData{
					Name:    name,
					Summary: ex.Summary,
					Value:   string(valueJSON),
				})
			}

			// Handle single example
			if media.Example != nil && len(media.Examples) == 0 {
				valueJSON, _ := json.MarshalIndent(media.Example, "", "  ")
				data.Examples = append(data.Examples, ExampleData{
					Name:  "example",
					Value: string(valueJSON),
				})
			}
		}

		result = append(result, data)
	}

	return result
}
