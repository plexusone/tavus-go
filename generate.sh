#!/bin/bash
# generate.sh - Generate Go API client from OpenAPI specification using ogen
#
# Prerequisites:
#   go install github.com/ogen-go/ogen/cmd/ogen@latest
#
# Usage:
#   ./generate.sh
#
# This script:
#   1. Runs ogen to generate Go code from the OpenAPI spec
#   2. Applies post-processing fixes for known ogen bugs
#   3. Runs go mod tidy to update dependencies
#   4. Verifies the build compiles

set -e

# Check if ogen is installed
if ! command -v ogen &> /dev/null; then
    echo "Error: ogen is not installed."
    echo "Install with: go install github.com/ogen-go/ogen/cmd/ogen@latest"
    exit 1
fi

# Check if OpenAPI spec exists
if [ ! -f "openapi/openapi.yaml" ]; then
    echo "Error: openapi/openapi.yaml not found."
    exit 1
fi

# Create internal/api directory if it doesn't exist
mkdir -p internal/api

# Generate API code
echo "Generating API code with ogen..."
ogen --package api --target internal/api --clean openapi/openapi.yaml

# Post-process: Fix ogen null handling bug (https://github.com/ogen-go/ogen/issues/1358)
if [ -f "internal/api/oas_json_gen.go" ]; then
    echo ""
    echo "Post-processing: Fixing Opt* null handling..."
    go run github.com/plexusone/ogen-tools/cmd/ogen-fixnull@latest internal/api/oas_json_gen.go 2>/dev/null || true
fi

# Post-process: Fix error body preservation
if [ -f "internal/api/oas_response_decoders_gen.go" ]; then
    echo ""
    echo "Post-processing: Fixing error body preservation..."
    go run github.com/plexusone/ogen-tools/cmd/ogen-fixerror@latest internal/api/oas_response_decoders_gen.go 2>/dev/null || true
fi

echo ""
echo "Running go mod tidy..."
go mod tidy

echo ""
echo "Verifying build..."
go build ./...

echo ""
echo "Done! API client generated successfully."
echo ""
echo "Generated files in internal/api/:"
ls -la internal/api/*.go 2>/dev/null | wc -l | xargs echo "  Total Go files:"
echo ""
echo "Next steps:"
echo "  1. Review changes in internal/api/"
echo "  2. Create/update SDK wrapper code in client.go"
echo "  3. Run tests: go test ./..."
echo "  4. Run linter: golangci-lint run"
