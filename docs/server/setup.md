# NeoArc Web Server Setup Guide

This guide provides detailed instructions for setting up and running the NeoArc Web Server.

## Prerequisites

Before you begin, ensure you have the following installed on your system:

*   **Python 3.8+**: The server is built with Python Flask.
    *   [Download Python](https://www.python.org/downloads/)
*   **`pip` or `uv`**: Python package installer. `uv` is a faster alternative.
    *   `pip` usually comes with Python.
    *   [Install uv](https://astral.sh/blog/uv-the-fast-python-package-installer-and-resolver)

## Installation Steps

1.  **Clone the Repository (if you haven't already):**
    ```bash
    git clone https://github.com/rkriad585/neoarc-cli.git # Assuming this is the repo
    cd neoarc-cli/neoarc-server
    ```

2.  **Navigate to the Server Directory:**
    ```bash
    cd neoarc-server
    ```

3.  **Create and Activate a Virtual Environment (Recommended):**
    Using a virtual environment isolates your project's dependencies from your system's global Python packages.

    *   **Windows:**
        ```bash
        python -m venv .venv
        .venv\Scripts\activate
        ```
    *   **Linux/macOS:**
        ```bash
        python3 -m venv .venv
        source .venv/bin/activate
        ```

4.  **Install Dependencies:**
    Install the required Python packages using `pip` or `uv`.

    *   **Using `pip`:**
        ```bash
        pip install flask werkzeug
        ```
    *   **Using `uv`:**
        ```bash
        uv pip install flask werkzeug
        ```

5.  **Configure Admin Credentials:**
    Open the `neoarc-server/config.py` file in your text editor. Locate the `ADMIN_CREDENTIALS` dictionary and update the `email` and `password` fields with your desired administrator login details.

    ```python
    # neoarc-server/config.py
    
    # ... other configurations ...
    
    ADMIN_CREDENTIALS = {
        "email": "your_admin_email@example.com", # <--- CHANGE THIS
        "password": "your_secure_admin_password" # <--- CHANGE THIS
    }
    
    # ... rest of the file ...
    ```
    **Important:** Change `SECRET_KEY` in `config.py` to a strong, unique value for production environments.

6.  **Run the Server:**
    Execute the `main.py` file to start the Flask development server.

    ```bash
    python main.py
    ```

    You should see output indicating that the Flask development server is running. By default, it will be accessible at `http://localhost:5000` (or the port specified in `config.py`).

    ```
     * Serving Flask app 'main'
     * Debug mode: on
    WARNING: This is a development server. Do not use it in a production deployment. Use a production WSGI server instead.
     * Running on http://0.0.0.0:59248 (Press CTRL+C to quit)
     * Restarting with stat
     * Debugger is active!
     * Debugger PIN: XXX-XXX-XXX
    ```
    (Note: The port might be `59248` as per `config.py` in the provided code, not `5000`.)

## Accessing the Web Server

Once the server is running:

*   **User Interface:** Open your web browser and navigate to `http://localhost:59248` (or the `HOST:PORT` configured in `config.py`). You will be presented with the login/registration page.
*   **Admin Panel:** Access the Super Admin dashboard by navigating to `http://localhost:59248/admin/login`. Use the credentials you set in `config.py`.

## Database Initialization

The `init_db()` function in `core/db.py` is called automatically when `main.py` starts if the `neoarc.db` file does not exist. This ensures that the necessary `users` and `aliases` tables are created.

## Next Steps

*   **Register a User:** Create a new user account through the web interface.
*   **Create Aliases:** Start defining your commands and scripts in the dashboard.
*   **Set up the CLI Client:** Proceed to the [Client Setup Guide](docs/client/setup.md) to configure your local execution environment.

Written by [Neorwc](https://github.com/rkriad585/neorwc-cli), Created by RK Riad Khan