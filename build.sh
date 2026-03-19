#!/bin/bash

# Build script for skill-seed
# Supports building both English and Chinese versions

set -e

VERSION=${VERSION:-"dev"}
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS="-X main.version=$VERSION -X main.buildTime=$BUILD_TIME"

echo "Building skill-seed..."

# Build English version
echo "Building English version..."
go build -ldflags "$LDFLAGS" -o bin/skill-seed ./cmd/skill-seed
echo "✓ English version built: bin/skill-seed"

# Build Chinese version
echo "Building Chinese version..."
go build -ldflags "$LDFLAGS" -tags cn -o bin/skill-seed-cn ./cmd/skill-seed
echo "✓ Chinese version built: bin/skill-seed-cn"

echo ""
echo "Build completed successfully!"
echo "Binaries are available in the bin/ directory"
echo ""
echo "To install:"
echo "  cp bin/skill-seed /usr/local/bin/  # English version"
echo "  cp bin/skill-seed-cn /usr/local/bin/  # Chinese version"
