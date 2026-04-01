#!/bin/bash
# Oculus pre-commit hook - Go best practices

set -e

echo "Running pre-commit checks..."

# Check formatting
UNFORMATTED=$(gofmt -l . 2>/dev/null | grep -v vendor | grep -v old-src | head -5)
if [ -n "$UNFORMATTED" ]; then
    echo "❌ Unformatted Go files:"
    echo "$UNFORMATTED"
    echo "Run: gofmt -w ."
    exit 1
fi
echo "✓ gofmt"

# Run vet
if ! go vet ./... 2>/dev/null; then
    echo "❌ go vet failed"
    exit 1
fi
echo "✓ go vet"

# Run tests
if ! go test ./... -short 2>/dev/null; then
    echo "❌ Tests failed"
    exit 1
fi
echo "✓ go test"

# Check build
if ! go build ./... 2>/dev/null; then
    echo "❌ Build failed"
    exit 1
fi
echo "✓ go build"

echo "All checks passed!"
