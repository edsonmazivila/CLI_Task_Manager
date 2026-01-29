#!/bin/bash
set -e

echo "=== CI Verification Script ==="
echo "This script simulates CI pipeline checks locally"
echo ""

# Exit codes
EXIT_SUCCESS=0
EXIT_FAILURE=1

# Track overall status
OVERALL_STATUS=$EXIT_SUCCESS

# Function to run a check
run_check() {
    local name=$1
    local command=$2
    
    echo ""
    echo "▶ Running: $name"
    echo "─────────────────────────────────────────"
    
    if eval "$command"; then
        echo "✓ PASSED: $name"
        return 0
    else
        echo "✗ FAILED: $name"
        OVERALL_STATUS=$EXIT_FAILURE
        return 1
    fi
}

# Verify Go installation
run_check "Go Installation" "go version"

# Verify CGO
run_check "CGO Support" "go env | grep CGO_ENABLED=1 || (export CGO_ENABLED=1 && go env | grep CGO_ENABLED)"

# Download dependencies
run_check "Dependency Download" "go mod download"

# Verify dependencies
run_check "Dependency Verification" "go mod verify"

# Format check
run_check "Code Formatting" "test -z \$(gofmt -l .)"

# Build check
run_check "Build Verification" "CGO_ENABLED=1 go build -o /tmp/task-verify ./cmd/task"

# Integration tests (comprehensive test coverage)
run_check "Integration Tests" "CGO_ENABLED=1 go test -race -timeout 60s ./tests/integration/..."

# Test coverage from integration tests
run_check "Test Coverage" "CGO_ENABLED=1 go test -race -coverprofile=/tmp/coverage.out -coverpkg=./... ./tests/integration/... && go tool cover -func=/tmp/coverage.out | grep total | awk '{print \$3}' | grep -E '^([7-9][0-9]|100)'"

# Race detection
run_check "Race Detection" "CGO_ENABLED=1 go test -race ./tests/integration/..."

# Vet check
run_check "Go Vet" "go vet ./..."

# Print summary
echo ""
echo "═══════════════════════════════════════════"
if [ $OVERALL_STATUS -eq $EXIT_SUCCESS ]; then
    echo "✓ ALL CHECKS PASSED"
    echo "═══════════════════════════════════════════"
else
    echo "✗ SOME CHECKS FAILED"
    echo "═══════════════════════════════════════════"
fi

exit $OVERALL_STATUS
