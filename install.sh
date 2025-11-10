#!/bin/bash
set -e

# Detect OS and architecture
detect_os_arch() {
    OS="$(uname -s)"
    ARCH="$(uname -m)"
    case "$OS" in
        Darwin)
            OS_NAME="darwin"
            ;;
        Linux)
            OS_NAME="linux"
            ;;
        *)
            echo "Unsupported OS: $OS"
            exit 1
            ;;
    esac
    case "$ARCH" in
        arm64|aarch64)
            ARCH_NAME="arm64"
            ;;
        x86_64|amd64)
            ARCH_NAME="amd64"
            ;;
        *)
            echo "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
}

detect_os_arch

BINARY_URL="https://github.com/luiz1361/catv/releases/latest/download/catv-${OS_NAME}-${ARCH_NAME}"

# Download binary
curl -L "$BINARY_URL" -o catv
chmod +x catv

# macOS: remove quarantine attribute
if [ "$OS_NAME" = "darwin" ]; then
    xattr -dr com.apple.quarantine catv || true
fi

echo "catv installed! Run ./catv --help for usage."
