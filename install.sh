#!/bin/bash

set -e

REPO="f-asai-monox/mcpify"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
STDIO_BINARY="mcp-server-stdio"
HTTP_BINARY="mcp-server-http"

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
    
    # Install stdio binary
    if [ -f "${temp_dir}/${STDIO_BINARY}" ]; then
        if [ -w "$INSTALL_DIR" ]; then
            mv "${temp_dir}/${STDIO_BINARY}" "$INSTALL_DIR/"
        else
            echo "Permission denied. Trying with sudo..."
            sudo mv "${temp_dir}/${STDIO_BINARY}" "$INSTALL_DIR/"
        fi
        chmod +x "${INSTALL_DIR}/${STDIO_BINARY}"
        echo "Successfully installed ${STDIO_BINARY} to ${INSTALL_DIR}/${STDIO_BINARY}"
    fi
    
    # Install http binary
    if [ -f "${temp_dir}/${HTTP_BINARY}" ]; then
        if [ -w "$INSTALL_DIR" ]; then
            mv "${temp_dir}/${HTTP_BINARY}" "$INSTALL_DIR/"
        else
            sudo mv "${temp_dir}/${HTTP_BINARY}" "$INSTALL_DIR/"
        fi
        chmod +x "${INSTALL_DIR}/${HTTP_BINARY}"
        echo "Successfully installed ${HTTP_BINARY} to ${INSTALL_DIR}/${HTTP_BINARY}"
    fi
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
    
    echo ""
    echo "Installation complete!"
    echo ""
    echo "Available commands:"
    echo "  mcp-server-stdio  - MCP server with stdio transport"
    echo "  mcp-server-http   - MCP server with HTTP transport"
    echo ""
    echo "Run 'mcp-server-stdio --help' or 'mcp-server-http --help' to get started."
    
    if ! command -v "$STDIO_BINARY" &> /dev/null; then
        echo ""
        echo "Note: Make sure ${INSTALL_DIR} is in your PATH."
    fi
}

main "$@"