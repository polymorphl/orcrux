#!/bin/bash

# Orcrux Release Test Script
# This script simulates the GitHub Actions workflow locally to test for issues

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

# Test environment
test_environment() {
    print_status "Testing environment..."
    
    # Check Go
    if ! command -v go &> /dev/null; then
        print_error "Go not found"
        return 1
    fi
    local go_version=$(go version | awk '{print $3}' | sed 's/go//')
    print_success "Go version: $go_version"
    
    # Check Node.js
    if ! command -v node &> /dev/null; then
        print_error "Node.js not found"
        return 1
    fi
    local node_version=$(node --version)
    print_success "Node.js version: $node_version"
    
    # Check npm
    if ! command -v npm &> /dev/null; then
        print_error "npm not found"
        return 1
    fi
    local npm_version=$(npm --version)
    print_success "npm version: $npm_version"
    
    # Check Wails CLI
    if ! command -v wails &> /dev/null; then
        print_warning "Wails CLI not found, installing..."
        go install github.com/wailsapp/wails/v2/cmd/wails@latest
    fi
    local wails_version=$(wails version)
    print_success "Wails version: $wails_version"
    
    return 0
}

# Test frontend build
test_frontend_build() {
    print_status "Testing frontend build..."
    
    cd frontend
    
    # Clean install
    print_status "Installing dependencies..."
    rm -rf node_modules package-lock.json
    npm ci
    
    # Build
    print_status "Building frontend..."
    npm run build
    
    # Check if build succeeded
    if [ ! -d "dist" ]; then
        print_error "Frontend build failed - dist directory not created"
        return 1
    fi
    
    print_success "Frontend build completed successfully"
    cd ..
    return 0
}

# Test Wails build for specific platform
test_wails_build() {
    local platform=$1
    local arch=$2
    
    print_status "Testing Wails build for $platform/$arch..."
    
    # Check cross-compilation limitations
    local current_os=$(uname | tr '[:upper:]' '[:lower:]')
    if [[ "$current_os" == "darwin" && "$platform" == "linux" ]]; then
        print_warning "Cross-compilation from macOS to Linux is not supported by Wails"
        print_warning "This test will be skipped. Use Linux or GitHub Actions for full testing."
        return 0
    fi
    
    # Clean previous builds
    rm -rf build/
    
    # Build
    local build_cmd="wails build -platform $platform/$arch -clean"
    print_status "Running: $build_cmd"
    
    if ! eval "$build_cmd"; then
        print_error "Wails build failed for $platform/$arch"
        return 1
    fi
    
    # Check if build succeeded
    if [ ! -d "build/bin" ]; then
        print_error "Build failed - build/bin directory not created"
        return 1
    fi
    
    # List build artifacts
    print_status "Build artifacts for $platform/$arch:"
    ls -la build/bin/
    
    print_success "Wails build for $platform/$arch completed successfully"
    return 0
}

# Test all platforms
test_all_platforms() {
    print_status "Testing builds for all platforms..."
    
    local platforms=("windows/amd64" "darwin/amd64")
    local failed_builds=()
    
    for platform in "${platforms[@]}"; do
        local platform_name=$(echo "$platform" | cut -d'/' -f1)
        local arch=$(echo "$platform" | cut -d'/' -f2)
        
        print_status "Testing $platform_name/$arch..."
        
        if test_wails_build "$platform_name" "$arch"; then
            print_success "✓ $platform_name/$arch build passed"
        else
            print_error "✗ $platform_name/$arch build failed"
            failed_builds+=("$platform")
        fi
        
        echo ""
    done
    
    # Summary
    if [ ${#failed_builds[@]} -eq 0 ]; then
        print_success "All platform builds passed! 🎉"
        return 0
    else
        print_error "Some builds failed: ${failed_builds[*]}"
        return 1
    fi
}

# Test artifact organization
test_artifact_organization() {
    print_status "Testing artifact organization..."
    
    # Create dist directory
    rm -rf dist/
    mkdir -p dist/
    
    # Organize artifacts
    local platforms=("windows" "darwin")
    
    for platform in "${platforms[@]}"; do
        local arch="amd64"
        local build_dir="build/orcrux-$platform-$arch"
        
        # Clean and rebuild for this platform
        rm -rf build/
        if wails build -platform "$platform/$arch" -clean; then
            mkdir -p "dist/orcrux-$platform-$arch"
            cp -r build/bin/* "dist/orcrux-$platform-$arch/"
            print_success "Organized artifacts for $platform/$arch"
        else
            print_error "Failed to build for $platform/$arch"
            return 1
        fi
    done
    
    # Show final structure
    print_status "Final artifact structure:"
    tree dist/ || ls -la dist/
    
    print_success "Artifact organization completed successfully"
    return 0
}

# Main test function
main() {
    print_status "Starting Orcrux release test..."
    echo ""
    
    local tests_passed=0
    local tests_failed=0
    
    # Test 1: Environment
    print_status "=== Test 1: Environment Check ==="
    if test_environment; then
        print_success "Environment test passed"
        ((tests_passed++))
    else
        print_error "Environment test failed"
        ((tests_failed++))
    fi
    echo ""
    
    # Test 2: Frontend Build
    print_status "=== Test 2: Frontend Build ==="
    if test_frontend_build; then
        print_success "Frontend build test passed"
        ((tests_passed++))
    else
        print_error "Frontend build test failed"
        ((tests_failed++))
    fi
    echo ""
    
    # Test 3: Platform Builds
    print_status "=== Test 3: Platform Builds ==="
    if test_all_platforms; then
        print_success "Platform builds test passed"
        ((tests_passed++))
    else
        print_error "Platform builds test failed"
        ((tests_failed++))
    fi
    echo ""
    
    # Test 4: Artifact Organization
    print_status "=== Test 4: Artifact Organization ==="
    if test_artifact_organization; then
        print_success "Artifact organization test passed"
        ((tests_passed++))
    else
        print_error "Artifact organization test failed"
        ((tests_failed++))
    fi
    echo ""
    
    # Final summary
    print_status "=== Test Summary ==="
    print_success "Tests passed: $tests_passed"
    if [ $tests_failed -gt 0 ]; then
        print_error "Tests failed: $tests_failed"
    fi
    
    if [ $tests_failed -eq 0 ]; then
        print_success "All tests passed! Your release workflow should work on GitHub. 🎉"
        print_status "Next steps:"
        print_status "1. Commit your changes"
        print_status "2. Create a tag: git tag v0.0.1"
        print_status "3. Push the tag: git push origin v0.0.1"
    else
        print_error "Some tests failed. Please fix the issues before pushing to GitHub."
        print_status "Check the error messages above for details."
    fi
}

# Run main function
main "$@"
