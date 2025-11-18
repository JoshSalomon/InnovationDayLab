#!/bin/bash

# Test runner script for Task Management Application
# Ensures tests are run from the correct directory

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_SRC_DIR="$SCRIPT_DIR/backend/src"

if [ ! -f "$BACKEND_SRC_DIR/go.mod" ]; then
    echo "Error: go.mod not found at $BACKEND_SRC_DIR/go.mod"
    echo "Please ensure you're running this from the project root directory."
    exit 1
fi

cd "$BACKEND_SRC_DIR" || exit 1

echo "Running tests from: $(pwd)"
echo ""

# Run tests with any additional arguments passed to this script
go test "$@" ./...
