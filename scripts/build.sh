#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR/.."

cd "$PROJECT_ROOT"

# Create build directory
mkdir -p build

# Build for Linux (amd64)
echo "Building for Linux..."
go build -o build/onTop-C2_linux ./src/main.go

echo "Build completed. Executable are in the build/ directory."
