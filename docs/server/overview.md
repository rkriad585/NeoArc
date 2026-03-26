# Server-Side Overview (Python Flask)

The NeoArc Web Server, codenamed "The Matrix," is built using the Python Flask microframework. It provides the web interface for users to manage aliases, authenticate, and interact with the system, as well as the API endpoint for the CLI client.

## Core Components

### `main.py` - Application Entry Point
This file initializes the Flask application, registers blueprints, configures settings, and defines core routes and error handlers.

*   **Flask App Initialization:** `app = Flask(__name__)`
*   **Blueprint Registration:** The `admin_bp` (for the admin panel) is registered with a configurable URL prefix (`config.ADMIN_ROUTE`).
*   **Secret Key:** `app.secret_key = config.SECRET_KEY` is crucial for session management and security.
*   **Upload Folder:** Configures `app.config['UPLOAD_FOLDER']` for profile pictures and ensures the directory exists.
*   **Database Initialization:** Calls `init_db()` from `core.db` to ensure the SQLite database and its tables are set up.
*   **Error Handlers:** Custom `404` (page not found) and `500` (internal server error) pages are rendered using `error.html` for a consistent UI.
*   **Main Routes:**
    *   `/`: Redirects to `dashboard` if logged in, otherwise shows `login.html`.
    *   `/login` (POST): Handles user authentication.
    *   `/register` (GET, POST): User registration.
    *   `/dashboard` (GET, POST): User's personal alias management and creation.
    *   `/profile` (GET, POST): User profile management, including avatar upload.
    *   `/search` (GET): Global search for public aliases.
    *   `/view/<alias_name>` (GET): Displays details of a specific alias.
    *   `/delete_alias/<int:id>` (GET): Deletes a user's alias.
    *   `/logout` (GET): Clears the user session.
    *   `/api/alias/<alias_name>` (GET): The API endpoint for the CLI client.

### `config.py` - Configuration Management
This file centralizes all configurable parameters for the Flask application.

*   **Server Settings:** `PORT`, `HOST`, `DEBUG`.
*   **File Uploads:** `UPLOAD_FOLDER`, `ALLOWED_EXTENSIONS` for profile pictures.
*   **Database Path:** `DB_PATH` for the SQLite database.
*   **Admin Route:** `ADMIN_ROUTE` prefix for the admin panel.
*   **Security:** `SECRET_KEY` (critical, should be changed in production).
*   **Admin Credentials:** `ADMIN_CREDENTIALS` dictionary for super admin login.

### `core/db.py` - Database Interaction Layer
Provides functions to interact with the SQLite database.

*   **`get_db()`:** Establishes a connection to `neoarc.db` and sets `row_factory` to `sqlite3.Row` for dictionary-like access to query results.
*   **`init_db()`:** Creates the `users` and `aliases` tables if they don't already exist.

### `core/admin.py` - Admin Panel Blueprint
This module defines the routes and logic for the Super Admin panel. It's registered as a Flask Blueprint to keep admin-related functionality separate.

*   **`admin_bp = Blueprint('admin', __name__)`**: Defines the blueprint.
*   **`require_admin()`**: A `before_request` handler that ensures only authenticated administrators can access admin routes (except `/admin/login`).
*   **Admin Routes:**
    *   `/admin/login` (GET, POST): Admin authentication.
    *   `/admin/` (Dashboard): Lists all users and their alias counts.
    *   `/admin/user/<int:user_id>`: Displays a specific user's details and aliases.
    *   `/admin/block/<int:user_id>`: Toggles a user's `is_blocked` status.
    *   `/admin/delete/<int:user_id>`: Deletes a user and all their associated aliases.
    *   `/admin/delete_alias/<int:alias_id>/<int:user_id>`: Deletes a specific alias from a user.
    *   `/admin/logout`: Clears the admin session.

## Templating and UI

NeoArc uses Jinja2 for server-side templating, rendering HTML pages with dynamic data.
*   **`templates/base.html`**: The main layout template, including TailwindCSS, jQuery, Highlight.js, and custom styling. It defines the overall look and feel, navigation, and common scripts.
*   **`templates/` (User-facing):** `login.html`, `register.html`, `dashboard.html`, `profile.html`, `search.html`, `view.html`, `error.html`. These implement the "Liquid Glass" UI.
*   **`templates/admin/` (Admin-facing):** `login.html`, `dashboard.html`, `view_user.html`. These provide a simpler, functional UI for administrative tasks.

## Static Files

The `static/` directory holds assets like the logo (`logo.svg`) and uploaded profile pictures (`static/uploads/`).

Written by [Neorwc](https://github.com/rkriad585/neorwc-cli), Created by RK Riad Khan