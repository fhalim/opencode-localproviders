#!/usr/bin/env just --justfile

# Build the binary
build:
    go build -o opencode-localproviders .

# Run tests
test:
    go test -v ./...

# Update LM Studio config from live API
lmstudio:
    ./opencode-localproviders --base-url http://localhost:1234/ --provider lmstudio

# Update Ollama config from live API
ollama:
    ./opencode-localproviders --base-url http://localhost:11434/ --provider ollama

# Dry-run for LM Studio (preview changes)
dry-lmstudio:
    ./opencode-localproviders --base-url http://localhost:1234/ --provider lmstudio --dry-run

# Clean build artifacts
clean:
    rm -f opencode-localproviders

# Format code
fmt:
    go fmt ./...

# Lint code
lint:
    golangci-lint run ./... || true

# Install locally
install: build
    cp localopenapi-opencode ~/.local/bin/

# Run in dev mode (compile and execute)
dev args:
    go run . {{args}}
