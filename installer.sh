#!/usr/bin/env bash
# installer.sh - Installer and Uninstaller for NeoArc on Linux / macOS
# Auto detects OS and architecture, downloads the release binary, sets up PATH, and supports self-uninstallation.

set -euo pipefail

# ── Configuration ────────────────────────────────────────────────────────────
PROJECT_NAME="neoarc"
PUBLISHER_NAME="rkriad585"
GITHUB_REPO="rkriad585/NeoArc"

# Paths
CONFIG_DIR="${HOME}/.config/neostore/${PROJECT_NAME}"
BIN_DIR="${CONFIG_DIR}/bin"
BINARY_PATH="${BIN_DIR}/${PROJECT_NAME}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

clean_profile() {
    local profile="$1"
    if [[ -f "$profile" ]]; then
        grep -v "neostore/${PROJECT_NAME}/bin" "$profile" > "${profile}.tmp" || true
        mv "${profile}.tmp" "$profile"
    fi
}

add_to_profile() {
    local profile="$1"
    if [[ -f "$profile" ]]; then
        if ! grep -q "neostore/${PROJECT_NAME}/bin" "$profile"; then
            echo "" >> "$profile"
            echo "# NeoArc PATH configuration" >> "$profile"
            echo "export PATH=\"\$HOME/.config/neostore/${PROJECT_NAME}/bin:\$PATH\"" >> "$profile"
        fi
    fi
}

# ── Check for Uninstallation Flag ───────────────────────────────────────────
IS_UNINSTALL=false
for arg in "$@"; do
    if [[ "$arg" == "--selfuninstall" || "$arg" == "-selfuninstall" || "$arg" == "--uninstall" || "$arg" == "-u" ]]; then
        IS_UNINSTALL=true
    fi
done

if [ "$IS_UNINSTALL" = true ]; then
    echo -e ""
    echo -e "${YELLOW}╔══════════════════════════════════════════════════╗${NC}"
    echo -e "${YELLOW}║              NeoArc Uninstaller                  ║${NC}"
    echo -e "${YELLOW}╚══════════════════════════════════════════════════╝${NC}"
    echo -e ""

    if [[ -f "$BINARY_PATH" ]]; then
        printf "  Removing binary: %s ... " "$BINARY_PATH"
        rm -f "$BINARY_PATH"
        echo -e "${GREEN}OK${NC}"
    fi

    if [[ -d "$CONFIG_DIR" ]]; then
        printf "  Removing config directory: %s ... " "$CONFIG_DIR"
        rm -rf "$CONFIG_DIR"
        echo -e "${GREEN}OK${NC}"
    fi

    echo -e "  Cleaning shell profile PATH configurations ... "
    PROFILES=(
        "${HOME}/.bashrc"
        "${HOME}/.zshrc"
        "${HOME}/.profile"
        "${HOME}/.bash_profile"
    )
    for p in "${PROFILES[@]}"; do
        if [[ -f "$p" ]]; then
            printf "    Updating %s ... " "$(basename "$p")"
            clean_profile "$p"
            echo -e "${GREEN}OK${NC}"
        fi
    done

    echo -e ""
    echo -e "${GREEN}  NeoArc has been successfully uninstalled from your system.${NC}"
    echo -e "${CYAN}  Please restart your shell or run: hash -r${NC}"
    echo -e ""
    exit 0
fi

# ── Installation / Update Flow ───────────────────────────────────────────────
echo -e ""
echo -e "${CYAN}╔══════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║               NeoArc Installer                   ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════════════╝${NC}"
echo -e ""

printf "  Checking latest version from GitHub ... "
VERSION_URL="https://raw.githubusercontent.com/${GITHUB_REPO}/main/.version"
if command -v curl >/dev/null 2>&1; then
    VERSION=$(curl -fsSL "$VERSION_URL" | tr -d '[:space:]')
elif command -v wget >/dev/null 2>&1; then
    VERSION=$(wget -qO- "$VERSION_URL" | tr -d '[:space:]')
