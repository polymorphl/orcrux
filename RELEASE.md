# Orcrux Release Guide

This guide explains how to create releases for Orcrux, including building binaries for all platforms and publishing them to GitHub.

## Prerequisites

Before creating a release, ensure you have:

- [Go 1.23+](https://golang.org/dl/) installed
- [Node.js 18+](https://nodejs.org/) installed
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) installed
- Access to the GitHub repository with push permissions

## Release Process

### 1. Prepare for Release

1. **Update version information**:
   - Update version in `frontend/package.json` if needed
   - Ensure all tests pass: `make test`
   - Commit any pending changes

2. **Create and push a new tag**:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

### 2. Automated Release (Recommended)

The GitHub Actions workflow will automatically:
- Build binaries for all platforms (Linux, Windows, macOS)
- Create a GitHub release with the tag
- Upload all build artifacts

**To trigger a release**:
1. Push a new tag: `git push origin v1.0.0`
2. The workflow will automatically start
3. Monitor the build progress in the Actions tab
4. The release will be created automatically

### 3. Manual Release

If you prefer to build manually or need to troubleshoot:

#### Using Makefile (Recommended)
```bash
# Build for all platforms
make build-all

# Create release package
make release
```

#### Using Build Scripts
```bash
# On macOS/Linux
chmod +x scripts/build.sh
./scripts/build.sh

# On Windows
scripts/build.bat
```

#### Using Wails CLI directly
```bash
# Build frontend first
cd frontend && npm run build && cd ..

# Build for each platform
wails build -platform linux/amd64 -clean
wails build -platform windows/amd64 -clean
wails build -platform darwin/amd64 -clean
```

### 4. Release Artifacts

Each release should include:

- **Linux (amd64)**: `orcrux-linux-amd64/` directory
- **Windows (amd64)**: `orcrux-windows-amd64/` directory  
- **macOS (amd64)**: `orcrux-darwin-amd64/` directory

### 5. GitHub Release

1. Go to the repository's Releases page
2. Click "Create a new release"
3. Select the tag you created
4. Add release notes describing:
   - New features
   - Bug fixes
   - Breaking changes
   - Installation instructions
5. Upload the build artifacts from the `dist/` directory
6. Publish the release

## Build Commands Reference

### Makefile Commands
```bash
make help          # Show available commands
make install-deps  # Install dependencies
make build         # Build for current platform
make build-all     # Build for all platforms
make clean         # Clean build artifacts
make test          # Run tests
make release       # Create release package
```

### Wails Commands
```bash
wails build                    # Build for current platform
wails build -platform linux/amd64    # Build for Linux
wails build -platform windows/amd64  # Build for Windows
wails build -platform darwin/amd64   # Build for macOS
wails build -clean            # Clean build before building
wails dev                     # Start development mode
```

## Troubleshooting

### Common Issues

1. **Build fails with dependency errors**:
   ```bash
   make clean
   make install-deps
   make build-all
   ```

2. **Cross-compilation issues**:
   - Ensure Go is properly configured for cross-compilation
   - Check that all required tools are installed

3. **Frontend build fails**:
   ```bash
   cd frontend
   rm -rf node_modules package-lock.json
   npm install
   npm run build
   ```

4. **Wails CLI not found**:
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

### Platform-Specific Notes

- **Linux**: Builds should work on most Linux distributions
- **Windows**: Cross-compilation from Linux/macOS works well
- **macOS**: Universal builds (Intel + Apple Silicon) are supported

## Release Checklist

- [ ] All tests pass
- [ ] Version numbers updated
- [ ] Changelog updated
- [ ] Code committed and pushed
- [ ] Tag created and pushed
- [ ] GitHub Actions workflow completed successfully
- [ ] Release notes written
- [ ] Release published on GitHub
- [ ] Release artifacts verified

## Support

If you encounter issues during the release process:

1. Check the GitHub Actions logs for detailed error information
2. Review the troubleshooting section above
3. Open an issue on GitHub with detailed error messages
4. Check the [Wails documentation](https://wails.io/docs/) for build issues
