#!/bin/bash

# Orcrux Code Signing Script
# This script handles code signing for macOS releases

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
        print_error "This script must be run on macOS for code signing"
        exit 1
    fi
    print_success "Running on macOS"
}

# Check if codesign is available
check_codesign() {
    if ! command -v codesign &> /dev/null; then
        print_error "codesign not found. Install Xcode Command Line Tools."
        exit 1
    fi
    print_success "codesign found"
}

# Generate free development certificate
generate_dev_cert() {
    print_status "Generating free development certificate..."
    
    # Try to create a basic development certificate
    if security create-keychain -p "" temp.keychain 2>/dev/null; then
        print_success "Temporary keychain created"
        
        # Try to generate a self-signed certificate
        if openssl req -new -newkey rsa:2048 -keyout temp.key -out temp.csr -subj "/CN=Mac Developer" 2>/dev/null; then
            print_success "Development certificate generated"
            print_status "You can now use: ./scripts/code-sign.sh sign 'Mac Developer'"
        else
            print_warning "Could not generate certificate automatically"
            print_status "Try opening Xcode once to generate certificates"
        fi
        
        # Clean up
        security delete-keychain temp.keychain 2>/dev/null
    else
        print_warning "Could not create temporary keychain"
        print_status "Try opening Xcode once to generate certificates"
    fi
}

# Check for developer identity
check_identity() {
    local identity=$1
    
    if [[ -z "$identity" ]]; then
        print_warning "No developer identity specified"
        print_status "Available identities:"
        security find-identity -v -p codesigning
        print_status "Usage: $0 sign <identity> [app_path]"
        print_status "Example: $0 sign 'Developer ID Application: Your Name'"
        exit 1
    fi
    
    # Verify identity exists
    if ! security find-identity -v -p codesigning | grep -q "$identity"; then
        print_error "Identity not found: $identity"
        print_status "Available identities:"
        security find-identity -v -p codesigning
        exit 1
    fi
    
    print_success "Identity verified: $identity"
}

# Sign the app
sign_app() {
    local identity=$1
    local app_path=${2:-"build/bin/orcrux.app"}
    
    print_status "Signing app: $app_path"
    
    if [[ ! -d "$app_path" ]]; then
        print_error "App not found at: $app_path"
        exit 1
    fi
    
    # Sign the app
    codesign --force --deep --sign "$identity" "$app_path"
    
    # Verify signature
    if codesign --verify --verbose "$app_path"; then
        print_success "App signed successfully"
    else
        print_error "Code signing verification failed"
        exit 1
    fi
}

# Notarize the app (requires Apple Developer Account - $99/year)
notarize_app() {
    local app_path=$1
    local apple_id=${2:-"$APPLE_ID"}
    local app_password=${3:-"$APPLE_APP_PASSWORD"}
    
    print_warning "Notarization requires a paid Apple Developer Account ($99/year)"
    print_status "For free development, use code signing only:"
    print_status "  ./scripts/code-sign.sh sign 'Mac Developer: Your Name'"
    print_status ""
    print_status "If you have a paid account, set these environment variables:"
    print_status "  APPLE_ID, APPLE_APP_PASSWORD, APPLE_TEAM_ID"
    
    if [[ -z "$apple_id" || -z "$app_password" ]]; then
        print_error "Apple Developer Account credentials required for notarization"
        return 1
    fi
    
    print_status "Notarizing app: $app_path"
    
    # Submit for notarization using modern notarytool
    local request_id=$(xcrun notarytool submit "$app_path" \
        --apple-id "$apple_id" \
        --password "$app_password" \
        --team-id "$APPLE_TEAM_ID" \
        --wait | grep "id:" | awk '{print $2}')
    
    if [[ -n "$request_id" ]]; then
        print_success "Notarization completed. Request ID: $request_id"
        print_status "App has been notarized successfully!"
    else
        print_error "Notarization failed"
        return 1
    fi
}

# Check notarization status
check_notarization() {
    local request_id=$1
    local apple_id=${2:-"$APPLE_ID"}
    local app_password=${3:-"$APPLE_APP_PASSWORD"}
    
    if [[ -z "$request_id" || -z "$apple_id" || -z "$app_password" ]]; then
        print_error "Missing required parameters for notarization check"
        return 1
    fi
    
    print_status "Checking notarization status for: $request_id"
    xcrun notarytool info "$request_id" \
        --apple-id "$apple_id" \
        --password "$app_password" \
        --team-id "$APPLE_TEAM_ID"
}

# Main function
main() {
    local command=$1
    local identity=$2
    local app_path=$3
    
    check_macos
    check_codesign
    
    case "$command" in
        "sign")
            check_identity "$identity"
            sign_app "$identity" "$app_path"
            ;;
        "notarize")
            local apple_id=${2:-"$APPLE_ID"}
            local app_password=${3:-"$APPLE_APP_PASSWORD"}
            local app_path=${4:-"orcrux-macos.dmg"}
            notarize_app "$app_path" "$apple_id" "$app_password"
            ;;
        "check")
            local request_id=$2
            local apple_id=${3:-"$APPLE_ID"}
            local app_password=${4:-"$APPLE_APP_PASSWORD"}
            check_notarization "$request_id" "$apple_id" "$app_password"
            ;;
        "list")
            print_status "Available developer identities:"
            security find-identity -v -p codesigning
            ;;
        "generate")
            generate_dev_cert
            ;;
        "help"|"--help"|"-h"|"")
            echo "Orcrux Code Signing Script"
            echo "=========================="
            echo ""
            echo "Usage: $0 <command> [options]"
            echo ""
            echo "Commands:"
            echo "  sign <identity> [app_path]     - Sign the app with developer identity"
            echo "  generate                       - Generate free development certificate"
            echo "  notarize [apple_id] [password] [app_path] - Submit app for notarization (paid)"
            echo "  check <request_id> [apple_id] [password] - Check notarization status"
            echo "  list                            - List available developer identities"
            echo "  help                           - Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0 generate                      # Generate free dev certificate"
            echo "  $0 sign 'Mac Developer: Your Name' # Sign with free certificate"
            echo "  $0 list                          # List available certificates"
            echo "  $0 notarize your@email.com app_password # Paid: Apple notarization"
            echo "  $0 check <request_id>            # Check notarization status"
            echo ""
            echo "Environment Variables:"
            echo "  APPLE_ID         - Your Apple ID email"
            echo "  APPLE_APP_PASSWORD - App-specific password"
            echo "  APPLE_TEAM_ID   - Your Apple Team ID (optional)"
            ;;
        *)
            print_error "Unknown command: $command"
            echo "Run '$0 help' for usage information"
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
