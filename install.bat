@echo off
setlocal enabledelayedexpansion

set REPO=f-asai-monox/mcpify
set INSTALL_DIR=%LOCALAPPDATA%\mcpify
set BINARY_NAME=mcp-server-stdio.exe
set HTTP_BINARY_NAME=mcp-server-http.exe

echo Installing mcpify...

:: Detect architecture
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    set ARCH=x86_64
) else if "%PROCESSOR_ARCHITECTURE%"=="x86" (
    if defined PROCESSOR_ARCHITEW6432 (
        set ARCH=x86_64
    ) else (
        echo Unsupported architecture: %PROCESSOR_ARCHITECTURE%
        exit /b 1
    )
) else (
    echo Unsupported architecture: %PROCESSOR_ARCHITECTURE%
    exit /b 1
)

:: Get latest release version
echo Fetching latest release...
for /f "tokens=2 delims=:, " %%a in ('curl -s https://api.github.com/repos/%REPO%/releases/latest ^| findstr /C:"tag_name"') do (
    set VERSION=%%a
    set VERSION=!VERSION:"=!
)

if "%VERSION%"=="" (
    echo Error: Could not fetch latest release version
    exit /b 1
)

echo Latest version: %VERSION%

:: Set download URL
set FILENAME=mcpify_Windows_%ARCH%.zip
set URL=https://github.com/%REPO%/releases/download/%VERSION%/%FILENAME%

:: Create install directory
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

:: Download
echo Downloading %FILENAME%...
curl -L -o "%TEMP%\%FILENAME%" "%URL%"
if errorlevel 1 (
    echo Error: Failed to download file
    exit /b 1
)

:: Extract
echo Extracting...
powershell -NoProfile -ExecutionPolicy Bypass -Command "Expand-Archive -Path '%TEMP%\%FILENAME%' -DestinationPath '%TEMP%\mcpify_extract' -Force"
if errorlevel 1 (
    echo Error: Failed to extract archive
    exit /b 1
)

:: Move binaries to install directory
echo Installing to %INSTALL_DIR%...
if exist "%TEMP%\mcpify_extract\%BINARY_NAME%" (
    move /Y "%TEMP%\mcpify_extract\%BINARY_NAME%" "%INSTALL_DIR%\" >nul
)
if exist "%TEMP%\mcpify_extract\%HTTP_BINARY_NAME%" (
    move /Y "%TEMP%\mcpify_extract\%HTTP_BINARY_NAME%" "%INSTALL_DIR%\" >nul
)

:: Clean up
if exist "%TEMP%\%FILENAME%" del "%TEMP%\%FILENAME%"
if exist "%TEMP%\mcpify_extract" rmdir /s /q "%TEMP%\mcpify_extract"

:: Add to PATH if not already present
echo %PATH% | findstr /C:"%INSTALL_DIR%" >nul
if errorlevel 1 (
    echo.
    echo Adding %INSTALL_DIR% to PATH...
    setx PATH "%PATH%;%INSTALL_DIR%" >nul
    echo.
    echo PATH updated. Please restart your terminal for changes to take effect.
)

echo.
echo Successfully installed mcpify to %INSTALL_DIR%
echo.
echo To get started, restart your terminal and run:
echo   mcp-server-stdio --help
echo.

endlocal