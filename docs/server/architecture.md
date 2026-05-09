# Server Architecture

The NeoArc server is built with Python Flask using a modular blueprint-based design.

## Blueprint Architecture

Registration is handled by `register_route_blueprints(app)` in `core/routes/__init__.py`:

```
main.py
  +-- admin_bp          (core/admin.py)          -- /admin/*
  +-- auth_bp           (core/routes/auth_routes.py)    -- /login, /register, /logout
  +-- alias_bp          (core/routes/alias_routes.py)   -- /dashboard, /edit, /delete, /search, /view
  +-- profile_bp        (core/routes/profile_routes.py) -- /profile
  +-- api_bp            (core/routes/api_routes.py)     -- /api/alias/<name>, /api/aliases
```

Each blueprint is independent, with its own route prefix and set of middleware.

## Request Flow

### Web Request (e.g., dashboard, login)

```
Browser Request
  -> Flask Routing (matches blueprint route)
    -> CSRF Check (POST/PUT/DELETE only; `csrf_required` decorator or inline check)
      -> Rate Limiter Check (login: 5/min, register: 3/5min)
        -> Authentication Check (session `user_id` or admin session)
          -> Business Logic (CRUD via `db_execute`/`db_commit`)
            -> Template Rendering (Jinja2 + Tailwind CSS)
              -> HTML Response
```

### API Request (`/api/alias/<name>`)

```
CLI Request
  -> `api_token_required` decorator (validates Bearer token via `secrets.compare_digest`)
    -> Rate Limiter Check (60/min per IP)
      -> AliasCache lookup (30s TTL, MD5 ETag)
        -> Cache HIT with If-None-Match matching    --> 304 Not Modified (empty body)
        -> Cache HIT without matching                --> 200 + ETag header
        -> Cache MISS --> DB query --> AliasCache set --> 200 + ETag header
```

## Session Management

- Sessions are set to `permanent = True` via a `before_request` handler with a 24-hour lifetime (`PERMANENT_SESSION_LIFETIME`).
- On successful login, `regenerate_session()` preserves session data while rotating the session ID to prevent session fixation.
- Logout clears the session entirely.

## Database Layer

- Every `get_db()` call opens a fresh `sqlite3.connect()` connection with per-connection PRAGMAs: WAL mode, `foreign_keys=ON`, `busy_timeout=5000`.
- Connections are scoped to the Flask application context via `g` (`get_conn()` in `core/helpers.py`).
- `db_execute()` and `db_commit()` wrap `get_conn()` with error logging.
- The `retry_on_locked` decorator retries operations 3 times on "database is locked" errors.

## Caching Layer

- `AliasCache` is an in-memory dict-based cache with 30-second TTL.
- ETags are MD5 hashes of the sorted JSON response data.
- Cache is invalidated on alias create, update, or delete operations (in both `alias_routes.py` and `api_routes.py`).

## Security

- **CSRF:** Token generated via `secrets.token_hex(32)`, stored in session, matched with constant-time `secrets.compare_digest()`. Required on all POST/PUT/DELETE web forms.
- **Password Hashing:** `werkzeug.security.generate_password_hash()` / `check_password_hash()` (Werkzeug's default pbkdf2:sha256).
- **API Authentication:** Bearer token verified with constant-time comparison.
- **Session Signing:** Flask session cookies signed with `config.SECRET_KEY`.
- **Input Validation:** Server-side validation of email, alias name, password strength, and image content (magic byte detection).
- **Rate Limiting:** In-memory per-IP rate limiting for login, registration, and API endpoints.

Written by [Neorwc](https://github.com/rkriad585/neorwc-cli), Created by RK Riad Khan
