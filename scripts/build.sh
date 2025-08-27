#!/bin/bash

# Orcrux Build Script
# This script builds the application for all supported platforms

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Wails CLI is installed
check_wails() {
    if ! command -v wails &> /dev/null; then
        print_error "Wails CLI not found. Installing..."
        go install github.com/wailsapp/wails/v2/cmd/wails@latest
    fi
    print_success "Wails CLI found"
}

# Build frontend
build_frontend() {
    print_status "Building frontend..."
    cd frontend
    npm ci
    npm run build
    cd ..
    print_success "Frontend built successfully"
}

# Build for specific platform
build_platform() {
    local platform=$1
    local arch=$2
    local extension=$3
    
    print_status "Building for $platform/$arch..."
    
    # Clean previous builds
    rm -rf build/
    
    # Build the application
    wails build -platform "$platform/$arch" -clean
    
    # Create output directory
    mkdir -p "dist/orcrux-$platform-$arch"
    
    # Copy build artifacts
    if [ "$platform" = "windows" ]; then
        cp -r build/bin/* "dist/orcrux-$platform-$arch/"
    elif [ "$platform" = "darwin" ]; then
        # For macOS, copy the app bundle and create DMG
        cp -r build/bin/* "dist/orcrux-$platform-$arch/"
        
        # Create DMG if we're on macOS
        if [[ "$(uname)" == "Darwin" ]]; then
            print_status "Creating DMG for macOS..."
            ./scripts/create-dmg.sh
            print_success "DMG created and added to artifacts"
        fi
    else
        cp -r build/bin/* "dist/orcrux-$platform-$arch/"
    fi
    
    print_success "Built for $platform/$arch"
}

# Main build function
main() {
    print_status "Starting Orcrux build process..."
    
    # Check prerequisites
    check_wails
    
    # Build frontend
    build_frontend
    
    # Create dist directory
    rm -rf dist/
    mkdir -p dist/
    
    # Build for all platforms
    build_platform "windows" "amd64" ".exe"
    build_platform "darwin" "amd64" ""
    
    print_success "All builds completed successfully!"
    print_status "Build artifacts are available in the 'dist/' directory"
    
    # List build artifacts
    echo ""
    print_status "Build artifacts:"
    ls -la dist/
}

# Run main function
main "$@"
