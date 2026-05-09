$ErrorActionPreference = "Stop"

$APP_NAME = "neoarc"
$VERSION = "1.0.0"
$COMMIT = "none"

try { $VERSION = git describe --tags --always } catch {}
try { $COMMIT = git rev-parse --short HEAD } catch {}

$BUILD_TIME = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")

$SERVER_ENV = Join-Path $PSScriptRoot "..\neoarc-server\.env"
$API_TOKEN = ""
if (Test-Path $SERVER_ENV) {
    $envContent = Get-Content $SERVER_ENV -Raw
    $match = [regex]::Match($envContent, 'NEOARC_API_TOKEN=(.+)')
    if ($match.Success) {
        $API_TOKEN = $match.Groups[1].Value.Trim()
        Write-Host "API token loaded from server .env"
    }
}

$OUTPUT_DIR = "bin"
New-Item -ItemType Directory -Force -Path $OUTPUT_DIR | Out-Null

$LDFLAGS = "-s -w -X neoarc/internal/cli.Version=$VERSION -X neoarc/internal/cli.BuildTime=$BUILD_TIME -X neoarc/internal/cli.Commit=$COMMIT"
if ($API_TOKEN) {
    $LDFLAGS += " -X neoarc/internal/cli.APIToken=$API_TOKEN"
}

$PLATFORMS = @(
    "windows/amd64",
    "darwin/amd64",
    "darwin/arm64",
    "linux/amd64",
    "linux/arm64"
)

Write-Host "Building $APP_NAME version $VERSION"

foreach ($PLATFORM in $PLATFORMS) {
    $GOOS, $GOARCH = $PLATFORM -split "/"
    $OUTPUT_NAME = "$APP_NAME-$GOOS-$GOARCH"
    if ($GOOS -eq "windows") {
        $OUTPUT_NAME += ".exe"
    }
    Write-Host "Building for $GOOS/$GOARCH"
    $env:CGO_ENABLED = "0"
    $env:GOOS = $GOOS
    $env:GOARCH = $GOARCH
    go build -trimpath -buildvcs=false -ldflags "$LDFLAGS" -o "$OUTPUT_DIR\$OUTPUT_NAME" ./cmd/neoarc
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Build failed for $GOOS/$GOARCH"
        exit 1
    }
}

Write-Host "----------------------------------------"
Write-Host "Build complete!"
Write-Host "Output directory: $OUTPUT_DIR"
