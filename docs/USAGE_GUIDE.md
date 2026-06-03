# NeoArc CLI Usage Guide

## Server Configuration

Before running aliases, point the CLI at your NeoArc server:

```bash
neoarc config http://localhost:59248
```

| Command | Description |
|---------|-------------|
| `neoarc config <server-url>` | Set server URL (default: `http://localhost:59248`) |
| `neoarc config insecure` | Enable TLS skip-verify |
| `neoarc config secure` | Disable TLS skip-verify |
| `neoarc config theme <name>` | Set color theme (`list` to show all) |
| `neoarc config-token <token>` | Set API token in local config file |
| `neoarc completion <shell>` | Generate shell completion script |
| `neoarc help` | Show help |
| `neoarc --config <path>` | Use a custom config file (global flag) |

## Alias Commands

| Command | Description |
|---------|-------------|
| `neoarc get <alias> [args...]` | Print alias code without executing |
| `neoarc run <alias> [args...]` | Fetch and execute alias with args |
| `neoarc <alias> [args...]` | Shorthand for `run` |

### Standalone Flags

| Flag | Description |
|------|-------------|
| `--install` | Download and install NeoArc to `~/.config/neostore/neoarc/bin/` |
| `--selfuninstall` | Remove NeoArc config, cache, and binary from the system |

### Flags

Flags must be placed before the alias name:

| Flag | Description |
|------|-------------|
| `--dry-run <alias>` | Print code without executing |
| `--yes <alias>` | Skip trust confirmation prompt |

### Argument Passing

All arguments after the alias name are forwarded to the executed script:

```bash
neoarc run myalias arg1 arg2
neoarc myalias --flag value
```

How args are accessed depends on the exec type:

| Type | Access |
|------|--------|
| bash/sh | `$1`, `$2`, `$@` |
| powershell | `$args[0]`, `$args[1]` |
| python | `sys.argv[1]`, `sys.argv[2]` |
| cmd | `%1`, `%2` |
| go | `os.Args[1]`, `os.Args[2]` |

### Shell Completion

Generate and install tab completion for your shell:

```bash
# Bash
neoarc completion bash > /etc/bash_completion.d/neoarc

# Zsh
neoarc completion zsh > /usr/local/share/zsh/site-functions/_neoarc

# PowerShell
neoarc completion powershell >> $PROFILE
```

Examples:

```bash
neoarc get sys_info
neoarc run sys_info
neoarc sys_info
neoarc --dry-run sys_info
neoarc --yes sys_info
neoarc completion bash
```

## Trust-on-First-Use

The first time you execute an alias, you are prompted to confirm:

```
Execute alias 'sys_info'? This will run code from the remote server. [y/N]:
```

Answering `y` or `yes` saves the alias to `trusted.json` so future runs proceed
without prompting. Answering anything else cancels execution.

## Theme Configuration

The CLI features 13 built-in color themes that affect help output, update messages, and confirmation prompts:

```bash
# Switch theme
neoarc config theme sunny_beach_day

# List available themes
neoarc config theme list
```

Available themes: `dark`, `light`, `sunny_beach_day`, `olive_garden_feast`, `summer_ocean_breeze`, `refreshing_summer_fun`, `black_gold_elegance`, `vibrant_color_fiesta`, `light_steel`, `golden_twilight`, `deep_sea`, `bright_green`, `vivid_nightfall`.

The selection is persisted in `config.toml`:
```toml
theme = "sunny_beach_day"
```

## Caching

Alias responses are cached locally for 30 seconds. Subsequent requests include
an `If-None-Match` header with the previous `ETag`. If the server responds
`304 Not Modified`, the cached value is used.

Cache file: `trusted.json` (same directory as `config.json`)

## Error Handling

Connection errors include a hint to check the configured server URL:

```
Error: connection failed: <underlying error>
Hint: Check server URL with: neoarc config <url> (default: http://localhost:59248)
```

Authentication failures (401/403) show:

```
Error: API authentication failed. Use 'neoarc config-token <token>' or rebuild the binary.
```

## Execution Types

The CLI infers how to execute code based on the alias's `exec_type` field:

| Type | Runner |
|------|--------|
| `bash` | `bash <tempfile>` |
| `powershell` | `powershell -ExecutionPolicy Bypass -File <tempfile>` |
| `python` | `python <tempfile>` |
| `go` | `go run <tempfile>` (auto-prepends `package main` if missing) |
| default (Windows) | `cmd /C <tempfile>` |
| default (Linux/macOS) | `sh <tempfile>` |

## Configuration File

**New location (v3+):**
- **Windows:** `%USERPROFILE%\.config\neostore\neoarc\config.toml`
- **Linux / macOS:** `~/.config/neostore/neoarc/config.toml`

*Legacy paths (`%APPDATA%\neoarc\` or `~/.neoarc/`) are auto-migrated on first run.*

### Config File Format (TOML)

```toml
server_url = "http://localhost:59248"
api_token = "your-token-here"
insecure_tls = false
```

### Custom Config Path

Use the `--config` global flag to specify a non-default config file:

```bash
neoarc --config /path/to/custom/config.toml run my-alias
neoarc --config /path/to/custom/config.toml config http://server:59248
```
