@echo off
REM Orcrux Build Script for Windows
REM This script builds the application for all supported platforms

setlocal enabledelayedexpansion

echo [INFO] Starting Orcrux build process...

REM Check if Go is installed
go version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Go is not installed or not in PATH
    exit /b 1
)

REM Check if Node.js is installed
node --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Node.js is not installed or not in PATH
    exit /b 1
)

REM Check if Wails CLI is installed
wails version >nul 2>&1
if errorlevel 1 (
    echo [INFO] Wails CLI not found. Installing...
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
)

echo [SUCCESS] Wails CLI found

REM Build frontend
echo [INFO] Building frontend...
cd frontend
call npm ci
if errorlevel 1 (
    echo [ERROR] Failed to install frontend dependencies
    exit /b 1
)

call npm run build
if errorlevel 1 (
    echo [ERROR] Failed to build frontend
    exit /b 1
)
cd ..

echo [SUCCESS] Frontend built successfully

REM Create dist directory
if exist dist rmdir /s /q dist
mkdir dist

REM Build for Windows
echo [INFO] Building for windows/amd64...
wails build -platform windows/amd64 -clean
if errorlevel 1 (
    echo [ERROR] Failed to build for Windows
    exit /b 1
)

mkdir dist\orcrux-windows-amd64
xcopy /e /i build\bin dist\orcrux-windows-amd64
echo [SUCCESS] Built for Windows

REM Build for Linux (cross-compilation)
echo [INFO] Building for linux/amd64...
wails build -platform linux/amd64 -clean
if errorlevel 1 (
    echo [ERROR] Failed to build for Linux
    exit /b 1
)

mkdir dist\orcrux-linux-amd64
xcopy /e /i build\bin dist\orcrux-linux-amd64
echo [SUCCESS] Built for Linux

REM Build for macOS (cross-compilation)
echo [INFO] Building for darwin/amd64...
wails build -platform darwin/amd64 -clean
if errorlevel 1 (
    echo [ERROR] Failed to build for macOS
    exit /b 1
)

mkdir dist\orcrux-darwin-amd64
xcopy /e /i build\bin dist\orcrux-darwin-amd64
echo [SUCCESS] Built for macOS

echo [SUCCESS] All builds completed successfully!
echo [INFO] Build artifacts are available in the 'dist' directory

echo.
echo [INFO] Build artifacts:
dir dist

pause
