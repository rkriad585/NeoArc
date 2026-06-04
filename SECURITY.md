# Security Policy

## Supported Versions

| Version | Supported          |
|---------|--------------------|
| 1.x     | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

NeoArc takes security seriously. If you discover a security vulnerability, please report it privately before disclosing it publicly.

**Do not** report security vulnerabilities via public GitHub issues. Instead, email:

**rkriad585@gmail.com**

You will receive a response within 48 hours acknowledging the report. We will keep you informed of the progress toward a fix and may request additional information.

### What to include

- A clear description of the vulnerability
- Steps to reproduce it
- Affected versions
- Any potential impact or exploit scenarios

### Scope

The following are in scope:
- The NeoArc server (Flask web application)
- The NeoArc CLI (Go client)
- The build and release pipeline

Out of scope:
- Third-party dependencies (report those to their respective maintainers)
- Infrastructure operated by third parties

## Security Best Practices

### For Production Deployments

1. **Use strong environment variables:**
   - Set a unique `NEOARC_SECRET_KEY` for session signing
   - Set a strong `NEOARC_API_TOKEN` for CLI authentication
   - Set a strong `NEOARC_ADMIN_PASSWORD`

2. **Run behind a reverse proxy** (nginx, Caddy) for TLS termination, rate limiting, and access logging.

3. **Keep dependencies updated:**
   ```bash
   cd neoarc-server && pip list --outdated
   cd neoarc-cli && go list -u -m all
   ```

4. **Restrict database file permissions** — the SQLite DB contains hashed passwords and session data.

5. **Use the admin panel sparingly** — change the admin route prefix from the default `/admin` via `NEOARC_ADMIN_ROUTE`.

### For CLI Users

1. Always verify the alias content before executing, especially when using `--yes` to skip the trust prompt.
2. The API token is stored in plaintext in `~/.config/neostore/neoarc/config.toml` — protect this file.
3. Only use `--install` and `--selfuninstall` from trusted installer scripts.
