# Release Guide

## Building Locally

### Prerequisites
- **Go 1.23+** installed
- **gcc** installed (for CGO/SQLite)
- **Docker** (for cross-platform builds from macOS)

### Build Commands

```bash
cd App

# Build macOS DMG (native, works on Mac)
bash release.sh macos

# Build all platforms (Docker required on macOS)
bash release.sh all

# Build individual platforms
bash release.sh linux-amd64
bash release.sh linux-arm64
bash release.sh windows-amd64
```

### Output
Artifacts are placed in `releases/` at the repo root:
- `Homestyle Plate_v1.0.0_macos.dmg` — macOS Disk Image
- `Homestyle Plate_v1.0.0_linux-amd64.tar.gz` — Linux x86_64
- `Homestyle Plate_v1.0.0_windows-amd64.zip` — Windows x64

## Cross-Compilation (from macOS)

Building Linux/Windows from macOS requires Docker or a cross-compiler:

### Option 1: Docker (recommended)
```bash
# Install Docker Desktop
brew install --cask docker
open -a Docker

# Then build
bash release.sh all
```

### Option 2: Cross-compiler
```bash
# For Windows:
brew install mingw-w64

# For Linux: Not practical from macOS — use Docker instead
```

## Automated Releases (CI/CD)

Push a version tag to trigger automated builds:

```bash
# Update version
echo "1.0.1" > App/VERSION
git add App/VERSION
git commit -m "Bump version to 1.0.1"

# Tag and push
git tag v1.0.1
git push origin main --tags
```

This triggers GitHub Actions which:
1. Builds for macOS (on macOS runner)
2. Builds for Linux (on Ubuntu runner)
3. Builds for Windows (on Ubuntu runner with mingw)
4. Creates a GitHub Release with all artifacts

## Releasing a New Version

1. **Update version**: Edit `App/VERSION`
2. **Test locally**: `cd App && bash release.sh all`
3. **Verify artifacts**: Check `releases/` directory
4. **Commit**: `git add App/VERSION && git commit -m "Release v1.0.1"`
5. **Tag**: `git tag v1.0.1 && git push origin main --tags`
6. **Done**: GitHub Actions creates the release automatically
