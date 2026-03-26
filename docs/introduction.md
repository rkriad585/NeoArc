# Introduction to NeoArc

NeoArc is an innovative, cross-platform **Command Obfuscation and Alias Execution System** designed to streamline the execution of complex scripts across various environments. It acts as a bridge between a centralized web-based "Execution Grid" and local terminal execution, allowing users to define and store scripts (Bash, PowerShell, Python, Go, etc.) as simple aliases. These aliases can then be triggered remotely on any configured device using a lightweight command-line interface (CLI) tool.

The system is built with a focus on both aesthetics and functionality, featuring a premium "Liquid Glass" / "Nothing OS" inspired user interface on the web server and a high-performance, polyglot CLI client.

## Core Concepts

*   **Execution Grid (Web Server):** The central hub where users manage their profiles, create and store aliases, and discover public scripts. It provides a secure, visually appealing interface for script management.
*   **Client (CLI Tool):** A lightweight, cross-platform binary that connects to the Execution Grid. It fetches aliases and executes the associated commands locally, automatically detecting the appropriate interpreter or compiler.
*   **Aliases:** User-defined names linked to specific commands or scripts. These can range from simple shell commands to multi-line scripts in various languages.
*   **Polyglot Execution:** The CLI's ability to intelligently determine the execution environment (Bash, PowerShell, Python, Go, CMD) based on the alias's `exec_type` and run the script accordingly.

## Key Features

### Web Server (The Matrix)

*   **Liquid Glass UI:** A modern, dark-mode interface leveraging TailwindCSS, frosted glass effects, and the distinctive DotGothic typography, creating a cyberpunk aesthetic.
*   **User System:** Comprehensive user management including secure registration, login, and profile customization with avatar uploads.
*   **Global Search:** A powerful feature allowing users to discover and inspect public aliases shared by the community within the Execution Grid.
*   **Node Inspector:** Provides detailed information about aliases, including syntax-highlighted code previews for better understanding and security review.
*   **Role-Based Access Control (RBAC):** A dedicated Super Admin panel for platform moderation, user management, and alias control.

### CLI Tool (The Client)

*   **Cross-Platform Compatibility:** Natively supported on Windows, Linux, macOS, and Termux (Android), ensuring broad accessibility.
*   **Polyglot Execution Engine:** Automatically identifies and executes code written in:
    *   `Bash` / `Sh`
    *   `PowerShell`
    *   `CMD` (Windows Command Prompt)
    *   `Python`
    *   `Golang` (Includes auto-compilation for Go scripts)
*   **Stealth & Performance:** Developed in Go, the CLI offers high performance, minimal resource usage, and single-binary deployment for ease of distribution and execution.

## Use Cases

*   **DevOps Automation:** Store and execute complex deployment scripts, server maintenance tasks, or CI/CD triggers with simple aliases.
*   **System Administration:** Centralize common administrative commands (e.g., log analysis, service restarts, system health checks) for quick access across multiple servers.
*   **Personal Productivity:** Create shortcuts for frequently used commands, development environment setups, or custom utility scripts.
*   **Security & Penetration Testing:** Manage and deploy custom exploit scripts or reconnaissance tools from a central repository.
*   **Educational Tool:** Share and learn from community-contributed scripts in a controlled environment.

NeoArc aims to empower users with a flexible, secure, and visually engaging platform for managing and executing their digital commands.

Written by [Neorwc](https://github.com/rkriad585/neorwc-cli), Created by RK Riad Khan