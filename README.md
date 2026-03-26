# NeoArc

![Logo](neoarc-server/static/logo.svg)

![NeoArc Banner](https://img.shields.io/badge/NeoArc-Execution_Grid-ea2b2b?style=for-the-badge&logo=terminator&logoColor=white) 
![Python](https://img.shields.io/badge/Server-Python_Flask-black?style=for-the-badge&logo=python) 
![Go](https://img.shields.io/badge/Client-Golang-00ADD8?style=for-the-badge&logo=go) 
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

**NeoArc** is a sophisticated, cross-platform **Command Obfuscation and Alias Execution System**. It bridges the gap between web-based storage and terminal execution, allowing users to store complex scripts (Bash, PowerShell, Python, Go) in a centralized "Execution Grid" and trigger them remotely on any device using a simple alias.

Featuring a premium **"Liquid Glass" / "Nothing OS" UI**, NeoArc offers a cyberpunk aesthetic with powerful administrative control.

---

## 📸 Features

### 🌐 Web Server (The Matrix)
*   **Liquid Glass UI:** A stunning dark-mode interface using TailwindCSS, frosted glass effects, and DotGothic typography.
*   **User System:** Secure Registration, Login, and Profile management with Avatar uploads.
*   **Global Search:** Discover public aliases created by other users in the grid.
*   **Node Inspector:** View detailed info on aliases, including syntax-highlighted code previews.
*   **Role-Based Access Control (RBAC):** Dedicated Super Admin panel to manage the platform.

### 💻 CLI Tool (The Client)
*   **Cross-Platform:** Runs natively on Windows, Linux, macOS, and Termux (Android).
*   **Polyglot Execution:** Automatically detects and runs code in:
    *   `Bash` / `Sh`
    *   `PowerShell`
    *   `CMD`
    *   `Python`
    *   `Golang` (Auto-compiles and runs)
*   **Stealth/Speed:** written in Go for high performance and single-binary deployment.

---

## 📂 Project Structure

```text
NEOARC/
├── neoarc-cli/                 # Golang Client
│   ├── build.ps1               # Windows Build Script
│   ├── build.sh                # Linux/Mac Build Script
│   ├── go.mod
│   └── main.go                 # CLI Logic
│
├── neoarc-server/              # Python Flask Server
│   ├── core/                   # Backend Modules
│   │   └── admin.py            # Super Admin Logic
│   ├── static/                 # Logo, Favicons
│   ├── templates/              # HTML Templates (Jinja2)
│   │   ├── admin/              # Admin Panel Templates
│   │   │   ├── dashboard.html
│   │   │   ├── login.html
│   │   │   └── view_user.html
│   │   ├── base.html           # Main Layout
│   │   ├── dashboard.html      # User Dashboard
│   │   ├── error.html          # Custom 404/500 Pages
│   │   ├── profile.html        # User Profile & Settings
│   │   ├── register.html
│   │   ├── search.html         # Global Search
│   │   └── view.html           # Alias Details
│   ├── config.py               # Admin Credentials Config
│   └── main.py                 # App Entry Point
│
└── README.md
```

---

## 🚀 Installation & Setup

### 1. Server Setup (Python)

Ensure you have **Python 3.8+** installed.

1.  Navigate to the server directory:
    ```bash
    cd neoarc-server
    ```

2.  Create and activate a virtual environment (recommended):
    ```bash
    # Windows
    python -m venv .venv
    .venv\Scripts\activate

    # Linux/macOS
    python3 -m venv .venv
    source .venv/bin/activate
    ```

3.  Install dependencies:
    ```bash
    pip install flask werkzeug
    # Or if you use uv:
    uv pip install flask werkzeug
    ```

4.  **Configuration:**
    Open `neoarc-server/config.py` and set your Admin credentials:
    ```python
    ADMIN_CREDENTIALS = {
        "email": "your_email@example.com",
        "password": "your_secure_password"
    }
    ```

5.  Run the server:
    ```bash
    python main.py
    ```
    *The server will start at `http://localhost:5000`.*

---

### 2. Client Setup (Golang)

Ensure you have **Go 1.20+** installed.

1.  Navigate to the CLI directory:
    ```bash
    cd neoarc-cli
    ```

2.  Build the binary:

    *   **Windows (PowerShell):**
        ```powershell
        .\build.ps1
        ```
    *   **Linux / macOS / Termux:**
        ```bash
        chmod +x build.sh
        ./build.sh
        ```

3.  Add the binary to your system PATH or move it to `/usr/local/bin` (Linux) or `C:\Windows\System32` (Windows) for global access.

---

## 🎮 Usage Guide

### Configuring the Client
Before running aliases, point the CLI to your web server:

```bash
neoarc config http://localhost:5000
# Or your public IP/Domain
neoarc config http://192.168.1.100:5000
```

### Executing Aliases
Once you have created an alias on the website (e.g., named `sys_info`), run it instantly:

```bash
# Run immediately
neoarc run sys_info

# Shorthand
neoarc sys_info

# View the code without running
neoarc get sys_info
```

---

## 🛡️ Admin Panel

NeoArc comes with a built-in Super Admin dashboard to moderate the grid.

1.  Navigate to: `http://localhost:5000/admin/login`
2.  Log in using credentials from `config.py`.

**Capabilities:**
*   **Dashboard:** View all users and their alias counts.
*   **User Inspector:** View full profile, email, and list of aliases.
*   **Moderation:**
    *   **Block User:** Prevent a user from logging in.
    *   **Delete User:** Purge a user and all their aliases from the database.
    *   **Delete Alias:** Remove specific malicious or broken aliases.

---

## 🛠️ Technology Stack

*   **Frontend:** HTML5, TailwindCSS (CDN), jQuery, Highlight.js, FontAwesome (via SVG).
*   **Backend:** Python Flask, Blueprint Architecture, SQLite3.
*   **Client:** Golang (`net/http`, `os/exec`).
*   **Design Language:** Nothing OS (Dot Matrix Typography & Monochrome/Red Palette).

---

## 🤝 Contributing

Contributions are welcome!
1.  Fork the Project.
2.  Create your Feature Branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your Changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the Branch (`git push origin feature/AmazingFeature`).
5.  Open a Pull Request.

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

---

<p align="center">
  <span style="font-family: monospace;">NEOARC v1.0.0</span>
</p>
