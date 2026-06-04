# NeoArc CLI

A Go command-line tool that connects to the NeoArc server to fetch and execute command aliases remotely. Supports multiple execution environments with caching, ETag support, trust-on-first-use, and TLS configuration.

## Features

- Remote alias execution across multiple environments: bash, PowerShell, cmd, Python, Go
- Local caching with 30-second TTL and ETag/304-based cache validation
- Trust-on-first-use (TOFU) prompting for untrusted aliases
- Configurable TLS (insecure mode for development)
- API token resolution: config file overrides compile-time token
- Cross-platform: Windows, macOS, Linux
- Dry-run mode and auto-yes flag
- Argument passing — all args after alias name are forwarded to the executed script
- Shell tab completion — `neoarc completion bash|zsh|powershell` generates completion scripts

## Installation

### Pre-built binaries

Download the latest binary for your platform from the releases page. Binaries are cross-compiled for:

- `windows/amd64`
- `windows/arm64`
- `darwin/amd64`, `darwin/arm64`
- `linux/amd64`, `linux/arm64`

### Build from source

```bash
# Unix
./build.sh

# Windows
.\build.ps1
```

Build output is placed in the `bin/` directory with the version, commit hash, and build time injected via ldflags.

## Usage

```
neoarc help
neoarc config <url>
neoarc config insecure
neoarc config secure
neoarc config-token <token>
neoarc completion <shell>
neoarc update
neoarc self-update
neoarc get <alias> [args...]
neoarc run <alias> [args...]
neoarc <alias> [args...]
```

### Commands

**`neoarc help`** — Display full usage information.

**`neoarc config <url>`** — Set the NeoArc server URL.

```
neoarc config https://neoarc.example.com
```

**`neoarc config insecure`** — Enable insecure TLS (skip certificate verification, useful for development with self-signed certs).

```
neoarc config insecure
```

**`neoarc config secure`** — Re-enable strict TLS verification.

```
neoarc config secure
```

**`neoarc config-token <token>`** — Set the API token in local config (takes priority over compile-time token).

```
neoarc config-token ncl_abc123...
```

**`neoarc update`** / **`neoarc self-update`** — Check for updates and replace the current binary with the latest release version.

```
neoarc update
neoarc self-update
neoarc update --proxy http://proxy:8080
neoarc self-update -p http://proxy:8080
```

The command fetches the latest version from GitHub, compares it to the current version, downloads the appropriate binary for your platform, and replaces the running executable (via `.old` rename on Windows, direct rename on Unix).

**`neoarc get <alias> [args...]`** — Fetch and display the alias command without executing it.

```
neoarc get deploy-staging
```

**`neoarc run <alias> [args...]`** — Fetch and execute an alias after trust confirmation (unless already trusted or `--yes` is set). All remaining args are forwarded to the executed script.

```
neoarc run deploy-staging --env production
```

**`neoarc <alias> [args...]`** — Shorthand for `neoarc run <alias>`.

```
neoarc deploy-staging --env production
```

**`neoarc completion <shell>`** — Generate a shell completion script for the given shell.

```
neoarc completion bash   # source this output in .bashrc
neoarc completion zsh    # source this output in .zshrc
neoarc completion powershell  # source this output in your PowerShell profile
```

Supported shells: `bash`, `zsh`, `powershell`.

### Argument Passing

All arguments after the alias name are forwarded to the executed script. How they are accessed depends on the `exec_type`:

| Exec type | Access pattern |
|-----------|---------------|
| `bash` / `sh` | `$1`, `$2`, `$@` |
| `powershell` | `$args[0]`, `$args[1]` |
| `python` | `sys.argv[1]`, `sys.argv[2]` |
| `cmd` | `%1`, `%2` |
| `go` | `os.Args[1]`, `os.Args[2]` |

Example with a bash alias:
```
neoarc run greet Alice Bob
# Inside the script: $1=Alice, $2=Bob
```

### Flags

Flags must be placed before the alias name:

