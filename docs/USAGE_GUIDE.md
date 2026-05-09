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
| `neoarc config-token <token>` | Set API token in local config file |
| `neoarc completion <shell>` | Generate shell completion script |
| `neoarc help` | Show help |

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

- **Windows:** `%APPDATA%\neoarc\config.json`
- **Linux / macOS:** `~/.neoarc/config.json`

Example config:

```json
{
  "server_url": "http://localhost:59248",
  "api_token": "your-token-here",
  "insecure_tls": false
}
```
