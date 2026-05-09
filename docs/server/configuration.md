# Server Configuration

All configuration is loaded from environment variables prefixed with `NEOARC_` via `python-dotenv`. A `.env` file in the `neoarc-server/` directory is loaded automatically. Missing tokens/keys with no env var set are auto-generated with a warning.

## Server Settings

### `NEOARC_PORT`
- **Description:** The port number the server listens on.
- **Type:** `int`
- **Default:** `59248`
- **Env:** `NEOARC_PORT=59248`

### `NEOARC_HOST`
- **Description:** The host address to bind to. `0.0.0.0` exposes to the network.
- **Type:** `str`
- **Default:** `0.0.0.0`
- **Env:** `NEOARC_HOST=0.0.0.0`

### `NEOARC_DEBUG`
- **Description:** Enables Flask debug mode. Must be `False` in production.
- **Type:** `bool`
- **Default:** `False`
- **Env:** `NEOARC_DEBUG=False`

## Security

### `NEOARC_SECRET_KEY`
- **Description:** Cryptographic key for signing session cookies. **Must be a long random string in production.**
- **Type:** `str`
- **Default:** Auto-generated via `secrets.token_hex(32)` if unset. A warning is emitted.
- **Env:** `NEOARC_SECRET_KEY=<64-char-hex>`
- **Generate:** `python -c "import secrets; print(secrets.token_hex(32))"`
- **Note:** If auto-generated, sessions will be invalidated on every server restart.

### `NEOARC_ADMIN_EMAIL`
- **Description:** Email address for the admin panel login.
- **Type:** `str`
- **Default:** `admin@neoarc.local`
- **Env:** `NEOARC_ADMIN_EMAIL=admin@neoarc.local`

### `NEOARC_ADMIN_PASSWORD`
- **Description:** Plain-text password for the admin panel. Automatically hashed at startup via `werkzeug.security.generate_password_hash()`.
- **Type:** `str`
- **Default:** `change_me`
- **Env:** `NEOARC_ADMIN_PASSWORD=change_me`
- **Warning:** A warning is emitted if the default is used.

### `NEOARC_ADMIN_PASSWORD_HASH`
- **Description:** Pre-computed password hash. If set, overrides `NEOARC_ADMIN_PASSWORD`.
- **Type:** `str`
- **Default:** Auto-generated from `NEOARC_ADMIN_PASSWORD` if unset.
- **Env:** `NEOARC_ADMIN_PASSWORD_HASH=<werkzeug-hash>`

### `NEOARC_API_TOKEN`
- **Description:** Bearer token for API authentication (used by the NeoArc CLI client). **Must be a long random string in production.**
- **Type:** `str`
- **Default:** Auto-generated via `secrets.token_hex(32)` if unset. A warning is emitted.
- **Env:** `NEOARC_API_TOKEN=<64-char-hex>`
- **Generate:** `python -c "import secrets; print(secrets.token_hex(32))"`
- **Note:** Build scripts inject this token into the neoarc binary at compile time.

## Database

### `NEOARC_DB_PATH`
- **Description:** Path to the SQLite database file, relative to `neoarc-server/`.
- **Type:** `str`
- **Default:** `neoarc.db`
- **Env:** `NEOARC_DB_PATH=neoarc.db`

## File Uploads

### `NEOARC_UPLOAD_FOLDER`
- **Description:** Directory for uploaded profile pictures, relative to `neoarc-server/`.
- **Type:** `str`
- **Default:** `static/uploads`
- **Env:** `NEOARC_UPLOAD_FOLDER=static/uploads`

### `NEOARC_MAX_UPLOAD_MB`
- **Description:** Maximum upload file size in megabytes.
- **Type:** `int`
- **Default:** `5`
- **Env:** `NEOARC_MAX_UPLOAD_MB=5`
- **Note:** Enforced by Flask's `MAX_CONTENT_LENGTH`.

### `ALLOWED_EXTENSIONS`
- **Description:** Allowed file extensions for uploads (not configurable via env var).
- **Type:** `set` of `str`
- **Hardcoded:** `{png, jpg, jpeg, gif, webp}`
- **Note:** `validate_image()` also performs content-based MIME detection via magic bytes, not just extension checking.

## Routing

### `NEOARC_ADMIN_ROUTE`
- **Description:** URL prefix for the admin panel. Changing this adds obscurity.
- **Type:** `str`
- **Default:** `/admin`
- **Env:** `NEOARC_ADMIN_ROUTE=/admin`

## Production Best Practices

- Set `NEOARC_SECRET_KEY` to a long random hex string (64 chars).
- Set `NEOARC_API_TOKEN` to a long random hex string (64 chars).
- Change `NEOARC_ADMIN_PASSWORD` from the default `change_me`.
- Set `NEOARC_DEBUG=False`.
- Never commit `.env` to version control.

Written by [Neorwc](https://github.com/rkriad585/neorwc-cli), Created by RK Riad Khan
