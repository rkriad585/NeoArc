# NeoArc CLI Setup Guide

This guide covers building and installing the NeoArc CLI client from source.

## Prerequisites

- **Go 1.25+** — [Download Go](https://go.dev/dl/)
- **Git** — for version info injection during build

## Build

### Automated (GitHub Actions)

The recommended approach is to push a version tag and let GitHub Actions build all binaries:

```bash
git tag v3.0.3
git push --tags
```

This triggers the `.github/workflows/release.yml` pipeline which:
- Cross-compiles for 6 platforms in parallel
- Injects version, commit, publisher metadata via ldflags
- Generates SHA-256 checksums
- Publishes a GitHub Release with auto-generated changelog

### Local Build Scripts

Run the platform-specific build script from the **repo root**.

### Windows (PowerShell)

```powershell
.\build.ps1
```

### Linux / macOS

```bash
chmod +x build.sh
./build.sh
```

Both scripts cross-compile for 6 platforms (Windows amd64/arm64, macOS amd64/arm64, Linux amd64/arm64), injecting version (from `.version`), commit hash, and build time via ldflags.

## Output

**Local build** artifacts are placed in `./bin/`. **Release builds** (via GitHub Actions) use the same naming convention:

| Platform       | Architecture | Binary name                    |
|----------------|-------------|--------------------------------|
| Windows        | AMD64       | `neoarc-windows-amd64.exe`     |
| Windows        | ARM64       | `neoarc-windows-arm64.exe`     |
| macOS (Intel)  | AMD64       | `neoarc-darwin-amd64`          |
| macOS (Silicon)| ARM64       | `neoarc-darwin-arm64`          |
| Linux          | AMD64       | `neoarc-linux-amd64`           |
| Linux          | ARM64       | `neoarc-linux-arm64`           |

## Compile-time API Token

Build scripts automatically read `NEOARC_API_TOKEN` from `../neoarc-server/.env`
and inject it via Go linker flags (`-X neoarc/internal/cli.APIToken=...`).

## API Token Resolution

At runtime, tokens are resolved in order:

1. **Local config** — set via `neoarc config-token <token>`
2. **Compile-time token** — injected during build from `neoarc-server/.env`
3. **Empty** — no token sent (server may reject unauthenticated requests)

## Configuration File Location

All platforms: `~/.config/neostore/neoarc/config.toml`
- **Windows:** `%USERPROFILE%\.config\neostore\neoarc\config.toml`
- **Linux / macOS:** `~/.config/neostore/neoarc/config.toml`

> **Legacy migration:** The CLI automatically migrates from old JSON config paths
> (`%APPDATA%\neoarc\config.json` or `~/.neoarc/config.json`) to the new TOML format on first run.

If no config file exists, the CLI creates one with a default server URL of `http://localhost:59248`.

The config file also supports a `theme` field — see [Theme Configuration](USAGE_GUIDE.md#theme-configuration) for details.

Use the interactive TUI editor to configure settings with a guided form:
```bash
neoarc edit
neoarc config edit
```

Available config commands:
```bash
neoarc config <server-url>
neoarc config insecure
neoarc config secure
neoarc config theme <name>
neoarc config theme edit
neoarc config-token <token>
neoarc edit
neoarc --config /path/to/config.toml run my-alias
```

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
