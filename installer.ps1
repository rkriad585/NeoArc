<#
.SYNOPSIS
  NeoArc Installer for Windows — downloads the latest release binary and installs globally.
.DESCRIPTION
  Detects the system architecture, downloads the correct NeoArc binary from GitHub,
  stores it in ~\.config\neostore\ncoarc\bin\neoarc.exe, and adds that directory
  to the user's PATH.
  Use --selfuninstall to remove NeoArc from the system.
.EXAMPLE
  .\installer.ps1                # Install NeoArc
  neoarc --selfuninstall  # Uninstall NeoArc
#>

$ErrorActionPreference = "Stop"

$PROJECT_NAME = "neoarc"
$REPO_OWNER   = "rkriad585"
$REPO         = "$REPO_OWNER/$PROJECT_NAME"
$INSTALL_DIR  = "$env:USERPROFILE\.config\neostore\$PROJECT_NAME\bin"
$INSTALL_PATH = "$INSTALL_DIR\$PROJECT_NAME.exe"

function Write-Step   { param([string]$Message) Write-Host ">>> $Message" -ForegroundColor Cyan }
function Write-Success { param([string]$Message) Write-Host "OK  $Message" -ForegroundColor Green }
function Write-Error   { param([string]$Message) Write-Host "ERR $Message" -ForegroundColor Red; exit 1 }

# ---- Resolve version from GitHub raw .version file ----
Write-Step "Resolving latest version..."
try {
    $ProgressPreference = "SilentlyContinue"
    $v = (Invoke-WebRequest -Uri "https://raw.githubusercontent.com/$REPO/main/.version" -UseBasicParsing).Content.Trim()
    $Version = $v
    Write-Success "Version: $Version"
} catch {
    Write-Error "Could not determine version from https://raw.githubusercontent.com/$REPO/main/.version"
}

# ---- Detect architecture using WMI ----
Write-Step "Detecting system architecture..."
$archCode = (Get-WmiObject Win32_Processor).Architecture
switch ($archCode) {
    0  { Write-Error "Unsupported architecture: x86 (32-bit)" }
    5  { Write-Error "Unsupported architecture: ARM (32-bit)" }
    6  { Write-Error "Unsupported architecture: IA64 (Itanium)" }
    9  { $Binary = "$PROJECT_NAME-windows-amd64.exe"; Write-Success "Architecture: AMD64" }
    12 { $Binary = "$PROJECT_NAME-windows-arm64.exe"; Write-Success "Architecture: ARM64" }
    default { Write-Error "Unknown architecture code: $archCode" }
}

$Url = "https://github.com/$REPO/releases/download/$Version/$Binary"

# ---- Handle uninstall ----
if ($args[0] -eq "--selfuninstall") {
    Write-Step "Uninstalling $PROJECT_NAME..."
    if (Test-Path $INSTALL_DIR) {
        Remove-Item -Recurse -Force $INSTALL_DIR
        Write-Success "Removed $INSTALL_DIR"
    } else {
        Write-Success "No install directory found."
    }

    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($currentPath -like "*$INSTALL_DIR*") {
        $newPath = ($currentPath -split ";" | Where-Object { $_ -ne $INSTALL_DIR }) -join ";"
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        Write-Success "Removed $INSTALL_DIR from user PATH"
    } else {
        Write-Success "Install directory not in PATH."
    }

    Write-Host ""
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "  $PROJECT_NAME has been uninstalled." -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    exit 0
}

# ---- Create install directory ----
Write-Step "Creating install directory: $INSTALL_DIR"
New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
Write-Success "Directory ready."

# ---- Download binary ----
Write-Step "Downloading $Binary from $Url"
try {
    $ProgressPreference = "SilentlyContinue"
    Invoke-WebRequest -Uri $Url -OutFile $INSTALL_PATH -UseBasicParsing
    Write-Success "Downloaded to $INSTALL_PATH"
} catch {
    Write-Error "Download failed: $_"
}

if (-not (Test-Path $INSTALL_PATH)) {
    Write-Error "Binary not found after download."
}

# ---- Add to PATH ----
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$INSTALL_DIR*") {
    Write-Step "Adding $INSTALL_DIR to user PATH"
    [Environment]::SetEnvironmentVariable("Path", "$currentPath;$INSTALL_DIR", "User")
    Write-Success "PATH updated (restart your terminal for changes to take effect)."
} else {
    Write-Success "$INSTALL_DIR already in PATH."
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "  $PROJECT_NAME installed successfully!" -ForegroundColor Green
Write-Host "  Binary : $INSTALL_PATH" -ForegroundColor Gray
Write-Host "  Version: $Version" -ForegroundColor Gray
Write-Host "  Usage  : $PROJECT_NAME help" -ForegroundColor Gray
Write-Host "  To uninstall, run:" -ForegroundColor Gray
Write-Host "    neoarc --selfuninstall" -ForegroundColor Gray
Write-Host "========================================" -ForegroundColor Green
