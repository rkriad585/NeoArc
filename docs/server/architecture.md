# Server Architecture (The Matrix)

The NeoArc server, codenamed "The Matrix," is built using Python Flask, providing a robust and scalable backend for managing user accounts and alias commands. It follows a modular design, leveraging Flask Blueprints for administrative functionalities and a centralized database for data persistence.

## Core Components

### 1. Flask Application (`main.py`)
-   **Entry Point:** `main.py` serves as the primary entry point for the Flask application. It initializes the Flask app, registers blueprints, and defines core routes for user interaction and API access.
-   **Configuration:** It imports settings from `config.py` for application-wide parameters like `SECRET_KEY`, `UPLOAD_FOLDER`, database path, and admin route prefix.
-   **Database Initialization:** Ensures the SQLite database is initialized and tables are created on startup if they don't exist, using `core.db.init_db()`.
-   **Error Handling:** Custom error handlers for 404 (Not Found) and 500 (Internal Server Error) are implemented, rendering user-friendly error pages.
-   **User Authentication & Authorization:** Integrates `werkzeug.security` for password hashing and `Flask-Bcrypt` (implied by `werkzeug.security` usage in the codebase context, though `Flask-Bcrypt` is explicitly listed in `site-packages`) for secure password management. User sessions are managed using Flask's built-in session mechanism.
-   **File Uploads:** Handles profile picture uploads, ensuring allowed file extensions.

### 2. Blueprints
-   **Modularization:** Flask Blueprints are used to organize different parts of the application. The primary example is `admin_bp` for the Super Admin panel.
-   **`admin_bp` (`core/admin.py`):** This blueprint encapsulates all administrative functionalities, including login, dashboard, user management (view, block, delete), and alias moderation. It has its own set of routes and template directory (`templates/admin`).

### 3. Database Layer (`core/db.py`)
-   **SQLite3:** NeoArc uses SQLite3 as its database, providing a lightweight and file-based solution suitable for development and smaller deployments.
-   **`get_db()`:** A utility function to establish a connection to the database and configure `row_factory` to return rows as dictionary-like objects.
-   **`init_db()`:** Creates the `users` and `aliases` tables if they do not already exist, defining their schema.

### 4. Configuration (`config.py`)
-   **Centralized Settings:** Stores all configurable parameters, including server host/port, debug mode, upload settings, database path, admin route, secret key, and admin credentials.
-   **Security:** Emphasizes the importance of changing `SECRET_KEY` and `ADMIN_CREDENTIALS` for production deployments.

### 5. Templating Engine (Jinja2)
-   **Dynamic Content:** Jinja2 is used for rendering dynamic HTML content.
-   **Template Structure:** Templates are organized in the `templates/` directory, with a `base.html` for common layout and specific templates for different views (e.g., `login.html`, `dashboard.html`, `profile.html`, `search.html`, `view.html`, `error.html`).
-   **Admin Templates:** A dedicated `templates/admin/` subdirectory houses templates for the admin panel.

### 6. Static Files
-   **Assets:** The `static/` directory stores static assets like the logo, favicons, and uploaded user profile pictures.

## Request Flow

1.  **Incoming Request:** A user's browser or the NeoArc CLI sends an HTTP request to the server.
2.  **Flask Routing:** `main.py` (or a registered blueprint) matches the request URL to a defined route.
3.  **Authentication/Authorization (if applicable):** For protected routes (e.g., dashboard, profile, admin panel), session data is checked to verify user login and roles.
4.  **Business Logic:** The corresponding Python function executes, interacting with the database (`core/db.py`) to fetch or store data (e.g., user details, alias commands).
5.  **Template Rendering:** For web requests, data is passed to a Jinja2 template, which generates the HTML response.
6.  **API Response:** For CLI requests (e.g., `/api/alias/<alias_name>`), data is formatted as JSON and returned.
7.  **Response to Client:** The server sends the HTML or JSON response back to the client.

## Security Considerations
-   **Password Hashing:** Passwords are not stored in plaintext; `werkzeug.security.generate_password_hash` and `check_password_hash` are used for secure storage and verification.
-   **Session Management:** Flask's `SECRET_KEY` is crucial for securing session cookies.
-   **Admin Access:** The admin panel is protected by separate credentials defined in `config.py`.
-   **File Uploads:** Basic validation for allowed file extensions is in place to prevent malicious uploads.

written by Neorwc