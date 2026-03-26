#!/usr/bin/env bash
set -e

APP_NAME="$1"
VERSION=$(git describe --tags --always 2>/dev/null || echo "1.0.0")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

OUTPUT_DIR="build"
mkdir -p $OUTPUT_DIR

LDFLAGS="-s -w \
-X main.Version=$VERSION \
-X main.BuildTime=$BUILD_TIME \
-X main.Commit=$COMMIT"

PLATFORMS=(
"windows/amd64"
"linux/amd64"
"linux/arm64"
"darwin/amd64"
"darwin/arm64"
"android/arm64"
)

echo "🚀 Building $APP_NAME version $VERSION"
echo "----------------------------------------"

for PLATFORM in "${PLATFORMS[@]}"
do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}

    OUTPUT_NAME="$APP_NAME-$GOOS-$GOARCH"
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME+=".exe"
    fi

    echo "🔨 Building for $GOOS/$GOARCH..."

    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH \
    go build \
    -trimpath \
    -buildvcs=false \
    -ldflags="$LDFLAGS" \
    -o "$OUTPUT_DIR/$OUTPUT_NAME"

done

echo "----------------------------------------"
echo "✅ Build complete!"
echo "📦 Output directory: $OUTPUT_DIR"