# NeoArc Web Server Setup Guide

## Prerequisites

- **Python 3.10+** (required by `pyproject.toml`)
- **pip** (Python package installer)

## Installation

1.  **Navigate to the Server Directory:**
    ```bash
    cd neoarc-server
    ```

2.  **Create and Activate a Virtual Environment (Recommended):**
    ```bash
    python -m venv .venv
    .venv\Scripts\activate   # Windows
    source .venv/bin/activate  # Linux/macOS
    ```

3.  **Install the Package:**
    Installs Flask, Werkzeug, python-dotenv, and waitress from `pyproject.toml`.
    ```bash
    pip install .
    ```

4.  **Configure Environment Variables:**
    Copy the example env file and fill in your values.
    ```bash
    copy .env.example .env       # Windows
    cp .env.example .env         # Linux/macOS
    ```

    Key variables to review:
    - `NEOARC_SECRET_KEY` - Long random string for session signing. Generate with:
      `python -c "import secrets; print(secrets.token_hex(32))"`
    - `NEOARC_PORT=59248` - Server port
    - `NEOARC_ADMIN_PASSWORD` - Admin panel password (default: `change_me`)
    - `NEOARC_API_TOKEN` - Token for CLI authentication. Generate with:
      `python -c "import secrets; print(secrets.token_hex(32))"`

## Running the Server

### Development
```bash
python main.py
```
Starts the Flask development server on port 59248 with debug mode.

### Production
```bash
python wsgi.py
```
Starts a Waitress production server on port 59248.

## Docker

```bash
docker build -t neoarc-server .
docker run -p 59248:59248 neoarc-server
```

The Docker image uses `python:3.12-alpine`, exposes port 59248, and runs `wsgi.py` by default.

## Running Tests

```bash
pytest tests/ -v
```

Runs the test suite (57 tests) against an isolated temporary SQLite database. Rate limiters are cleared per test via the `client` fixture.

## Accessing the Server

- **Web Interface:** `http://localhost:59248`
- **Admin Panel:** `http://localhost:59248/admin/login`
- **API:** `http://localhost:59248/api/alias/<alias_name>`

Written by [Neorwc](https://github.com/rkriad585/neorwc-cli), Created by RK Riad Khan
