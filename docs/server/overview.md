# Server-Side Overview (Python Flask)

The NeoArc Web Server provides the web interface for users to manage aliases, authenticate, and interact with the system, plus an API endpoint for the CLI client. The entire project lives under `neoarc-server/`.

## Core Components

### `main.py` - Application Entry Point

Initializes the Flask application, registers blueprints, configures settings, and defines core routes and error handlers.

- **Flask App Initialization:** `app = Flask(__name__)`
- **Blueprint Registration:** `admin_bp` is registered with a configurable URL prefix (`config.ADMIN_ROUTE`). The remaining blueprints are registered via `register_route_blueprints(app)` which is a factory function in `core/routes/__init__.py`.
- **Secret Key:** `app.secret_key = config.SECRET_KEY` for session management and security.
- **Upload Folder:** Configures `UPLOAD_FOLDER` and `MAX_CONTENT_LENGTH`, ensures the upload directory exists.
- **Permanent Sessions:** A `before_request` handler sets `session.permanent = True` with a 24-hour lifetime.
- **Database Initialization:** Calls `init_db()` from `core.db` to ensure the SQLite database and its tables and indexes are set up.
- **CSRF Token:** `csrf_token` is injected into all Jinja templates via `app.jinja_env.globals`.
- **Error Handlers:** Custom `404`, `500`, and `413` (file too large) error handlers for a consistent cyber-themed UI.
- **Route:** `/` redirects to dashboard if logged in, otherwise shows `login.html`.

### `core/routes/__init__.py` - Blueprint Factory

`register_route_blueprints(app)` imports and registers four blueprints:
- `auth_bp` from `core/routes/auth_routes.py`
- `alias_bp` from `core/routes/alias_routes.py`
- `profile_bp` from `core/routes/profile_routes.py`
- `api_bp` from `core/routes/api_routes.py`

### `core/routes/auth_routes.py` - Authentication Blueprint

Defines `/login` (POST, CSRF-protected), `/register` (GET, POST), and `/logout` (GET). Rate-limited login at 5 attempts/minute, register at 3 attempts/5 minutes.

### `core/routes/alias_routes.py` - Alias Management Blueprint

Defines `/dashboard` (GET, POST), `/edit/<id>` (GET, POST), `/delete/<id>` (GET), `/search` (GET, paginated, 18 per page), `/view/<alias_name>` (GET). Inline CSRF check for POST operations. Invalidates `alias_cache` on create/update/delete.

### `core/routes/profile_routes.py` - Profile Blueprint

Defines `/profile` (GET, POST) for viewing and updating user profile, including avatar upload with content-based MIME detection.

### `core/routes/api_routes.py` - API Blueprint

Defines two endpoints for the CLI client:
- `GET /api/alias/<alias_name>` — Fetch a single alias by name. Protected by `api_token_required`, rate-limited (60/min), uses `AliasCache` with ETag support.
- `GET /api/aliases` — List all alias names (used by shell completion scripts). Protected by `api_token_required`, rate-limited (60/min).

### `core/admin.py` - Admin Panel Blueprint

Defined as `admin_bp` with a configurable URL prefix. Routes:
- `/login` (GET, POST) - Admin authentication against `config.ADMIN_EMAIL` and `config.ADMIN_PASSWORD_HASH`
- `/` - Dashboard listing all users with alias counts
- `/user/<int:user_id>` - View user details and aliases
- `/block/<int:user_id>` - Toggle blocked status
- `/delete/<int:user_id>` - Delete user and aliases
- `/delete_alias/<int:alias_id>/<int:user_id>` - Delete specific alias
- `/logout` - Clear admin session

### `core/auth.py` - Security Utilities

- `generate_csrf_token()` - Generates and stores a 32-byte hex CSRF token in the session
- `csrf_required` decorator - Validates CSRF token on POST/PUT/DELETE requests
- `api_token_required` decorator - Validates Bearer token from `Authorization` header against `config.API_TOKEN`
- `regenerate_session()` - Preserves session data while rotating the session ID

### `core/cache.py` - Alias Cache

`AliasCache` class with a 30-second TTL. Uses MD5 hash of sorted JSON as ETag. Exposes `get()`, `set()`, and `invalidate()` methods. A module-level singleton `alias_cache` is used throughout.

### `core/rate_limiter.py` - Rate Limiter

