# NeoArc Server

Flask backend for the NeoArc command obfuscation and alias execution system. Users register, create command aliases with execution types, and a Go CLI fetches and executes them remotely.

## Features

- **User authentication** — Registration, login, session management with CSRF protection
- **Alias management** — Create, edit, delete, and search command aliases with custom execution types
- **REST API** — Token-protected endpoints for CLI clients to fetch aliases by name and list all alias names
- **Profile management** — Avatar upload with magic-byte validation, editable user details
- **Rate limiting** — In-memory throttling for login (5/min), registration (3/5min), and API (60/min)
- **Caching** — 30-second TTL alias cache with MD5 ETags for conditional requests
- **Admin panel** — User management (view, block, delete) behind static credential login
- **Input validation** — Regex-based email/alias checks, password strength rules, image MIME detection
- **Cyber-themed UI** — Tailwind CSS templates with flash messages, syntax highlighting, and responsive design
- **SQLite storage** — WAL mode, foreign keys, busy timeout for safe concurrent access

## Project Structure

```
neoarc-server/
├── main.py                  # Flask app entry point
├── wsgi.py                  # Production entry point (waitress)
├── config.py                # Environment-based configuration
├── pyproject.toml           # Project metadata and dependencies
├── Dockerfile               # Multi-stage Docker build
├── .env.example             # Environment variable template
├── neoarc.db                # SQLite database (created at runtime)
├── static/
│   └── uploads/             # User avatar uploads
├── templates/
│   ├── base.html            # Layout with navbar, search, flash messages
│   ├── login.html           # User login page
│   ├── register.html        # User registration page
│   ├── dashboard.html       # Alias dashboard
│   ├── edit_alias.html      # Create/edit alias form
│   ├── view.html            # Alias detail view
│   ├── search.html          # Search results
│   ├── profile.html         # User profile with avatar upload
│   ├── error.html           # Error page template
│   └── admin/
│       ├── login.html       # Admin login page
│       ├── dashboard.html   # Admin user list
│       └── view_user.html   # Admin user detail
├── core/
│   ├── __init__.py
│   ├── db.py                # SQLite initialization and schema
│   ├── helpers.py           # g-scoped DB connection, execute, commit
│   ├── auth.py              # CSRF tokens, API token auth, session helpers
│   ├── admin.py             # Admin blueprint (login, dashboard, user ops)
│   ├── rate_limiter.py      # In-memory rate limiter
│   ├── cache.py             # AliasCache with TTL and ETag generation
│   ├── validation.py        # Email, alias, password, image validation
│   └── routes/
│       ├── __init__.py      # Blueprint registration
│       ├── auth_routes.py   # /login, /register, /logout
│       ├── alias_routes.py  # /dashboard, /edit, /delete, /search, /view
│       ├── profile_routes.py# /profile with avatar upload
│       └── api_routes.py    # /api/alias/<name>, /api/aliases
└── tests/
    ├── __init__.py
    ├── conftest.py          # Temp DB, env setup, per-test cleanup
    ├── test_routes.py       # 37 route tests (auth, aliases, profile, isolation)
    ├── test_api.py          # 9 API tests (token, ETag, rate limiting)
    └── test_security.py     # 11 security tests (CSRF, rate limits, validation, blocks)
```

## Quick Start

### Prerequisites

- Python 3.10 or later
- pip

### Installation

```bash
git clone <repo> && cd neoarc-server
python -m venv .venv
source .venv/bin/activate    # or .venv\Scripts\activate on Windows
pip install .
```

### Configuration

Copy the environment template and edit the values:

```bash
cp .env.example .env
```

At minimum, set `NEOARC_SECRET_KEY` to a random string. The server will auto-generate one if omitted (with a warning).

### Run in Development

```bash
python main.py
```

The server starts by default on `http://0.0.0.0:59248`.

### Run in Production

```bash
python wsgi.py
```

