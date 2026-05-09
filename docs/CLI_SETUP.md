# NeoArc CLI Setup Guide

This guide covers building and installing the NeoArc CLI client from source.

## Prerequisites

- **Go 1.25+** — [Download Go](https://go.dev/dl/)
- **Git** — for version info injection during build

## Build

Navigate to the `neoarc-cli/` directory and run the platform-specific build script.

### Windows (PowerShell)

```powershell
.\build.ps1
```

### Linux / macOS

```bash
chmod +x build.sh
./build.sh
```

## Output

Build artifacts are placed in the `bin/` directory (not `build/`):

| Platform       | Binary name                    |
|----------------|--------------------------------|
| Windows x86_64 | `neoarc-windows-amd64.exe`     |
| macOS x86_64   | `neoarc-darwin-amd64`          |
| macOS ARM64    | `neoarc-darwin-arm64`          |
| Linux x86_64   | `neoarc-linux-amd64`           |
| Linux ARM64    | `neoarc-linux-arm64`           |

## Compile-time API Token

Build scripts automatically read `NEOARC_API_TOKEN` from `../neoarc-server/.env`
and inject it via Go linker flags (`-X neoarc/internal/cli.APIToken=...`).

## API Token Resolution

At runtime, tokens are resolved in order:

1. **Local config** — set via `neoarc config-token <token>`
2. **Compile-time token** — injected during build from `neoarc-server/.env`
3. **Empty** — no token sent (server may reject unauthenticated requests)

## Configuration File Location

- **Windows:** `%APPDATA%\neoarc\config.json`
- **Linux / macOS:** `~/.neoarc/config.json`

If no config file exists, the CLI creates one with a default server URL of `http://localhost:59248`.

## Verify Installation

```bash
neoarc help
```

You should see the NeoArc help output. If the binary is not on your PATH,
move it to a directory in your PATH or add the `bin/` folder to your PATH.

## One-Line Installation

**Unix (Linux / macOS):**
```bash
curl -fsSL https://raw.githubusercontent.com/rkriad585/NeoArc/main/installer.sh | sh
```

**Windows (PowerShell):**
```powershell
iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/rkriad585/NeoArc/main/installer.ps1'))
```

Installs the binary to `~/.config/neostore/neoarc/bin/neoarc` and adds it to your PATH.

## Uninstall

```bash
./installer.sh --selfuninstall          # Unix
.\installer.ps1 --selfuninstall         # Windows
neoarc --selfuninstall                  # CLI (if accessible)
```

## Next Steps

- [Usage Guide](USAGE_GUIDE.md) — configure the CLI and run aliases