else
    echo -e "${RED}FAILED${NC}"
    echo "Error: curl or wget is required to run this installer." >&2
    exit 1
fi
echo -e "${GREEN}${VERSION}${NC}"

printf "  Detecting platform ... "
OS_NAME=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS_NAME" in
    linux*)  OS="linux" ;;
    darwin*) OS="darwin" ;;
    *)
        echo -e "${RED}FAILED${NC}"
        echo "Error: Unsupported operating system: ${OS_NAME}" >&2
        exit 1
        ;;
esac

ARCH_NAME=$(uname -m)
case "$ARCH_NAME" in
    x86_64)  ARCH="amd64" ;;
    amd64)   ARCH="amd64" ;;
    arm64)   ARCH="arm64" ;;
    aarch64) ARCH="arm64" ;;
    *)
        echo -e "${RED}FAILED${NC}"
        echo "Error: Unsupported CPU architecture: ${ARCH_NAME}" >&2
        exit 1
        ;;
esac
echo -e "${GREEN}${OS}-${ARCH}${NC}"

DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/${PROJECT_NAME}-${OS}-${ARCH}"
echo -e "  Downloading binary from ${DOWNLOAD_URL} ...\n"

mkdir -p "$BIN_DIR"

if command -v curl >/dev/null 2>&1; then
    curl -L --progress-bar -o "$BINARY_PATH" "$DOWNLOAD_URL"
elif command -v wget >/dev/null 2>&1; then
    wget --show-progress -O "$BINARY_PATH" "$DOWNLOAD_URL"
else
    echo -e "  ${RED}Failed to download binary. No curl or wget found.${NC}" >&2
    exit 1
fi

if [[ ! -f "$BINARY_PATH" ]]; then
    echo -e "  ${RED}Downloaded file not found. Download failed.${NC}" >&2
    exit 1
fi

chmod +x "$BINARY_PATH"
SIZE=$(du -h "$BINARY_PATH" | cut -f1)
echo -e "  ${GREEN}Successfully downloaded ${PROJECT_NAME} (${SIZE})${NC}"

printf "  Configuring shell profile PATH configurations ... "
PROFILES=(
    "${HOME}/.bashrc"
    "${HOME}/.zshrc"
    "${HOME}/.profile"
    "${HOME}/.bash_profile"
)
UPDATED_PROFILES=()

ACTIVE_SHELL=$(basename "$SHELL")
if [ "$ACTIVE_SHELL" = "bash" ]; then
    touch "${HOME}/.bashrc"
    add_to_profile "${HOME}/.bashrc"
    UPDATED_PROFILES+=(".bashrc")
elif [ "$ACTIVE_SHELL" = "zsh" ]; then
    touch "${HOME}/.zshrc"
    add_to_profile "${HOME}/.zshrc"
    UPDATED_PROFILES+=(".zshrc")
else
    for p in "${PROFILES[@]}"; do
        if [[ -f "$p" ]]; then
            add_to_profile "$p"
            UPDATED_PROFILES+=("$(basename "$p")")
        fi
    done
fi

if [ ${#UPDATED_PROFILES[@]} -eq 0 ]; then
    touch "${HOME}/.profile"
    add_to_profile "${HOME}/.profile"
    UPDATED_PROFILES+=(".profile")
fi

echo -e "${GREEN}OK${NC} (Updated: ${UPDATED_PROFILES[*]})"

# ── Success Banner ────────────────────────────────────────────────────────────
echo -e ""
echo -e "${GREEN}╔══════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║         NeoArc successfully installed!           ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════════╝${NC}"
echo -e ""
echo -e "  Installation Path: ${BINARY_PATH}"
echo -e "  Version          : ${VERSION}"
echo -e ""
echo -e "${YELLOW}  Please RESTART your terminal window or run:${NC}"
echo -e "${CYAN}  source ~/${UPDATED_PROFILES[0]}${NC}"
echo -e "${YELLOW}  to start using NeoArc immediately!${NC}"
echo -e ""