Uses [waitress](https://docs.pylonsproject.org/projects/waitress/) as the production WSGI server.

## API Reference

### `GET /api/alias/<name>`

Returns a command alias by its name. Protected by Bearer token authentication.

**Headers:**

| Header | Required | Description |
|--------|----------|-------------|
| `Authorization` | Yes | `Bearer <token>` (the `NEOARC_API_TOKEN` value) |

**Response (200):**

```json
{
  "success": true,
  "alias": {
    "alias_name": "myalias",
    "command": "ls -la",
    "exec_type": "shell"
  }
}
```

**Response (404):**

```json
{
  "success": false,
  "message": "Alias not found"
}
```

**Response (401):**

```json
{
  "success": false,
  "message": "Unauthorized"
}
```

**Caching:**

Responses are cached in memory for 30 seconds. The server returns an `ETag` header (MD5 of the response body). Clients may send a `If-None-Match` header; the server responds with `304 Not Modified` when the cache is still fresh.

**Rate limiting:**

This endpoint is throttled at 60 requests per minute per client IP.

### Example (curl)

```bash
curl -H "Authorization: Bearer $NEOARC_API_TOKEN" http://localhost:59248/api/alias/myalias
```

### `GET /api/aliases`

Returns a list of all alias names. Protected by Bearer token authentication.

**Headers:**

| Header | Required | Description |
|--------|----------|-------------|
| `Authorization` | Yes | `Bearer <token>` (the `NEOARC_API_TOKEN` value) |

**Response (200):**

```json
{
  "success": true,
  "aliases": ["myalias", "deploy", "sys_info"]
}
```

This endpoint is used by shell completion generators (`neoarc completion bash|zsh|powershell`) and the hidden `neoarc _list_aliases` command.

## Configuration Reference

All configuration is via environment variables (or a `.env` file in the project root).

| Variable | Default | Description |
|----------|---------|-------------|
| `NEOARC_SECRET_KEY` | auto-generated | Flask session signing key |
| `NEOARC_PORT` | `59248` | Server listen port |
| `NEOARC_HOST` | `0.0.0.0` | Server bind address |
| `NEOARC_DEBUG` | `False` | Enable Flask debug mode |
| `NEOARC_DB_PATH` | `neoarc.db` | SQLite database file path |
| `NEOARC_UPLOAD_FOLDER` | `static/uploads` | Avatar upload directory |
| `NEOARC_MAX_UPLOAD_MB` | `5` | Maximum upload size in MB |
| `NEOARC_ADMIN_ROUTE` | `/admin` | Admin panel URL prefix |
| `NEOARC_ADMIN_EMAIL` | `admin@neoarc.local` | Admin login email |
| `NEOARC_ADMIN_PASSWORD` | `change_me` | Admin login password (used when `NEOARC_ADMIN_PASSWORD_HASH` is not set) |
| `NEOARC_ADMIN_PASSWORD_HASH` | *none* | Pre-computed Werkzeug password hash (overrides `NEOARC_ADMIN_PASSWORD`) |
| `NEOARC_API_TOKEN` | auto-generated | Bearer token for API authentication |

## Docker Deployment

```bash
docker build -t neoarc-server .
docker run -d \
  -p 59248:59248 \
  -v neoarc-data:/app/neoarc.db \
  -e NEOARC_SECRET_KEY=your-secret \
  -e NEOARC_API_TOKEN=your-api-token \
  -e NEOARC_ADMIN_PASSWORD=your-admin-password \
  neoarc-server
```

The Dockerfile uses a multi-stage build:
1. **Builder stage** — pip installs dependencies into a dist directory
2. **Runtime stage** — `python:3.12-alpine` with gcc (for waitress), copies dependencies and application code, runs `python wsgi.py` on port 59248

## Running Tests

66 tests covering routes, API, and security.

```bash
# Install with test dependencies (currently pytest and its ecosystem)
pip install pytest pytest-flask

# Run all tests
pytest tests/

# Run with verbose output
pytest tests/ -v

# Run a specific test file
pytest tests/test_api.py
```

Test categories:

| File | Count | Coverage |
|------|-------|----------|
| `test_routes.py` | 46 | Registration, login, alias CRUD, search, profile, user isolation, password reset |
| `test_api.py` | 9 | Token auth, ETag caching, rate limiting |
| `test_security.py` | 11 | CSRF protection, rate limits, input validation, blocked user handling |

Tests use a temporary SQLite database, set required environment variables in `conftest.py`, and perform per-test cleanup.

## Architecture Notes

- **Flask blueprints** organize routes into auth, alias, profile, and API modules, plus a separate admin blueprint
- **SQLite with WAL mode** — Write-Ahead Logging enables concurrent reads without blocking writes; `busy_timeout=5000` prevents locked-database errors under light concurrency
- **CSRF + API token** — Web forms use session-based CSRF tokens; the API uses a static Bearer token (suitable for CLI clients, not user-facing mobile/web apps)
- **Rate limiting** is purely in-memory (no Redis dependency) and resets on server restart — adequate for single-instance deployments
- **Caching layer** wraps alias lookups with a 30-second TTL; ETags let well-behaved clients avoid unnecessary transfers
- **Image validation** uses magic bytes (not file extension) to detect PNG, JPEG, GIF, and WebP — prevents MIME-type spoofing on avatar uploads
- **Session lifetime** is 24 hours with `PERMANENT_SESSION_LIFETIME`; each request renews the session via `make_session_permanent`

## License

MIT
