#!/bin/bash
set -e

if ! command -v go &> /dev/null; then
    echo "Go not found. Please install Go 1.22.2 or newer."
    exit 1
fi

REQUIRED_VERSION="1.22.2"
INSTALLED_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
if [[ "$(printf '%s\n' "$REQUIRED_VERSION" "$INSTALLED_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]]; then
    echo "Go version found ($INSTALLED_VERSION) is older than $REQUIRED_VERSION. Please update Go."
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR/.."

cd "$PROJECT_ROOT"
mkdir -p build

echo "Installing Go dependencies..."
go mod tidy

echo "Building for Linux..."
go build -o -ldflags="-s -w" -o build/TheDarkMark_linux ./src/main.go

echo "Build completed. Executables are in the build/ directory."
