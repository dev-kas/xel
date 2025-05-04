#!/bin/sh
# Xel Installer
# This script installs Xel, a runtime for VirtLang
# Usage: curl -fsSL https://raw.githubusercontent.com/dev-kas/xel/master/scripts/install.sh | sh

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

# Check for required tools
check_dependencies() {
    print_info "Checking dependencies..."
    
    if ! command -v curl >/dev/null 2>&1; then
        print_error "curl is required but not installed. Please install curl and try again."
        exit 1
    fi
    
    if ! command -v tar >/dev/null 2>&1 && ! command -v unzip >/dev/null 2>&1; then
        print_error "Either tar or unzip is required but neither is installed. Please install one of them and try again."
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

# Determine installation directory
get_install_dir() {
    if [ "$OS" = "windows" ]; then
        # For Windows, prefer a location in the user's path
        if [ -n "$LOCALAPPDATA" ]; then
            INSTALL_DIR="$LOCALAPPDATA/Xel"
        else
            INSTALL_DIR="$HOME/Xel"
        fi
    else
        # For Unix-like systems, use /usr/local/bin if possible, otherwise use $HOME/.local/bin
        if [ -d "/usr/local/bin" ] && [ -w "/usr/local/bin" ]; then
            INSTALL_DIR="/usr/local/bin"
        else
            INSTALL_DIR="$HOME/.local/bin"
            # Ensure the directory exists
            mkdir -p "$INSTALL_DIR"
        fi
    fi
    
    print_info "Installation directory: $INSTALL_DIR"
}

# Download the latest release
download_release() {
    print_info "Downloading Xel..."
    
    # Determine binary name based on platform
    if [ "$OS" = "windows" ]; then
        BINARY_NAME="xel-${OS}-${ARCH}.exe"
    else
        BINARY_NAME="xel-${OS}-${ARCH}"
    fi
    
    # GitHub release URL
    RELEASE_URL="https://github.com/dev-kas/xel/releases/latest/download/${BINARY_NAME}"
    
    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    
    # Download the binary
    if ! curl -fsSL "$RELEASE_URL" -o "$TMP_DIR/$BINARY_NAME"; then
        print_error "Failed to download Xel. Please check your internet connection and try again."
        rm -rf "$TMP_DIR"
        exit 1
    fi
    
    print_success "Download complete!"
}

# Install the binary
install_binary() {
    print_info "Installing Xel..."
    
    # Determine the final binary name
    if [ "$OS" = "windows" ]; then
        FINAL_BINARY_NAME="xel.exe"
    else
        FINAL_BINARY_NAME="xel"
    fi
    
    # Copy the binary to the installation directory
    cp "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$FINAL_BINARY_NAME"
    
    # Make the binary executable (not needed for Windows)
    if [ "$OS" != "windows" ]; then
        chmod +x "$INSTALL_DIR/$FINAL_BINARY_NAME"
    fi
    
    # Clean up temporary directory
    rm -rf "$TMP_DIR"
    
    print_success "Xel has been installed to $INSTALL_DIR/$FINAL_BINARY_NAME"
}

# Check if the installation directory is in PATH
check_path() {
    if [ "$OS" != "windows" ]; then
        # For Unix-like systems
        if ! echo "$PATH" | tr ':' '\n' | grep -q "^$INSTALL_DIR$"; then
            print_info "NOTE: $INSTALL_DIR is not in your PATH."
            print_info "To add it, run:"
            print_info "  echo 'export PATH=\"$INSTALL_DIR:\$PATH\"' >> ~/.bashrc"
            print_info "  source ~/.bashrc"
            print_info "Or for zsh:"
            print_info "  echo 'export PATH=\"$INSTALL_DIR:\$PATH\"' >> ~/.zshrc"
            print_info "  source ~/.zshrc"
        fi
    else
        # For Windows, just inform the user
        print_info "NOTE: Make sure $INSTALL_DIR is in your PATH."
    fi
}

# Verify the installation
verify_installation() {
    if [ "$OS" = "windows" ]; then
        BINARY_PATH="$INSTALL_DIR/xel.exe"
    else
        BINARY_PATH="$INSTALL_DIR/xel"
    fi
    
    if [ -f "$BINARY_PATH" ]; then
        print_success "Installation successful!"
        
        # Display the installed version
        if [ -x "$BINARY_PATH" ]; then
            VERSION=$("$BINARY_PATH" --version 2>/dev/null | grep -o "v[0-9]*\.[0-9]*\.[0-9]*" || echo "unknown")
            print_info "Installed Xel version: $VERSION"
        fi
        
        print_info "You can now use Xel by running: xel"
        print_info "For help, run: xel --help"
        print_info "To check version, run: xel --version"
    else
        print_error "Installation failed. The binary was not found at $BINARY_PATH."
        exit 1
    fi
}

# Main installation process
main() {
    print_info "=== Xel Installer ==="
    check_dependencies
    detect_platform
    get_install_dir
    download_release
    install_binary
    check_path
    verify_installation
    print_success "=== Installation Complete ==="
}

# Run the installer
main