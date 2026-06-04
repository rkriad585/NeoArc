# Docker Deployment Guide

NeoArc provides Docker support for both the server and CLI components.

## Prerequisites

- Docker Engine 24+ (or Docker Desktop)
- Docker Compose v2+

## Server (Production Stack)

### Quick Start with Compose

```bash
# From the repo root
docker compose build
docker compose up -d
```

This starts the NeoArc server with:
- SQLite database persisted in a `neoarc-data` volume
- Avatar uploads stored in a `neoarc-uploads` volume
- Configurable via environment variables

### Required Environment Variables

Set these in a `.env` file in the repo root or export them in your shell:

```bash
NEOARC_SECRET_KEY=your-secret-key
NEOARC_API_TOKEN=your-api-token
NEOARC_ADMIN_PASSWORD=your-admin-password
```

### Standalone Server Container

```bash
cd neoarc-server
docker build -t neoarc-server .
docker run -d \
  --name neoarc-server \
  -p 59248:59248 \
  -v neoarc-data:/app/data \
  -v neoarc-uploads:/app/static/uploads \
  -e NEOARC_SECRET_KEY=... \
  -e NEOARC_API_TOKEN=... \
  -e NEOARC_ADMIN_PASSWORD=... \
  neoarc-server
```

### Configuration Reference

| Variable | Default | Description |
|----------|---------|-------------|
| `NEOARC_PORT` | `59248` | Server listen port |
| `NEOARC_HOST` | `0.0.0.0` | Server bind address |
| `NEOARC_DEBUG` | `False` | Enable Flask debug mode |
| `NEOARC_SECRET_KEY` | auto-generated | Flask session signing key |
| `NEOARC_API_TOKEN` | auto-generated | Bearer token for CLI auth |
| `NEOARC_DB_PATH` | `/app/data/neoarc.db` | SQLite database path |
| `NEOARC_UPLOAD_FOLDER` | `/app/static/uploads` | Avatar upload directory |
| `NEOARC_MAX_UPLOAD_MB` | `5` | Max upload size in MB |
| `NEOARC_ADMIN_EMAIL` | `admin@neoarc.local` | Admin login email |
| `NEOARC_ADMIN_PASSWORD` | `change_me` | Admin password |
| `NEOARC_ADMIN_ROUTE` | `/admin` | Admin panel URL prefix |

## CLI (Docker Build)

Build the CLI binary inside a Docker container (useful for CI or if you don't have Go installed):

```bash
# From the repo root
docker build -t neoarc-cli .
docker run --rm neoarc-cli --help
```

### Build Arguments

| Argument | Default | Description |
|----------|---------|-------------|
| `VERSION` | `v0.0.0` | Version string (injected via ldflags) |
| `COMMIT` | `unknown` | Git commit SHA |
| `BUILD_TIME` | `unknown` | Build timestamp |

### Extract the Binary

```bash
docker build -t neoarc-cli .
docker create --name tmp neoarc-cli
docker cp tmp:/usr/local/bin/neoarc ./neoarc
docker rm tmp
```

## Common Tasks

### View Logs

```bash
docker compose logs -f
```

### Stop and Remove

```bash
docker compose down
docker volume rm neoarc_neoarc-data neoarc_neoarc-uploads
```

### Rebuild After Changes

```bash
docker compose build --no-cache
docker compose up -d
```
