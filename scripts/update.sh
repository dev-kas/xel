#!/bin/sh
# Xel Updater
# This script updates Xel to the latest version
# Usage: curl -fsSL https://raw.githubusercontent.com/dev-kas/xel/main/scripts/update.sh | sh

set -e # Exit immediately if a command exits with a non-zero status

# Print colored output
print_info() {
    printf "\033[0;34m%s\033[0m\n" "$1"
}

print_success() {
    printf "\033[0;32m%s\033[0m\n" "$1"
}

print_error() {
    printf "\033[0;31m%s\033[0m\n" "$1" >&2
}

print_warning() {
    printf "\033[0;33m%s\033[0m\n" "$1"
}

# Check for required tools
check_dependencies() {
    print_info "Checking dependencies..."
    
    if ! command -v curl >/dev/null 2>&1; then
        print_error "curl is required but not installed. Please install curl and try again."
        exit 1
    fi
}

# Detect operating system and architecture
detect_platform() {
    print_info "Detecting platform..."
    
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    # Normalize architecture names
    case "$ARCH" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
    
    # Normalize OS names
    case "$OS" in
        darwin)
            OS="darwin"
            ;;
        linux)
            OS="linux"
            ;;
        msys*|mingw*|cygwin*|windows*)
            OS="windows"
            ;;
        *)
            print_error "Unsupported operating system: $OS"
            exit 1
            ;;
    esac
    
    print_info "Detected platform: $OS-$ARCH"
}

# Find Xel binary
find_binary() {
    print_info "Looking for Xel installation..."
    
    if [ "$OS" = "windows" ]; then
        BINARY_NAME="xel.exe"
    else
        BINARY_NAME="xel"
    fi
    
    # Try to find the binary in PATH
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        BINARY_PATH=$(command -v "$BINARY_NAME")
        INSTALL_DIR=$(dirname "$BINARY_PATH")
        print_info "Found Xel at: $BINARY_PATH"
        return 0
    fi
    
    # Common installation locations
    COMMON_LOCATIONS=""
    if [ "$OS" = "windows" ]; then
        if [ -n "$LOCALAPPDATA" ]; then
            COMMON_LOCATIONS="$COMMON_LOCATIONS $LOCALAPPDATA/Xel/$BINARY_NAME"
        fi
        COMMON_LOCATIONS="$COMMON_LOCATIONS $HOME/Xel/$BINARY_NAME"
    else
        COMMON_LOCATIONS="$COMMON_LOCATIONS /usr/local/bin/$BINARY_NAME $HOME/.local/bin/$BINARY_NAME"
    fi
    
    # Check common locations
    for LOCATION in $COMMON_LOCATIONS; do
        if [ -f "$LOCATION" ]; then
            BINARY_PATH="$LOCATION"
            INSTALL_DIR=$(dirname "$BINARY_PATH")
            print_info "Found Xel at: $BINARY_PATH"
            return 0
        fi
    done
    
    print_error "Could not find Xel installation. Please install Xel first."
    print_info "You can install Xel with: curl -fsSL https://raw.githubusercontent.com/dev-kas/xel/main/scripts/install.sh | sh"
    exit 1
}

# Check current version
check_current_version() {
    print_info "Checking current version..."
    
    if [ -x "$BINARY_PATH" ]; then
        CURRENT_VERSION=$("$BINARY_PATH" --version 2>/dev/null || echo "unknown")
        print_info "Current version: $CURRENT_VERSION"
    else
        print_warning "Could not determine current version."
        CURRENT_VERSION="unknown"
    fi
}

# Download the latest release
download_release() {
    print_info "Downloading latest version of Xel..."
    
    # Determine binary name based on platform
    if [ "$OS" = "windows" ]; then
        DOWNLOAD_BINARY_NAME="xel-${OS}-${ARCH}.exe"
    else
        DOWNLOAD_BINARY_NAME="xel-${OS}-${ARCH}"
    fi
    
    # GitHub release URL
    RELEASE_URL="https://github.com/dev-kas/xel/releases/latest/download/${DOWNLOAD_BINARY_NAME}"
    
    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    
    # Download the binary
    if ! curl -fsSL "$RELEASE_URL" -o "$TMP_DIR/$DOWNLOAD_BINARY_NAME"; then
        print_error "Failed to download Xel. Please check your internet connection and try again."
        rm -rf "$TMP_DIR"
        exit 1
    fi
    
    print_success "Download complete!"
}

# Update the binary
update_binary() {
    print_info "Updating Xel..."
    
    # Backup the current binary
    if [ -f "$BINARY_PATH" ]; then
        cp "$BINARY_PATH" "$BINARY_PATH.backup"
        print_info "Created backup at: $BINARY_PATH.backup"
    fi
    
    # Determine the downloaded binary name
    if [ "$OS" = "windows" ]; then
        DOWNLOAD_BINARY_NAME="xel-${OS}-${ARCH}.exe"
    else
        DOWNLOAD_BINARY_NAME="xel-${OS}-${ARCH}"
    fi
    
    # Copy the new binary
    if ! cp "$TMP_DIR/$DOWNLOAD_BINARY_NAME" "$BINARY_PATH"; then
        print_error "Failed to update Xel. You may need administrator privileges."
        print_error "Try running this script with sudo or as an administrator."
        
        # Restore from backup if it exists
        if [ -f "$BINARY_PATH.backup" ]; then
            print_info "Restoring from backup..."
            cp "$BINARY_PATH.backup" "$BINARY_PATH"
            rm "$BINARY_PATH.backup"
        fi
        
        rm -rf "$TMP_DIR"
        exit 1
    fi
    
    # Make the binary executable (not needed for Windows)
    if [ "$OS" != "windows" ]; then
        chmod +x "$BINARY_PATH"
    fi
    
    # Clean up
    rm -rf "$TMP_DIR"
    if [ -f "$BINARY_PATH.backup" ]; then
        rm "$BINARY_PATH.backup"
    fi
    
    print_success "Xel has been updated successfully!"
}

# Verify the update
verify_update() {
    print_info "Verifying update..."
    
    if [ -x "$BINARY_PATH" ]; then
        NEW_VERSION=$("$BINARY_PATH" --version 2>/dev/null || echo "unknown")
        print_info "Updated to version: $NEW_VERSION"
        
        if [ "$NEW_VERSION" = "$CURRENT_VERSION" ] && [ "$NEW_VERSION" != "unknown" ]; then
            print_warning "You already have the latest version of Xel."
        else
            print_success "Xel has been updated successfully!"
        fi
    else
        print_error "Update verification failed. Please check your installation."
        exit 1
    fi
}

# Main update process
main() {
    print_info "=== Xel Updater ==="
    check_dependencies
    detect_platform
    find_binary
    check_current_version
    download_release
    update_binary
    verify_update
    print_success "=== Update Complete ==="
}

# Run the updater
main