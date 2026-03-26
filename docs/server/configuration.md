# Server Configuration (`config.py`)

The `config.py` file centralizes all configurable parameters for the NeoArc server. It's essential to review and customize these settings, especially for production deployments, to ensure security and optimal performance.

## Configuration Parameters

### `PORT`
-   **Description:** The port number on which the Flask server will listen for incoming connections.
-   **Type:** `int`
-   **Default:** `59248`
-   **Example:** `PORT = 8000`

### `HOST`
-   **Description:** The host address the server will bind to.
    -   `"0.0.0.0"` makes the server accessible from any IP address on the network.
    -   `"127.0.0.1"` or `"localhost"` restricts access to the local machine only.
-   **Type:** `str`
-   **Default:** `"0.0.0.0"`
-   **Example:** `HOST = "127.0.0.1"`

### `DEBUG`
-   **Description:** Enables or disables Flask's debug mode.
    -   `True`: Activates the debugger and reloader, providing detailed error messages and automatic code reloading on changes. **Should always be `False` in production.**
    -   `False`: Disables debug features.
-   **Type:** `bool`
-   **Default:** `True`
-   **Example:** `DEBUG = False`

### `UPLOAD_FOLDER`
-   **Description:** The directory where uploaded files (e.g., user profile pictures) will be stored. This path is relative to the `neoarc-server` directory.
-   **Type:** `str`
-   **Default:** `'static/uploads'`
-   **Example:** `UPLOAD_FOLDER = 'data/uploads'`

### `ALLOWED_EXTENSIONS`
-   **Description:** A set of allowed file extensions for uploads. This helps prevent users from uploading potentially malicious file types.
-   **Type:** `set` of `str`
-   **Default:** `{'png', 'jpg', 'jpeg', 'gif', 'webp'}`
-   **Example:** `ALLOWED_EXTENSIONS = {'png', 'jpg'}`

### `DB_PATH`
-   **Description:** The path to the SQLite database file. This path is relative to the `neoarc-server` directory.
-   **Type:** `str`
-   **Default:** `'neoarc.db'`
-   **Example:** `DB_PATH = 'data/neoarc_prod.db'`

### `ADMIN_ROUTE`
-   **Description:** The URL prefix for the Super Admin panel. Changing this can add a layer of obscurity.
-   **Type:** `str`
-   **Default:** `'/admin'`
-   **Example:** `ADMIN_ROUTE = '/control_panel_xyz'`

### `SECRET_KEY`
-   **Description:** A cryptographic key used by Flask to sign session cookies and other security-related operations.
-   **Type:** `str`
-   **Default:** `'neoarc_super_secret_key_change_me'`
-   **Security Warning:** **THIS MUST BE CHANGED TO A LONG, RANDOM, AND COMPLEX STRING IN PRODUCTION.** Failure to do so will compromise the security of user sessions.
-   **Example:** `SECRET_KEY = 'your_very_long_and_random_secret_key_here_1234567890abcdef'`

### `ADMIN_CREDENTIALS`
-   **Description:** A dictionary containing the email and password for the Super Admin account. These credentials are used to log into the `/admin` panel.
-   **Type:** `dict`
-   **Default:**
    ```python
    ADMIN_CREDENTIALS = {
        "email": "rkriad585@gamil.com",
        "password": "riad"
    }
    ```
-   **Security Warning:** **THESE CREDENTIALS MUST BE CHANGED TO SECURE VALUES IN PRODUCTION.** Do not use default or easily guessable passwords.
-   **Example:**
    ```python
    ADMIN_CREDENTIALS = {
        "email": "admin@yourdomain.com",
        "password": "a_very_strong_admin_password_123!"
    }
    ```

## Best Practices for Production

-   **Change Defaults:** Always change `SECRET_KEY` and `ADMIN_CREDENTIALS` to strong, unique values.
-   **Disable Debug Mode:** Set `DEBUG = False` to prevent sensitive information from being exposed.
-   **Secure File Storage:** Ensure `UPLOAD_FOLDER` has appropriate file system permissions to prevent unauthorized access or execution of uploaded files.
-   **Environment Variables:** For highly sensitive information like `SECRET_KEY` and `ADMIN_CREDENTIALS`, consider loading them from environment variables rather than hardcoding them in `config.py`. This prevents them from being committed to version control.

written by Neorwc