#!/bin/bash

# shuru-hoja installer script
# Production-safe, enterprise-grade filesystem analyzer

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Version
VERSION="1.0.0"
BINARY_NAME="shuru-hoja"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc"
COMPLETION_DIR="/etc/bash_completion.d"

# Logging function
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check if running as root (for system installation)
    if [[ $EUID -eq 0 ]]; then
        log_warning "Running as root. Installation will be system-wide."
    fi
    
    # Check for curl/wget
    if command -v curl &> /dev/null; then
        DOWNLOAD_CMD="curl -fsSL"
    elif command -v wget &> /dev/null; then
        DOWNLOAD_CMD="wget -qO-"
    else
        log_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi
    
    # Check architecture
    ARCH=$(uname -m)
    case $ARCH in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            log_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
    
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
}

# Detect package manager
detect_package_manager() {
    if command -v apt-get &> /dev/null; then
        echo "apt"
    elif command -v yum &> /dev/null; then
        echo "yum"
    elif command -v dnf &> /dev/null; then
        echo "dnf"
    elif command -v pacman &> /dev/null; then
        echo "pacman"
    elif command -v zypper &> /dev/null; then
        echo "zypper"
    else
        echo "unknown"
    fi
}

# Install dependencies
install_dependencies() {
    local pkg_manager=$(detect_package_manager)
    
    log_info "Installing dependencies using $pkg_manager..."
    
    case $pkg_manager in
        apt)
            apt-get update
            apt-get install -y ca-certificates curl wget
            ;;
        yum|dnf)
            $pkg_manager install -y ca-certificates curl wget
            ;;
        pacman)
            pacman -Sy --noconfirm ca-certificates curl wget
            ;;
        zypper)
            zypper install -y ca-certificates curl wget
            ;;
        *)
            log_warning "Unknown package manager. Please install ca-certificates, curl, and wget manually."
            ;;
    esac
}

# Download binary
download_binary() {
    local url="https://github.com/shuru-hoja/releases/download/v${VERSION}/shuru-hoja_${VERSION}_${OS}_${ARCH}.tar.gz"
    
    log_info "Downloading shuru-hoja v${VERSION}..."
    
    # Create temp directory
    TEMP_DIR=$(mktemp -d)
    trap "rm -rf $TEMP_DIR" EXIT
    
    # Download and extract
    $DOWNLOAD_CMD "$url" | tar -xz -C "$TEMP_DIR"
    
    # Verify binary
    if [[ ! -f "$TEMP_DIR/$BINARY_NAME" ]]; then
        log_error "Download failed or binary not found in archive"
        exit 1
    fi
    
    # Make binary executable
    chmod +x "$TEMP_DIR/$BINARY_NAME"
    
    # Move to install directory
    log_info "Installing to $INSTALL_DIR..."
    if [[ $EUID -eq 0 ]]; then
        mv "$TEMP_DIR/$BINARY_NAME" "$INSTALL_DIR/"
    else
        sudo mv "$TEMP_DIR/$BINARY_NAME" "$INSTALL_DIR/"
    fi
    
    # Verify installation
    if command -v "$BINARY_NAME" &> /dev/null; then
        log_success "Binary installed successfully"
    else
        log_error "Installation failed. $BINARY_NAME not found in PATH"
        exit 1
    fi
}

# Install configuration
install_config() {
    log_info "Installing configuration..."
    
    CONFIG_FILE="$CONFIG_DIR/shuruhoja.conf"
    SAMPLE_CONFIG="# shuru-hoja configuration
# Production-safe filesystem analyzer

# Scanning options
max_workers = 100
skip_paths = /proc,/sys,/dev,/run

# Detection thresholds
log_file_age_days = 30
orphan_dir_age_days = 90
orphan_dir_min_size_gb = 1
cache_dir_min_size_mb = 100
duplicate_min_size_mb = 10

# Risk levels
critical_size_gb = 10
caution_size_gb = 1

# Output options
output_format = table
color_output = true
max_results = 50
"
    
    if [[ ! -f "$CONFIG_FILE" ]]; then
        if [[ $EUID -eq 0 ]]; then
            echo "$SAMPLE_CONFIG" > "$CONFIG_FILE"
        else
            echo "$SAMPLE_CONFIG" | sudo tee "$CONFIG_FILE" > /dev/null
        fi
        log_success "Configuration file created at $CONFIG_FILE"
    else
        log_info "Configuration file already exists at $CONFIG_FILE"
    fi
}

# Install shell completion
install_completion() {
    log_info "Installing shell completion..."
    
    COMPLETION_SCRIPT="# bash completion for shuru-hoja
_shuru_hoja_completion() {
    local cur prev opts
    COMPREPLY=()
    cur=\"\${COMP_WORDS[COMP_CWORD]}\"
    prev=\"\${COMP_WORDS[COMP_CWORD-1]}\"
    opts=\"--help --version --config --json --quiet --version\"
    
    if [[ \${cur} == -* ]] ; then
        COMPREPLY=( \$(compgen -W \"\${opts}\" -- \${cur}) )
        return 0
    fi
}
complete -F _shuru_hoja_completion shuru-hoja"
    
    if [[ -d "$COMPLETION_DIR" ]]; then
        if [[ $EUID -eq 0 ]]; then
            echo "$COMPLETION_SCRIPT" > "$COMPLETION_DIR/shuru-hoja"
        else
            echo "$COMPLETION_SCRIPT" | sudo tee "$COMPLETION_DIR/shuru-hoja" > /dev/null
        fi
        log_success "Bash completion installed"
    else
        log_warning "Bash completion directory not found. Skipping completion installation."
    fi
}

# Verify installation
verify_installation() {
    log_info "Verifying installation..."
    
    # Check if binary is in PATH
    if command -v "$BINARY_NAME" &> /dev/null; then
        log_success "$BINARY_NAME is available in PATH"
    else
        log_error "$BINARY_NAME not found in PATH"
        exit 1
    fi
    
    # Test run with --version
    if $BINARY_NAME --version &> /dev/null; then
        log_success "Binary is functional"
    else
        log_warning "Binary version check failed (might be expected)"
    fi
}

# Main installation flow
main() {
    echo -e "${BLUE}╔════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║     shuru-hoja Installation Script v$VERSION     ║${NC}"
    echo -e "${BLUE}║  Production-Safe Filesystem Analyzer           ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════════╝${NC}"
    echo ""
    
    check_prerequisites
    install_dependencies
    download_binary
    install_config
    install_completion
    verify_installation
    
    echo ""
    echo -e "${GREEN}✓ Installation completed successfully!${NC}"
    echo ""
    echo -e "${YELLOW}To start analyzing your filesystem, run:${NC}"
    echo -e "  ${BLUE}$BINARY_NAME${NC}"
    echo ""
    echo -e "${YELLOW}The tool runs in read-only mode and will NEVER delete files.${NC}"
    echo -e "${YELLOW}It only provides recommendations for manual cleanup.${NC}"
    echo ""
}

# Run main function
main "$@"