- `--dry-run` — Print the command that would be executed without running it.
- `--yes` — Skip the trust prompt and auto-approve execution.

```
neoarc --dry-run deploy-staging
neoarc --yes deploy-staging
```

## Install

Install NeoArc to `~/.config/neostore/neoarc/bin/` automatically:

```bash
neoarc --install
```

This downloads the latest release binary for your platform, copies it to the install directory, and adds the directory to your PATH. If the download fails (e.g., the release tag doesn't exist yet), it falls back to copying the current running binary.

## Uninstall

Remove NeoArc config, cache, and binary from your system:

```bash
# Via the CLI (self-uninstall)
neoarc --selfuninstall
```

This deletes the config directory (`~/.config/neostore/neoarc/`), removes the binary on Linux/macOS, and creates a deferred delete script on Windows. PATH removal instructions are printed after the operation.

You can also use the installer scripts:
```bash
# Linux / macOS
./installer.sh --selfuninstall

# Windows
.\installer.ps1 --selfuninstall
```

## Configuration

The CLI stores configuration in `~/.config/neostore/neoarc/` on all platforms. Old paths (`%APPDATA%/neoarc/` on Windows, `~/.neoarc/` on Unix) are auto-migrated on first run.

### config.json

```json
{
  "server_url": "https://neoarc.example.com",
  "api_token": "ncl_...",
  "insecure_tls": false
}
```

### API token resolution

1. If a token is set in `config.json`, it is used.
2. Otherwise, the compile-time token (injected via ldflags at build) is used.
3. If neither is available, requests are made without authentication.

## Caching

After a successful alias fetch, the result is cached locally with a 30-second TTL. Subsequent requests within the TTL return the cached entry immediately.

On cache hit after the TTL expires, the CLI sends the cached ETag via the `If-None-Match` header. If the server responds with `304 Not Modified`, the cached entry is reused and its TTL is refreshed. The ETag is computed from the server response (MD5 header) or falls back to a SHA256 hash of the response content.

Cache is stored in `cache.json` in the config directory.

## Trust store

When executing an alias for the first time, the CLI prompts for confirmation:

```
Execute alias 'deploy-staging'? [y/N]
```

On approval, the alias is added to the trust store (`trusted.json`). Subsequent executions skip the prompt automatically.

Use `--yes` to bypass the prompt for all aliases in a single invocation.

## Architecture

```
+-----------+     GET /api/alias/{name}     +-------------+
| NeoArc CLI |     GET /api/aliases (list)   | NeoArc Server |
|            | ----------------------------> |              |
+-----------+  <---------------------------- +-------------+
     |                    |
     | 304 Not Modified   |
     | (cached entry)     |
     v                    v
  Local cache          Local trust store
  (30s TTL + ETag)     (trusted aliases)
     |
     v
  Execute with args:
  bash | powershell | cmd | python | go run
```

The CLI sends authenticated requests to the NeoArc server's `/api/alias/{name}` endpoint. The server returns the alias metadata including the command string and execution environment. The CLI then writes the command to a temporary file with the appropriate extension and executes it via the corresponding interpreter.

Supported exec types and their handlers:

| Exec type | Interpreter | Temp file |
|-----------|-------------|-----------|
| `bash` | `bash` | `.sh` |
| `powershell` | `powershell` | `.ps1` |
| `cmd` | `cmd /c` (Windows) / `sh -c` (Unix) | `.bat` |
| `python` | `python` | `.py` |
| `go` | `go run` | `.go` |

## Build from source

### Test Suite

```bash
go test ./internal/cli/... -v    # 40+ unit tests
go test ./tests/... -v           # 3 integration tests (builds binary)
```

### Prerequisites

- Go 1.25+
- Git

### Unix

```bash
chmod +x build.sh
./build.sh
```

The script reads `NEOARC_API_TOKEN` from `../neoarc-server/.env`, reads version from the root `.version` file, and cross-compiles for all supported targets.

### Windows

```powershell
.\build.ps1
```

Same logic as the Unix variant using PowerShell equivalents.

## License

MIT
