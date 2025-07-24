.PHONY: build test clean lint fmt deps

# Variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt

# Default target
all: test build

# Build the project
build:
	$(GOBUILD) -v ./...

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f coverage.out coverage.html

# Format code
fmt:
	$(GOFMT) -s -w .

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Verify dependencies
verify:
	$(GOMOD) verify

# Update dependencies
update:
	$(GOGET) -u ./...
	$(GOMOD) tidy

# Run go vet
vet:
	$(GOCMD) vet ./...

# Check for potential issues
check: fmt vet lint test

# Install golangci-lint (if not already installed)
install-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2

# Test SSH connection (requires TEST_SSH_* environment variables)
test-ssh:
	@echo "Testing SSH connection..."
	@if [ -z "$(TEST_SSH_HOST)" ]; then echo "TEST_SSH_HOST not set"; exit 1; fi
	@if [ -z "$(TEST_SSH_USER)" ]; then echo "TEST_SSH_USER not set"; exit 1; fi
	@if [ -z "$(TEST_SSH_PASS)" ]; then echo "TEST_SSH_PASS not set"; exit 1; fi
	@echo "SSH Test: $(TEST_SSH_HOST) with user $(TEST_SSH_USER)"
