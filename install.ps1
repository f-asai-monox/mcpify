# PowerShell installation script for mcpify

$ErrorActionPreference = "Stop"

$REPO = "f-asai-monox/mcpify"
$INSTALL_DIR = "$env:LOCALAPPDATA\mcpify"
$BINARY_NAME = "mcp-server-stdio.exe"
$HTTP_BINARY_NAME = "mcp-server-http.exe"

Write-Host "Installing mcpify..." -ForegroundColor Green

# Detect architecture
$ARCH = ""
if ([Environment]::Is64BitOperatingSystem) {
    $ARCH = "x86_64"
} else {
    Write-Host "Error: 32-bit Windows is not supported" -ForegroundColor Red
    exit 1
}

# Get latest release version
Write-Host "Fetching latest release..." -ForegroundColor Yellow
try {
    $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO/releases/latest"
    $VERSION = $release.tag_name
} catch {
    Write-Host "Error: Could not fetch latest release version" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
}

Write-Host "Latest version: $VERSION" -ForegroundColor Cyan

# Set download URL
$FILENAME = "mcpify_Windows_$ARCH.zip"
$URL = "https://github.com/$REPO/releases/download/$VERSION/$FILENAME"

# Create install directory
if (!(Test-Path $INSTALL_DIR)) {
    New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
}

# Download
$tempFile = "$env:TEMP\$FILENAME"
Write-Host "Downloading $FILENAME..." -ForegroundColor Yellow
try {
    Invoke-WebRequest -Uri $URL -OutFile $tempFile -UseBasicParsing
} catch {
    Write-Host "Error: Failed to download file" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
}

# Extract
Write-Host "Extracting..." -ForegroundColor Yellow
$extractPath = "$env:TEMP\mcpify_extract"
try {
    Expand-Archive -Path $tempFile -DestinationPath $extractPath -Force
} catch {
    Write-Host "Error: Failed to extract archive" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    Remove-Item $tempFile -Force -ErrorAction SilentlyContinue
    exit 1
}

# Move binaries to install directory
Write-Host "Installing to $INSTALL_DIR..." -ForegroundColor Yellow
if (Test-Path "$extractPath\$BINARY_NAME") {
    Move-Item -Path "$extractPath\$BINARY_NAME" -Destination "$INSTALL_DIR\" -Force
}
if (Test-Path "$extractPath\$HTTP_BINARY_NAME") {
    Move-Item -Path "$extractPath\$HTTP_BINARY_NAME" -Destination "$INSTALL_DIR\" -Force
}

# Clean up
Remove-Item $tempFile -Force -ErrorAction SilentlyContinue
Remove-Item $extractPath -Recurse -Force -ErrorAction SilentlyContinue

# Add to PATH if not already present
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$INSTALL_DIR*") {
    Write-Host "`nAdding $INSTALL_DIR to PATH..." -ForegroundColor Yellow
    [Environment]::SetEnvironmentVariable(
        "Path",
        "$currentPath;$INSTALL_DIR",
        "User"
    )
    Write-Host "PATH updated. Please restart your terminal for changes to take effect." -ForegroundColor Green
}

Write-Host "`nSuccessfully installed mcpify to $INSTALL_DIR" -ForegroundColor Green
Write-Host "`nTo get started, restart your terminal and run:" -ForegroundColor Cyan
Write-Host "  mcp-server-stdio --help" -ForegroundColor White
Write-Host ""