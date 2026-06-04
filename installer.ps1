# installer.ps1 - Installer and Uninstaller for NeoArc on Windows
# Auto detects architecture, downloads the release binary, sets up PATH, and supports self-uninstallation.

[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

# ── Configuration ────────────────────────────────────────────────────────────
$ProjectName    = "neoarc"
$PublisherName  = "rkriad585"
$GitHubRepo     = "rkriad585/NeoArc"

# Paths
$ConfigDir      = Join-Path $env:USERPROFILE ".config\neostore\$ProjectName"
$BinDir         = Join-Path $ConfigDir "bin"
$BinaryPath     = Join-Path $BinDir "${ProjectName}.exe"

# ── Check for Uninstallation Flag ───────────────────────────────────────────
$IsUninstall = $false
if ($args -contains "--selfuninstall" -or $args -contains "-selfuninstall" -or $args -contains "--uninstall" -or $args -contains "-u") {
    $IsUninstall = $true
}

if ($IsUninstall) {
    Write-Host ""
    Write-Host "╔══════════════════════════════════════════════════╗" -ForegroundColor Yellow
    Write-Host "║              NeoArc Uninstaller                  ║" -ForegroundColor Yellow
    Write-Host "╚══════════════════════════════════════════════════╝" -ForegroundColor Yellow
    Write-Host ""

    # 1. Remove binary and directories
    if (Test-Path $BinaryPath) {
        Write-Host "  Removing binary: $BinaryPath ... " -NoNewline
        Remove-Item $BinaryPath -Force
        Write-Host "OK" -ForegroundColor Green
    }

    if (Test-Path $BinDir) {
        Write-Host "  Removing bin directory: $BinDir ... " -NoNewline
        Remove-Item $BinDir -Recurse -Force -ErrorAction SilentlyContinue
        Write-Host "OK" -ForegroundColor Green
    }

    if (Test-Path $ConfigDir) {
        Write-Host "  Removing config directory: $ConfigDir ... " -NoNewline
        Remove-Item $ConfigDir -Recurse -Force -ErrorAction SilentlyContinue
        Write-Host "OK" -ForegroundColor Green
    }

    # 2. Clean up PATH environment variable (User level)
    Write-Host "  Updating User PATH environment variable ... " -NoNewline
    $userPath = [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::User)
    if ($userPath) {
        $pathEntries = $userPath -split ';'
        $cleanEntries = $pathEntries | Where-Object { $_.Trim().TrimEnd('\') -ine $BinDir.TrimEnd('\') }
        $newUserPath = $cleanEntries -join ';'

        [Environment]::SetEnvironmentVariable("Path", $newUserPath, [EnvironmentVariableTarget]::User)
        $env:PATH = ($env:PATH -split ';' | Where-Object { $_.Trim().TrimEnd('\') -ine $BinDir.TrimEnd('\') }) -join ';'
        Write-Host "OK" -ForegroundColor Green
    } else {
        Write-Host "Skipped (PATH empty)" -ForegroundColor Yellow
    }

    Write-Host ""
    Write-Host "  NeoArc has been successfully uninstalled from your system." -ForegroundColor Green
    Write-Host "  Please restart any open terminal windows to apply the changes." -ForegroundColor Cyan
    Write-Host ""
    exit 0
}

# ── Installation / Update Flow ───────────────────────────────────────────────
Write-Host ""
Write-Host "╔══════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║               NeoArc Installer                   ║" -ForegroundColor Cyan
Write-Host "╚══════════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

# 1. Resolve Version from GitHub
Write-Host "  Checking latest version from GitHub ... " -NoNewline
$versionUrl = "https://raw.githubusercontent.com/${GitHubRepo}/main/.version"
try {
    $version = (Invoke-RestMethod -Uri $versionUrl -UseBasicParsing).Trim()
    Write-Host "$version" -ForegroundColor Green
} catch {
    Write-Host "FAILED" -ForegroundColor Red
    Write-Error "Could not fetch version from GitHub. Please check your internet connection."
    exit 1
}

# 2. Detect System Architecture (Win32_Processor mapping)
Write-Host "  Detecting system architecture ... " -NoNewline
$procArch = $null
try {
    $procArch = (Get-CimInstance Win32_Processor -ErrorAction SilentlyContinue | Select-Object -First 1).Architecture
} catch {
    try {
        $procArch = (Get-WmiObject Win32_Processor -ErrorAction SilentlyContinue | Select-Object -First 1).Architecture
    } catch {}
}

if ($null -eq $procArch) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "AMD64" -or $env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
        $arch = "amd64"
    } elseif ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
        $arch = "arm64"
    } else {
        $arch = "amd64"
    }
} else {
    switch ($procArch) {
        9  { $arch = "amd64" }
        12 { $arch = "arm64" }
        default { $arch = "amd64" }
    }
}
Write-Host "windows-$arch" -ForegroundColor Green

# 3. Download the Release Binary
$downloadUrl = "https://github.com/${GitHubRepo}/releases/download/${version}/${ProjectName}-windows-${arch}.exe"
Write-Host "  Downloading binary from $downloadUrl ...`n"

if (-not (Test-Path $BinDir)) {
    New-Item -ItemType Directory -Path $BinDir -Force | Out-Null
}

try {
    $oldProgressPreference = $ProgressPreference
    $ProgressPreference = 'SilentlyContinue'
    Invoke-WebRequest -Uri $downloadUrl -OutFile $BinaryPath -UseBasicParsing
    $ProgressPreference = $oldProgressPreference
    $size = [math]::Round((Get-Item $BinaryPath).Length / 1MB, 2)
    Write-Host "  Successfully downloaded neoarc.exe (${size} MB)" -ForegroundColor Green
} catch {
    $ProgressPreference = $oldProgressPreference
    Write-Host "  Failed to download binary." -ForegroundColor Red
    Write-Error "Please ensure a release exists for version $version."
    exit 1
}

# 4. Update PATH Environment Variable
Write-Host "  Configuring PATH environment variable ... " -NoNewline
$userPath = [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::User)
$pathEntries = if ($userPath) { $userPath -split ';' } else { @() }

$binDirNormalized = $BinDir.TrimEnd('\')
$alreadyInPath = $false
foreach ($entry in $pathEntries) {
    if ($entry.Trim().TrimEnd('\') -ieq $binDirNormalized) {
        $alreadyInPath = $true
        break
    }
}

if (-not $alreadyInPath) {
    $newPathEntries = $pathEntries + $BinDir
    $newUserPath = $newPathEntries -join ';'
    [Environment]::SetEnvironmentVariable("Path", $newUserPath, [EnvironmentVariableTarget]::User)
    $env:PATH = $env:PATH + ";" + $BinDir
    Write-Host "Added to PATH" -ForegroundColor Green
} else {
    Write-Host "Already in PATH" -ForegroundColor Yellow
}

# ── Success Banner ────────────────────────────────────────────────────────────
Write-Host ""
Write-Host "╔══════════════════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║         NeoArc successfully installed!           ║" -ForegroundColor Green
Write-Host "╚══════════════════════════════════════════════════╝" -ForegroundColor Green
Write-Host ""
Write-Host "  Installation Path: $BinaryPath" -ForegroundColor Cyan
Write-Host "  Version          : $version" -ForegroundColor Cyan
Write-Host ""
Write-Host "  Please RESTART your terminal/PowerShell window to start using NeoArc." -ForegroundColor Yellow
Write-Host "  You can then run: neoarc --help" -ForegroundColor Yellow
Write-Host ""
