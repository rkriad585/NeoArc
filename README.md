# NeoArc

![Logo](neoarc-server/static/logo.svg)

![NeoArc Banner](https://img.shields.io/badge/NeoArc-Execution_Grid-ea2b2b?style=for-the-badge&logo=terminator&logoColor=white) 
![Python](https://img.shields.io/badge/Server-Python_Flask-black?style=for-the-badge&logo=python) 
![Go](https://img.shields.io/badge/Client-Golang-00ADD8?style=for-the-badge&logo=go) 
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
[![Auto Build & Release](https://github.com/rkriad585/NeoArc/actions/workflows/release.yml/badge.svg)](https://github.com/rkriad585/NeoArc/actions/workflows/release.yml)

**Cross-Platform Command Obfuscation & Alias Execution System**

NeoArc is a cyber/Matrix-themed ecosystem for storing, managing, and executing command aliases and scripts across platforms. It consists of a Flask web server (the dashboard) and a Go CLI client that fetches and executes aliases remotely. Designed for developers, sysadmins, and anyone who wants a centralized, authenticated command store with a dystopian cyberpunk aesthetic.

---

## Architecture

```
┌─────────────────┐       HTTPS / Bearer Auth        ┌──────────────────────┐
│  NeoArc Server   │ ◄──────────────────────────────► │   NeoArc CLI (Go)    │
│  (Flask + WSGI)  │   GET /api/alias/<name>          │  get / run / config  │
│                  │   GET /api/aliases (list)         │  completion (tab)    │
│                  │       ETag caching, 304s         │                      │
│  SQLite (WAL)    │                                  │  Local cache + trust │
│  Web Dashboard   │                                  │  Update + selfuninst │
│  Admin Panel     │                                  │  6 platform builds   │
└─────────────────┘                                  └──────────────────────┘
```

Aliases are created and managed via the web dashboard, stored in a SQLite database, and served to authenticated CLI clients through a REST API protected by Bearer tokens. The CLI caches responses locally (30s TTL) and prompts users to confirm execution of remote code (with a persistent trust store for repeat aliases).

---

## Features

| Category | Details |
|---|---|
| **Web Dashboard** | Create, edit, delete, search, and view aliases. Cyber-themed dark UI with Tailwind CSS |
| **CLI Client** | `neoarc get`, `neoarc run`, `neoarc config`, `neoarc config-token`, `neoarc edit` (TUI form), `neoarc update`, `--config` flag. 13 color themes. Supports bash, PowerShell, cmd, Python, and Go exec types |
| **Authentication** | Session-based auth with CSRF protection. Rate-limited login (5/min) and registration (3/5min) |
| **API Security** | Bearer token authentication, rate-limited (60/min per IP), ETag caching with 304 responses |
| **Admin Panel** | User management — view, block/unblock, delete users. Configurable admin route prefix |
| **Profile Management** | Update full name, email. Upload avatar with MIME-signature validation (PNG, JPG, GIF, WebP) |
| **Caching** | Server-side in-memory alias cache (30s TTL, ETag). Client-side file-based cache (30s TTL) |
| **Trust System** | CLI prompts before executing remote code. Persistent trust store per alias with timestamp |
| **Database** | SQLite with WAL mode, foreign keys, busy timeout, retry-on-lock logic |
| **Argument Passing** | CLI passes args after alias name through to the executed script (`$1`, `$args[0]`, `sys.argv[1]`) |
| **Tab Completion** | `neoarc completion bash|zsh|powershell` generates shell completions for commands, flags, and alias names |
| **Self-Update** | `neoarc update` fetches latest release, compares version, and replaces the binary (deferred batch on Windows) |
| **Self-Uninstall** | `neoarc --selfuninstall` removes config, cache, and binary from the system |
| **Cross-Platform CLI** | Pre-built binaries for Windows (amd64/arm64), macOS (amd64/arm64), Linux (amd64/arm64) |
| **Docker** | Multi-stage Alpine build. Compose file with persistent volumes for DB and uploads |
| **Exec Types** | `bash`, `powershell`, `cmd`, `python`, `go` — mapped to appropriate runtime per OS |
| **Error Pages** | Themed 404, 500 error pages with matrix-flavored messages |

---

## Project Structure

```
NeoArc/
├── .version                    # Version file (v3.0.3)
├── build.ps1                   # Root-level Windows cross-compile script
├── build.sh                    # Root-level Unix cross-compile script
├── installer.ps1               # One-line Windows installer
├── installer.sh                # One-line Unix installer
├── docker-compose.yml          # Production Docker stack
│
├── .github/workflows/          # GitHub Actions automation
│   └── release.yml             # Auto Build & Release pipeline
│
├── neoarc-server/              # Flask web application
│   ├── core/
│   │   ├── routes/             # Blueprint modules
│   │   │   ├── alias_routes.py     # Dashboard, CRUD, search
│   │   │   ├── api_routes.py       # REST API endpoint
│   │   │   ├── auth_routes.py      # Login, register, logout
│   │   │   └── profile_routes.py   # Profile management
│   │   ├── admin.py            # Admin panel (Blueprint)
│   │   ├── auth.py             # CSRF, API token, session helpers
│   │   ├── cache.py            # In-memory alias cache with ETag
│   │   ├── db.py               # SQLite setup (WAL, indexes)
│   │   ├── helpers.py          # DB connection helpers
│   │   ├── rate_limiter.py     # In-memory sliding-window rate limiter
│   │   └── validation.py       # Input validation + image MIME detection
│   ├── templates/              # Jinja2 templates (Tailwind CSS)
│   ├── tests/                  # Pytest suite (66 tests)
│   │   ├── test_api.py
│   │   ├── test_routes.py
│   │   └── test_security.py
│   ├── config.py               # Environment-based configuration
│   ├── main.py                 # Flask application factory
│   ├── wsgi.py                 # Waitress production entry point
│   ├── pyproject.toml          # Python package metadata
│   └── Dockerfile              # Multi-stage Alpine build
├── neoarc-cli/                 # Go CLI client
│   ├── cmd/neoarc/main.go      # Entry point (version vars injected via ldflags)
│   ├── internal/
│   │   ├── cli/
│   │   │   ├── cli.go              # Core logic (fetch, execute, config, cache, trust, update)
│   │   │   ├── completion.go       # Shell completion generators (bash/zsh/powershell)
│   │   │   └── cli_test.go         # Unit tests (40+)
│   │   └── config/
│   │       └── config.go           # Cross-platform config directory helpers (TOML)
│   ├── tests/
│   │   └── integration_test.go # Integration tests (3)
│   ├── go.mod
│   └── go.sum
│
└── static/uploads/             # Avatar uploads directory
```

---

## Quick Start

### Server (Python 3.10+)

```bash
cd neoarc-server
pip install -e .
python main.py
```

The server starts on `http://0.0.0.0:59248` by default. Open it in a browser, register an account, and start creating aliases.

### CLI (Go 1.25+)

```bash
cd neoarc-cli
go build -o neoarc ./cmd/neoarc
./neoarc config http://localhost:59248
./neoarc config-token <api-token-from-server>
./neoarc run my-alias
```

The API token is shown once on server startup (unless `NEOARC_API_TOKEN` is set). You can also bake it into the binary at build time (see build scripts).

---

## Client Configuration System

The CLI uses a **cross-platform configuration system** built into `internal/config/`.

### Config Directory Locations

| Platform | Config Path |
|----------|-------------|
| Windows | `%USERPROFILE%\.config\neostore\neoarc\` |
| Linux/macOS | `~/.config/neostore/neoarc/` |

Config is stored as TOML: `~/.config/neostore/neoarc/config.toml`

### Log Files

Logs are stored in the same config directory:
- **Windows:** `%USERPROFILE%\.config\neostore\neoarc\history.log`
- **Linux/macOS:** `~/.config/neostore/neoarc/history.log`

### Output / Saved Files

Downloaded or generated files are saved to:
- **Windows:** `%USERPROFILE%\Downloads\neostore\neoarc\`
- **Linux/macOS:** `~/Downloads/neostore/neoarc/`

### Helper Functions (`internal/config/config.go`)

| Function | Description |
|----------|-------------|
| `ConfigDir()` | Returns the config directory path |
| `EnsureConfigDir()` | Creates the config directory if it doesn't exist |
| `ConfigFile(name)` | Returns full path to a named config file |
| `LogFile(name)` | Returns full path to a log file |
| `SaveDir()` | Returns the output/save directory path |
| `HomeDir()` | Returns the user's home directory (OS-aware) |

### Custom Config File

Use `--config <path>` to point the CLI at a non-default config file:
```bash
neoarc --config /path/to/custom/config.toml run my-alias
```

### Config File Format (TOML)

```toml
server_url = "http://localhost:59248"
api_token = "your-api-token"
insecure_tls = false
```

### Config File Resolution Order

1. `--config <path>` flag (if provided)
2. Default: `~/.config/neostore/neoarc/config.toml`

### Migration from Legacy Config

On first run, the CLI automatically migrates config from old paths:
- Old Windows: `%APPDATA%\neoarc\config.json`
- Old Unix: `~/.neoarc/config.json`

Old JSON config files are converted to the new TOML format automatically.

---

## Theme Configuration

The CLI includes 13 built-in color themes that can be switched at runtime. Themes apply color to help output, update messages, and confirmation prompts.

### Commands

| Command | Description |
|---------|-------------|
| `neoarc config theme list` | List all available themes |
| `neoarc config theme <name>` | Switch to a named theme |

### Available Themes

| Name | Description | Colors (Hex) |
|------|-------------|--------------|
| `dark` | Dark Theme | `#0f172a #111827 #1e293b #334155 #64748b #e2e8f0` |
| `light` | Light Theme | `#ffffff #f8fafc #e2e8f0 #cbd5e1 #475569 #0f172a` |
| `sunny_beach_day` | Sunny Beach Day **(Default)** | `#264653 #2a9d8f #e9c46a #f4a261 #e76f51` |
| `olive_garden_feast` | Olive Garden Feast | `#606c38 #283618 #fefae0 #dda15e #bc6c25` |
| `summer_ocean_breeze` | Summer Ocean Breeze | `#e63946 #f1faee #a8dadc #457b9d #1d3557` |
| `refreshing_summer_fun` | Refreshing Summer Fun | `#8ecae6 #219ebc #023047 #ffb703 #fb8500` |
| `black_gold_elegance` | Black & Gold Elegance | `#000000 #14213d #fca311 #e5e5e5 #ffffff` |
| `vibrant_color_fiesta` | Vibrant Color Fiesta | `#ffbe0b #fb5607 #ff006e #8338ec #3a86ff` |
| `light_steel` | Light Steel | `#f8f9fa #e9ecef #dee2e6 #ced4da #adb5bd #6c757d #495057 #343a40 #212529` |
| `golden_twilight` | Golden Twilight | `#000814 #001d3d #003566 #ffc300 #ffd60a` |
| `deep_sea` | Deep Sea | `#0d1b2a #1b263b #415a77 #778da9 #e0e1dd` |
| `bright_green` | Bright Green | `#004b23 #006400 #007200 #008000 #38b000 #70e000 #9ef01a #ccff33` |
| `vivid_nightfall` | Vivid Nightfall | `#10002b #240046 #3c096c #5a189a #7b2cbf #9d4edd #c77dff #e0aaff` |

### Theme Persistence

The active theme is saved to `config.toml` and persists across sessions:

```toml
server_url = "http://localhost:59248"
api_token = "your-api-token"
insecure_tls = false
theme = "vibrant_color_fiesta"
```

### TUI Theme Picker

Open an interactive theme selector:

```bash
neoarc config theme edit
```

### TUI Config Editor

Open an interactive form to edit all config settings at once:

```bash
neoarc edit
# or
neoarc config edit
```

The editor shows a banner at startup and provides input fields for:
- Server URL
- API Token
- Insecure TLS toggle
- Theme selector (with preview)

### Examples

```bash
# List all themes
neoarc config theme list

# Switch to dark theme
neoarc config theme dark

# Switch via interactive picker
neoarc config theme edit

# Edit all config via TUI
neoarc edit
```

## Server Configuration Reference

All server configuration is via environment variables prefixed with `NEOARC_`.

| Variable | Default | Description |
|---|---|---|
| `NEOARC_PORT` | `59248` | Server listen port |
| `NEOARC_HOST` | `0.0.0.0` | Server bind address |
| `NEOARC_DEBUG` | `False` | Enable Flask debug mode |
| `NEOARC_SECRET_KEY` | auto-generated | Flask session signing key (set for persistence across restarts) |
| `NEOARC_API_TOKEN` | auto-generated | Bearer token for CLI authentication |
| `NEOARC_DB_PATH` | `neoarc.db` | SQLite database file path |
| `NEOARC_UPLOAD_FOLDER` | `static/uploads` | Avatar upload directory |
| `NEOARC_MAX_UPLOAD_MB` | `5` | Max upload file size in MB |
| `NEOARC_ADMIN_EMAIL` | `admin@neoarc.local` | Admin panel login email |
| `NEOARC_ADMIN_PASSWORD` | `change_me` | Admin panel password |
| `NEOARC_ADMIN_ROUTE` | `/admin` | Admin panel URL prefix |

Set these via a `.env` file in the `neoarc-server/` directory or pass them directly to Docker.

---

## Build Instructions

### Server

```bash
cd neoarc-server
pip install -e .
```

Or install from `pyproject.toml`:

```bash
cd neoarc-server
pip install .
```

Run with Waitress for production:

```bash
python wsgi.py
```

### CLI

#### Automated (Recommended)

Push a tag and let GitHub Actions build all 6 platform binaries automatically:

```bash
git tag v3.0.3
git push --tags
```

See the [Release Workflow](#release-workflow-automated) section above for details.

#### Local Build Scripts

The repository includes build scripts for cross-compiling locally:

**Unix:**
```bash
chmod +x build.sh
./build.sh
```

**Windows (PowerShell):**
```powershell
.\build.ps1
```

Both scripts:
- Auto-detect the server's `.env` file to bake in the API token
- Build for Windows/amd64, Windows/arm64, macOS/amd64, macOS/arm64, Linux/amd64, Linux/arm64
- Inject version (from `.version` file), commit hash, and build time via linker flags
- Output binaries to `./bin/`

#### Manual Single-Platform Build

```bash
VERSION="v$(cat .version | tr -d '[:space:]')"
cd neoarc-cli
go build -ldflags="-s -w -X neoarc/internal/cli.Version=$VERSION" -o neoarc ./cmd/neoarc
```

---

## Release Workflow (Automated)

NeoArc includes a **GitHub Actions release workflow** (`.github/workflows/release.yml`) that automatically builds and publishes binaries whenever you push a tag.

### How It Works

```text
prepare ──────────────────────────────────────────────┐
   └─ fetches .version from GitHub raw URL            │
   └─ extracts commit SHA + pre-release flag          │
                                                      ▼
build (6 parallel) ──────────────────────────────► release ──► notify-on-failure
   windows/amd64                                      ▲
   windows/arm64                                      │
   linux/amd64                                        │
   linux/arm64    (cross-compiled w/ aarch64-gcc)     │
   darwin/amd64   (macos-latest Intel runner)         │
   darwin/arm64   (macos-latest M2 native runner)     │
                                                      │
changelog ────────────────────────────────────────────┘
   └─ groups commits: feat / fix / perf / docs / other
```

### Binary Output Names

| Platform | Architecture | File |
|----------|-------------|------|
| Windows | AMD64 | `neoarc-windows-amd64.exe` |
| Windows | ARM64 | `neoarc-windows-arm64.exe` |
| Linux | AMD64 | `neoarc-linux-amd64` |
| Linux | ARM64 | `neoarc-linux-arm64` |
| macOS | AMD64 (Intel) | `neoarc-darwin-amd64` |
| macOS | ARM64 (Apple Silicon) | `neoarc-darwin-arm64` |

### Metadata Injected at Build Time

| Variable | Source |
|----------|--------|
| `main.Version` | `.version` file (e.g. `v3.0.3`) |
| `main.Commit` | Short git SHA (8 chars) |
| `main.PublisherName` | `rkriad585` |
| `main.PublisherEmail` | `rkriad585@gmail.com` |

### How to Publish a Release

```bash
# 1. Update .version file
echo "v3.0.4" > .version

# 2. Commit and push
git add .version
git commit -m "Release v3.0.4"
git push

# 3. Tag and push
git tag v3.0.4
git push --tags
```

The workflow will:
1. Read the version from `.version` (with fallback to the tag itself)
2. Cross-compile for all 6 platforms in parallel
3. Generate SHA-256 checksums
4. Build a changelog from commit history (grouped by feat/fix/perf/docs)
5. Publish a GitHub Release with all assets
6. Send a failure notification if any step fails

> **Security note:** The release workflow only injects public metadata (`Version`, `Commit`, `PublisherName`, `PublisherEmail`). The `NEOARC_API_TOKEN` is **not** baked into release binaries — doing so would leak the server's API secret to every downloader. Each user must set their own token at runtime:
> ```bash
> neoarc config-token <your-token>
> ```
> This saves the token locally to `~/.config/neostore/neoarc/config.toml`. The local `build.sh`/`build.ps1` scripts can optionally auto-inject the token for private/internal builds by reading `neoarc-server/.env`.
>
> **Note:** The workflow also supports manual trigger via the GitHub Actions UI (`workflow_dispatch`).

---

## Docker Deployment

Build and run with Docker Compose:

```bash
docker compose build
docker compose up -d
```

Required environment variables (set in `.env` or shell):

```bash
NEOARC_SECRET_KEY=your-secret-key
NEOARC_API_TOKEN=your-api-token
NEOARC_ADMIN_PASSWORD=your-admin-password
```

The compose file mounts two persistent volumes:
- `neoarc-data` — SQLite database at `/app/data/`
- `neoarc-uploads` — Avatar uploads at `/app/static/uploads/`

The server is exposed on port `59248`.

### Standalone Docker

```bash
cd neoarc-server
docker build -t neoarc-server .
docker run -d \
  -p 59248:59248 \
  -v neoarc-data:/app/data \
  -v neoarc-uploads:/app/static/uploads \
  -e NEOARC_SECRET_KEY=... \
  -e NEOARC_API_TOKEN=... \
  -e NEOARC_ADMIN_PASSWORD=... \
  neoarc-server
```

---

## Installation (One-Line)

Download and install the latest NeoArc binary with a single command:

**Unix (Linux / macOS):**
```bash
curl -fsSL https://raw.githubusercontent.com/rkriad585/NeoArc/main/installer.sh | sh
```

**Windows (PowerShell):**
```powershell
# PowerShell 5+ (recommended)
irm https://raw.githubusercontent.com/rkriad585/NeoArc/main/installer.ps1 | iex

# Legacy fallback
iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/rkriad585/NeoArc/main/installer.ps1'))
```

### What the installer does

1. Fetches the latest version from the repository's `.version` file
2. Auto-detects your OS and CPU architecture
3. Downloads the correct pre-built binary from GitHub Releases
4. Installs it to `~/.config/neostore/neoarc/bin/neoarc` (or `.exe` on Windows)
5. Adds that directory to your `PATH` so `neoarc` is available globally

### Supported platforms

| OS | Architectures |
|----|--------------|
| Windows | AMD64, ARM64 |
| Linux | AMD64, ARM64 |
| macOS | AMD64 (Intel), ARM64 (Apple Silicon) |

### Uninstall

Use any of these methods (they all do the same thing):

```bash
# Via installer script (local copy)
./installer.sh --selfuninstall          # Linux / macOS
.\installer.ps1 --selfuninstall         # Windows

# Via installer script (one-liner, no download needed)
curl -fsSL https://raw.githubusercontent.com/rkriad585/NeoArc/main/installer.sh | sh -s -- --selfuninstall       # Linux / macOS
irm https://raw.githubusercontent.com/rkriad585/NeoArc/main/installer.ps1 | iex "-" "--selfuninstall"            # Windows PowerShell

# Via CLI itself (if still accessible)
neoarc --selfuninstall
```

Short aliases: `-u`, `--uninstall`, `-selfuninstall` also work.

**What uninstall removes:**
- The NeoArc binary (`~/.config/neostore/neoarc/bin/neoarc`)
- All config, cache, and trust data (`~/.config/neostore/neoarc/`)
- PATH entries from shell profiles (Unix) or User environment variables (Windows)

---

## CLI Usage

```
neoarc get <alias> [args...]       — Print alias code to stdout
neoarc run <alias> [args...]       — Fetch and execute alias with args
neoarc <alias> [args...]           — Shorthand for run
neoarc config <server-url>         — Set NeoArc server URL
neoarc config-token <tok>          — Set API token
neoarc config insecure             — Skip TLS certificate verification
neoarc config secure               — Re-enable TLS verification
neoarc config theme <name>         — Set color theme (use 'list' for all)
neoarc edit                        — Open TUI configuration editor
neoarc version                     — Show the installed version
neoarc help                        — Show help
neoarc completion <shell>          — Generate shell completion (bash|zsh|powershell)
neoarc update                      — Self-update to the latest version

Flags:
  -v, --version                    — Show the installed version
  -h, --help                       — Show this help menu
  --install                        — Download and install NeoArc to ~/.config/neostore/neoarc/bin/
  --selfuninstall                  — Remove NeoArc config, cache, and binary
  --config <path>                  — Use a custom config file path (before any command)
  --dry-run                        — Print alias code without executing (before alias)
  --yes                            — Skip execution confirmation prompt (before alias)

Args after <alias> are passed through to the executed command:
  bash/sh:    $1, $2, $@
  powershell: $args[0], $args[1]
  python:     sys.argv[1], sys.argv[2]
```

---

## Test Suite

### Server (66 tests)

```bash
cd neoarc-server
pip install -e ".[test]"   # if pytest is in extras, else pip install pytest
pytest tests/ -v
```

Tests cover: registration, login, alias CRUD, search, profile updates, API auth, ETag caching, rate limiting, CSRF protection, input validation, session protection, blocked users, cross-user isolation, password reset.

### CLI (40+ unit + 3 integration)

```bash
cd neoarc-cli
go test ./internal/cli/... -v    # unit tests
go test ./tests/... -v           # integration tests (builds binary)
```

---

## Contributing

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature/my-feature`).
3. Make changes and ensure tests pass.
4. Commit with a descriptive message.
5. Open a pull request.

Code style: Python follows PEP 8, Go follows `gofmt`. All new features should include tests.

---
