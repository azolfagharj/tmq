# Makefile for tmq - TOML Query Tool

.PHONY: all build test test-coverage lint clean help

# Default target
all: test lint build

# Build the binary
build:
	go build -ldflags="-s -w" -o bin/tmq ./cmd/tmq

# Run tests
test:
	go test ./... -race -v

# Run tests with coverage
test-coverage:
	go test ./... -race -coverprofile=coverage.out -covermode=atomic

# Run linter
lint:
	@echo "Running Go linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "golangci-lint not found. Running basic checks..."; \
		go vet ./...; \
		if [ "$$(go fmt ./... | wc -l)" -gt 0 ]; then \
			echo "Code is not formatted. Run: go fmt ./..."; \
			exit 1; \
		fi; \
		echo "Basic lint checks passed (golangci-lint not available)"; \
	fi

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out coverage.html

# Install development dependencies
install-deps:
	go mod download
	# Install golangci-lint if not present
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2; \
	fi

# Generate coverage report
coverage-html: test-coverage
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Show help
help:
	@echo "Available targets:"
	@echo "  all           - Run tests, lint, and build"
	@echo "  build         - Build the binary"
	@echo "  test          - Run tests with race detection"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  lint          - Run golangci-lint"
	@echo "  clean         - Clean build artifacts"
	@echo "  install-deps  - Install development dependencies"
	@echo "  coverage-html - Generate HTML coverage report"
	@echo "  help          - Show this help message"

#