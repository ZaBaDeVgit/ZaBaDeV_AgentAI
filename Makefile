.PHONY: all build test lint clean install dev run

# Variables
BINARY_NAME=zabadev
GO_VERSION=1.21
MAIN_PACKAGE=./cmd/zabadev

# Default target
all: clean lint test build

# Build the binary
build:
	go build -o bin/$(BINARY_NAME) $(MAIN_PACKAGE)

# Run tests
test:
	go test ./... -v

# Run tests with coverage
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Lint code
lint:
	gofmt -s -w .
	go vet ./...
	golangci-lint run ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install dependencies
install:
	go mod download
	go mod tidy

# Development setup
dev: install
	go install github.com/air-verse/air@latest

# Run in development mode with hot reload
run-dev:
	air

# Run the application
run:
	go run $(MAIN_PACKAGE)

# Cross-platform build
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)

# Docker build
docker-build:
	docker build -t $(BINARY_NAME) .

# Docker run
docker-run:
	docker run --rm $(BINARY_NAME)

# CI pipeline
ci: install lint test build