#!/bin/sh

## Don't run this script directly. Use the command `make build` instead.
echo "Fetching latest version information"

# Fetch the latest to ensure we have the latest tag
git fetch origin

# Get the version from the git tag, ff not available, default to v0.0.0
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")

# All the ldflags for the build
LDFLAGS="-X main.Version=$VERSION"

echo "Building app $VERSION"

OUTPUT="bin/pismo"

go build -v -ldflags="$LDFLAGS" -o $OUTPUT ./cmd || exit 1

echo "Build successful. Output file: $OUTPUT"

