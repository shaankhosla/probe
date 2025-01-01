#!/bin/bash

# Constants
REPO_OWNER="shaankhosla"  
REPO_NAME="probe"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="probe"        
LATEST_RELEASE_URL="https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest"

# Colors for fancy output
GREEN='\033[0;32m'
RED='\033[0;31m'
NO_COLOR='\033[0m'

# Function to detect platform (OS and architecture)
detect_platform() {
    OS=$(uname | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case $ARCH in
        x86_64) ARCH="amd64" ;;
        arm64 | aarch64) ARCH="arm64" ;;
        *) echo -e "${RED}Unsupported architecture: $ARCH${NO_COLOR}"; exit 1 ;;
    esac

    if [ "$OS" == "darwin" ]; then
        GOOS="darwin"
    elif [ "$OS" == "linux" ]; then
        GOOS="linux"
    else
        echo -e "${RED}Unsupported operating system: $OS${NO_COLOR}"
        exit 1
    fi

    echo "$GOOS $ARCH"
}

# Function to get the latest release download URL for the binary
get_download_url() {
    PLATFORM=$(detect_platform)
    GOOS=$(echo $PLATFORM | cut -d' ' -f1)
    GOARCH=$(echo $PLATFORM | cut -d' ' -f2)

    DOWNLOAD_URL=$(curl -sL $LATEST_RELEASE_URL | grep "browser_download_url" | grep "$GOOS-$GOARCH" | cut -d '"' -f 4)

    if [[ -z "$DOWNLOAD_URL" ]]; then
        echo -e "${RED}Error: Failed to find a matching binary for your system ($GOOS-$GOARCH).${NO_COLOR}"
        exit 1
    fi

    echo $DOWNLOAD_URL
}

# Function to install the binary
install_binary() {
    echo -e "${GREEN}Detecting platform...${NO_COLOR}"
    PLATFORM=$(detect_platform)
    echo -e "${GREEN}Platform detected: $PLATFORM${NO_COLOR}"

    echo -e "${GREEN}Fetching download URL for the latest release...${NO_COLOR}"
    DOWNLOAD_URL=$(get_download_url)
    echo -e "${GREEN}Download URL: $DOWNLOAD_URL${NO_COLOR}"

    TMP_FILE=$(mktemp)
    echo -e "${GREEN}Downloading the binary...${NO_COLOR}"
    curl -L "$DOWNLOAD_URL" -o "$TMP_FILE"

    echo -e "${GREEN}Installing $BINARY_NAME to $INSTALL_DIR...${NO_COLOR}"
    chmod +x "$TMP_FILE"
    sudo mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"

    echo -e "${GREEN}Installation complete! Run '$BINARY_NAME --help' to get started.${NO_COLOR}"
}

# Main script
install_binary
