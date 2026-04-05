#!/bin/bash
# Keyarch installer for macOS and Linux
# Usage: curl -fsSL https://raw.githubusercontent.com/vinayydv3695/keyarch/main/scripts/install.sh | bash

set -e

REPO="vinayydv3695/keyarch"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="keyarch"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$OS" in
        linux*)
            OS="linux"
            ;;
        darwin*)
            OS="darwin"
            ;;
        *)
            error "Unsupported operating system: $OS"
            ;;
    esac
    
    case "$ARCH" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            error "Unsupported architecture: $ARCH"
            ;;
    esac
    
    info "Detected platform: $OS-$ARCH"
}

# Get latest version
get_latest_version() {
    VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$VERSION" ]; then
        error "Failed to get latest version"
    fi
    info "Latest version: $VERSION"
}

# Download and install
install() {
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/keyarch-$OS-$ARCH"
    
    info "Downloading from: $DOWNLOAD_URL"
    
    TMP_DIR=$(mktemp -d)
    TMP_FILE="$TMP_DIR/keyarch"
    
    if command -v curl &> /dev/null; then
        curl -fsSL "$DOWNLOAD_URL" -o "$TMP_FILE"
    elif command -v wget &> /dev/null; then
        wget -q "$DOWNLOAD_URL" -O "$TMP_FILE"
    else
        error "curl or wget is required"
    fi
    
    chmod +x "$TMP_FILE"
    
    # Install to INSTALL_DIR
    if [ -w "$INSTALL_DIR" ]; then
        mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
    else
        info "Requesting sudo access to install to $INSTALL_DIR"
        sudo mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    rm -rf "$TMP_DIR"
    
    info "Installed $BINARY_NAME to $INSTALL_DIR/$BINARY_NAME"
}

# Verify installation
verify() {
    if command -v keyarch &> /dev/null; then
        info "Installation successful!"
        echo ""
        echo "Run 'keyarch' to start typing!"
    else
        warn "Installation completed but 'keyarch' is not in PATH"
        warn "You may need to add $INSTALL_DIR to your PATH"
    fi
}

main() {
    echo ""
    echo "  Keyarch Installer"
    echo "  Feature-rich TUI typing test"
    echo ""
    
    detect_platform
    get_latest_version
    install
    verify
    
    echo ""
}

main
