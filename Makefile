.PHONY: build test clean install fmt lint help test-policy

# Default target
all: build

# Build the binary
build:
	go build -o yspec ./cmd/yspec

# Install the binary to GOPATH/bin
install:
	go install ./cmd/yspec

# Run tests
test:
	go test ./...

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -f yspec

# Run with example files
example: build
	./yspec examples/simple-before.yaml examples/simple-after.yaml

# Test policy integration with conftest
test-policy: build
	./examples/test-policy.sh

# Show help
help:
	@echo "Available targets:"
	@echo "  build       - Build the yspec binary"
	@echo "  install     - Install to GOPATH/bin"
	@echo "  test        - Run tests"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  clean       - Clean build artifacts"
	@echo "  example     - Run with example files"
	@echo "  test-policy - Test policy integration with conftest"
	@echo "  help        - Show this help"