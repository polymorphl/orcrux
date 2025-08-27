#!/bin/bash

# Orcrux DMG Creation Script
# This script creates a DMG file from the macOS app bundle

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

# Check if we're on macOS
check_macos() {
    if [[ "$(uname)" != "Darwin" ]]; then
        print_error "This script must be run on macOS to create DMG files"
        exit 1
    fi
    print_success "Running on macOS"
}

# Check if hdiutil is available
check_hdiutil() {
    if ! command -v hdiutil &> /dev/null; then
        print_error "hdiutil not found. This script requires macOS system tools."
        exit 1
    fi
    print_success "hdiutil found"
}

# Create DMG from app bundle
create_dmg() {
    local app_path=$1
    local dmg_name=$2
    local volume_name=$3
    
    print_status "Creating DMG: $dmg_name"
    
    # Create temporary directory for DMG contents
    local temp_dir=$(mktemp -d)
    local app_name=$(basename "$app_path")
    
    # Copy app to temp directory
    cp -R "$app_path" "$temp_dir/"
    
    # Create DMG
    hdiutil create -volname "$volume_name" -srcfolder "$temp_dir" -ov -format UDZO "$dmg_name"
    
    # Clean up temp directory
    rm -rf "$temp_dir"
    
    print_success "DMG created: $dmg_name"
}

# Main function
main() {
    local app_path="build/bin/orcrux.app"
    local dmg_name="orcrux-macos.dmg"
    local volume_name="Orcrux"
    
    print_status "Starting DMG creation for Orcrux..."
    
    # Check prerequisites
    check_macos
    check_hdiutil
    
    # Check if app exists
    if [[ ! -d "$app_path" ]]; then
        print_error "App bundle not found at: $app_path"
        print_status "Please build the macOS app first with: wails build -platform darwin/amd64"
        exit 1
    fi
    
    # Create DMG
    create_dmg "$app_path" "$dmg_name" "$volume_name"
    
    # Show result
    print_success "DMG creation completed!"
    print_status "DMG file: $dmg_name"
    print_status "Size: $(du -h "$dmg_name" | cut -f1)"
    
    # Move to dist directory if it exists
    if [[ -d "dist" ]]; then
        mkdir -p "dist/orcrux-darwin-amd64"
        mv "$dmg_name" "dist/orcrux-darwin-amd64/"
        print_status "DMG moved to dist/orcrux-darwin-amd64/"
    fi
    
    # Also copy to build/bin for GitHub Actions artifacts
    if [[ -d "build/bin" ]]; then
        cp "$dmg_name" "build/bin/"
        print_status "DMG copied to build/bin for artifacts"
    fi
}

# Run main function
main "$@"
