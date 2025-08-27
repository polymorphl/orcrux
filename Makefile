# Orcrux Makefile
# Provides easy commands for building and managing the application

.PHONY: help build clean build-all build-linux build-windows build-darwin test install-deps release

# Default target
help:
	@echo "Orcrux Build System"
	@echo "==================="
	@echo ""
	@echo "Available commands:"
	@echo "  help          - Show this help message"
	@echo "  install-deps  - Install all dependencies"
	@echo "  build         - Build for current platform"
	@echo "  build-all     - Build for all platforms"

	@echo "  build-windows - Build for Windows (amd64)"
	@echo "  build-darwin  - Build for macOS (amd64)"
	@echo "  create-dmg    - Create macOS DMG installer"
	@echo "  sign-app      - Code sign macOS app"
	@echo "  notarize      - Submit app for Apple notarization"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  release       - Create release package"
	@echo ""

# Install dependencies
install-deps:
	@echo "Installing dependencies..."
	@cd frontend && npm ci
	@go mod download
	@go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Build frontend
build-frontend:
	@echo "Building frontend..."
	@cd frontend && npm run build

# Build for current platform
build: build-frontend
	@echo "Building for current platform..."
	@wails build -clean

# Build for all platforms
build-all: build-frontend
	@echo "Building for all platforms..."
	@mkdir -p dist
	@wails build -platform windows/amd64 -clean
	@mkdir -p dist/orcrux-windows-amd64
	@cp -r build/bin/* dist/orcrux-windows-amd64/
	@wails build -platform darwin/amd64 -clean
	@mkdir -p dist/orcrux-darwin-amd64
	@cp -r build/bin/* dist/orcrux-darwin-amd64/
	@echo "All builds completed successfully!"

# Build for Windows
build-windows: build-frontend
	@echo "Building for Windows..."
	@wails build -platform windows/amd64 -clean

# Build for macOS
build-darwin: build-frontend
	@echo "Building for macOS..."
	@wails build -platform darwin/amd64 -clean
	@echo "Creating DMG for macOS..."
	@./scripts/create-dmg.sh

# Create DMG from existing app bundle
create-dmg:
	@echo "Creating DMG from existing macOS app..."
	@./scripts/create-dmg.sh

# Code sign the macOS app
sign-app:
	@echo "Code signing macOS app..."
	@./scripts/code-sign.sh help
	@echo ""
	@echo "To sign with your identity, run:"
	@echo "  ./scripts/code-sign.sh sign 'Developer ID Application: Your Name'"

# Notarize the app with Apple
notarize:
	@echo "Submitting app for Apple notarization..."
	@./scripts/code-sign.sh help
	@echo ""
	@echo "To notarize, run:"
	@echo "  ./scripts/code-sign.sh notarize your@email.com app_password"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf build/
	@rm -rf dist/
	@rm -rf frontend/dist/

# Run tests
test:
	@echo "Running tests..."
	@go test ./...
	@cd frontend && npm test

# Create release package
release: build-all
	@echo "Creating release package..."
	@cd dist && tar -czf orcrux-release-$(shell date +%Y%m%d-%H%M%S).tar.gz *
	@echo "Release package created in dist/"

# Development mode
dev:
	@echo "Starting development mode..."
	@wails dev

# Install Wails CLI
install-wails:
	@echo "Installing Wails CLI..."
	@go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Check prerequisites
check-prereqs:
	@echo "Checking prerequisites..."
	@command -v go >/dev/null 2>&1 || { echo "Go is required but not installed. Aborting." >&2; exit 1; }
	@command -v node >/dev/null 2>&1 || { echo "Node.js is required but not installed. Aborting." >&2; exit 1; }
	@command -v npm >/dev/null 2>&1 || { echo "npm is required but not installed. Aborting." >&2; exit 1; }
	@echo "All prerequisites are satisfied."
