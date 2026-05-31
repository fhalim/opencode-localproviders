#!/usr/bin/env just --justfile

BIN_NAME := "opencode-localproviders"
BIN_EXT := if os_family() == "windows" { ".exe" } else { "" }
BIN_PATH := BIN_NAME + BIN_EXT
BIN_EXEC := "." / BIN_PATH

# Build the binary
build:
    go build -o {{BIN_PATH}} .

# Run tests
test:
    go test -v ./...

# Update LM Studio config from live API
lmstudio: build
    {{BIN_EXEC}} --base-url http://localhost:1234/ --provider lmstudio

# Update Ollama config from live API
ollama: build
    {{BIN_EXEC}} --base-url http://localhost:11434/ --provider ollama

# Dry-run for LM Studio (preview changes)
dry-lmstudio: build
    {{BIN_EXEC}} --base-url http://localhost:1234/ --provider lmstudio --dry-run

# Clean build artifacts
clean:
    go clean
    @if [ -f {{BIN_PATH}} ]; then rm {{BIN_PATH}}; fi

# Format code
fmt:
    go fmt ./...

# Lint code
lint:
    golangci-lint run ./... || true

# Install locally
install: build
    mkdir -p ~/.local/bin
    cp {{BIN_PATH}} ~/.local/bin/

# Run in dev mode (compile and execute)
dev args:
    go run . {{args}}
