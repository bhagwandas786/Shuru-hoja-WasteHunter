#!/bin/bash

set -e

echo "Building shuru hoja..."
echo "======================"

# Check for Go
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed"
    exit 1
fi

# Create build directory
mkdir -p build

# Build
echo "Building binary..."
go build -ldflags "-s -w" -o build/shuru-hoja ./cmd/shuru-hoja

# Check if build succeeded
if [ $? -eq 0 ]; then
    echo "Build successful!"
    echo "Binary created at: build/shuru-hoja"
    
    # Test the binary
    echo "Testing binary..."
    ./build/shuru-hoja --version 2>/dev/null && echo "✓ Binary works" || echo "⚠ Binary test inconclusive"
    
    # Show file info
    echo ""
    echo "Binary information:"
    ls -lh build/shuru-hoja
    file build/shuru-hoja
else
    echo "Build failed!"
    exit 1
fi
