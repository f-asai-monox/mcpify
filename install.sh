#!/bin/bash

set -e

REPO="f-asai-monox/mcpify"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
BINARY_NAME="mcpify"

detect_os() {
    OS=""
    case "$(uname -s)" in
        Linux*)     OS=Linux;;
        Darwin*)    OS=Darwin;;
        CYGWIN*|MINGW*|MSYS*) OS=Windows;;
        *)          echo "Unsupported OS: $(uname -s)"; exit 1;;
    esac
    echo "$OS"
}

detect_arch() {
    ARCH=""
    case "$(uname -m)" in
        x86_64|amd64)   ARCH=x86_64;;
        aarch64|arm64)  ARCH=arm64;;
        *)              echo "Unsupported architecture: $(uname -m)"; exit 1;;
    esac
    echo "$ARCH"
}

get_latest_release() {
    curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
}

download_and_install() {
    local version=$1
    local os=$2
    local arch=$3
    
    local filename="${BINARY_NAME}_${os}_${arch}.tar.gz"
    if [ "$os" = "Windows" ]; then
        filename="${BINARY_NAME}_${os}_${arch}.zip"
    fi
    
    local url="https://github.com/${REPO}/releases/download/${version}/${filename}"
    
    echo "Downloading ${BINARY_NAME} ${version} for ${os}/${arch}..."
    
    local temp_dir=$(mktemp -d)
    trap "rm -rf $temp_dir" EXIT
    
    curl -L -o "${temp_dir}/${filename}" "$url"
    
    echo "Extracting..."
    if [ "$os" = "Windows" ]; then
        unzip -q "${temp_dir}/${filename}" -d "$temp_dir"
    else
        tar -xzf "${temp_dir}/${filename}" -C "$temp_dir"
    fi
    
    echo "Installing to ${INSTALL_DIR}..."
    
    if [ -w "$INSTALL_DIR" ]; then
        mv "${temp_dir}/${BINARY_NAME}" "$INSTALL_DIR/"
    else
        echo "Permission denied. Trying with sudo..."
        sudo mv "${temp_dir}/${BINARY_NAME}" "$INSTALL_DIR/"
    fi
    
    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    
    echo "Successfully installed ${BINARY_NAME} to ${INSTALL_DIR}/${BINARY_NAME}"
}

main() {
    echo "Installing mcpify..."
    
    OS=$(detect_os)
    ARCH=$(detect_arch)
    VERSION=$(get_latest_release)
    
    if [ -z "$VERSION" ]; then
        echo "Error: Could not fetch latest release version"
        exit 1
    fi
    
    download_and_install "$VERSION" "$OS" "$ARCH"
    
    if command -v "$BINARY_NAME" &> /dev/null; then
        echo "Installation complete! Run '${BINARY_NAME} --help' to get started."
    else
        echo "Installation complete! Make sure ${INSTALL_DIR} is in your PATH."
    fi
}

main "$@"