`RateLimiter` class with configurable `max_attempts` and `window_seconds`. Uses a `defaultdict` of timestamps per key. Key limits:
- Login: 5 requests per 60 seconds
- Register: 3 requests per 300 seconds
- API: 60 requests per 60 seconds

### `core/validation.py` - Input Validation

Regex-based validators:
- `validate_email()` - Standard email format check
- `validate_alias_name()` - Alphanumeric, hyphens, underscores (1-100 chars)
- `validate_full_name()` - 1-200 characters
- `validate_password()` - Minimum 8 chars, must not match username/email
- `allowed_file()` - Extension check against `config.ALLOWED_EXTENSIONS`
- `validate_image()` - Content-based MIME detection via magic bytes (PNG, JPG, GIF, WebP)

### `core/helpers.py` - Database Helpers

Flask `g`-scoped connection management:
- `get_conn()` - Returns the database connection from `g._db`, creating it via `get_db()` if needed
- `db_execute(query, params)` - Executes a query with error logging
- `db_commit()` - Commits with error logging

### `core/db.py` - Database Layer

- `get_db()` - Creates a new `sqlite3.connect()` with `row_factory=sqlite3.Row`, enables WAL journal mode, `foreign_keys=ON`, and `busy_timeout=5000`
- `init_db()` - Creates `users` and `aliases` tables plus four indexes (`idx_aliases_alias_name`, `idx_aliases_user_id`, `idx_users_username`, `idx_users_email`), migrates `created_at` column if missing
- `retry_on_locked` decorator - Retries up to 3 times with exponential backoff on "database is locked"

### `config.py` - Configuration Management

Environment-based configuration using `NEOARC_*` environment variables with auto-generated fallbacks:
- Server: `NEOARC_PORT`, `NEOARC_HOST`, `NEOARC_DEBUG`
- Security: `NEOARC_SECRET_KEY` (auto-generated via `secrets.token_hex(32)` if unset, with warning)
- Admin: `NEOARC_ADMIN_EMAIL`, `NEOARC_ADMIN_PASSWORD` / `NEOARC_ADMIN_PASSWORD_HASH` (auto-hashed with warning)
- API: `NEOARC_API_TOKEN` (auto-generated if unset, with warning)
- Storage: `NEOARC_DB_PATH`, `NEOARC_UPLOAD_FOLDER`, `NEOARC_MAX_UPLOAD_MB`
- Routing: `NEOARC_ADMIN_ROUTE`

### `wsgi.py` - Production Server

Starts a Waitress production WSGI server on the configured `HOST` and `PORT`.

## Templating and UI

NeoArc uses Jinja2 for server-side templating with a cyber-themed "Liquid Glass" aesthetic.

**Templates:**

| Template | Description |
|---|---|
| `templates/base.html` | Main layout with Tailwind CSS, DotGothic16 font, highlight.js, jQuery |
| `templates/login.html` | User login page |
| `templates/register.html` | User registration page |
| `templates/dashboard.html` | Alias management dashboard |
| `templates/edit_alias.html` | Edit alias form |
| `templates/view.html` | View alias details with syntax highlighting |
| `templates/search.html` | Global alias search with pagination |
| `templates/profile.html` | User profile management |
| `templates/error.html` | Error pages (404, 500) |
| `templates/admin/login.html` | Admin login page |
| `templates/admin/dashboard.html` | Admin user management dashboard |
| `templates/admin/view_user.html` | Admin view user details |

## Static Files

- `static/logo.svg` - NeoArc logo
- `static/uploads/` - User-uploaded profile pictures

## Test Suite

Tests are located in `tests/` and use `pytest` with a temporary SQLite database:

| File | Tests | Description |
|---|---|---|
| `tests/conftest.py` | Fixtures | Creates isolated temp DB, sets env vars, provides `client` fixture |
| `tests/helpers.py` | Utilities | `register_user()`, `login()`, `create_alias()`, `extract_csrf()` |
| `tests/test_routes.py` | 37 | Route behavior, registration, login, CRUD, search, pagination, authorization |
| `tests/test_api.py` | 9 | API auth, missing/invalid tokens, ETag caching, rate limiting, whitespace trimming |
| `tests/test_security.py` | 11 | CSRF protection, rate limiting, input validation, blocked user access |

**Total: 57 tests**

Written by [Neorwc](https://github.com/rkriad585/neorwc-cli), Created by RK Riad Khan
