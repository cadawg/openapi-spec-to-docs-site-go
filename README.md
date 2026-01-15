# OpenAPI Documentation Server

Self-hosted API documentation generator for OpenAPI 3.0 specifications.

## Features

- Clean, modern UI with Tailwind CSS
- Server-side rendering with Go templates
- Request/response examples with syntax highlighting
- Single-page navigation with smooth scrolling
- Copy-to-clipboard for code examples
- Single binary deployment

## Installation

### Prerequisites

- Go 1.21 or higher

### Build

```bash
go build -o openapi-docs
```

## Usage

### Start Server

```bash
# Use default spec file (./spec/api.json) on port 8080
./openapi-docs

# Custom spec file and port
./openapi-docs --spec path/to/spec.json --port 3000
```

### Command-Line Options

- `--spec` or `-s`: Path to OpenAPI spec file (default: `./spec/api.json`)
- `--port` or `-p`: Server port (default: `8080`)

### Examples

```bash
# Default configuration
./openapi-docs

# Custom spec file
./openapi-docs --spec myapi.json

# Custom port
./openapi-docs --port 3000

# Both options
./openapi-docs --spec myapi.json --port 3000
```

## Development

### Run Tests

```bash
go test ./...
```

### Run Locally

```bash
go run main.go
```

Then open http://localhost:8080 in your browser.

## Project Structure

```
openapi-docs/
├── main.go              # HTTP server and entry point
├── parser/              # OpenAPI spec parsing
│   ├── openapi.go
│   └── openapi_test.go
├── templates/           # Template data model
│   ├── model.go
│   ├── model_test.go
│   └── index.html      # HTML template (embedded)
└── spec/               # OpenAPI specifications
    └── api.json        # Default spec file
```

## License

MIT
