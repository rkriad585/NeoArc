#!/usr/bin/env sh
# NeoArc Installer for Linux / macOS
# Detects architecture, downloads the correct binary from GitHub,
# installs to ~/.config/neostore/neoarc/bin/neoarc, and adds it to PATH.
# Pass --selfuninstall to remove NeoArc from the system.

set -eu

PROJECT_NAME="neoarc"
REPO_OWNER="rkriad585"
REPO="${REPO_OWNER}/${PROJECT_NAME}"
INSTALL_DIR="${HOME}/.config/neostore/${PROJECT_NAME}/bin"
BINARY="${PROJECT_NAME}"
INSTALL_PATH="${INSTALL_DIR}/${BINARY}"

# ---- Resolve version from GitHub raw .version file ----
echo ">>> Resolving latest version..."
if command -v curl >/dev/null 2>&1; then
    VERSION="$(curl -fsSL "https://raw.githubusercontent.com/${REPO}/main/.version" 2>/dev/null | tr -d '[:space:]')"
elif command -v wget >/dev/null 2>&1; then
    VERSION="$(wget -qO- "https://raw.githubusercontent.com/${REPO}/main/.version" 2>/dev/null | tr -d '[:space:]')"
else
    echo "Error: Neither curl nor wget found."
    exit 1
fi

if [ -z "${VERSION}" ]; then
    echo "Error: Could not determine version from GitHub."
    exit 1
fi
echo "OK   Version: ${VERSION}"

# ---- Detect OS + architecture ----
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

echo ">>> Detecting system architecture..."
case "${OS}" in
    linux)
        case "${ARCH}" in
            x86_64|amd64) DOWNLOAD="${PROJECT_NAME}-linux-amd64"    ;;
            aarch64|arm64) DOWNLOAD="${PROJECT_NAME}-linux-arm64"   ;;
            *)            echo "Unsupported architecture: ${ARCH}"; exit 1 ;;
        esac
        echo "OK   Architecture: ${ARCH}"
        ;;
    darwin)
        case "${ARCH}" in
            x86_64|amd64) DOWNLOAD="${PROJECT_NAME}-darwin-amd64"   ;;
            arm64)        DOWNLOAD="${PROJECT_NAME}-darwin-arm64"   ;;
            *)            echo "Unsupported architecture: ${ARCH}"; exit 1 ;;
        esac
        echo "OK   Architecture: ${ARCH}"
        ;;
    *)
        echo "Unsupported OS: ${OS}. Use installer.ps1 on Windows."
        exit 1
        ;;
esac

URL="https://github.com/${REPO}/releases/download/${VERSION}/${DOWNLOAD}"

# ---- Handle uninstall ----
if [ "${1:-}" = "--selfuninstall" ]; then
    echo ">>> Uninstalling ${PROJECT_NAME}..."
    if [ -d "${INSTALL_DIR}" ]; then
        rm -rf "${INSTALL_DIR}"
        echo "OK   Removed ${INSTALL_DIR}"
    else
        echo "OK   No install directory found."
    fi

    for RC in "${HOME}/.zshrc" "${HOME}/.bashrc" "${HOME}/.bash_profile" "${HOME}/.profile"; do
        if [ -f "${RC}" ]; then
            if grep -qsF "${INSTALL_DIR}" "${RC}" 2>/dev/null; then
                grep -vF "${INSTALL_DIR}" "${RC}" | grep -v "# Added by ${PROJECT_NAME} installer" > "${RC}.tmp" && mv "${RC}.tmp" "${RC}"
                echo "OK   Removed PATH entry from ${RC}"
            fi
        fi
    done

    echo ""
    echo "========================================"
    echo "  ${PROJECT_NAME} has been uninstalled."
    echo "========================================"
    exit 0
fi

# ---- Create install directory ----
echo ">>> Creating install directory: ${INSTALL_DIR}"
mkdir -p "${INSTALL_DIR}"

# ---- Download binary ----
echo ">>> Downloading ${DOWNLOAD} from ${URL}"
if command -v curl >/dev/null 2>&1; then
    curl -fsSL -o "${INSTALL_PATH}" "${URL}"
elif command -v wget >/dev/null 2>&1; then
    wget -q -O "${INSTALL_PATH}" "${URL}"
else
    echo "Error: Neither curl nor wget found."
    exit 1
fi

if [ ! -f "${INSTALL_PATH}" ]; then
    echo "Error: Download failed — binary not found at ${INSTALL_PATH}"
    exit 1
fi

chmod +x "${INSTALL_PATH}"
echo "OK   Installed to ${INSTALL_PATH}"

# ---- Add to PATH (via shell rc) ----
LINE="export PATH=\"\${PATH}:${INSTALL_DIR}\""

add_to_rc() {
    RC="$1"
    if [ ! -f "${RC}" ]; then
        touch "${RC}"
    fi
    if ! grep -qsF "${INSTALL_DIR}" "${RC}"; then
        echo "" >> "${RC}"
        echo "# Added by ${PROJECT_NAME} installer" >> "${RC}"
        echo "${LINE}" >> "${RC}"
        echo "OK   Added ${INSTALL_DIR} to PATH in ${RC}"
    else
        echo "OK   ${INSTALL_DIR} already in ${RC}"
    fi
}

case "${SHELL}" in
    *zsh) add_to_rc "${HOME}/.zshrc" ;;
    *bash)
        if [ "${OS}" = "darwin" ]; then
            add_to_rc "${HOME}/.bash_profile"
        else
            add_to_rc "${HOME}/.bashrc"
        fi
        ;;
    *) add_to_rc "${HOME}/.profile" ;;
esac

echo ""
echo "========================================"
echo "  ${PROJECT_NAME} installed successfully!"
echo "  Binary : ${INSTALL_PATH}"
echo "  Version: ${VERSION}"
echo "  Usage  : ${PROJECT_NAME} help"
echo "  To uninstall, run:"
echo "    ${PROJECT_NAME} --selfuninstall"
echo "========================================"
echo ""
echo "Run the following now (or restart your terminal):"
case "${SHELL}" in
    *zsh) echo "  source ~/.zshrc" ;;
    *bash)
        if [ "${OS}" = "darwin" ]; then
            echo "  source ~/.bash_profile"
        else
            echo "  source ~/.bashrc"
        fi
        ;;
    *) echo "  export PATH=\"\${PATH}:${INSTALL_DIR}\"" ;;
esac
