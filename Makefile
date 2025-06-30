# Makefile for envsubst

# Variables
BINARY_NAME=envsubst
MODULE_NAME=github.com/allex/envsubst
CMD_DIR=./cmd/envsubst
BUILD_DIR=./bin
GO_FILES=$(shell find . -name '*.go' -not -path './vendor/*')

# Build info
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

# Default target
.DEFAULT_GOAL := help

## build: Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

## install: Install the binary to GOPATH/bin
.PHONY: install
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) $(CMD_DIR)

## test: Run all tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v -race ./...

## test-coverage: Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report saved to coverage.html"

## benchmark: Run benchmarks
.PHONY: benchmark
benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

## clean: Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

## fmt: Format Go code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

## vet: Run go vet
.PHONY: vet
vet:
	@echo "Running go vet..."
	go vet ./...

## mod-tidy: Tidy module dependencies
.PHONY: mod-tidy
mod-tidy:
	@echo "Tidying module dependencies..."
	go mod tidy

## mod-verify: Verify module dependencies
.PHONY: mod-verify
mod-verify:
	@echo "Verifying module dependencies..."
	go mod verify

## lint: Run golangci-lint (if available)
.PHONY: lint
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "Running golangci-lint..."; \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, skipping..."; \
	fi

## check: Run all checks (fmt, vet, test)
.PHONY: check
check: fmt vet test

## ci: Run CI pipeline locally
.PHONY: ci
ci: mod-tidy fmt vet lint test

## run-example: Run the example
.PHONY: run-example
run-example: build
	@echo "Running example..."
	cd _example && ../$(BUILD_DIR)/$(BINARY_NAME) -i config.yaml

## git-release: Create a git release with patch version bump
.PHONY: git-release
git-release:
	@echo "Creating git release..."
	git release -t v -r patch

## help: Show this help message
.PHONY: help
help: Makefile
	@echo "Available targets:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
