#!/bin/bash

# Build script for grow-check
# Supports building both English and Chinese versions

set -e

VERSION=${VERSION:-"dev"}
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS="-X main.version=$VERSION -X main.buildTime=$BUILD_TIME"

echo "Building grow-check..."

# Build English version
echo "Building English version..."
go build -ldflags "$LDFLAGS" -o bin/grow-check ./cmd/grow-check
echo "✓ English version built: bin/grow-check"

# Build Chinese version
echo "Building Chinese version..."
go build -ldflags "$LDFLAGS" -tags cn -o bin/grow-check-cn ./cmd/grow-check
echo "✓ Chinese version built: bin/grow-check-cn"

echo ""
echo "Build completed successfully!"
echo "Binaries are available in the bin/ directory"
echo ""
echo "To install:"
echo "  cp bin/grow-check /usr/local/bin/  # English version"
echo "  cp bin/grow-check-cn /usr/local/bin/  # Chinese version"
