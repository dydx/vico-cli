#!/bin/bash

# vico-cli CLI Installer
# This script detects your OS/architecture and installs the appropriate binary

set -e

# Default installation directory
INSTALL_DIR="/usr/local/bin"
REPO_OWNER="dydx"
REPO_NAME="vico-cli"
GITHUB_RELEASES_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Installing vico-cli CLI...${NC}"

# Version handling
if [ -z "$1" ]; then
  VERSION=$(curl -s https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
  echo -e "Installing latest version: ${GREEN}${VERSION}${NC}"
else
  VERSION=$1
  echo -e "Installing specified version: ${GREEN}${VERSION}${NC}"
fi

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map architecture to release binary architecture
if [ "$ARCH" = "x86_64" ]; then
  ARCH="amd64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
  ARCH="arm64"
elif [[ "$ARCH" == armv* ]]; then
  echo -e "${RED}Error: ARM 32-bit architectures are not supported${NC}"
  exit 1
fi

# Handle OS detection
case "$OS" in
  darwin)
    BINARY_NAME="vico-cli-darwin-${ARCH}"
    ;;
  linux)
    BINARY_NAME="vico-cli-linux-${ARCH}"
    ;;
  mingw*|msys*|cygwin*|windows*)
    OS="windows"
    BINARY_NAME="vico-cli-windows-${ARCH}.exe"
    INSTALL_DIR="$HOME/bin"
    REPO_NAME="vico-cli.exe"
    ;;
  *)
    echo -e "${RED}Unsupported operating system: $OS${NC}"
    exit 1
    ;;
esac

# Create a temporary directory
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

# Download URL
DOWNLOAD_URL="${GITHUB_RELEASES_URL}/download/${VERSION}/${BINARY_NAME}"
echo -e "Downloading binary for ${YELLOW}${OS}/${ARCH}${NC} from ${BLUE}${DOWNLOAD_URL}${NC}"

# Download the binary
if ! curl -L -s --fail "$DOWNLOAD_URL" -o "$TMP_DIR/$BINARY_NAME"; then
  echo -e "${RED}Failed to download binary. Please check if the version exists at:${NC}"
  echo -e "${BLUE}${GITHUB_RELEASES_URL}${NC}"
  exit 1
fi

# For non-Windows, make binary executable
if [ "$OS" != "windows" ]; then
  chmod +x "$TMP_DIR/$BINARY_NAME"
fi

# Create installation directory if it doesn't exist
if [ ! -d "$INSTALL_DIR" ]; then
  echo -e "Creating installation directory: ${YELLOW}${INSTALL_DIR}${NC}"
  mkdir -p "$INSTALL_DIR"
fi

# Move binary to installation directory
if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$REPO_NAME"
else
  echo -e "${YELLOW}Elevated permissions required to install to $INSTALL_DIR${NC}"
  sudo mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$REPO_NAME"
fi

# Verify installation
if [ -x "$INSTALL_DIR/$REPO_NAME" ] || [ "$OS" = "windows" ]; then
  echo -e "${GREEN}Successfully installed vico-cli CLI to ${INSTALL_DIR}/${REPO_NAME}${NC}"
  echo -e "\nTo use the CLI, you need to set up your credentials:"
  echo -e "${YELLOW}export vico-cli_EMAIL=\"your.email@example.com\"${NC}"
  echo -e "${YELLOW}export vico-cli_PASSWORD=\"your-password\"${NC}"
  echo -e "\nYou can now run: ${GREEN}${REPO_NAME} devices list${NC}"
else
  echo -e "${RED}Installation failed. Please try again or install manually.${NC}"
  exit 1
fi

exit 0