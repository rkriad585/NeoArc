# Development Guide

This guide covers setting up a NeoArc development environment, building from source, and contributing changes.

## Prerequisites

- **Go 1.25+** — [Download](https://go.dev/dl/)
- **Python 3.8+** — [Download](https://python.org/downloads/)
- **Git** — for version info and release management
- **Docker** (optional) — for containerized builds and server deployment
- **Make** (optional) — for using the project Makefile

## Project Layout

```
NeoArc/
├── neoarc-cli/          # Go CLI client
│   ├── cmd/neoarc/      # Entry point
│   └── internal/        # Library code
├── neoarc-server/       # Python Flask server
│   ├── core/            # Application logic
│   └── tests/           # Pytest suite
├── Dockerfile           # CLI multi-stage Docker build
├── docker-compose.yml   # Server Docker stack
├── Makefile             # Common dev commands
├── CMakeLists.txt       # CMake integration (optional)
└── build.sh / .ps1      # Cross-compile scripts
```

## Quick Start

### Server

```bash
cd neoarc-server
python -m venv .venv
source .venv/bin/activate   # or .venv\Scripts\activate on Windows
pip install -e .
cp .env.example .env        # edit as needed
python main.py
```

Server starts at `http://localhost:59248`.

### CLI

```bash
cd neoarc-cli
go build -o bin/neoarc ./cmd/neoarc/
./bin/neoarc config http://localhost:59248
./bin/neoarc config-token <api-token>
./bin/neoarc help
```

## Using the Makefile

The root `Makefile` provides convenience targets:

| Target | Description |
|--------|-------------|
| `make build` | Build CLI for current platform |
| `make build-all` | Cross-compile for all 6 platforms (runs `build.sh`) |
| `make test` | Run all tests (CLI + server) |
| `make lint` | Run `go vet` and `flake8` |
| `make format` | Format Go and Python code |
| `make clean` | Remove build artifacts |
| `make install` | Build and install CLI locally |
| `make release` | Print release instructions |
| `make docker` | Build CLI Docker image |
| `make docker-server` | Build server Docker image |

Example workflow:

```bash
make build          # quick CLI build
make test           # verify everything passes
make docker         # build Docker image
```

## Docker Development

### Build the CLI in Docker

```bash
make docker
# or
docker build -t neoarc-cli \
  --build-arg VERSION=$(cat .version) \
  --build-arg COMMIT=$(git rev-parse --short HEAD) \
  --build-arg BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  .
```

### Run the Server in Docker

```bash
make docker-server
# or
docker compose up -d
```

See [DOCKER.md](DOCKER.md) for detailed Docker usage.

## Cross-Compilation

Build for all supported platforms with a single command:

```bash
make build-all        # Unix
.\build.ps1           # Windows PowerShell
```

Output goes to `./bin/`:

| File | Platform |
|------|----------|
| `neoarc-windows-amd64.exe` | Windows AMD64 |
| `neoarc-windows-arm64.exe` | Windows ARM64 |
| `neoarc-darwin-amd64` | macOS Intel |
| `neoarc-darwin-arm64` | macOS Apple Silicon |
| `neoarc-linux-amd64` | Linux AMD64 |
| `neoarc-linux-arm64` | Linux ARM64 |

## Testing

### CLI Tests

```bash
# Unit tests
cd neoarc-cli && go test ./internal/... -v

# Integration tests
cd neoarc-cli && go test ./tests/... -v
```

### Server Tests

```bash
cd neoarc-server
pip install -e ".[test]"
python -m pytest tests/ -v
```

## Code Style

- **Go:** Run `gofmt -l -w .` before committing. Follow [Effective Go](https://go.dev/doc/effective_go).
- **Python:** Follow [PEP 8](https://peps.python.org/pep-0008/). Use `black` for formatting and `flake8` for linting.
- **General:** Keep functions small, write tests, avoid unnecessary comments.

## Release Process

1. Update `.version`:
   ```bash
   echo "v1.5.1" > .version
   ```
2. Commit and tag:
   ```bash
   git add .version && git commit -m "Release v1.5.1"
   git tag v1.5.1 && git push --tags
   ```
3. The GitHub Actions workflow builds all binaries and publishes a release.

## IDE Integration

The `CMakeLists.txt` at the repo root enables basic IDE integration for editors that support CMake (CLion, VS Code with CMake Tools, etc.). Use the Makefile targets through CMake:

```bash
cmake -B build
cmake --build build --target build
cmake --build build --target test
```
