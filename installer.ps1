<#
.SYNOPSIS
  NeoArc Installer for Windows — downloads the latest release binary and installs globally.
.DESCRIPTION
  Detects the system architecture, downloads the correct NeoArc binary from GitHub,
  stores it in ~\.config\neostore\neoarc\bin\neoarc.exe, and adds that directory
  to the user's PATH.
.EXAMPLE
  .\installer.ps1
#>

$ErrorActionPreference = "Stop"
$Repo = "rkriad585/NeoArc"
$Version = "v1.0.0"
$InstallDir = "$env:USERPROFILE\.config\neostore\neoarc\bin"

function Write-Step {
    param([string]$Message)
    Write-Host ">>> $Message" -ForegroundColor Cyan
}

function Write-Success {
    param([string]$Message)
    Write-Host "OK  $Message" -ForegroundColor Green
}

function Write-Error {
    param([string]$Message)
    Write-Host "ERR $Message" -ForegroundColor Red
    exit 1
}

# ---- Detect architecture ----
$arch = $env:PROCESSOR_ARCHITECTURE
if ($arch -eq "AMD64") {
    $Binary = "neoarc-windows-amd64.exe"
} elseif ($arch -eq "ARM64") {
    $Binary = "neoarc-windows-amd64.exe"
    Write-Host "Warn: No native ARM64 binary — using x86_64 (will run under emulation)" -ForegroundColor Yellow
} else {
    Write-Error "Unsupported architecture: $arch"
}

$Url = "https://github.com/$Repo/releases/download/$Version/$Binary"
$InstallPath = "$InstallDir\neoarc.exe"

# ---- Create install directory ----
Write-Step "Creating install directory: $InstallDir"
New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
Write-Success "Directory ready."

# ---- Download binary ----
Write-Step "Downloading $Binary from $Url"
try {
    $ProgressPreference = "SilentlyContinue"
    Invoke-WebRequest -Uri $Url -OutFile $InstallPath -UseBasicParsing
    Write-Success "Downloaded to $InstallPath"
} catch {
    Write-Error "Download failed: $_"
}

# ---- Verify file ----
if (-not (Test-Path $InstallPath)) {
    Write-Error "Binary not found after download."
}

# ---- Add to PATH ----
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$InstallDir*") {
    Write-Step "Adding $InstallDir to user PATH"
    [Environment]::SetEnvironmentVariable("Path", "$currentPath;$InstallDir", "User")
    Write-Success "PATH updated (log out and back in, or restart your terminal for changes to take effect)."
} else {
    Write-Success "$InstallDir already in PATH."
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "  NeoArc installed successfully!" -ForegroundColor Green
Write-Host "  Binary : $InstallPath" -ForegroundColor Gray
Write-Host "  Version: $Version" -ForegroundColor Gray
Write-Host "  Usage  : neoarc help" -ForegroundColor Gray
Write-Host "========================================" -ForegroundColor Green
