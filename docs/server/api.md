# Server API Endpoints

The NeoArc Web Server exposes a single public API endpoint primarily for the NeoArc CLI client to fetch alias commands. This API is designed to be simple and efficient.

## 1. Get Alias Command

This endpoint allows the CLI client to retrieve the command and execution type associated with a given alias name.

*   **URL:** `/api/alias/<alias_name>`
*   **Method:** `GET`
*   **Description:** Fetches the command and execution type for a specified alias.
*   **Path Parameters:**
    *   `alias_name` (string, required): The unique name of the alias to retrieve.

*   **Response (JSON):**

    *   **Success (HTTP 200 OK):**
        ```json
        {
            "success": true,
            "alias": "sys_info",
            "command": "uname -a && lsb_release -a",
            "exec_type": "bash"
        }
        ```
        *   `success` (boolean): `true` indicates the alias was found.
        *   `alias` (string): The name of the alias requested.
        *   `command` (string): The actual script or command associated with the alias.
        *   `exec_type` (string): The type of interpreter required to execute the command (e.g., "bash", "powershell", "python", "go", "cmd").

    *   **Alias Not Found (HTTP 404 Not Found):**
        ```json
        {
            "success": false,
            "message": "Alias not found"
        }
        ```
        *   `success` (boolean): `false` indicates the alias was not found in the database.
        *   `message` (string): A descriptive error message.

### Example Usage (CLI Perspective)

When the NeoArc CLI client executes `neoarc sys_info`, it performs an HTTP GET request similar to:

```bash
GET http://localhost:5000/api/alias/sys_info
```

The server then responds with the JSON data, which the CLI parses to execute the `command` using the specified `exec_type`.

Written by [Neorwc](https://github.com/rkriad585/neorwc-cli), Created by RK Riad Khan