#!/bin/bash
set -e

echo "=== Running Task Manager Test Suite ==="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print section headers
print_section() {
    echo ""
    echo -e "${GREEN}==== $1 ====${NC}"
    echo ""
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
print_section "Checking Prerequisites"

if ! command_exists go; then
    echo -e "${RED}Error: Go is not installed${NC}"
    exit 1
fi

if ! command_exists gcc; then
    echo -e "${RED}Error: GCC is not installed${NC}"
    echo "Install with: sudo apt-get install build-essential"
    exit 1
fi

echo "✓ Go version: $(go version)"
echo "✓ GCC version: $(gcc --version | head -n1)"

# Set up environment
export CGO_ENABLED=1

# Run go mod tidy
print_section "Installing Dependencies"
go mod download
go mod tidy

# Run unit tests
print_section "Running Unit Tests"
go test -v -race -coverprofile=coverage.out ./internal/...

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Unit tests passed${NC}"
else
    echo -e "${RED}✗ Unit tests failed${NC}"
    exit 1
fi

# Run integration tests
print_section "Running Integration Tests"
go test -v -race -timeout 60s ./tests/integration/...

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Integration tests passed${NC}"
else
    echo -e "${RED}✗ Integration tests failed${NC}"
    exit 1
fi

# Generate coverage report
print_section "Coverage Report"
go tool cover -func=coverage.out | tail -n 1

# Run benchmarks
print_section "Running Benchmarks"
go test -bench=. -benchmem -run=^$ ./tests/integration/... | tee benchmark.txt

# Build the application
print_section "Building Application"
go build -o bin/task ./cmd/task

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Build successful${NC}"
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi

# Verify binary
print_section "Verifying Binary"
./bin/task --help > /dev/null 2>&1

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Binary verification passed${NC}"
else
    echo -e "${RED}✗ Binary verification failed${NC}"
    exit 1
fi

# Run linter if available
if command_exists golangci-lint; then
    print_section "Running Linter"
    golangci-lint run --timeout=5m
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Linting passed${NC}"
    else
        echo -e "${YELLOW}⚠ Linting found issues${NC}"
    fi
else
    echo -e "${YELLOW}⚠ golangci-lint not installed, skipping linting${NC}"
fi

# Run security checks if available
if command_exists gosec; then
    print_section "Running Security Scan"
    gosec -quiet ./...
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Security scan passed${NC}"
    else
        echo -e "${YELLOW}⚠ Security scan found issues${NC}"
    fi
else
    echo -e "${YELLOW}⚠ gosec not installed, skipping security scan${NC}"
fi

# Summary
print_section "Test Suite Complete"
echo -e "${GREEN}All tests passed successfully!${NC}"
echo ""
echo "Next steps:"
echo "  - Review coverage report: go tool cover -html=coverage.out"
echo "  - Check benchmarks: cat benchmark.txt"
echo "  - Install binary: make install"
