#!/usr/bin/env bash
set -e

APP_NAME="neoarc"

# Read version from .version file
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
VERSION_FILE="${SCRIPT_DIR}/.version"
if [ -f "$VERSION_FILE" ]; then
    VERSION="$(cat "$VERSION_FILE" | tr -d '[:space:]')"
else
    VERSION="v0.0.0"
fi

COMMIT=$(git rev-parse --short HEAD)
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

OUTPUT_DIR="${SCRIPT_DIR}/bin"
mkdir -p "$OUTPUT_DIR"

CLI_DIR="${SCRIPT_DIR}/neoarc-cli"

LDFLAGS="-s -w \
-X neoarc/internal/cli.Version=${VERSION} \
-X neoarc/internal/cli.BuildTime=${BUILD_TIME} \
-X neoarc/internal/cli.Commit=${COMMIT} \
-X neoarc/internal/version.Version=${VERSION} \
-X neoarc/internal/version.BuildTime=${BUILD_TIME} \
-X neoarc/internal/version.Commit=${COMMIT}"

PLATFORMS=(
    "windows/amd64"
    "windows/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
)

echo "Building ${APP_NAME} version ${VERSION}"
echo "----------------------------------------"

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS="${PLATFORM%/*}"
    GOARCH="${PLATFORM#*/}"
    OUTPUT_NAME="${APP_NAME}-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME+=".exe"
    fi
    OUTPUT_PATH="${OUTPUT_DIR}/${OUTPUT_NAME}"

    echo "Building for ${GOOS}/${GOARCH}..."

    CGO_ENABLED=0 GOOS="${GOOS}" GOARCH="${GOARCH}" \
        go build -ldflags "${LDFLAGS}" -o "${OUTPUT_PATH}" "${CLI_DIR}/cmd/neoarc/"
done

echo "----------------------------------------"
echo "Build complete!"
echo "Output directory: ${OUTPUT_DIR}"
ls -1 "${OUTPUT_DIR}"
