# NeoArc: Usage Guide

This guide covers how to interact with the NeoArc system, both through its web interface (the "Execution Grid") and the command-line interface (the "Client").

## 1. Web Server (The Matrix) Usage

The web server provides a user-friendly interface for managing your aliases, discovering public ones, and accessing administrative features.

### 1.1. Accessing the Web UI

Open your web browser and navigate to the server's address (e.g., `http://localhost:59248`).

### 1.2. User Authentication

*   **Registration (`/register`):**
    *   Click on "Register here" from the login page or directly go to `/register`.
    *   Fill in your `Full Name`, `Username`, `Secure Email`, and `Password`.
    *   Click "REGISTER" to create your account.
*   **Login (`/login`):**
    *   Enter your `Username` and `Password` on the login page.
    *   Click "Login" to access your dashboard.
*   **Logout (`/logout`):**
    *   Click the "Logout" button (usually in the navigation bar) to end your session.

### 1.3. Dashboard (`/dashboard`)

This is your personal "Execution Grid" where you manage your aliases.

*   **Creating a New Alias:**
    1.  In the "Compile" section (left column), fill in:
        *   **Alias Name:** A unique, short name for your command (e.g., `sys_info`, `update_packages`).
        *   **Execution Type (Optional):** Select the language/shell for your command (e.g., `BASH`, `POWERSHELL`, `PYTHON`, `GOLANG`). If left as `AUTO`, the CLI will attempt to infer the type or default to `sh`/`cmd`.
        *   **Code / Command:** Paste or type your script/command here.
    2.  Click the "Compile" button to save your alias.
*   **Viewing Existing Aliases:**
    *   Your created aliases are listed in the "Active Nodes" section (right column).
    *   Each alias card shows its name, execution type, and a code preview.
*   **Copying Alias Command for CLI:**
    *   Click the copy icon on an alias card to copy the `neoarc <alias_name>` command to your clipboard.
*   **Deleting an Alias:**
    *   Click the trash can icon on an alias card. You will be prompted to confirm deletion.

### 1.4. Profile Management (`/profile`)

Update your personal information and avatar.

*   **Avatar Upload:**
    *   Click "CLICK TO SELECT IMAGE..." to upload a profile picture.
*   **Update Details:**
    *   Modify your `FULL NAME` and `SECURE EMAIL`. Your `USERNAME` is immutable.
    *   Click "Update" to save changes.
*   **System Stats:** View your active alias count and system status.

### 1.5. Global Node Search (`/search`)

Discover public aliases created by other users.

*   **Search:** Enter keywords in the search bar to find aliases.
*   **View Node:** Click on any alias in the search results to view its full code and creator details.

### 1.6. Node Inspector (`/view/<alias_name>`)

Provides a detailed view of a specific alias, including syntax-highlighted code.

*   **Copy Command:** Use the "Copy" button to get the `neoarc <alias_name>` command for CLI execution.

## 2. CLI Tool (The Client) Usage

The NeoArc CLI is your gateway to executing aliases from your terminal.

### 2.1. Configuring the Client

Before you can run any aliases, you must tell the CLI where your NeoArc server is located.

```bash
neoarc config http://localhost:59248
# Or your public IP/Domain:
# neoarc config http://192.168.1.100:59248
# neoarc config https://your.domain.com
```
This command saves the server URL to a local configuration file (`config.json` in `~/.neoarc` or `%APPDATA%\neoarc`).

### 2.2. Executing Aliases

Once configured, you can run any alias you've created or found on the web server.

*   **Run an Alias:**
    ```bash
    neoarc run <alias_name>
    # Shorthand:
    neoarc <alias_name>
    ```
    Example:
    ```bash
    neoarc sys_info
    ```
    The CLI will fetch the command from the server and execute it on your local machine.

*   **View Alias Code (without running):**
    To inspect the code of an alias before executing it:
    ```bash
    neoarc get <alias_name>
    ```
    This will print the command's code and its execution type to your terminal.

### 2.3. Help Menu

To see a list of available CLI commands:
```bash
neoarc help
```

## 3. Admin Panel Usage

The Super Admin panel provides powerful moderation capabilities for the NeoArc platform.

### 3.1. Accessing the Admin Panel

Navigate to `http://localhost:59248/admin/login` in your web browser.

### 3.2. Admin Login

*   Use the `email` and `password` configured in `neoarc-server/config.py` under `ADMIN_CREDENTIALS`.
*   Click "ACCESS CONTROL" to log in.

### 3.3. Admin Dashboard (`/admin/`)

*   **User List:** View all registered users, their email, the number of aliases they've created, and their current status (Active/Blocked).
*   **Actions per User:**
    *   **View:** Click "View" to see a user's detailed profile and a list of all their aliases.
    *   **Block/Unblock:** Click "Block" to prevent a user from logging in. The button will change to "Unblock" for blocked users.
    *   **Delete:** Click "Delete" to permanently remove a user and all their associated aliases from the database. **Use with caution!**

### 3.4. View User Details (`/admin/user/<user_id>`)

*   **User Information:** See the user's full name, username, email, and profile picture.
*   **User Aliases:** A table listing all aliases created by that specific user, including alias name, type, and a command preview.
*   **Delete User Alias:** From this view, you can delete individual aliases created by the user.

Written by [Neorwc](https://github.com/rkriad585/neorwc-cli), Created by RK Riad Khan