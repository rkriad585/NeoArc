# NeoArc: Project Overview

NeoArc is an innovative, cross-platform Command Obfuscation and Alias Execution System designed to streamline the execution of complex scripts across various environments. It provides a centralized web-based "Execution Grid" for storing scripts (Bash, PowerShell, Python, Go) and a lightweight, high-performance Go CLI tool for triggering them remotely via simple aliases.

## Key Concepts

*   **Execution Grid (Web Server):** The central hub where users manage their aliases, view public aliases, and interact with the system through a "Liquid Glass" UI.
*   **Client (CLI Tool):** A fast, single-binary Go application that connects to the Execution Grid to fetch and execute commands on the local machine.
*   **Aliases:** User-defined names for complex scripts or commands, stored on the server and executed via the CLI.
*   **Polyglot Execution:** The CLI's ability to automatically detect and run scripts written in Bash, PowerShell, CMD, Python, or Go.

## Core Features

### Web Server (The Matrix)
*   **Liquid Glass UI:** A modern, dark-mode interface built with TailwindCSS, featuring frosted glass effects and DotGothic typography for a cyberpunk aesthetic.
*   **Secure User System:** Comprehensive user management including registration, login, profile updates, and avatar uploads.
*   **Global Search:** Enables users to discover and utilize public aliases shared by the community.
*   **Node Inspector:** Provides detailed views of aliases, including syntax-highlighted code previews.
*   **Role-Based Access Control (RBAC):** A dedicated Super Admin panel for platform moderation and user management.

### CLI Tool (The Client)
*   **Cross-Platform Compatibility:** Natively supported on Windows, Linux, macOS, and Termux (Android).
*   **Versatile Script Execution:** Capable of running scripts in Bash/Sh, PowerShell, CMD, Python, and Go (with automatic compilation).
*   **Optimized for Performance:** Developed in Go for speed, efficiency, and single-binary deployment, ensuring minimal overhead.

## Technology Stack

*   **Frontend:** HTML5, TailwindCSS (CDN), jQuery, Highlight.js, FontAwesome.
*   **Backend:** Python Flask, Blueprint Architecture, SQLite3.
*   **Client:** Golang (`net/http`, `os/exec`).
*   **Design Language:** Inspired by Nothing OS (Dot Matrix Typography & Monochrome/Red Palette).

## Project Structure

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
│   │   ├── admin.py            # Super Admin Logic
│   │   └── db.py               # Database Connection & Schema
│   ├── static/                 # Logo, Favicons, Uploads
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
│   ├── .python-version         # Python version specification
│   ├── config.py               # Admin Credentials Config
│   └── main.py                 # App Entry Point
│
└── README.md                   # Main project README
```

Written by [Neorwc](https://github.com/rkriad585/neorwc-cli), Created by RK Riad Khan