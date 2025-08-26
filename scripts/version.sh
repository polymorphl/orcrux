#!/bin/bash

# Orcrux Version Management Script
# This script helps manage version numbers and create releases

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

# Get current version from package.json
get_current_version() {
    local version
    version=$(grep '"version"' frontend/package.json | sed 's/.*"version": "\([^"]*\)".*/\1/')
    echo "$version"
}

# Update version in package.json
update_version() {
    local new_version=$1
    local platform=$(uname)
    
    if [[ "$platform" == "Darwin" ]]; then
        # macOS
        sed -i '' "s/\"version\": \"[^\"]*\"/\"version\": \"$new_version\"/" frontend/package.json
    else
        # Linux
        sed -i "s/\"version\": \"[^\"]*\"/\"version\": \"$new_version\"/" frontend/package.json
    fi
}

# Show current version
show_version() {
    local version=$(get_current_version)
    print_status "Current version: $version"
}

# Bump version
bump_version() {
    local bump_type=$1
    local current_version=$(get_current_version)
    local new_version
    
    if [[ "$bump_type" == "major" ]]; then
        new_version=$(echo "$current_version" | awk -F. '{print $1+1 ".0.0"}')
    elif [[ "$bump_type" == "minor" ]]; then
        new_version=$(echo "$current_version" | awk -F. '{print $1 "." $2+1 ".0"}')
    elif [[ "$bump_type" == "patch" ]]; then
        new_version=$(echo "$current_version" | awk -F. '{print $1 "." $2 "." $3+1}')
    else
        print_error "Invalid bump type. Use: major, minor, or patch"
        exit 1
    fi
    
    print_status "Bumping version from $current_version to $new_version"
    update_version "$new_version"
    print_success "Version updated to $new_version"
}

# Create and push tag
create_tag() {
    local version=$1
    local tag="v$version"
    
    print_status "Creating tag: $tag"
    
    # Check if tag already exists
    if git tag -l | grep -q "^$tag$"; then
        print_warning "Tag $tag already exists"
        read -p "Do you want to delete and recreate it? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            git tag -d "$tag"
            git push origin ":refs/tags/$tag" 2>/dev/null || true
        else
            print_error "Tag creation cancelled"
            exit 1
        fi
    fi
    
    # Create tag
    git add frontend/package.json
    git commit -m "Bump version to $version" || true
    git tag "$tag"
    
    print_status "Pushing tag to remote..."
    git push origin "$tag"
    
    print_success "Tag $tag created and pushed successfully!"
    print_status "GitHub Actions will now build and release automatically"
}

# Show help
show_help() {
    echo "Orcrux Version Management Script"
    echo "================================"
    echo ""
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  version                    - Show current version"
    echo "  bump [major|minor|patch]  - Bump version number"
    echo "  release [version]         - Create release tag and push"
    echo "  help                      - Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 version                    # Show current version"
    echo "  $0 bump patch                 # Bump patch version (0.0.1 -> 0.0.2)"
    echo "  $0 bump minor                 # Bump minor version (0.1.0 -> 0.2.0)"
    echo "  $0 bump major                 # Bump major version (1.0.0 -> 2.0.0)"
    echo "  $0 release                    # Create release with current version"
    echo "  $0 release 1.0.0             # Create release with specific version"
    echo ""
}

# Main function
main() {
    case "$1" in
        "version")
            show_version
            ;;
        "bump")
            if [[ -z "$2" ]]; then
                print_error "Bump type required. Use: major, minor, or patch"
                exit 1
            fi
            bump_version "$2"
            ;;
        "release")
            if [[ -n "$2" ]]; then
                # Use specified version
                update_version "$2"
                create_tag "$2"
            else
                # Use current version
                local version=$(get_current_version)
                create_tag "$version"
            fi
            ;;
        "help"|"--help"|"-h"|"")
            show_help
            ;;
        *)
            print_error "Unknown command: $1"
            show_help
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
