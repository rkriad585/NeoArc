# NeoArc CLI Setup Guide

This guide provides instructions for setting up and building the NeoArc Command Line Interface (CLI) client, which is written in Golang. The CLI allows you to execute aliases stored on your NeoArc server.

## Prerequisites

*   **Go 1.20+:** Ensure Go is installed on your system. You can download it from [go.dev/dl/](https://go.dev/dl/).
*   **Git:** Required if you are cloning the repository.

## Installation Steps

1.  **Clone the Repository (if not already done):**
    If you haven't already, clone the NeoArc project repository:
    ```bash
    git clone https://github.com/rkriad585/neoarc-cli.git # Assuming this is the repo
    cd neoarc-cli/neoarc-cli
    ```
    *(Note: The provided codebase context does not include a git repository, so this is a common assumption for project setup.)*

2.  **Navigate to the CLI Directory:**
    Change your current directory to the `neoarc-cli` folder:
    ```bash
    cd neoarc-cli
    ```

3.  **Build the Binary:**
    NeoArc provides platform-specific build scripts to simplify the compilation process.

    *   **On Windows (using PowerShell):**
        Open PowerShell and run:
        ```powershell
        .\build.ps1 neoarc
        ```
        This script will compile the `main.go` file and create an executable named `neoarc-windows-amd64.exe` (and other platform binaries) in a `build/` directory.

    *   **On Linux / macOS / Termux (using Bash):**
        First, make the build script executable:
        ```bash
        chmod +x build.sh
        ```
        Then, run the script:
        ```bash
        ./build.sh neoarc
        ```
        This script will compile the `main.go` file and create an executable named `neoarc-linux-amd64` (and other platform binaries) in a `build/` directory.

    Upon successful completion, you will see a `build/` directory containing the compiled binaries for various operating systems and architectures.

4.  **Add the Binary to Your System PATH (for global access):**
    To run `neoarc` commands from any directory in your terminal, you need to add the compiled binary to your system's PATH environment variable.

    *   **Option A: Move to a standard system directory:**
        *   **Linux/macOS:** Move the `neoarc-<os>-<arch>` binary to `/usr/local/bin/`.
            ```bash
            sudo mv build/neoarc-linux-amd64 /usr/local/bin/neoarc # Or neoarc-darwin-amd64 for macOS
            ```
        *   **Windows:** Move `neoarc-windows-amd64.exe` to `C:\Windows\System32` or any directory already in your PATH.
            ```powershell
            Move-Item -Path "build\neoarc-windows-amd64.exe" -Destination "C:\Windows\System32\neoarc.exe"
            ```

    *   **Option B: Add the `build` directory to your PATH:**
        This allows you to keep the binary in the `build` folder.
        *   **Linux/macOS (Bash/Zsh):** Add the following line to your `~/.bashrc`, `~/.zshrc`, or `~/.profile` file:
            ```bash
            export PATH="$PATH:/path/to/your/neoarc-cli/build"
            ```
            Replace `/path/to/your/neoarc-cli` with the actual path to your `neoarc-cli` directory. After saving, run `source ~/.bashrc` (or your respective file) or restart your terminal.
        *   **Windows:**
            1.  Search for "Environment Variables" in the Start Menu and select "Edit the system environment variables".
            2.  Click "Environment Variables..."
            3.  Under "System variables" or "User variables for <YourUser>", find the `Path` variable and click "Edit...".
            4.  Click "New" and add the full path to your `neoarc-cli\build` directory (e.g., `C:\Users\YourUser\Documents\neoarc-cli\build`).
            5.  Click "OK" on all windows to save changes. Restart your terminal.

5.  **Verify Installation:**
    Open a new terminal window and type:
    ```bash
    neoarc help
    ```
    You should see the NeoArc help message, confirming that the CLI is correctly installed and accessible.

## Next Steps

*   **Configure the Client:** Before executing aliases, you must configure the CLI to point to your NeoArc web server. Refer to the [Usage Guide](USAGE_GUIDE.md) for instructions on `neoarc config`.

Written by [Neorwc](https://github.com/rkriad585/neorwc-cli), Created by RK Riad Khan