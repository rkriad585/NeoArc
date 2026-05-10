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

---

## Project Structure

```text
NEOARC/
|
+-- neoarc-cli/                     # Golang Client
|   +-- bin/                        # Precompiled binaries
|   +-- cmd/neoarc/main.go          # CLI Entry point
|   +-- internal/cli/
|   |   +-- cli.go                  # CLI Logic (config, cache, trust store, execution)
|   |   +-- completion.go           # Shell completion generators (bash/zsh/powershell)
|   |   +-- cli_test.go             # CLI unit tests
|   +-- tests/integration_test.go   # Integration tests
|   +-- go.mod
|   +-- go.sum
|   +-- build.ps1                   # Windows Build Script
|   +-- build.sh                    # Linux/Mac Build Script
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

Ensure you have **Go 1.20+** installed.

1.  Navigate to the CLI directory:
    ```bash
    cd neoarc-cli
    ```

2.  Build the binary:

    - **Windows (PowerShell):**
        ```powershell
        .\build.ps1
        ```
    - **Linux / macOS / Termux:**
        ```bash
        chmod +x build.sh
        ./build.sh
        ```

3.  Add the binary to your system PATH or move it to `/usr/local/bin` (Linux) or `C:\Windows\System32` (Windows) for global access.

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
  <span style="font-family: monospace;">NEOARC v3.0.3</span>
</p>

Written by [Neorwc](https://github.com/rkriad585/neorwc-cli), Created by RK Riad Khan
