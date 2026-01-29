.PHONY: build clean install test run help

# Binary name
BINARY_NAME=task

# Build directory
BUILD_DIR=bin

# Go parameters
GOCMD=go
GOBUILD=CGO_ENABLED=1 $(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/task
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Install the binary to system path
install: build
	@echo "Installing $(BINARY_NAME)..."
	@sudo mv $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	CGO_ENABLED=1 $(GOTEST) -v -race -timeout 60s ./tests/integration/...

# Run all tests (unit + integration)
test-all: test test-integration
	@echo "All tests complete"

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	CGO_ENABLED=1 $(GOTEST) -bench=. -benchmem -run=^$$ ./tests/integration/...

# Run CI verification locally
ci-verify:
	@echo "Running CI verification..."
	@./scripts/ci-verify.sh

# Run full test suite
test-suite:
	@echo "Running full test suite..."
	@./scripts/run-tests.sh

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies downloaded"

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...
	@echo "Code formatted"

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	golangci-lint run
	@echo "Linting complete"

# Check for vulnerabilities (requires govulncheck)
vuln:
	@echo "Checking for vulnerabilities..."
	govulncheck ./...
	@echo "Vulnerability check complete"

# Run the application
run: build
	@$(BUILD_DIR)/$(BINARY_NAME)

# Build for multiple platforms
build-all: build-linux build-darwin build-windows
	@echo "All builds complete"

build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux ./cmd/task

build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-macos ./cmd/task

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=1 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME).exe ./cmd/task

# Display help
help:
	@echo "Task Manager - Makefile Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build         Build the application"
	@echo "  clean         Remove build artifacts"
	@echo "  install       Install the binary to /usr/local/bin"
	@echo "  test          Run tests"
	@echo "  test-coverage Run tests with coverage"
	@echo "  test-integration Run integration tests"
	@echo "  test-all      Run all tests (unit + integration)"
	@echo "  bench         Run benchmarks"
	@echo "  ci-verify     Run CI verification locally"
	@echo "  test-suite    Run full test suite"
	@echo "  deps          Download and tidy dependencies"
	@echo "  fmt           Format code"
	@echo "  lint          Run linter (requires golangci-lint)"
	@echo "  vuln          Check for vulnerabilities (requires govulncheck)"
	@echo "  run           Build and run the application"
	@echo "  build-all     Build for all platforms"
	@echo "  build-linux   Build for Linux"
	@echo "  build-darwin  Build for macOS"
	@echo "  build-windows Build for Windows"
	@echo "  help          Display this help message"
