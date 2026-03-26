param([string]$APP_NAME)

$ErrorActionPreference = "Stop"

if (-not $APP_NAME) {
    Write-Host "Usage: .\build.ps1 <appname>"
    exit 1
}

$VERSION = "1.0.0"
$COMMIT = "none"

try { $VERSION = git describe --tags --always } catch {}
try { $COMMIT = git rev-parse --short HEAD } catch {}

$BUILD_TIME = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")

$OUTPUT_DIR = "build"
New-Item -ItemType Directory -Force -Path $OUTPUT_DIR | Out-Null

$LDFLAGS = "-s -w -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.Commit=$COMMIT"

$PLATFORMS = @(
    "windows/amd64",
    "linux/amd64",
    "linux/arm64",
    "darwin/amd64",
    "darwin/arm64",
    "android/arm64"
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

    go build -trimpath -buildvcs=false -ldflags "$LDFLAGS" -o "$OUTPUT_DIR\$OUTPUT_NAME"

    if ($LASTEXITCODE -ne 0) {                 
        Write-Error "Build failed for $GOOS/$GOARCH"
        exit 1
    }
}

Write-Host "----------------------------------------"
Write-Host "✅ Build complete!"
Write-Host "📦 Output directory: $OUTPUT_DIR"
