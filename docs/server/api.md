# Server API

The NeoArc Web Server exposes two API endpoints for the NeoArc CLI client to retrieve alias commands.

## Authentication

All API requests require a Bearer token in the `Authorization` header. The token is configured via `NEOARC_API_TOKEN` in your `.env` file (auto-generated if unset).

## Endpoint

### GET /api/alias/<alias_name>

Fetches the command and execution type for a specified alias.

**Path Parameters:**
- `alias_name` (string, required) - The unique name of the alias. Whitespace is automatically trimmed.

**Headers:**
- `Authorization: Bearer <token>` (required)
- `If-None-Match: <etag>` (optional, for conditional requests)

**Rate Limit:** 60 requests per minute per IP address.

**Success Response (200 OK):**
```json
{
    "success": true,
    "alias": "sys_info",
    "command": "uname -a && lsb_release -a",
    "exec_type": "bash"
}
```

**Not Modified (304):**
Returned when the `If-None-Match` header matches the current ETag. Empty body.

**Error Responses:**

| Status | Condition | Body |
|---|---|---|
| 401 | Missing or malformed Authorization header | `{"success": false, "message": "Missing or invalid Authorization header"}` |
| 403 | Invalid API token | `{"success": false, "message": "Invalid API token"}` |
| 404 | Alias not found in database | `{"success": false, "message": "Alias not found"}` |
| 429 | Rate limit exceeded | `{"success": false, "message": "Rate limit exceeded"}` |
| 500 | Internal server error | `{"success": false, "message": "Internal server error"}` or `{"success": false, "message": "Database error"}` |

## Caching

Responses include an `ETag` header computed as an MD5 hash of the response JSON. The `AliasCache` layer caches results in memory for 30 seconds. The cache is invalidated when an alias is created, updated, or deleted.

### GET /api/aliases

Returns a list of all alias names. Protected by Bearer token authentication.

**Headers:**
- `Authorization: Bearer <token>` (required)

**Rate Limit:** 60 requests per minute per IP address.

**Success Response (200 OK):**
```json
{
    "success": true,
    "aliases": ["deploy", "sys_info", "greet"]
}
```

This endpoint is used by the CLI's shell completion generators (`neoarc completion bash|zsh|powershell`) and the hidden `neoarc _list_aliases` command.

## CLI Usage Examples

```bash
curl -H "Authorization: Bearer <token>" http://localhost:59248/api/alias/my_alias

curl -H "Authorization: Bearer <token>" http://localhost:59248/api/aliases
```

With ETag support:
```bash
curl -H "Authorization: Bearer <token>" \
     -H "If-None-Match: <etag>" \
     http://localhost:59248/api/alias/my_alias
```

Written by [Neorwc](https://github.com/rkriad585/neorwc-cli), Created by RK Riad Khan
