$ErrorActionPreference = "Stop"

$APP_NAME = "neoarc"

# Read version from .version file
$VERSION_FILE = Join-Path $PSScriptRoot ".version"
if (Test-Path $VERSION_FILE) {
    $VERSION = "$(Get-Content $VERSION_FILE -Raw | ForEach-Object { $_.Trim() })"
} else {
    $VERSION = "v0.0.0"
}

$COMMIT = git rev-parse --short HEAD

$BUILD_TIME = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")

$OUTPUT_DIR = Join-Path $PSScriptRoot "bin"
New-Item -ItemType Directory -Force -Path $OUTPUT_DIR | Out-Null

$CLI_DIR = Join-Path $PSScriptRoot "neoarc-cli"

$LDFLAGS = "-s -w -X neoarc/internal/cli.Version=$VERSION -X neoarc/internal/cli.BuildTime=$BUILD_TIME -X neoarc/internal/cli.Commit=$COMMIT -X neoarc/internal/version.Version=$VERSION -X neoarc/internal/version.BuildTime=$BUILD_TIME -X neoarc/internal/version.Commit=$COMMIT"

$PLATFORMS = @(
    "windows/amd64",
    "windows/arm64",
    "darwin/amd64",
    "darwin/arm64",
    "linux/amd64",
    "linux/arm64"
)

Write-Host "Building $APP_NAME version $VERSION"
Write-Host "----------------------------------------"

foreach ($PLATFORM in $PLATFORMS) {
    $GOOS, $GOARCH = $PLATFORM -split "/"
    $OUTPUT_NAME = "$APP_NAME-$GOOS-$GOARCH"
    if ($GOOS -eq "windows") {
        $OUTPUT_NAME += ".exe"
    }
    $OUTPUT_PATH = Join-Path $OUTPUT_DIR $OUTPUT_NAME

    Write-Host "Building for $GOOS/$GOARCH..."

    $env:CGO_ENABLED = "0"
    $env:GOOS = $GOOS
    $env:GOARCH = $GOARCH

    Push-Location $CLI_DIR
    & go build -ldflags "$LDFLAGS" -o "$OUTPUT_PATH" ./cmd/neoarc/
    $exitCode = $LASTEXITCODE
    Pop-Location

    if ($exitCode -ne 0) {
        Write-Host "ERR Failed to build for $GOOS/$GOARCH" -ForegroundColor Red
        exit 1
    }
}

Remove-Item Env:\CGO_ENABLED, Env:\GOOS, Env:\GOARCH -ErrorAction SilentlyContinue

Write-Host "----------------------------------------"
Write-Host "Build complete!"
Write-Host "Output directory: $OUTPUT_DIR"
Get-ChildItem -Path $OUTPUT_DIR | ForEach-Object { Write-Host "  $($_.Name)" }
