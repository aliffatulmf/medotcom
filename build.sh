#!/bin/bash

set -e

BINARY_NAME="you"
OUTPUT_DIR="build"
MIN_GO_VERSION="1.18.0"

# Check if Go is installed
if ! command -v go &> /dev/null
then
    echo "Go could not be found, please install it."
    exit
fi

# Check if the correct Go version is installed
CURRENT_GO_VERSION=$(go version | awk '{print substr($3,3)}')
if [ "$(printf '%s\n' "$MIN_GO_VERSION" "$CURRENT_GO_VERSION" | sort -V | head -n1)" = "$CURRENT_GO_VERSION" ]; then
  echo "Go version must be greater than $MIN_GO_VERSION"
  exit
fi

# Build the application
echo "Building..."
go build -o $OUTPUT_DIR/$BINARY_NAME main.go

echo "Build complete"