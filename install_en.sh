#!/bin/bash

# Display help information
show_help() {
    echo "Usage: $0 [options]"
    echo "Options:"
    echo "  -v, --version VERSION  Specify the version to install"
    echo "  -m, --mirror MIRROR    Specify mirror source (available: ghproxy, wgetla, default)"
    echo "  -h, --help            Show help information"
    exit 0
}

# Get system information
case "$(uname -s)" in
    Linux*)
        OS="linux"
        ;;
    Darwin*)
        OS="darwin"
        ;;
    MINGW64*|MSYS*|CYGWIN*)
        OS="windows"
        ;;
    *)
        echo "Unsupported operating system: $(uname -s)"
        exit 1
        ;;
esac

# Get system architecture
case "$(uname -m)" in
    x86_64|amd64)
        ARCH="x86_64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $(uname -m)"
        exit 1
        ;;
esac

# Set download file extension and installation path based on operating system
case $OS in
    linux|darwin)
        EXT=""
        INSTALL_PATH="/usr/bin/wen"
        ;;
    windows)
        EXT=".exe"
        INSTALL_PATH="/usr/bin/wen.exe"
        ;;
    *)
        echo "Unsupported operating system: $OS"
        exit 1
        ;;
esac

# If version is not specified, get the latest version
if [ -z "$VERSION" ]; then
    echo "Getting latest version..."
    VERSION=$(curl -s https://api.github.com/repos/zdt1013/wen-ai-cli/releases/latest | grep -oP '"tag_name": "\K(.*)(?=")')
    
    if [ -z "$VERSION" ]; then
        echo "Unable to get latest version"
        exit 1
    fi
    echo "Will install latest version: $VERSION"
else
    echo "Will install specified version: $VERSION"
fi

# Define mirror domain list
declare -A MIRRORS=(
    ["ghproxy"]="https://ghproxy.imciel.com/"
    ["wgetla"]="https://wget.la/"
    ["default"]=""
)

# Parse command line arguments
MIRROR="default"
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -m|--mirror)
            if [[ -n "${MIRRORS[$2]}" ]]; then
                MIRROR="$2"
            else
                echo "Invalid mirror source: $2"
                echo "Available mirror sources: ${!MIRRORS[*]}"
                exit 1
            fi
            shift 2
            ;;
        -h|--help)
            show_help
            ;;
        *)
            echo "Unknown option: $1"
            show_help
            ;;
    esac
done

# Process version number for directory and filename
RELEASE_VERSION=$VERSION
VERSION=${VERSION#v}

# Build download URL
GITHUB_URL="https://github.com/zdt1013/wen-ai-cli/releases/download/${RELEASE_VERSION}/wen-ai-cli_${VERSION}_${OS}_${ARCH}${EXT}"

# Ensure target directory exists
mkdir -p $(dirname $INSTALL_PATH)

# Download based on selected mirror
if [ "$MIRROR" == "default" ]; then
    echo "Using default download method..."
    curl -L -o $INSTALL_PATH $GITHUB_URL
    
    if [ $? -ne 0 ]; then
        echo "Default download failed, trying mirror sources..."
        for mirror_name in "${!MIRRORS[@]}"; do
            if [ "$mirror_name" != "default" ]; then
                mirror_url="${MIRRORS[$mirror_name]}${GITHUB_URL}"
                echo "Trying mirror source: $mirror_name ($mirror_url)"
                curl -L -o $INSTALL_PATH $mirror_url
                if [ $? -eq 0 ]; then
                    echo "Download successful using mirror $mirror_name"
                    break
                fi
            fi
        done
    fi
else
    mirror_url="${MIRRORS[$MIRROR]}${GITHUB_URL}"
    echo "Using mirror $MIRROR for download: $mirror_url"
    curl -L -o $INSTALL_PATH $mirror_url
fi

if [ $? -ne 0 ]; then
    echo "All download methods failed"
    exit 1
fi

# Set execution permissions
chmod +x $INSTALL_PATH

echo "Installation complete!"
echo "wen-ai-cli has been installed to $INSTALL_PATH" 