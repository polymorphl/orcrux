# Orcrux Build System

This document provides a comprehensive overview of the Orcrux build system, designed to generate binaries for all platforms and automate GitHub releases.

## 🏗️ System Overview

The build system consists of multiple components working together:

- **GitHub Actions Workflow** - Automated CI/CD pipeline
- **Makefile** - Cross-platform build commands
- **Build Scripts** - Platform-specific build automation
- **Version Management** - Automated version bumping and tagging
- **Release Automation** - GitHub releases with artifacts

## 🚀 Quick Start

### 1. Automated Release (Recommended)
```bash
# Bump version and create release
./scripts/version.sh bump patch
./scripts/version.sh release
```

### 2. Manual Build
```bash
# Build for all platforms
make build-all

# Or use build scripts
./scripts/build.sh        # macOS/Linux
scripts/build.bat         # Windows
```

## 📋 Available Commands

### Makefile Commands
```bash
make help          # Show all available commands
make install-deps  # Install dependencies
make build         # Build for current platform
make build-all     # Build for all platforms
make build-windows # Build for Windows (amd64)
make build-darwin  # Build for macOS (amd64)
make clean         # Clean build artifacts
make test          # Run tests
make release       # Create release package
```

### Version Management
```bash
./scripts/version.sh version           # Show current version
./scripts/version.sh bump patch        # Bump patch version
./scripts/version.sh bump minor        # Bump minor version
./scripts/version.sh bump major        # Bump major version
./scripts/version.sh release           # Create release tag
./scripts/version.sh release 1.0.0    # Create release with specific version
```

### Build Scripts
```bash
./scripts/build.sh        # Build all platforms (macOS/Linux)
scripts/build.bat         # Build all platforms (Windows)
```

## 🔄 GitHub Actions Workflow

### Trigger
- **Automatic**: Push a tag starting with `v` (e.g., `v1.0.0`)
- **Manual**: Use the workflow dispatch trigger

### Process
1. **Matrix Build**: Builds simultaneously on Ubuntu, Windows, and macOS
2. **Frontend Build**: Installs dependencies and builds React app
3. **Wails Build**: Compiles Go backend and packages with frontend
4. **Artifact Upload**: Uploads build artifacts for each platform
5. **Release Creation**: Automatically creates GitHub release with artifacts

### Platforms Supported
- **Windows (amd64)**: Windows Server 2022
- **macOS (amd64)**: macOS 13

## 🛠️ Build Process

### 1. Frontend Build
```bash
cd frontend
npm ci              # Install dependencies
npm run build       # Build production bundle
cd ..
```

### 2. Wails Build
```bash
wails build -platform <platform>/<arch> -clean
```

### 3. Artifact Organization
```
dist/
├── orcrux-windows-amd64/   # Windows binaries
└── orcrux-darwin-amd64/    # macOS binaries
```

## 📦 Release Artifacts

Each release includes:

- **Executable**: Platform-specific binary
- **Resources**: Icons, metadata, and bundled assets
- **Installers**: Platform-specific installation packages (if applicable)

### File Structure
```
orcrux-<platform>-<arch>/
├── orcrux(.exe)           # Main executable
├── Info.plist             # macOS metadata
├── iconfile.icns          # macOS icon
├── icon.ico               # Windows icon
└── Resources/             # Additional resources
```

## 🔧 Configuration

### Wails Configuration (`wails.json`)
```json
{
  "name": "orcrux",
  "outputfilename": "orcrux",
  "frontend:install": "npm install",
  "frontend:build": "npm run build"
}
```

### GitHub Actions Configuration (`.github/workflows/release.yml`)
- **Go Version**: 1.23
- **Node Version**: 18
- **Matrix Strategy**: Ubuntu, Windows, macOS
- **Artifact Retention**: 90 days

## 🚨 Troubleshooting

### Common Issues

1. **Build Fails with Dependency Errors**
   ```bash
   make clean
   make install-deps
   make build-all
   ```

2. **Cross-Compilation Issues**
   - Ensure Go is properly configured
   - Check that all required tools are installed
   - Verify platform/architecture combinations

3. **Frontend Build Fails**
   ```bash
   cd frontend
   rm -rf node_modules package-lock.json
   npm install
   npm run build
   ```

4. **Wails CLI Not Found**
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

### Debug Commands
```bash
# Check prerequisites
make check-prereqs

# Clean and rebuild
make clean
make build-all

# Check Wails version
wails version

# Check Go version
go version

# Check Node version
node --version
```

## 📚 Additional Resources

- [RELEASE.md](RELEASE.md) - Detailed release instructions
- [Wails Documentation](https://wails.io/docs/) - Framework documentation
- [Go Cross-Compilation](https://golang.org/doc/install/source#environment) - Go build guide

## 🤝 Contributing to Build System

### Adding New Platforms
1. Update matrix strategy in `.github/workflows/release.yml`
2. Add build commands to Makefile
3. Update build scripts
4. Test cross-compilation

### Modifying Build Process
1. Update relevant build scripts
2. Modify Makefile targets
3. Update GitHub Actions workflow
4. Test on all platforms

### Best Practices
- Always test builds on target platforms
- Use semantic versioning for releases
- Keep build artifacts organized
- Document any platform-specific requirements
