#!/bin/bash
set -e
echo "Building Go application..."
go build -tags netgo -ldflags '-s -w' -o app ./cmd/main.go
echo "Build completed successfully!"
