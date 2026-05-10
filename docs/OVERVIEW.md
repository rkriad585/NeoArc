# NeoArc: Project Overview

NeoArc is an innovative, cross-platform Command Obfuscation and Alias Execution System designed to streamline the execution of complex scripts across various environments. It provides a centralized web-based "Execution Grid" for storing scripts (Bash, PowerShell, Python, Go) and a lightweight, high-performance Go CLI tool for triggering them remotely via simple aliases.

## Key Concepts

- **Execution Grid (Web Server):** The central hub where users manage their aliases, view public aliases, and interact with the system through a "Liquid Glass" UI.
- **Client (CLI Tool):** A fast, single-binary Go application that connects to the Execution Grid to fetch and execute commands on the local machine.
- **Aliases:** User-defined names for complex scripts or commands, stored on the server and executed via the CLI.
- **Polyglot Execution:** The CLI's ability to automatically detect and run scripts written in Bash, PowerShell, CMD, Python, or Go.

## Core Features

### Web Server (The Matrix)
- **Liquid Glass UI:** A modern, dark-mode interface built with TailwindCSS, featuring frosted glass effects and DotGothic typography for a cyberpunk aesthetic.
- **Secure User System:** Comprehensive user management including registration, login, profile updates, and avatar uploads.
- **Global Search:** Enables users to discover and utilize public aliases shared by the community.
- **Node Inspector:** Provides detailed views of aliases, including syntax-highlighted code previews.
- **Role-Based Access Control (RBAC):** A dedicated Super Admin panel for platform moderation and user management.
- **CSRF Protection:** All POST/PUT/DELETE routes require a CSRF token, preventing cross-site request forgery attacks.
- **API Token Authentication:** CLI requests must present a Bearer API token for server-to-client authentication.
- **Rate Limiting:** Login, registration, and API endpoints are rate-limited to prevent brute-force and abuse.
- **ETag Caching:** Server-side AliasCache (30s TTL) uses MD5-based ETags for conditional responses (HTTP 304).
- **Session Fixation Prevention:** Session regeneration on login prevents fixation attacks.
- **Image Upload Validation:** Avatar uploads verify both file extension and MIME magic bytes.

### CLI Tool (The Client)
- **Cross-Platform Compatibility:** Natively supported on Windows, Linux, macOS, and Termux (Android).
- **Versatile Script Execution:** Capable of running scripts in Bash/Sh, PowerShell, CMD, Python, and Go (with automatic compilation).
- **Optimized for Performance:** Developed in Go for speed, efficiency, and single-binary deployment, ensuring minimal overhead.
- **Client-Side Alias Cache:** Locally caches fetched aliases with a 30-second TTL, reducing network requests.
- **Trust-on-First-Use (TOFU):** Prompts the user for confirmation before executing a remote alias for the first time; saves trusted aliases locally.
- **ETag Support:** Sends `If-None-Match` headers to the server; uses cached data on HTTP 304 responses.
- **Dry-Run Mode:** `--dry-run` flag prints alias code without executing.
- **Yes Mode:** `--yes` flag skips the execution confirmation prompt for automation.
- **Argument Passing:** Arguments after the alias name are forwarded to the executed script (`$1`, `$args[0]`, `sys.argv[1]`, etc.).
- **Shell Tab Completion:** `neoarc completion bash|zsh|powershell` generates completions for subcommands, flags, and alias names.

## Technology Stack

- **Frontend:** HTML5, TailwindCSS (CDN), jQuery, Highlight.js, FontAwesome.
- **Backend:** Python Flask, Blueprint Architecture, SQLite3, waitress (production).
- **Client:** Golang (`net/http`, `os/exec`).
- **Design Language:** Inspired by Nothing OS (Dot Matrix Typography & Monochrome/Red Palette).
- **Configuration:** Environment-based via `.env` file (loaded by `python-dotenv`).

## Project Structure

```text
NEOARC/
|
+-- neoarc-cli/                     # Golang Client
|   +-- bin/                        # Precompiled binaries
|   +-- cmd/neoarc/main.go          # CLI entry point
|   +-- internal/cli/
|   |   +-- cli.go                  # CLI core logic (config, cache, trust, execution)
|   |   +-- completion.go           # Shell completion generators (bash/zsh/powershell)
|   |   +-- cli_test.go
|   +-- tests/
|   |   +-- integration_test.go
|   +-- go.mod / go.sum
|
+-- .version
+-- build.ps1 / build.sh
+-- installer.ps1 / installer.sh
+-- neoarc-server/                  # Python Flask Server
|   +-- core/
|   |   +-- admin.py                # Super admin logic
|   |   +-- auth.py                 # CSRF generation, csrf_required, api_token_required decorators, session regeneration
|   |   +-- cache.py                # AliasCache class with ETag support, 30s TTL
|   |   +-- db.py                   # SQLite init and connection management
|   |   +-- helpers.py              # Flask g-scoped db_execute and db_commit helpers
|   |   +-- rate_limiter.py         # In-memory RateLimiter (configurable max_attempts/window)
|   |   +-- validation.py           # Email, alias name, password, image MIME validation
|   |   +-- routes/                 # Flask Blueprints
|   |       +-- __init__.py         # register_route_blueprints()
|   |       +-- alias_routes.py
|   |       +-- api_routes.py
|   |       +-- auth_routes.py
|   |       +-- profile_routes.py
|   +-- templates/
|   +-- tests/
|   +-- .env.example
|   +-- config.py                   # All settings sourced from .env
|   +-- Dockerfile
|   +-- main.py
|   +-- pyproject.toml
|   +-- wsgi.py                     # Production entry point (waitress)
|
+-- docker-compose.yml
+-- static/
+-- README.md
```

Written by [Neorwc](https://github.com/rkriad585/neorwc-cli), Created by RK Riad Khan
