#!/usr/bin/env bash
set -e

APP_NAME="neoarc"
VERSION=$(git describe --tags --always)
COMMIT=$(git rev-parse --short HEAD)
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

SERVER_ENV="$(dirname "$0")/../neoarc-server/.env"
API_TOKEN=""
if [ -f "$SERVER_ENV" ]; then
    API_TOKEN=$(grep -E '^NEOARC_API_TOKEN=' "$SERVER_ENV" | cut -d '=' -f2- | tr -d ' \t\r\n')
    if [ -n "$API_TOKEN" ]; then
        echo "API token loaded from server .env"
    fi
fi

OUTPUT_DIR="bin"
mkdir -p "$OUTPUT_DIR"

LDFLAGS="-s -w \
-X neoarc/internal/cli.Version=$VERSION \
-X neoarc/internal/cli.BuildTime=$BUILD_TIME \
-X neoarc/internal/cli.Commit=$COMMIT"

if [ -n "$API_TOKEN" ]; then
    LDFLAGS="$LDFLAGS -X neoarc/internal/cli.APIToken=$API_TOKEN"
fi

PLATFORMS=(
"windows/amd64"
"darwin/amd64"
"darwin/arm64"
"linux/amd64"
"linux/arm64"
)

echo "Building $APP_NAME version $VERSION"
echo "----------------------------------------"

for PLATFORM in "${PLATFORMS[@]}"
do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}
    OUTPUT_NAME="$APP_NAME-$GOOS-$GOARCH"
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME+=".exe"
    fi
    echo "Building for $GOOS/$GOARCH..."
    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH \
    go build \
    -trimpath \
    -buildvcs=false \
    -ldflags="$LDFLAGS" \
    -o "$OUTPUT_DIR/$OUTPUT_NAME" ./cmd/neoarc
done

echo "----------------------------------------"
echo "Build complete!"
echo "Output directory: $OUTPUT_DIR"
