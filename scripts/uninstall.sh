#!/bin/sh
# Xel Uninstaller
# This script uninstalls Xel from your system
# Usage: curl -fsSL https://raw.githubusercontent.com/dev-kas/xel/master/scripts/uninstall.sh | sh
# Or with forced uninstallation: curl -fsSL https://raw.githubusercontent.com/dev-kas/xel/master/scripts/uninstall.sh | sh -s -- -y

set -e # Exit immediately if a command exits with a non-zero status

# Parse command line arguments
FORCE_UNINSTALL=false

# Check if -y is passed as an argument
for arg in "$@"; do
    if [ "$arg" = "-y" ] || [ "$arg" = "--yes" ]; then
        FORCE_UNINSTALL=true
        break
    fi
done

# Also check if the first argument contains -y (for cases where sh -s -- -y gets mangled)
if [ -n "$1" ] && echo "$1" | grep -q -- "-y"; then
    FORCE_UNINSTALL=true
fi

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

# Detect operating system
detect_platform() {
    print_info "Detecting platform..."
    
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    
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
    
    print_info "Detected OS: $OS"
}

# Find Xel binary
find_binary() {
    print_info "Looking for Xel installation..."
    
    if [ "$OS" = "windows" ]; then
        BINARY_NAME="xel.exe"
    else
        BINARY_NAME="xel"
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
    
    # Try to find the binary in PATH
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        BINARY_PATH=$(command -v "$BINARY_NAME")
        print_info "Found Xel at: $BINARY_PATH"
        return 0
    fi
    
    # Check common locations
    for LOCATION in $COMMON_LOCATIONS; do
        if [ -f "$LOCATION" ]; then
            BINARY_PATH="$LOCATION"
            print_info "Found Xel at: $BINARY_PATH"
            return 0
        fi
    done
    
    print_error "Could not find Xel installation. It may have been installed in a non-standard location."
    print_error "Please manually remove the Xel binary from your system."
    exit 1
}

# Uninstall Xel
uninstall_binary() {
    print_info "Uninstalling Xel..."
    
    # Skip confirmation if force flag is set
    if [ "$FORCE_UNINSTALL" = "true" ]; then
        print_info "Force uninstall flag detected. Proceeding without confirmation."
        CONFIRM="y"
    else
        # Check if script is being piped from curl
        if [ -t 0 ]; then
            # Terminal is interactive, ask for confirmation
            printf "Are you sure you want to uninstall Xel from %s? [y/N] " "$BINARY_PATH"
            # Use read with a timeout to ensure it waits for input
            if [ "$OS" = "darwin" ]; then
                # macOS doesn't support -t option for read
                read -r CONFIRM
            else
                # Linux and other systems
                read -r -t 300 CONFIRM  # 5-minute timeout
            fi
        else
            # Being piped from curl, provide clear instructions and exit with helpful message
            print_warning "This script is being run via curl pipe."
            print_warning "For non-interactive uninstallation, use:"
            print_warning "curl -fsSL https://raw.githubusercontent.com/dev-kas/xel/master/scripts/uninstall.sh | sh -s -- -y"
            print_warning ""
            print_warning "Uninstallation cancelled. Please run again with the -y flag as shown above."
            exit 0
        fi
        
        if [ "$CONFIRM" != "y" ] && [ "$CONFIRM" != "Y" ]; then
            print_info "Uninstallation cancelled."
            exit 0
        fi
    fi
    
    # Remove the binary
    if rm "$BINARY_PATH"; then
        print_success "Xel has been uninstalled successfully."
    else
        print_error "Failed to remove Xel. You may need administrator privileges."
        print_error "Try running this script with sudo or as an administrator."
        exit 1
    fi
}

# Main uninstallation process
main() {
    print_info "=== Xel Uninstaller ==="
    detect_platform
    find_binary
    uninstall_binary
    print_success "=== Uninstallation Complete ==="
}

# Run the uninstaller
main