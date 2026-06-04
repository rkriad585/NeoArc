# NeoArc Server Setup Guide

This guide details the steps to set up and run the NeoArc web server, which acts as the central "Execution Grid" for managing aliases and user accounts.

## Prerequisites

*   **Python 3.8+:** Ensure Python is installed on your system. You can download it from [python.org](https://www.python.org/downloads/).
*   **Internet Connection:** Required to download Python packages.

## Quick Start with Docker

```bash
git clone https://github.com/rkriad585/NeoArc.git
cd NeoArc
docker compose up -d
```

This starts the server on port 59248. Set required env vars in a `.env` file first (see [Configuration](#configuration) below).

## Installation Steps

1.  **Clone the Repository (if not already done):**
    If you haven't already, clone the NeoArc project repository:
    ```bash
    git clone https://github.com/rkriad585/NeoArc.git
    cd NeoArc/neoarc-server
    ```

2.  **Navigate to the Server Directory:**
    Change your current directory to the `neoarc-server` folder:
    ```bash
    cd neoarc-server
    ```

3.  **Create and Activate a Virtual Environment (Recommended):**
    Using a virtual environment isolates your project's dependencies from your global Python installation.

    *   **On Windows:**
        ```powershell
        python -m venv .venv
        .venv\Scripts\activate
        ```
    *   **On Linux/macOS:**
        ```bash
        python3 -m venv .venv
        source .venv/bin/activate
        ```
    You should see `(.venv)` or similar in your terminal prompt, indicating the virtual environment is active.

4.  **Install Dependencies:**
    Install the required Python packages using `pip` (or `uv` if you prefer):
    ```bash
    pip install flask werkzeug
    # Or if you use uv:
    # uv pip install flask werkzeug
    ```
    These packages are essential for the Flask web application.

5.  **Configure Admin Credentials:**
    Open the `neoarc-server/config.py` file in your text editor. Locate the `ADMIN_CREDENTIALS` dictionary and update the `email` and `password` fields with your desired administrator login details.

    ```python
    # neoarc-server/config.py
    ADMIN_CREDENTIALS = {
        "email": "your_admin_email@example.com", # <--- CHANGE THIS
        "password": "your_secure_admin_password" # <--- CHANGE THIS
    }
    ```
    **Important:** Change the `SECRET_KEY` as well for production environments.
    ```python
    # neoarc-server/config.py
    SECRET_KEY = 'a_very_long_and_random_string_for_security' # <--- CHANGE THIS
    ```

6.  **Initialize the Database:**
    The `init_db()` function in `core/db.py` will create the `neoarc.db` SQLite database and its tables (`users`, `aliases`) if they don't already exist. This is automatically called when `main.py` runs for the first time if the database file is not found.

7.  **Run the Server:**
    Start the Flask development server:
    ```bash
    python main.py
    ```
    You should see output similar to:
    ```
     * Serving Flask app 'main'
     * Debug mode: on
     * Running on http://0.0.0.0:59248 (Press CTRL+C to quit)
    ```
    The server will typically start on `http://localhost:59248` (or the port specified in `config.py`).

## Post-Installation

*   **Access the Web UI:** Open your web browser and navigate to `http://localhost:59248` (or the configured host and port).
*   **Register a User:** You can register new users through the `/register` endpoint.
*   **Access Admin Panel:** The Super Admin panel is available at `http://localhost:59248/admin/login`. Use the credentials you set in `config.py` to log in.

## Client Configuration

After setting up the server, configure the CLI to connect to it:

```bash
neoarc config http://your-server:59248
neoarc config-token <api-token-from-server>
```

All CLI config is stored in `~/.config/neostore/neoarc/config.toml`.

See the [CLI Setup Guide](CLI_SETUP.md) for more details.

## Troubleshooting

*   **`ModuleNotFoundError`:** Ensure you have activated your virtual environment and installed all dependencies (`pip install flask werkzeug`).
*   **Port in Use:** If the server fails to start due to the port being in use, you can change the `PORT` variable in `config.py`.
*   **Database Issues:** If you encounter database errors, ensure the `neoarc.db` file is writable by the server process. You can try deleting `neoarc.db` (if it's a fresh setup and no data needs to be preserved) and restarting the server to re-initialize it.

Written by [Neorwc](https://github.com/rkriad585/neorwc-cli), Created by RK Riad Khan