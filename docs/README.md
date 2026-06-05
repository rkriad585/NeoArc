# NeoArc

![NeoArc Banner](https://img.shields.io/badge/NeoArc-Execution_Grid-ea2b2b?style=for-the-badge&logo=terminator&logoColor=white)
![Python](https://img.shields.io/badge/Server-Python_Flask-black?style=for-the-badge&logo=python)
![Go](https://img.shields.io/badge/Client-Golang-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

**NeoArc** is a sophisticated, cross-platform **Command Obfuscation and Alias Execution System**. It bridges the gap between web-based storage and terminal execution, allowing users to store complex scripts (Bash, PowerShell, Python, Go) in a centralized "Execution Grid" and trigger them remotely on any device using a simple alias.

Featuring a premium **"Liquid Glass" / "Nothing OS" UI**, NeoArc offers a cyberpunk aesthetic with powerful administrative control.

---

## Features

### Web Server (The Matrix)
- **Liquid Glass UI:** A stunning dark-mode interface using TailwindCSS, frosted glass effects, and DotGothic typography.
- **User System:** Secure Registration, Login, and Profile management with Avatar uploads.
- **Global Search:** Discover public aliases created by other users in the grid.
- **Node Inspector:** View detailed info on aliases, including syntax-highlighted code previews.
- **Role-Based Access Control (RBAC):** Dedicated Super Admin panel to manage the platform.

### CLI Tool (The Client)
- **Cross-Platform:** Runs natively on Windows, Linux, macOS, and Termux (Android).
- **Polyglot Execution:** Automatically detects and runs code in:
    - `Bash` / `Sh`
    - `PowerShell`
    - `CMD`
    - `Python`
    - `Golang` (Auto-compiles and runs)
- **Argument Passing:** CLI arguments after the alias name are forwarded to the script (`$1`, `$args[0]`, `sys.argv[1]`, etc.)
- **Tab Completion:** `neoarc completion bash|zsh|powershell` generates shell-specific completions for commands, flags, and alias names
- **Stealth/Speed:** written in Go for high performance and single-binary deployment.
- **Themes:** 13 built-in color themes with runtime switching via `neoarc config theme <name>`.
- **TUI Editor:** `neoarc edit` opens an interactive configuration form

---

## Project Structure

```text
NEOARC/
|
+-- .version                        # Version file (v1.5.0)
+-- build.ps1                       # Root-level Windows Build Script
+-- build.sh                        # Root-level Unix Build Script
+-- installer.ps1                   # One-line Windows installer
+-- installer.sh                    # One-line Unix installer
+-- Dockerfile                      # Multi-stage CLI Docker build
+-- docker-compose.yml              # Server Docker stack
+-- Makefile                        # Common dev/build commands
+-- CMakeLists.txt                  # CMake IDE integration
+-- SECURITY.md                     # Security policy
+-- CODE_OF_CONDUCT.md              # Contributor covenant
+-- CONTRIBUTING.md                 # Contribution guidelines
+-- LICENSE                         # MIT license
|
+-- neoarc-cli/                     # Golang Client
|   +-- cmd/neoarc/main.go          # CLI Entry point
|   +-- internal/
|   |   +-- cli/
|   |   |   +-- cli.go              # CLI Logic (config, cache, trust store, execution, update)
|   |   |   +-- completion.go       # Shell completion generators (bash/zsh/powershell)
|   |   |   +-- form.go             # Interactive TUI configuration editor
|   |   |   +-- cli_test.go         # CLI unit tests (40+)
|   |   +-- config/
|   |   |   +-- config.go           # Cross-platform config directory helpers (TOML)
|   |   |   +-- theme.go            # Built-in color themes (13 themes)
|   |   +-- banner/
|   |   |   +-- banner.go           # Startup banner display
|   |   +-- version/
|   |   |   +-- version.go          # Version info helpers
|   +-- tests/integration_test.go   # Integration tests
|   +-- go.mod / go.sum
|
+-- neoarc-server/                  # Python Flask Server
|   +-- core/                       # Backend Modules
|   |   +-- admin.py                # Super Admin Logic
|   |   +-- auth.py                 # CSRF protection, API token auth, session management
|   |   +-- cache.py                # AliasCache with server-side ETag (30s TTL)
|   |   +-- db.py                   # Database Connection & Schema
|   |   +-- helpers.py              # g-scoped DB connection helpers
|   |   +-- rate_limiter.py         # In-memory rate limiter middleware
|   |   +-- validation.py           # Email, alias name, password, image MIME validation
|   |   +-- routes/                 # Blueprint-based route modules
|   |       +-- __init__.py         # Blueprint registration
|   |       +-- alias_routes.py     # Alias CRUD routes
|   |       +-- api_routes.py       # CLI-facing API endpoints
|   |       +-- auth_routes.py      # Login/register routes
|   |       +-- profile_routes.py   # Profile management routes
|   +-- static/                     # Logo, Favicons, Uploads
|   +-- templates/                  # HTML Templates (Jinja2)
|   |   +-- admin/                  # Admin Panel Templates
|   |   |   +-- dashboard.html
|   |   |   +-- login.html
|   |   |   +-- view_user.html
|   |   +-- base.html
|   |   +-- dashboard.html
|   |   +-- edit_alias.html
|   |   +-- error.html
|   |   +-- login.html
|   |   +-- profile.html
|   |   +-- register.html
|   |   +-- search.html
|   |   +-- view.html
|   +-- tests/
|   |   +-- conftest.py
|   |   +-- helpers.py
|   |   +-- test_api.py
|   |   +-- test_routes.py
|   |   +-- test_security.py
|   +-- .env.example                # Environment configuration template
|   +-- config.py                   # Reads all config from .env (not inline creds)
|   +-- Dockerfile
|   +-- main.py                     # App entry point
|   +-- pyproject.toml              # Python package definition & dependencies
|   +-- wsgi.py                     # Production WSGI entry point (waitress)
|
+-- docker-compose.yml
+-- .dockerignore
+-- .gitignore
+-- static/
+-- README.md
```

---

## Installation & Setup

### 1. Server Setup (Python)

Ensure you have **Python 3.10+** installed.

1.  Navigate to the server directory:
    ```bash
    cd neoarc-server
    ```

2.  Create and activate a virtual environment (recommended):
    ```bash
    # Windows
    python -m venv .venv
    .venv\Scripts\activate

    # Linux/macOS
    python3 -m venv .venv
    source .venv/bin/activate
    ```

3.  Install dependencies (uses pyproject.toml):
    ```bash
    pip install .
    # Or for editable development install:
    pip install -e .
    # Or if you use uv:
    uv pip install -e .
    ```

4.  **Configuration:**
    Copy `.env.example` to `.env` and set your values:
    ```bash
    cp .env.example .env
    ```
    Then edit `.env` with your preferred settings (secret key, admin password, API token, etc.).

5.  Run the server:
    ```bash
    python main.py
    ```
    *The server will start at `http://localhost:59248`.*

### 2. Client Setup (Golang)

Ensure you have **Go 1.25+** installed.

Build the binary from the repo root:

- **Windows (PowerShell):**
    ```powershell
    .\build.ps1
    ```
- **Linux / macOS / Termux:**
    ```bash
    chmod +x build.sh
    ./build.sh
    ```

Both scripts cross-compile for 6 platforms (Windows amd64/arm64, macOS amd64/arm64, Linux amd64/arm64), injection version (from `.version`), commit hash, and build time via ldflags. Output is placed in `./bin/`.

### One-Line Install (recommended)

Download and install the latest pre-built binary without any toolchain:

**Unix (Linux / macOS):**
```bash
curl -fsSL https://raw.githubusercontent.com/rkriad585/NeoArc/main/installer.sh | sh
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/rkriad585/NeoArc/main/installer.ps1 | iex
```

The installer auto-detects your OS and architecture, downloads the correct binary from GitHub Releases, and adds it to your PATH.

### Manual Build

```bash
cd neoarc-cli
go build -o neoarc ./cmd/neoarc
```

### Uninstall

```bash
# Via installer script
./installer.sh --selfuninstall          # Linux / macOS
.\installer.ps1 --selfuninstall         # Windows

# Via CLI itself
neoarc --selfuninstall
```

---

## CLI Configuration System

All CLI configuration, cache, and trust data live under `~/.config/neostore/neoarc/`:

| File | Format | Purpose |
|------|--------|---------|
| `config.toml` | TOML | Server URL, API token, TLS settings |
| `alias_cache.json` | JSON | Cached alias responses (30s TTL) |
| `trusted.json` | JSON | Approved aliases (TOFU trust store) |

### Config Directory Locations

- **Windows:** `%USERPROFILE%\.config\neostore\neoarc\`
- **Linux/macOS:** `~/.config/neostore/neoarc/`

### Saved Output Directory

Downloaded or generated files are saved to:
- **Windows:** `%USERPROFILE%\Downloads\neostore\neoarc\`
- **Linux/macOS:** `~/Downloads/neostore/neoarc/`

### Custom Config Path

Use the `--config` global flag to override:
```bash
neoarc --config /path/to/custom.toml run my-alias
```

### Theme Configuration

The CLI includes 13 built-in color themes for help output, update messages, and confirmation prompts. Switch themes at runtime:

```bash
neoarc config theme sunny_beach_day
neoarc config theme list          # List all available themes
```

Theme names: `dark`, `light`, `sunny_beach_day`, `olive_garden_feast`, `summer_ocean_breeze`, `refreshing_summer_fun`, `black_gold_elegance`, `vibrant_color_fiesta`, `light_steel`, `golden_twilight`, `deep_sea`, `bright_green`, `vivid_nightfall`

The theme is persisted in `config.toml` as `theme = "name"`.

### Config File Example

```toml
server_url = "http://localhost:59248"
api_token = "your-api-token-here"
insecure_tls = false
theme = "sunny_beach_day"
```

### Legacy Migration

Old paths (`%APPDATA%\neoarc\` or `~/.neoarc/`) are detected and migrated to the new
`~/.config/neostore/neoarc/` directory automatically on first run. JSON configs are
converted to TOML format during migration.

---

## Usage Guide

### Configuring the Client
Before running aliases, point the CLI to your web server:

```bash
neoarc config http://localhost:59248
# Or your public IP/Domain
neoarc config http://192.168.1.100:59248
```

### Executing Aliases
Once you have created an alias on the website (e.g., named `sys_info`), run it instantly:

```bash
# Run immediately
neoarc run sys_info

# Shorthand
neoarc sys_info

# View the code without running
neoarc get sys_info

# Dry-run (print code without executing)
neoarc --dry-run sys_info

# Skip execution confirmation prompt
neoarc --yes sys_info

# Pass arguments to the alias
neoarc run greet Alice Bob

# Generate shell completion
neoarc completion bash > /etc/bash_completion.d/neoarc
```

---

## Admin Panel

NeoArc comes with a built-in Super Admin dashboard to moderate the grid.

1.  Navigate to: `http://localhost:59248/admin/login`
2.  Log in using credentials configured in `.env`.

**Capabilities:**
- **Dashboard:** View all users and their alias counts.
- **User Inspector:** View full profile, email, and list of aliases.
- **Moderation:**
    - **Block User:** Prevent a user from logging in.
    - **Delete User:** Purge a user and all their aliases from the database.
    - **Delete Alias:** Remove specific malicious or broken aliases.

---

## Technology Stack

- **Frontend:** HTML5, TailwindCSS (CDN), jQuery, Highlight.js, FontAwesome (via SVG).
- **Backend:** Python Flask, Blueprint Architecture, SQLite3, waitress (production).
- **Client:** Golang (`net/http`, `os/exec`).
- **Design Language:** Nothing OS (Dot Matrix Typography & Monochrome/Red Palette).

---

## Documentation Index

| Document | Description |
|----------|-------------|
| [README.md](README.md) | This file — project overview and setup |
| [introduction.md](introduction.md) | Introduction to NeoArc concepts |
| [OVERVIEW.md](OVERVIEW.md) | Detailed feature overview |
| [ARCHITECTURE.md](ARCHITECTURE.md) | System architecture and data flow |
| [CLI_SETUP.md](CLI_SETUP.md) | CLI build and installation guide |
| [SERVER_SETUP.md](SERVER_SETUP.md) | Server setup and configuration |
| [USAGE_GUIDE.md](USAGE_GUIDE.md) | CLI usage reference |
| [DOCKER.md](DOCKER.md) | Docker deployment guide |
| [DEVELOPMENT.md](DEVELOPMENT.md) | Development workflow and build guide |
| [server/](server/) | Server-specific documentation |

---

## Contributing

Contributions are welcome!
1.  Fork the Project.
2.  Create your Feature Branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your Changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the Branch (`git push origin feature/AmazingFeature`).
5.  Open a Pull Request.

## License

Distributed under the MIT License. See `LICENSE` for more information.

---

<p align="center">
  <span style="font-family: monospace;">NEOARC v1.5.0</span>
</p>

Written by [Neorwc](https://github.com/rkriad585/neorwc-cli), Created by RK Riad Khan
