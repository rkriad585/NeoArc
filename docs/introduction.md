# Introduction to NeoArc

NeoArc is an innovative, cross-platform **Command Obfuscation and Alias Execution System** designed to streamline the execution of complex scripts across various environments. It acts as a bridge between a centralized web-based "Execution Grid" and local terminal execution, allowing users to define and store scripts (Bash, PowerShell, Python, Go, etc.) as simple aliases. These aliases can then be triggered remotely on any configured device using a lightweight command-line interface (CLI) tool.

The system is built with a focus on both aesthetics and functionality, featuring a premium "Liquid Glass" / "Nothing OS" inspired user interface on the web server and a high-performance, polyglot CLI client.

## Core Concepts

- **Execution Grid (Web Server):** The central hub where users manage their profiles, create and store aliases, and discover public scripts. It provides a secure, visually appealing interface for script management.
- **Client (CLI Tool):** A lightweight, cross-platform binary that connects to the Execution Grid. It fetches aliases and executes the associated commands locally, automatically detecting the appropriate interpreter or compiler.
- **Aliases:** User-defined names linked to specific commands or scripts. These can range from simple shell commands to multi-line scripts in various languages.
- **Polyglot Execution:** The CLI's ability to intelligently determine the execution environment (Bash, PowerShell, Python, Go, CMD) based on the alias's `exec_type` and run the script accordingly.

## Key Features

### Web Server (The Matrix)

- **Liquid Glass UI:** A modern, dark-mode interface leveraging TailwindCSS, frosted glass effects, and the distinctive DotGothic typography, creating a cyberpunk aesthetic.
- **User System:** Comprehensive user management including secure registration, login, and profile customization with avatar uploads.
- **Global Search:** A powerful feature allowing users to discover and inspect public aliases shared by the community within the Execution Grid.
- **Node Inspector:** Provides detailed information about aliases, including syntax-highlighted code previews for better understanding and security review.
- **Role-Based Access Control (RBAC):** A dedicated Super Admin panel for platform moderation, user management, and alias control.
- **Environment-Based Configuration:** All settings (secret key, admin credentials, database path, upload limits, API token) are configured via a `.env` file, loaded by `python-dotenv`. No hardcoded credentials in source code.
- **CSRF Protection:** All POST, PUT, and DELETE routes require a valid CSRF token (generated per-session, compared via constant-time comparison). The `@csrf_required` decorator enforces this on every form submission.
- **API Bearer Token Authentication:** CLI requests must include an `Authorization: Bearer <token>` header. The token is validated via `secrets.compare_digest` for timing-attack resistance.
- **ETag-Based Caching:** The server-side `AliasCache` (30-second TTL) computes an MD5-based ETag for each alias JSON response. CLI clients send `If-None-Match` headers; the server returns HTTP 304 when the cache is still fresh.
- **Rate Limiting:** Three rate limiters protect login, registration, and API endpoints. Each is an in-memory sliding-window counter with configurable `max_attempts` and `window_seconds`. Defaults: 5 attempts per 60 seconds.
- **Session Fixation Prevention:** On successful login, the session is cleared and regenerated (`regenerate_session`), invalidating any session ID set before authentication.
- **Image Upload MIME Validation:** Avatar uploads verify both the file extension (against an allowlist) and the file content (magic bytes: PNG, JPEG, GIF, WebP). Files with wrong MIME types are rejected even if the extension matches.

### CLI Tool (The Client)

- **Cross-Platform Compatibility:** Natively supported on Windows, Linux, macOS, and Termux (Android), ensuring broad accessibility.
- **Polyglot Execution Engine:** Automatically identifies and executes code written in:
    - `Bash` / `Sh`
    - `PowerShell`
    - `CMD` (Windows Command Prompt)
    - `Python`
    - `Golang` (Includes auto-compilation for Go scripts)
- **Stealth & Performance:** Developed in Go, the CLI offers high performance, minimal resource usage, and single-binary deployment for ease of distribution and execution.
- **Cross-Platform Config System:** All configuration, cache, and trust data are stored in `~/.config/neostore/neoarc/` using TOML format. Legacy JSON configs are auto-migrated on first run.
- **Client-Side Alias Cache:** Fetched aliases are cached locally in `alias_cache.json` with a 30-second TTL. Subsequent fetches within the TTL use the cached version without a network request.
- **Trust-on-First-Use (TOFU):** Before executing an alias for the first time, the CLI prompts the user for confirmation. Approved aliases are stored in `trusted.json` and skip the prompt on future runs.
- **ETag Conditional Requests:** The CLI stores ETags from server responses and sends `If-None-Match` headers. On HTTP 304, the cached alias is used without re-downloading.
- **Dry-Run Mode:** The `--dry-run` flag prints the alias code to stdout without executing it.
- **Yes Mode:** The `--yes` flag skips the execution confirmation prompt, enabling automated/scripted usage.
- **Argument Passing:** Arguments after the alias name are forwarded verbatim to the executed script (`$1`, `$args[0]`, `sys.argv[1]`, etc. depending on the exec type).
- **Shell Tab Completion:** `neoarc completion bash|zsh|powershell` generates completion scripts for subcommands, flags, and alias names fetched from the server via `GET /api/aliases`.

## Use Cases

- **DevOps Automation:** Store and execute complex deployment scripts, server maintenance tasks, or CI/CD triggers with simple aliases.
- **System Administration:** Centralize common administrative commands (e.g., log analysis, service restarts, system health checks) for quick access across multiple servers.
- **Personal Productivity:** Create shortcuts for frequently used commands, development environment setups, or custom utility scripts.
- **Security & Penetration Testing:** Manage and deploy custom exploit scripts or reconnaissance tools from a central repository.
- **Educational Tool:** Share and learn from community-contributed scripts in a controlled environment.

NeoArc aims to empower users with a flexible, secure, and visually engaging platform for managing and executing their digital commands.

Written by [Neorwc](https://github.com/rkriad585/neorwc-cli), Created by RK Riad Khan
