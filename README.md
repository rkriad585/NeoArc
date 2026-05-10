# NeoArc

![Logo](neoarc-server/static/logo.svg)

![NeoArc Banner](https://img.shields.io/badge/NeoArc-Execution_Grid-ea2b2b?style=for-the-badge&logo=terminator&logoColor=white) 
![Python](https://img.shields.io/badge/Server-Python_Flask-black?style=for-the-badge&logo=python) 
![Go](https://img.shields.io/badge/Client-Golang-00ADD8?style=for-the-badge&logo=go) 
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

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
| **CLI Client** | `neoarc get`, `neoarc run`, `neoarc config`, `neoarc config-token`, `neoarc update`. Supports bash, PowerShell, cmd, Python, and Go exec types |
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
│   ├── cmd/neoarc/main.go      # Entry point
│   ├── internal/cli/
│   │   ├── cli.go              # Core logic (fetch, execute, config, cache, trust, update)
│   │   ├── completion.go       # Shell completion generators (bash/zsh/powershell)
│   │   └── cli_test.go         # Unit tests (40+)
│   ├── tests/
│   │   └── integration_test.go # Integration tests (3)
│   └── go.mod
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

## Configuration Reference

All configuration is via environment variables prefixed with `NEOARC_`.

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

The CLI includes build scripts for cross-compiling to all platforms.

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

Manual single-platform build (version from .version, or leave empty for GitHub fallback):

```bash
VERSION="v$(cat .version | tr -d '[:space:]')"
cd neoarc-cli
go build -ldflags="-s -w -X neoarc/internal/cli.Version=$VERSION" -o neoarc ./cmd/neoarc
```

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
iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/rkriad585/NeoArc/main/installer.ps1'))
```

The binary is installed to `~/.config/neostore/neoarc/bin/neoarc` and the directory is added to your PATH.

### Uninstall

```bash
# Using the installer script
./installer.sh --selfuninstall          # Linux / macOS
.\installer.ps1 --selfuninstall         # Windows

# Or using the CLI itself (if still accessible)
neoarc --selfuninstall
```

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
neoarc help                        — Show help
neoarc completion <shell>          — Generate shell completion (bash|zsh|powershell)
neoarc update                      — Self-update to the latest version

Standalone flags:
  --install                        — Download and install NeoArc to ~/.config/neostore/neoarc/bin/
  --selfuninstall                  — Remove NeoArc config, cache, and binary

Flags (place before alias):
  --dry-run                        — Print alias code without executing
  --yes                            — Skip execution confirmation prompt

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
